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

func TestStoreAuthenticatesUsers(t *testing.T) {
	store, err := tenants.NewStoreWithUsers(map[string]policy.Policy{
		"admin-token": {
			SubjectID:                 "user:admin",
			TenantID:                  "team-a",
			Role:                      "admin",
			AllowedTools:              []string{"codex"},
			AllowedWorkspaces:         []string{"repo-*"},
			AllowedCredentialProfiles: []string{"default"},
			AllowTerminal:             true,
			MaxJobDuration:            time.Minute,
		},
	}, []tenants.UserRequest{
		{Username: "admin", Password: "secret", Token: "admin-token"},
	})
	if err != nil {
		t.Fatalf("create store: %v", err)
	}

	token, p, ok := store.AuthenticateUser("admin", "secret")
	if !ok {
		t.Fatal("expected user authentication to succeed")
	}
	if token != "admin-token" || p.SubjectID != "user:admin" || !p.IsAdmin() {
		t.Fatalf("unexpected user session: token=%q policy=%#v", token, p)
	}
	if _, _, ok := store.AuthenticateUser("admin", "wrong"); ok {
		t.Fatal("expected wrong password to fail")
	}
}

func TestStoreDefaultsNewUserTenantPolicy(t *testing.T) {
	store, err := tenants.NewStoreWithUsers(nil, []tenants.UserRequest{
		{Username: "Team B", Password: "secret"},
	})
	if err != nil {
		t.Fatalf("create store: %v", err)
	}

	token, p, ok := store.AuthenticateUser("team b", "secret")
	if !ok {
		t.Fatal("expected user authentication to succeed")
	}
	if token == "" {
		t.Fatal("expected generated user token")
	}
	if p.TenantID != "team-b" {
		t.Fatalf("expected tenant to default from username, got %q", p.TenantID)
	}
	if p.SubjectID != "tenant-user:team-b" {
		t.Fatalf("expected subject to default from tenant, got %q", p.SubjectID)
	}
	if p.Role != "tenant" {
		t.Fatalf("expected tenant role, got %q", p.Role)
	}
	assertStrings(t, p.AllowedWorkspaces, []string{"repo-*"})
	assertStrings(t, p.AllowedCredentialProfiles, []string{"team-default"})
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
