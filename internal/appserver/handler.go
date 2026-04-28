package appserver

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"agent-runtime/internal/policy"
	"agent-runtime/internal/tools"

	"github.com/gorilla/websocket"
)

type PolicyStore interface {
	Lookup(token string) (policy.Policy, bool)
}

type ToolResolver interface {
	Resolve(name string) (tools.Tool, bool)
}

type Options struct {
	Policies                 PolicyStore
	Tools                    ToolResolver
	ResolveWorkspace         func(tenantID string, workspaceID string) (string, error)
	ResolveCredentialProfile func(tenantID string, profileID string) (string, error)
}

type Handler struct {
	policies                 PolicyStore
	tools                    ToolResolver
	resolveWorkspace         func(tenantID string, workspaceID string) (string, error)
	resolveCredentialProfile func(tenantID string, profileID string) (string, error)
	upgrader                 websocket.Upgrader
}

type sessionSpec struct {
	ToolName          string
	Executable        string
	TenantID          string
	WorkspaceID       string
	CredentialProfile string
	WorkingDir        string
	Env               map[string]string
	MaxDuration       time.Duration
}

func NewHandler(options Options) *Handler {
	return &Handler{
		policies:                 options.Policies,
		tools:                    options.Tools,
		resolveWorkspace:         options.ResolveWorkspace,
		resolveCredentialProfile: options.ResolveCredentialProfile,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeHTTPError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	spec, err := h.authorize(r)
	if err != nil {
		writeHTTPError(w, http.StatusForbidden, err.Error())
		return
	}

	client, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer client.Close()

	_ = h.proxy(r.Context(), client, spec)
	_ = client.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(time.Second),
	)
}

func (h *Handler) authorize(r *http.Request) (sessionSpec, error) {
	if h.policies == nil {
		return sessionSpec{}, fmt.Errorf("policy store is not configured")
	}
	token := tokenFromRequest(r)
	if token == "" {
		return sessionSpec{}, fmt.Errorf("missing app-server token")
	}
	p, ok := h.policies.Lookup(token)
	if !ok {
		return sessionSpec{}, fmt.Errorf("token is not authorized")
	}

	query := r.URL.Query()
	toolName := strings.TrimSpace(r.PathValue("tool"))
	if toolName == "" {
		toolName = strings.TrimSpace(query.Get("tool"))
	}
	if toolName == "" {
		toolName = "codex"
	}
	if toolName != "codex" {
		return sessionSpec{}, fmt.Errorf("app-server is currently supported for codex only")
	}
	tenantID := query.Get("tenant")
	if tenantID == "" {
		tenantID = p.TenantID
	}
	workspaceID := query.Get("workspace")
	if workspaceID == "" {
		return sessionSpec{}, fmt.Errorf("workspace is required")
	}
	profileID := query.Get("credential_profile")
	if profileID == "" {
		profileID = query.Get("profile")
	}
	if profileID == "" {
		return sessionSpec{}, fmt.Errorf("credential profile is required")
	}
	if err := p.AuthorizeJob(policy.JobRequest{
		TenantID:          tenantID,
		Tool:              toolName,
		WorkspaceID:       workspaceID,
		CredentialProfile: profileID,
		RequestedDuration: p.MaxJobDuration,
		WantsTerminal:     true,
	}); err != nil {
		return sessionSpec{}, err
	}
	if h.tools == nil {
		return sessionSpec{}, fmt.Errorf("tool registry is not configured")
	}
	tool, ok := h.tools.Resolve(toolName)
	if !ok {
		return sessionSpec{}, fmt.Errorf("tool %q is not registered", toolName)
	}
	if tool.Path == "" {
		return sessionSpec{}, fmt.Errorf("tool %q has no executable path", toolName)
	}
	if h.resolveWorkspace == nil || h.resolveCredentialProfile == nil {
		return sessionSpec{}, fmt.Errorf("app-server resolvers are not configured")
	}
	workingDir, err := h.resolveWorkspace(tenantID, workspaceID)
	if err != nil {
		return sessionSpec{}, fmt.Errorf("resolve workspace: %w", err)
	}
	credentialRoot, err := h.resolveCredentialProfile(tenantID, profileID)
	if err != nil {
		return sessionSpec{}, fmt.Errorf("resolve credential profile: %w", err)
	}

	env := map[string]string{
		"HOME":                  credentialRoot,
		"XDG_CONFIG_HOME":       filepath.Join(credentialRoot, ".config"),
		"NPM_CONFIG_PREFIX":     "/data/npm-global",
		"PATH":                  tools.RuntimePath(credentialRoot, os.Getenv("PATH")),
		"AGENT_RUNTIME_TENANT":  tenantID,
		"AGENT_RUNTIME_SUBJECT": p.SubjectID,
	}
	if tool.CredentialEnv != "" {
		credentialPath := credentialRoot
		if tool.CredentialSubdir != "" {
			credentialPath = filepath.Join(credentialRoot, tool.CredentialSubdir)
		}
		env[tool.CredentialEnv] = credentialPath
	}

	return sessionSpec{
		ToolName:          toolName,
		Executable:        tool.Path,
		TenantID:          tenantID,
		WorkspaceID:       workspaceID,
		CredentialProfile: profileID,
		WorkingDir:        workingDir,
		Env:               env,
		MaxDuration:       p.MaxJobDuration,
	}, nil
}

