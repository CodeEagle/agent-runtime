package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"agent-runtime/internal/api"
	"agent-runtime/internal/files"
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
		t.Fatalf("expected terminal and inline CLI manager, got %q", body)
	}
	if !strings.Contains(body, "Installed CLIs") || strings.Contains(body, "Repositories") || strings.Contains(body, `data-manager-tab`) || strings.Contains(body, `repositories-panel`) {
		t.Fatalf("expected CLI manager to show a single installed CLI list, got %q", body)
	}
	if strings.Contains(body, `data-view="tools-view"`) || strings.Contains(body, `id="tools-view"`) {
		t.Fatalf("expected CLI manager to live inside terminal view, got %q", body)
	}
	if !strings.Contains(body, "/api/files") || !strings.Contains(body, "File Explorer") {
		t.Fatalf("expected UI to include tenant file explorer, got %q", body)
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
	for _, marker := range []string{
		"https://claude.ai/favicon.ico",
		"https://avatars.githubusercontent.com/u/14957082",
		"https://avatars.githubusercontent.com/u/161781182",
		"https://opencode.ai/favicon-96x96-v3.png",
		"https://img.alicdn.com/imgextra",
		"https://www.kimi.com/favicon.ico",
		"https://docs.qoder.com/mintlify-assets",
	} {
		if !strings.Contains(body, marker) {
			t.Fatalf("expected UI to include official CLI logo or source %q", marker)
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
	req.Header.Set("Authorization", "Bearer token-1")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", res.Code, res.Body.String())
	}

	req = httptest.NewRequest(http.MethodDelete, "/api/tools/claude", nil)
	req.Header.Set("Authorization", "Bearer token-1")
	res = httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d: %s", res.Code, res.Body.String())
	}
}

func TestServerListsTenants(t *testing.T) {
	handler := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/tenants", nil)
	req.Header.Set("Authorization", "Bearer token-1")
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

func TestServerScopesTenantsByRole(t *testing.T) {
	registry := tools.NewRegistry(nil)
	tenantStore := tenants.NewStore(map[string]policy.Policy{
		"admin-token": {
			SubjectID:                 "admin:test",
			TenantID:                  "team-a",
			Role:                      "admin",
			AllowedCredentialProfiles: []string{"*"},
			AllowTerminal:             true,
		},
		"team-a-token": {
			SubjectID:                 "tenant:team-a",
			TenantID:                  "team-a",
			Role:                      "tenant",
			AllowedCredentialProfiles: []string{"default"},
		},
		"team-b-token": {
			SubjectID:                 "tenant:team-b",
			TenantID:                  "team-b",
			Role:                      "tenant",
			AllowedCredentialProfiles: []string{"default"},
		},
	})
	handler := api.NewServer(api.Options{Tools: registry, Tenants: tenantStore})

	req := httptest.NewRequest(http.MethodGet, "/api/tenants", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("expected admin status 200, got %d: %s", res.Code, res.Body.String())
	}
	var adminBody struct {
		Tenants []tenants.Summary `json:"tenants"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &adminBody); err != nil {
		t.Fatalf("decode admin response: %v", err)
	}
	if len(adminBody.Tenants) != 2 {
		t.Fatalf("expected admin to see two tenants, got %#v", adminBody.Tenants)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/tenants", nil)
	req.Header.Set("Authorization", "Bearer team-a-token")
	res = httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("expected tenant status 200, got %d: %s", res.Code, res.Body.String())
	}
	var tenantBody struct {
		Tenants []tenants.Summary `json:"tenants"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &tenantBody); err != nil {
		t.Fatalf("decode tenant response: %v", err)
	}
	if len(tenantBody.Tenants) != 1 || tenantBody.Tenants[0].ID != "team-a" {
		t.Fatalf("expected tenant to see only team-a, got %#v", tenantBody.Tenants)
	}
	if len(tenantBody.Tenants[0].Subjects) != 1 || tenantBody.Tenants[0].Subjects[0] != "tenant:team-a" {
		t.Fatalf("expected tenant summary to expose only current subject, got %#v", tenantBody.Tenants[0].Subjects)
	}
}

func TestServerListsFilesWithinTenantBoundary(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "team-a", "workspaces", "repo-main"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "team-a", "workspaces", "repo-main", "README.md"), []byte("ok"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "team-b", "workspaces"), 0o755); err != nil {
		t.Fatal(err)
	}
	tenantStore := tenants.NewStore(map[string]policy.Policy{
		"team-a-token": {
			SubjectID: "tenant:team-a",
			TenantID:  "team-a",
			Role:      "tenant",
		},
	})
	handler := api.NewServer(api.Options{
		Tenants: tenantStore,
		Files:   files.NewExplorer(root),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/files?tenant=team-a&space=workspaces&path=/", nil)
	req.Header.Set("Authorization", "Bearer team-a-token")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", res.Code, res.Body.String())
	}
	var body files.Listing
	if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode listing: %v", err)
	}
	if len(body.Entries) != 1 || body.Entries[0].Name != "repo-main" {
		t.Fatalf("expected repo-main entry, got %#v", body.Entries)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/files?tenant=team-b&space=workspaces&path=/", nil)
	req.Header.Set("Authorization", "Bearer team-a-token")
	res = httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusForbidden {
		t.Fatalf("expected status 403 for cross-tenant access, got %d: %s", res.Code, res.Body.String())
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
			Role:                      "admin",
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
