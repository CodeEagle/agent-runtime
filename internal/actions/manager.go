package actions

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"agent-runtime/internal/policy"
	"agent-runtime/internal/tools"
)

type Status string

const (
	StatusQueued    Status = "queued"
	StatusRunning   Status = "running"
	StatusSucceeded Status = "succeeded"
	StatusFailed    Status = "failed"
	StatusTimedOut  Status = "timed_out"
	StatusCanceled  Status = "canceled"
)

type EventType string

const (
	EventStatus EventType = "status"
	EventStdout EventType = "stdout"
	EventStderr EventType = "stderr"
	EventInput  EventType = "input"
	EventExit   EventType = "exit"
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

type CreateRequest struct {
	Token             string `json:"-"`
	TenantID          string `json:"tenant"`
	Tool              string `json:"tool"`
	Action            string `json:"action"`
	WorkspaceID       string `json:"workspace"`
	CredentialProfile string `json:"credential_profile"`
	TimeoutSeconds    int    `json:"timeout_seconds,omitempty"`
}

type InputRequest struct {
	Data string `json:"data"`
}

type Action struct {
	ID                string    `json:"id"`
	TenantID          string    `json:"tenant"`
	SubjectID         string    `json:"subject"`
	Tool              string    `json:"tool"`
	Action            string    `json:"action"`
	Command           string    `json:"command"`
	WorkspaceID       string    `json:"workspace"`
	WorkingDir        string    `json:"cwd"`
	CredentialProfile string    `json:"credential_profile"`
	Status            Status    `json:"status"`
	ExitCode          int       `json:"exit_code"`
	Error             string    `json:"error,omitempty"`
	AuthURLs          []string  `json:"auth_urls,omitempty"`
	AuthCodes         []string  `json:"auth_codes,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	StartedAt         time.Time `json:"started_at,omitempty"`
	FinishedAt        time.Time `json:"finished_at,omitempty"`
	Events            []Event   `json:"events,omitempty"`
}

type Event struct {
	Type    EventType `json:"type"`
	Message string    `json:"message"`
	At      time.Time `json:"at"`
}

type Manager struct {
	policies                 PolicyStore
	tools                    ToolResolver
	resolveWorkspace         func(tenantID string, workspaceID string) (string, error)
	resolveCredentialProfile func(tenantID string, profileID string) (string, error)

	mu      sync.RWMutex
	actions map[string]Action
	events  map[string][]Event
	runners map[string]*runner
}

type runner struct {
	cancel context.CancelFunc
	stdin  io.WriteCloser
}

type spec struct {
	action            Action
	tool              tools.Tool
	env               map[string]string
	timeout           time.Duration
	command           string
	credentialRoot    string
	credentialProfile string
}

var (
	urlPattern    = regexp.MustCompile(`https?://[^\s"'<>]+`)
	codePattern   = regexp.MustCompile(`(?i)(?:enter|copy|paste|use|code|one-time code|verification code)[^A-Z0-9]{0,80}([A-Z0-9]{4,}(?:[- ][A-Z0-9]{3,})+)`)
	secretPattern = regexp.MustCompile(`(?i)(sk-[A-Za-z0-9_-]{20,}|sk-ant-[A-Za-z0-9_-]{20,}|(?:oauth|auth|access|refresh)[_-]?token[=:]\s*[A-Za-z0-9._~+/=-]{16,})`)
)

func NewManager(options Options) *Manager {
	return &Manager{
		policies:                 options.Policies,
		tools:                    options.Tools,
		resolveWorkspace:         options.ResolveWorkspace,
		resolveCredentialProfile: options.ResolveCredentialProfile,
		actions:                  make(map[string]Action),
		events:                   make(map[string][]Event),
		runners:                  make(map[string]*runner),
	}
}

func (m *Manager) Create(_ context.Context, req CreateRequest) (Action, error) {
	if req.TimeoutSeconds < 0 {
		return Action{}, fmt.Errorf("timeout_seconds must be positive")
	}
	if m.policies == nil {
		return Action{}, fmt.Errorf("policy store is not configured")
	}
	p, ok := m.policies.Lookup(req.Token)
	if !ok {
		return Action{}, fmt.Errorf("token is not authorized")
	}
	req.Action = strings.TrimSpace(req.Action)
	if req.Action == "" {
		return Action{}, fmt.Errorf("action is required")
	}
	timeout := defaultTimeout(req.Action)
	if req.TimeoutSeconds > 0 {
		timeout = time.Duration(req.TimeoutSeconds) * time.Second
	}
	if req.TimeoutSeconds == 0 && p.MaxJobDuration > 0 && timeout > p.MaxJobDuration {
		timeout = p.MaxJobDuration
	}
	if err := p.AuthorizeJob(policy.JobRequest{
		TenantID:          req.TenantID,
		Tool:              req.Tool,
		WorkspaceID:       req.WorkspaceID,
		CredentialProfile: req.CredentialProfile,
		RequestedDuration: timeout,
	}); err != nil {
		return Action{}, err
	}
	if m.tools == nil {
		return Action{}, fmt.Errorf("tool registry is not configured")
	}
	tool, ok := m.tools.Resolve(req.Tool)
	if !ok {
		return Action{}, fmt.Errorf("tool %q is not registered", req.Tool)
	}
	if m.resolveWorkspace == nil {
		return Action{}, fmt.Errorf("workspace resolver is not configured")
	}
	workingDir, err := m.resolveWorkspace(req.TenantID, req.WorkspaceID)
	if err != nil {
		return Action{}, fmt.Errorf("resolve workspace: %w", err)
	}
	if m.resolveCredentialProfile == nil {
		return Action{}, fmt.Errorf("credential resolver is not configured")
	}
	credentialRoot, err := m.resolveCredentialProfile(req.TenantID, req.CredentialProfile)
	if err != nil {
		return Action{}, fmt.Errorf("resolve credential profile: %w", err)
	}
	env := actionEnv(tool, credentialRoot, req.TenantID, p.SubjectID)
	command, err := commandFor(tool, req.Action)
	if err != nil {
		return Action{}, err
	}

	now := time.Now()
	action := Action{
		ID:                newID(),
		TenantID:          req.TenantID,
		SubjectID:         p.SubjectID,
		Tool:              tool.Name,
		Action:            req.Action,
		Command:           displayCommand(command),
		WorkspaceID:       req.WorkspaceID,
		WorkingDir:        workingDir,
		CredentialProfile: req.CredentialProfile,
		Status:            StatusQueued,
		CreatedAt:         now,
	}
	m.mu.Lock()
	m.actions[action.ID] = action
	m.events[action.ID] = []Event{{Type: EventStatus, Message: string(StatusQueued), At: now}}
	m.mu.Unlock()

	go m.run(spec{
		action:            action,
		tool:              tool,
		env:               env,
		timeout:           timeout,
		command:           command,
		credentialRoot:    credentialRoot,
		credentialProfile: req.CredentialProfile,
	})
	return action, nil
}

func (m *Manager) Get(id string) (Action, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	action, ok := m.actions[id]
	if !ok {
		return Action{}, false
	}
	action.Events = append([]Event(nil), m.events[id]...)
	return action, true
}

func (m *Manager) Input(id string, data string) error {
	m.mu.RLock()
	action, ok := m.actions[id]
	runner := m.runners[id]
	m.mu.RUnlock()
	if !ok {
		return fmt.Errorf("action not found")
	}
	if action.Status != StatusRunning || runner == nil || runner.stdin == nil {
		return fmt.Errorf("action is not accepting input")
	}
	if _, err := io.WriteString(runner.stdin, data+"\n"); err != nil {
		return fmt.Errorf("send input: %w", err)
	}
	m.appendEvent(id, Event{Type: EventInput, Message: "input sent\n", At: time.Now()})
	return nil
}

func (m *Manager) Cancel(id string) error {
	m.mu.RLock()
	_, ok := m.actions[id]
	runner := m.runners[id]
	m.mu.RUnlock()
	if !ok {
		return fmt.Errorf("action not found")
	}
	if runner == nil || runner.cancel == nil {
		return fmt.Errorf("action is not running")
	}
	runner.cancel()
	return nil
}

func (m *Manager) run(runSpec spec) {
	started := time.Now()
	m.update(runSpec.action.ID, func(action *Action) {
		action.Status = StatusRunning
		action.StartedAt = started
	})
	m.appendEvent(runSpec.action.ID, Event{Type: EventStatus, Message: string(StatusRunning), At: started})

	ctx, cancel := context.WithTimeout(context.Background(), runSpec.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/bin/sh", "-lc", runSpec.command)
	cmd.Dir = runSpec.action.WorkingDir
	cmd.Env = mergeEnv(os.Environ(), runSpec.env)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.WaitDelay = 5 * time.Second
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		m.finish(runSpec.action.ID, StatusFailed, -1, err.Error())
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		m.finish(runSpec.action.ID, StatusFailed, -1, err.Error())
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		m.finish(runSpec.action.ID, StatusFailed, -1, err.Error())
		return
	}
	if err := cmd.Start(); err != nil {
		m.finish(runSpec.action.ID, StatusFailed, -1, err.Error())
		return
	}

	m.mu.Lock()
	m.runners[runSpec.action.ID] = &runner{cancel: cancel, stdin: stdin}
	m.mu.Unlock()

	done := make(chan struct{}, 2)
	go m.scan(runSpec.action.ID, stdout, EventStdout, done)
	go m.scan(runSpec.action.ID, stderr, EventStderr, done)
	err = cmd.Wait()
	<-done
	<-done
	_ = stdin.Close()

	m.mu.Lock()
	delete(m.runners, runSpec.action.ID)
	m.mu.Unlock()

	status := StatusSucceeded
	exitCode := 0
	errorMessage := ""
	if ctx.Err() == context.DeadlineExceeded {
		status = StatusTimedOut
		exitCode = -1
		errorMessage = ctx.Err().Error()
	} else if ctx.Err() == context.Canceled {
		status = StatusCanceled
		exitCode = -1
		errorMessage = ctx.Err().Error()
	} else if err != nil {
		status = StatusFailed
		errorMessage = err.Error()
		exitCode = -1
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitCode = exitError.ExitCode()
		}
	}
	m.finish(runSpec.action.ID, status, exitCode, errorMessage)
}

