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

func TestStoreDeleteUserRemovesOwnedPolicy(t *testing.T) {
	store, err := tenants.NewStoreWithUsers(nil, []tenants.UserRequest{
		{
			Username:                  "Team B",
			Password:                  "secret",
			TenantID:                  "team-b",
			SubjectID:                 "tenant-user:team-b",
			AllowedTools:              []string{"codex"},
			AllowedWorkspaces:         []string{"repo-*"},
			AllowedCredentialProfiles: []string{"team-default"},
			AllowTerminal:             true,
		},
	})
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	users := store.ListUsers()
	if len(users) != 1 {
		t.Fatalf("expected one user, got %#v", users)
	}
	if got := store.List(); len(got) != 1 || got[0].ID != "team-b" || !got[0].AllowTerminal {
		t.Fatalf("expected team-b tenant before delete, got %#v", got)
	}

	if err := store.DeleteUser(users[0].ID); err != nil {
		t.Fatalf("delete user: %v", err)
	}
	if got := store.ListUsers(); len(got) != 0 {
		t.Fatalf("expected no users after delete, got %#v", got)
	}
	if got := store.List(); len(got) != 0 {
		t.Fatalf("expected deleted user's tenant policy to be removed, got %#v", got)
	}
}

func TestStoreDeleteUserKeepsSharedPolicy(t *testing.T) {
	store, err := tenants.NewStoreWithUsers(map[string]policy.Policy{
		"shared-token": {
			SubjectID:     "tenant-user:shared",
			TenantID:      "team-shared",
			Role:          "tenant",
			AllowTerminal: true,
		},
	}, []tenants.UserRequest{
		{Username: "one", Password: "secret", Token: "shared-token"},
		{Username: "two", Password: "secret", Token: "shared-token"},
	})
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	users := store.ListUsers()
	if len(users) != 2 {
		t.Fatalf("expected two users, got %#v", users)
	}

	if err := store.DeleteUser(users[0].ID); err != nil {
		t.Fatalf("delete user: %v", err)
	}
	if got := store.List(); len(got) != 1 || got[0].ID != "team-shared" {
		t.Fatalf("expected shared policy to remain, got %#v", got)
	}
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
