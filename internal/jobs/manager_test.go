package jobs_test

import (
	"context"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"agent-runtime/internal/jobs"
	"agent-runtime/internal/policy"
	"agent-runtime/internal/tools"
)

func TestManagerCreatesAndRunsAuthorizedJob(t *testing.T) {
	executor := &recordingExecutor{}
	registry := tools.NewRegistry([]tools.Tool{
		{
			Name:             "codex",
			Path:             "/usr/bin/codex",
			Version:          "test-version",
			CredentialEnv:    "CODEX_HOME",
			CredentialSubdir: ".codex",
		},
	})
	workspaceRoot := filepath.Join(t.TempDir(), "repo-main")
	credentialRoot := filepath.Join(t.TempDir(), "team-default")
	manager := jobs.NewManager(jobs.Options{
		Policies: jobs.StaticPolicyStore{
			"token-1": policy.Policy{
				SubjectID:                 "service-account:cc-connect",
				TenantID:                  "team-a",
				AllowedTools:              []string{"codex"},
				AllowedWorkspaces:         []string{"repo-*"},
				AllowedCredentialProfiles: []string{"team-default"},
				MaxJobDuration:            time.Minute,
			},
		},
		Tools: registry,
		ResolveWorkspace: func(tenantID string, workspaceID string) (string, error) {
			if tenantID != "team-a" || workspaceID != "repo-main" {
				t.Fatalf("unexpected workspace lookup %s/%s", tenantID, workspaceID)
			}
			return workspaceRoot, nil
		},
		ResolveCredentialProfile: func(tenantID string, profileID string) (string, error) {
			if tenantID != "team-a" || profileID != "team-default" {
				t.Fatalf("unexpected credential lookup %s/%s", tenantID, profileID)
			}
			return credentialRoot, nil
		},
		Executor: executor,
	})

	job, err := manager.Create(context.Background(), jobs.CreateRequest{
		Token:             "token-1",
		TenantID:          "team-a",
		Tool:              "codex",
		Args:              []string{"exec", "fix tests"},
		WorkspaceID:       "repo-main",
		CredentialProfile: "team-default",
		Env:               map[string]string{"EXTRA": "1"},
		Timeout:           30 * time.Second,
	})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}

	got := waitForJobStatus(t, manager, job.ID, jobs.StatusSucceeded)
	if got.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", got.ExitCode)
	}

	spec := executor.lastSpec(t)
	if spec.Executable != "/usr/bin/codex" {
		t.Fatalf("expected executable to come from registry, got %q", spec.Executable)
	}
	if spec.WorkingDir != workspaceRoot {
		t.Fatalf("expected workspace root %q, got %q", workspaceRoot, spec.WorkingDir)
	}
	if spec.Env["CODEX_HOME"] != filepath.Join(credentialRoot, ".codex") {
		t.Fatalf("expected CODEX_HOME to point at credential profile subdir, got %q", spec.Env["CODEX_HOME"])
	}
	if spec.Env["EXTRA"] != "1" {
		t.Fatalf("expected caller env to be preserved")
	}

	events := manager.Events(job.ID)
	if len(events) == 0 {
		t.Fatal("expected job events to be recorded")
	}
	if !containsEvent(events, jobs.EventStdout, "ok") {
		t.Fatalf("expected stdout event containing ok, got %#v", events)
	}
}

func TestManagerRejectsUnauthorizedJobBeforeExecution(t *testing.T) {
	executor := &recordingExecutor{}
	manager := jobs.NewManager(jobs.Options{
		Policies: jobs.StaticPolicyStore{
			"token-1": policy.Policy{
				SubjectID:                 "service-account:worker",
				TenantID:                  "team-a",
				AllowedTools:              []string{"codex"},
				AllowedWorkspaces:         []string{"repo-main"},
				AllowedCredentialProfiles: []string{"team-default"},
				MaxJobDuration:            time.Minute,
			},
		},
		Tools: tools.NewRegistry([]tools.Tool{
			{Name: "claude", Path: "/usr/bin/claude", Version: "test"},
		}),
		ResolveWorkspace: func(string, string) (string, error) {
			t.Fatal("workspace should not be resolved for unauthorized job")
			return "", nil
		},
		ResolveCredentialProfile: func(string, string) (string, error) {
			t.Fatal("credential profile should not be resolved for unauthorized job")
			return "", nil
		},
		Executor: executor,
	})

	_, err := manager.Create(context.Background(), jobs.CreateRequest{
		Token:             "token-1",
		TenantID:          "team-a",
		Tool:              "claude",
		WorkspaceID:       "repo-main",
		CredentialProfile: "team-default",
		Timeout:           time.Second,
	})
	if err == nil {
		t.Fatal("expected authorization error")
	}
	if !strings.Contains(err.Error(), "tool") {
		t.Fatalf("expected tool authorization error, got %v", err)
	}
	if executor.calls() != 0 {
		t.Fatalf("executor should not be called for unauthorized job")
	}
}

func TestManagerRejectsUnknownBearerToken(t *testing.T) {
	manager := jobs.NewManager(jobs.Options{
		Policies: jobs.StaticPolicyStore{},
		Tools:    tools.NewRegistry(nil),
		Executor: &recordingExecutor{},
	})

	_, err := manager.Create(context.Background(), jobs.CreateRequest{
		Token:             "missing",
		TenantID:          "team-a",
		Tool:              "codex",
		WorkspaceID:       "repo-main",
		CredentialProfile: "team-default",
		Timeout:           time.Second,
	})
	if err == nil {
		t.Fatal("expected token error")
	}
	if !strings.Contains(err.Error(), "token") {
		t.Fatalf("expected token error, got %v", err)
	}
}

type recordingExecutor struct {
	mu    sync.Mutex
	specs []jobs.ExecutionSpec
}

func (e *recordingExecutor) Run(ctx context.Context, spec jobs.ExecutionSpec, sink jobs.EventSink) jobs.ExecutionResult {
	e.mu.Lock()
	e.specs = append(e.specs, spec)
	e.mu.Unlock()

	sink.Emit(jobs.Event{Type: jobs.EventStdout, Message: "ok\n", At: time.Now()})
	return jobs.ExecutionResult{ExitCode: 0}
}

func (e *recordingExecutor) calls() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.specs)
}

func (e *recordingExecutor) lastSpec(t *testing.T) jobs.ExecutionSpec {
	t.Helper()
	e.mu.Lock()
	defer e.mu.Unlock()
	if len(e.specs) == 0 {
		t.Fatal("executor was not called")
	}
	return e.specs[len(e.specs)-1]
}

func waitForJobStatus(t *testing.T, manager *jobs.Manager, id string, status jobs.Status) jobs.Job {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		job, ok := manager.Get(id)
		if !ok {
			t.Fatalf("job %s was not found", id)
		}
		if job.Status == status {
			return job
		}
		if job.Status == jobs.StatusFailed {
			t.Fatalf("job failed unexpectedly: %#v", job)
		}
		time.Sleep(10 * time.Millisecond)
	}
	job, _ := manager.Get(id)
	t.Fatalf("job did not reach %s: %#v", status, job)
	return jobs.Job{}
}

func containsEvent(events []jobs.Event, eventType jobs.EventType, text string) bool {
	for _, event := range events {
		if event.Type == eventType && strings.Contains(event.Message, text) {
			return true
		}
	}
	return false
}
