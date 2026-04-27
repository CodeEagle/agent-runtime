package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"agent-runtime/internal/api"
	"agent-runtime/internal/jobs"
	"agent-runtime/internal/policy"
	"agent-runtime/internal/tenants"
	"agent-runtime/internal/terminal"
	"agent-runtime/internal/tools"

	"github.com/gorilla/websocket"
)

func TestServerListsTools(t *testing.T) {
	handler := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/tools", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", res.Code, res.Body.String())
	}

	var body struct {
		Tools []tools.Tool `json:"tools"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Tools) != 1 || body.Tools[0].Name != "codex" {
		t.Fatalf("unexpected tools response: %#v", body.Tools)
	}
}

func TestServerServesWebUIAtRoot(t *testing.T) {
	handler := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", res.Code, res.Body.String())
	}
	if got := res.Header().Get("Content-Type"); !strings.Contains(got, "text/html") {
		t.Fatalf("expected HTML content type, got %q", got)
	}
	body := res.Body.String()
	if !strings.Contains(body, "Agent Runtime") {
		t.Fatalf("expected UI title in response, got %q", body)
	}
	if !strings.Contains(body, "/api/tools") {
		t.Fatalf("expected UI to reference API endpoints, got %q", body)
	}
	if strings.Contains(body, "Create Job") {
		t.Fatalf("expected UI to hide job creation controls, got %q", body)
	}
	if !strings.Contains(body, "Terminal") || !strings.Contains(body, "CLI Manager") {
		t.Fatalf("expected terminal and CLI manager tabs, got %q", body)
	}
	if !strings.Contains(body, "/assets/xterm/xterm.js") {
		t.Fatalf("expected UI to load vendored xterm assets, got %q", body)
	}
	if strings.Contains(body, `id="tool-path"`) || strings.Contains(body, `id="tool-name"`) {
		t.Fatalf("expected UI to use official install sources instead of manual path registration")
	}
	for _, marker := range []string{
		"https://claude.ai/install.sh",
		"npm install -g @openai/codex",
		"npm install -g @google/gemini-cli",
		"https://opencode.ai/install",
		"https://gitee.com/iflow-ai/iflow-cli/raw/main/install.sh",
		"https://code.kimi.com/install.sh",
		"https://qoder.com/install",
	} {
		if !strings.Contains(body, marker) {
			t.Fatalf("expected UI to include official install source %q", marker)
		}
	}
}

func TestServerServesVendoredXtermAsset(t *testing.T) {
	handler := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/assets/xterm/xterm.css", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", res.Code, res.Body.String())
	}
	if got := res.Header().Get("Content-Type"); !strings.Contains(got, "text/css") {
		t.Fatalf("expected CSS content type, got %q", got)
	}
}

func TestServerAddsAndDeletesTools(t *testing.T) {
	handler := newTestServer(t)

	payload := []byte(`{
		"name": "claude",
		"path": "/usr/bin/claude",
		"version": "test",
		"credential_env": "CLAUDE_CONFIG_DIR",
		"credential_subdir": ".claude"
	}`)
	req := httptest.NewRequest(http.MethodPost, "/api/tools", bytes.NewReader(payload))
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", res.Code, res.Body.String())
	}

	req = httptest.NewRequest(http.MethodDelete, "/api/tools/claude", nil)
	res = httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d: %s", res.Code, res.Body.String())
	}
}

func TestServerListsTenants(t *testing.T) {
	handler := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/tenants", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", res.Code, res.Body.String())
	}

	var body struct {
		Tenants []tenants.Summary `json:"tenants"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Tenants) != 1 || body.Tenants[0].ID != "team-a" {
		t.Fatalf("unexpected tenants response: %#v", body.Tenants)
	}
	if !body.Tenants[0].AllowTerminal {
		t.Fatalf("expected terminal permission to be exposed: %#v", body.Tenants[0])
	}
}

func TestServerRejectsTerminalWithoutToken(t *testing.T) {
	handler := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/terminal?tenant=team-a&workspace=repo-main&credential_profile=team-default", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", res.Code, res.Body.String())
	}
}