func (m *Manager) scan(id string, reader io.Reader, eventType EventType, done chan<- struct{}) {
	defer func() { done <- struct{}{} }()
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		m.appendOutput(id, eventType, scanner.Text()+"\n")
	}
	if err := scanner.Err(); err != nil && !strings.Contains(err.Error(), "file already closed") {
		m.appendOutput(id, EventStderr, err.Error()+"\n")
	}
}

func (m *Manager) appendOutput(id string, eventType EventType, message string) {
	clean := sanitizeOutput(message)
	m.captureAuthHints(id, clean)
	m.appendEvent(id, Event{Type: eventType, Message: clean, At: time.Now()})
}

func (m *Manager) captureAuthHints(id string, message string) {
	urls := urlPattern.FindAllString(message, -1)
	codes := codePattern.FindAllStringSubmatch(message, -1)
	if len(urls) == 0 && len(codes) == 0 {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	action := m.actions[id]
	for _, rawURL := range urls {
		action.AuthURLs = appendUnique(action.AuthURLs, strings.TrimRight(rawURL, ".,);]"))
	}
	for _, match := range codes {
		if len(match) > 1 {
			action.AuthCodes = appendUnique(action.AuthCodes, strings.TrimSpace(match[1]))
		}
	}
	m.actions[id] = action
}

func (m *Manager) finish(id string, status Status, exitCode int, errorMessage string) {
	finished := time.Now()
	m.update(id, func(action *Action) {
		action.Status = status
		action.ExitCode = exitCode
		action.Error = sanitizeOutput(errorMessage)
		action.FinishedAt = finished
	})
	m.appendEvent(id, Event{Type: EventExit, Message: string(status), At: finished})
}

func (m *Manager) update(id string, update func(*Action)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	action := m.actions[id]
	update(&action)
	m.actions[id] = action
}

func (m *Manager) appendEvent(id string, event Event) {
	if event.At.IsZero() {
		event.At = time.Now()
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events[id] = append(m.events[id], event)
}

func defaultTimeout(action string) time.Duration {
	switch action {
	case "install":
		return 20 * time.Minute
	case "auth":
		return 15 * time.Minute
	default:
		return 2 * time.Minute
	}
}

func actionEnv(tool tools.Tool, credentialRoot string, tenantID string, subjectID string) map[string]string {
	env := map[string]string{
		"HOME":                  credentialRoot,
		"XDG_CONFIG_HOME":       filepath.Join(credentialRoot, ".config"),
		"NPM_CONFIG_PREFIX":     "/data/npm-global",
		"PATH":                  tools.RuntimePath(credentialRoot, os.Getenv("PATH")),
		"AGENT_RUNTIME_TENANT":  tenantID,
		"AGENT_RUNTIME_SUBJECT": subjectID,
	}
	if tool.CredentialEnv != "" {
		credentialPath := credentialRoot
		if tool.CredentialSubdir != "" {
			credentialPath = filepath.Join(credentialRoot, tool.CredentialSubdir)
		}
		env[tool.CredentialEnv] = credentialPath
	}
	return env
}

func commandFor(tool tools.Tool, action string) (string, error) {
	if action == "install" {
		command, ok := installCommands[tool.Name]
		if !ok {
			return "", fmt.Errorf("install is not configured for tool %q", tool.Name)
		}
		return command, nil
	}
	if action == "verify" {
		return shellQuote(tool.Path) + " --version", nil
	}
	if action == "auth" {
		if command, ok := authCommands[tool.Name]; ok {
			if strings.Contains(command, "{{exe}}") {
				command = strings.ReplaceAll(command, "{{exe}}", shellQuote(tool.Path))
			}
			return command, nil
		}
		return shellQuote(tool.Path), nil
	}
	return "", fmt.Errorf("unsupported CLI action %q", action)
}

var installCommands = map[string]string{
	"claude":   "curl -fsSL https://claude.ai/install.sh | bash",
	"codex":    "npm install -g @openai/codex@latest",
	"gemini":   "npm install -g @google/gemini-cli@latest",
	"opencode": "curl -fsSL https://opencode.ai/install | bash",
	"iflow":    "bash -c \"$(curl -fsSL https://gitee.com/iflow-ai/iflow-cli/raw/main/install.sh)\"",
	"kimi":     "curl -LsSf https://code.kimi.com/install.sh | bash",
	"qoder":    "curl -fsSL https://qoder.com/install | bash",
}

var authCommands = map[string]string{
	"claude":   "{{exe}} setup-token",
	"codex":    "{{exe}} login --device-auth",
	"gemini":   "{{exe}}",
	"opencode": "{{exe}} auth login",
	"iflow":    "{{exe}}",
	"kimi":     "{{exe}}",
	"qoder":    "{{exe}}",
}

func displayCommand(command string) string {
	return sanitizeOutput(command)
}

func sanitizeOutput(value string) string {
	if value == "" {
		return value
	}
	return secretPattern.ReplaceAllString(value, "[redacted]")
}

func shellQuote(value string) string {
	if value == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
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

func envKey(item string) string {
	for i, char := range item {
		if char == '=' {
			return item[:i]
		}
	}
	return item
}

func appendUnique(values []string, value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return values
	}
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func newID() string {
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err == nil {
		return hex.EncodeToString(bytes[:])
	}
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
