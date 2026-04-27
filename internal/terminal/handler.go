package terminal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"agent-runtime/internal/policy"
	"agent-runtime/internal/tools"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

type PolicyStore interface {
	Lookup(token string) (policy.Policy, bool)
}

type ToolLister interface {
	List() []tools.Tool
}

type Options struct {
	Policies                 PolicyStore
	Tools                    ToolLister
	ResolveWorkspace         func(tenantID string, workspaceID string) (string, error)
	ResolveCredentialProfile func(tenantID string, profileID string) (string, error)
	Shell                    string
}

type Handler struct {
	policies                 PolicyStore
	tools                    ToolLister
	resolveWorkspace         func(tenantID string, workspaceID string) (string, error)
	resolveCredentialProfile func(tenantID string, profileID string) (string, error)
	shell                    string
	upgrader                 websocket.Upgrader
}

func NewHandler(options Options) *Handler {
	return &Handler{
		policies:                 options.Policies,
		tools:                    options.Tools,
		resolveWorkspace:         options.ResolveWorkspace,
		resolveCredentialProfile: options.ResolveCredentialProfile,
		shell:                    options.Shell,
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

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	h.runSession(r, conn, spec)
}

type sessionSpec struct {
	TenantID          string
	SubjectID         string
	WorkspaceID       string
	CredentialProfile string
	WorkingDir        string
	CredentialRoot    string
	Cols              int
	Rows              int
}

func (h *Handler) authorize(r *http.Request) (sessionSpec, error) {
	if h.policies == nil {
		return sessionSpec{}, fmt.Errorf("policy store is not configured")
	}
	token := tokenFromRequest(r)
	if token == "" {
		return sessionSpec{}, fmt.Errorf("missing terminal token")
	}
	p, ok := h.policies.Lookup(token)
	if !ok {
		return sessionSpec{}, fmt.Errorf("token is not authorized")
	}

	query := r.URL.Query()
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
	if err := p.AuthorizeTerminal(policy.TerminalRequest{
		TenantID:          tenantID,
		WorkspaceID:       workspaceID,
		CredentialProfile: profileID,
	}); err != nil {
		return sessionSpec{}, err
	}
	if h.resolveWorkspace == nil || h.resolveCredentialProfile == nil {
		return sessionSpec{}, fmt.Errorf("terminal resolvers are not configured")
	}
	workingDir, err := h.resolveWorkspace(tenantID, workspaceID)
	if err != nil {
		return sessionSpec{}, fmt.Errorf("resolve workspace: %w", err)
	}
	credentialRoot, err := h.resolveCredentialProfile(tenantID, profileID)
	if err != nil {
		return sessionSpec{}, fmt.Errorf("resolve credential profile: %w", err)
	}

	return sessionSpec{
		TenantID:          tenantID,
		SubjectID:         p.SubjectID,
		WorkspaceID:       workspaceID,
		CredentialProfile: profileID,
		WorkingDir:        workingDir,
		CredentialRoot:    credentialRoot,
		Cols:              clampDimension(query.Get("cols"), 120, 40, 240),
		Rows:              clampDimension(query.Get("rows"), 32, 12, 80),
	}, nil
}

func (h *Handler) runSession(r *http.Request, conn *websocket.Conn, spec sessionSpec) {
	shell := h.shell
	if shell == "" {
		shell = os.Getenv("SHELL")
	}
	if shell == "" {
		shell = "/bin/sh"
	}

	ctx := r.Context()
	cmd := exec.CommandContext(ctx, shell)
	cmd.Dir = spec.WorkingDir
	cmd.Env = h.sessionEnv(spec)

	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{
		Cols: uint16(spec.Cols),
		Rows: uint16(spec.Rows),
	})
	if err != nil {
		writeWS(conn, nil, "error", fmt.Sprintf("start shell: %v", err))
		return
	}
	defer func() {
		_ = ptmx.Close()
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		_ = cmd.Wait()
	}()

	var writeMu sync.Mutex
	go func() {
		buffer := make([]byte, 32768)
		for {
			n, readErr := ptmx.Read(buffer)
			if n > 0 {
				if err := writeWS(conn, &writeMu, "output", string(buffer[:n])); err != nil {
					return
				}
			}
			if readErr != nil {
				_ = writeWS(conn, &writeMu, "exit", "")
				_ = conn.Close()
				return
			}
		}
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return
		}
		var message wsMessage
		if err := json.Unmarshal(data, &message); err != nil {
			continue
		}
		switch message.Type {
		case "input":
			_, _ = ptmx.Write([]byte(message.Data))
		case "resize":
			cols := clampInt(message.Cols, 40, 240)
			rows := clampInt(message.Rows, 12, 80)
			_ = pty.Setsize(ptmx, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})
		}
	}
}

func (h *Handler) sessionEnv(spec sessionSpec) []string {
	overrides := map[string]string{
		"HOME":                             spec.CredentialRoot,
		"XDG_CONFIG_HOME":                  filepath.Join(spec.CredentialRoot, ".config"),
		"NPM_CONFIG_PREFIX":                "/data/npm-global",
		"PATH":                             tools.RuntimePath(spec.CredentialRoot, os.Getenv("PATH")),
		"TERM":                             "xterm-256color",
		"COLORTERM":                        "truecolor",
		"AGENT_RUNTIME_TENANT":             spec.TenantID,
		"AGENT_RUNTIME_SUBJECT":            spec.SubjectID,
		"AGENT_RUNTIME_WORKSPACE":          spec.WorkspaceID,
		"AGENT_RUNTIME_CREDENTIAL_PROFILE": spec.CredentialProfile,
	}
	_ = os.MkdirAll(overrides["XDG_CONFIG_HOME"], 0o700)
	if h.tools != nil {
		for _, tool := range h.tools.List() {
			credentialPath := spec.CredentialRoot
			if tool.CredentialSubdir != "" {
				credentialPath = filepath.Join(spec.CredentialRoot, tool.CredentialSubdir)
				_ = os.MkdirAll(credentialPath, 0o700)
			}
			if tool.CredentialEnv != "" {
				overrides[tool.CredentialEnv] = credentialPath
			}
		}
	}

	env := os.Environ()
	seen := make(map[string]bool, len(overrides))
	for i, item := range env {
		key, _, ok := strings.Cut(item, "=")
		if !ok {
			continue
		}
		if value, exists := overrides[key]; exists {
			env[i] = key + "=" + value
			seen[key] = true
		}
	}
	for key, value := range overrides {
		if !seen[key] {
			env = append(env, key+"="+value)
		}
	}
	return env
}

type wsMessage struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
	Cols int    `json:"cols,omitempty"`
	Rows int    `json:"rows,omitempty"`
}

func writeWS(conn *websocket.Conn, mu *sync.Mutex, messageType string, data string) error {
	if mu != nil {
		mu.Lock()
		defer mu.Unlock()
	}
	return conn.WriteJSON(wsMessage{Type: messageType, Data: data})
}

func tokenFromRequest(r *http.Request) string {
	if token := strings.TrimSpace(r.URL.Query().Get("token")); token != "" {
		return token
	}
	const prefix = "Bearer "
	header := r.Header.Get("Authorization")
	if strings.HasPrefix(header, prefix) {
		return strings.TrimSpace(strings.TrimPrefix(header, prefix))
	}
	return ""
}

func clampDimension(raw string, fallback int, min int, max int) int {
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return clampInt(parsed, min, max)
}

func clampInt(value int, min int, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func writeHTTPError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