func TestServerRunsTerminalWebSocket(t *testing.T) {
	handler := newTestServer(t)
	server := httptest.NewServer(handler)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/terminal/ws?token=token-1&tenant=team-a&workspace=repo-main&credential_profile=team-default"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial terminal websocket: %v", err)
	}
	defer conn.Close()

	if err := conn.WriteJSON(map[string]string{"type": "input", "data": "printf agent-runtime-ok\\n\r"}); err != nil {
		t.Fatalf("write terminal command: %v", err)
	}

	deadline := time.Now().Add(2 * time.Second)
	var output strings.Builder
	for time.Now().Before(deadline) {
		if err := conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond)); err != nil {
			t.Fatalf("set read deadline: %v", err)
		}
		var message struct {
			Type string `json:"type"`
			Data string `json:"data"`
		}
		if err := conn.ReadJSON(&message); err != nil {
			continue
		}
		if message.Type == "output" {
			output.WriteString(message.Data)
		}
		if strings.Contains(output.String(), "agent-runtime-ok") {
			return
		}
	}
	t.Fatalf("terminal output did not contain marker, got %q", output.String())
}

func TestServerRejectsJobWithoutBearerToken(t *testing.T) {
	handler := newTestServer(t)

	req := httptest.NewRequest(http.MethodPost, "/api/jobs", strings.NewReader(`{"tenant":"team-a"}`))
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d: %s", res.Code, res.Body.String())
	}
}

func TestServerCreatesAndFetchesAuthorizedJob(t *testing.T) {
	handler := newTestServer(t)

	payload := []byte(`{
		"tenant": "team-a",
		"tool": "codex",
		"args": ["exec", "fix tests"],
		"workspace": "repo-main",
		"credential_profile": "team-default",
		"timeout_seconds": 30
	}`)
	req := httptest.NewRequest(http.MethodPost, "/api/jobs", bytes.NewReader(payload))
	req.Header.Set("Authorization", "Bearer token-1")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusAccepted {
		t.Fatalf("expected status 202, got %d: %s", res.Code, res.Body.String())
	}

	var created jobs.Job
	if err := json.Unmarshal(res.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created job: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected job id, got %#v", created)
	}

	fetched := waitForHTTPJob(t, handler, created.ID, jobs.StatusSucceeded)
	if fetched.Tool != "codex" {
		t.Fatalf("expected codex job, got %#v", fetched)
	}
}

func newTestServer(t *testing.T) http.Handler {
	t.Helper()
	registry := tools.NewRegistry([]tools.Tool{
		{Name: "codex", Path: "/usr/bin/codex", Version: "test", CredentialEnv: "CODEX_HOME", CredentialSubdir: ".codex"},
	})
	tenantStore := tenants.NewStore(map[string]policy.Policy{
		"token-1": {
			SubjectID:                 "service-account:test",
			TenantID:                  "team-a",
			AllowedTools:              []string{"codex"},
			AllowedWorkspaces:         []string{"repo-*"},
			AllowedCredentialProfiles: []string{"team-default"},
			AllowTerminal:             true,
			MaxJobDuration:            time.Minute,
		},
	})
	resolveWorkspace := func(string, string) (string, error) {
		return t.TempDir(), nil
	}
	resolveCredentialProfile := func(string, string) (string, error) {
		return t.TempDir(), nil
	}
	manager := jobs.NewManager(jobs.Options{
		Policies:                 tenantStore,
		Tools:                    registry,
		ResolveWorkspace:         resolveWorkspace,
		ResolveCredentialProfile: resolveCredentialProfile,
		Executor:                 immediateExecutor{},
	})
	terminalHandler := terminal.NewHandler(terminal.Options{
		Policies:                 tenantStore,
		Tools:                    registry,
		ResolveWorkspace:         resolveWorkspace,
		ResolveCredentialProfile: resolveCredentialProfile,
		Shell:                    "/bin/sh",
	})

	return api.NewServer(api.Options{Jobs: manager, Tools: registry, Tenants: tenantStore, Terminal: terminalHandler})
}

type immediateExecutor struct{}

func (immediateExecutor) Run(ctx context.Context, spec jobs.ExecutionSpec, sink jobs.EventSink) jobs.ExecutionResult {
	sink.Emit(jobs.Event{Type: jobs.EventStdout, Message: "ok\n", At: time.Now()})
	return jobs.ExecutionResult{ExitCode: 0}
}

func waitForHTTPJob(t *testing.T, handler http.Handler, id string, status jobs.Status) jobs.Job {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		req := httptest.NewRequest(http.MethodGet, "/api/jobs/"+id, nil)
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		if res.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", res.Code, res.Body.String())
		}
		var job jobs.Job
		if err := json.Unmarshal(res.Body.Bytes(), &job); err != nil {
			t.Fatalf("decode fetched job: %v", err)
		}
		if job.Status == status {
			return job
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("job %s did not reach %s", id, status)
	return jobs.Job{}
}
