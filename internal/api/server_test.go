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
	"agent-runtime/internal/tools"
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
	manager := jobs.NewManager(jobs.Options{
		Policies: jobs.StaticPolicyStore{
			"token-1": policy.Policy{
				SubjectID:                 "service-account:test",
				TenantID:                  "team-a",
				AllowedTools:              []string{"codex"},
				AllowedWorkspaces:         []string{"repo-*"},
				AllowedCredentialProfiles: []string{"team-default"},
				MaxJobDuration:            time.Minute,
			},
		},
		Tools: registry,
		ResolveWorkspace: func(string, string) (string, error) {
			return t.TempDir(), nil
		},
		ResolveCredentialProfile: func(string, string) (string, error) {
			return t.TempDir(), nil
		},
		Executor: immediateExecutor{},
	})

	return api.NewServer(api.Options{Jobs: manager, Tools: registry})
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
