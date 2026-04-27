package tenants_test

import (
	"testing"
	"time"

	"agent-runtime/internal/policy"
	"agent-runtime/internal/tenants"
)

func TestStoreListsTenantSummariesWithoutTokens(t *testing.T) {
	store := tenants.NewStore(map[string]policy.Policy{
		"token-a": {
			SubjectID:                 "service:a",
			TenantID:                  "team-a",
			AllowedTools:              []string{"codex"},
			AllowedWorkspaces:         []string{"repo-*"},
			AllowedCredentialProfiles: []string{"default"},
			AllowTerminal:             true,
			MaxJobDuration:            time.Minute,
		},
		"token-b": {
			SubjectID:                 "service:b",
			TenantID:                  "team-a",
			AllowedTools:              []string{"claude", "codex"},
			AllowedWorkspaces:         []string{"docs"},
			AllowedCredentialProfiles: []string{"default", "ops"},
		},
	})

	got := store.List()
	if len(got) != 1 {
		t.Fatalf("expected one tenant, got %#v", got)
	}
	if got[0].ID != "team-a" {
		t.Fatalf("unexpected tenant id: %#v", got[0])
	}
	if !got[0].AllowTerminal {
		t.Fatalf("expected terminal to be allowed when any token permits it: %#v", got[0])
	}
	assertStrings(t, got[0].AllowedTools, []string{"claude", "codex"})
	assertStrings(t, got[0].WorkspacePatterns, []string{"docs", "repo-*"})
	assertStrings(t, got[0].CredentialProfiles, []string{"default", "ops"})
}

func assertStrings(t *testing.T, got []string, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("expected %#v, got %#v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected %#v, got %#v", want, got)
		}
	}
}