func (h *Handler) proxy(parent context.Context, client *websocket.Conn, spec sessionSpec) error {
	ctx := parent
	cancel := func() {}
	if spec.MaxDuration > 0 {
		ctx, cancel = context.WithTimeout(parent, spec.MaxDuration)
	} else {
		ctx, cancel = context.WithCancel(parent)
	}
	defer cancel()

	cmd := exec.CommandContext(ctx, resolveExecutable(spec.Executable, mergeEnv(os.Environ(), spec.Env)), "app-server", "--listen", "ws://127.0.0.1:0", "--session-source", "agent-runtime")
	cmd.Dir = spec.WorkingDir
	cmd.Env = mergeEnv(os.Environ(), spec.Env)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	listenURL := make(chan string, 1)
	logs := make(chan string, 64)
	if err := cmd.Start(); err != nil {
		return err
	}
	defer func() {
		cancel()
		_ = cmd.Wait()
	}()

	go discard(stdout)
	go scanForListenURL(stderr, listenURL, logs)

	upstreamURL, err := waitForListenURL(ctx, listenURL, logs)
	if err != nil {
		return err
	}
	upstream, _, err := websocket.DefaultDialer.DialContext(ctx, upstreamURL, nil)
	if err != nil {
		return fmt.Errorf("connect codex app-server: %w", err)
	}
	defer upstream.Close()

	errs := make(chan error, 2)
	go copyWS(upstream, client, errs)
	go copyWS(client, upstream, errs)
	err = <-errs
	cancel()
	return err
}

func scanForListenURL(reader io.Reader, out chan<- string, logs chan<- string) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	re := regexp.MustCompile(`listening on:\s*(ws://[^\s]+)`)
	for scanner.Scan() {
		line := scanner.Text()
		select {
		case logs <- line:
		default:
		}
		if match := re.FindStringSubmatch(line); len(match) == 2 {
			select {
			case out <- match[1]:
			default:
			}
		}
	}
}

func waitForListenURL(ctx context.Context, listenURL <-chan string, logs <-chan string) (string, error) {
	timer := time.NewTimer(15 * time.Second)
	defer timer.Stop()
	recent := make([]string, 0, 8)
	for {
		select {
		case value := <-listenURL:
			if _, err := url.Parse(value); err != nil {
				return "", fmt.Errorf("invalid app-server URL %q: %w", value, err)
			}
			return value, nil
		case line := <-logs:
			if line != "" {
				recent = append(recent, line)
				if len(recent) > 8 {
					recent = recent[1:]
				}
			}
		case <-timer.C:
			return "", fmt.Errorf("codex app-server did not become ready: %s", strings.Join(recent, " | "))
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

func copyWS(dst *websocket.Conn, src *websocket.Conn, errs chan<- error) {
	for {
		messageType, data, err := src.ReadMessage()
		if err != nil {
			errs <- err
			return
		}
		if err := dst.WriteMessage(messageType, data); err != nil {
			errs <- err
			return
		}
	}
}

func discard(reader io.Reader) {
	_, _ = io.Copy(io.Discard, reader)
}

func tokenFromRequest(r *http.Request) string {
	header := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if strings.HasPrefix(header, prefix) {
		return strings.TrimSpace(strings.TrimPrefix(header, prefix))
	}
	return strings.TrimSpace(r.URL.Query().Get("token"))
}

func mergeEnv(base []string, extra map[string]string) []string {
	out := make([]string, 0, len(base)+len(extra))
	seen := make(map[string]int, len(base)+len(extra))
	for _, item := range base {
		key := envKey(item)
		seen[key] = len(out)
		out = append(out, item)
	}
	for key, value := range extra {
		item := key + "=" + value
		if index, ok := seen[key]; ok {
			out[index] = item
			continue
		}
		seen[key] = len(out)
		out = append(out, item)
	}
	return out
}

func resolveExecutable(name string, env []string) string {
	if strings.ContainsAny(name, `/\`) {
		return name
	}
	pathValue := envValue(env, "PATH")
	for _, dir := range filepath.SplitList(pathValue) {
		if dir == "" {
			dir = "."
		}
		candidate := filepath.Join(dir, name)
		info, err := os.Stat(candidate)
		if err == nil && !info.IsDir() && info.Mode()&0o111 != 0 {
			return candidate
		}
	}
	return name
}

func envValue(env []string, key string) string {
	prefix := key + "="
	for i := len(env) - 1; i >= 0; i-- {
		if strings.HasPrefix(env[i], prefix) {
			return strings.TrimPrefix(env[i], prefix)
		}
	}
	return ""
}

func envKey(item string) string {
	for i, char := range item {
		if char == '=' {
			return item[:i]
		}
	}
	return item
}

func writeHTTPError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = fmt.Fprintf(w, `{"error":%q}`+"\n", message)
}

func IsPortOpen(address string) bool {
	conn, err := net.DialTimeout("tcp", address, 200*time.Millisecond)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
