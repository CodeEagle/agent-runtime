package credentials_test

import (
	"os"
	"path/filepath"
	"testing"

	"agent-runtime/internal/credentials"
)

func TestResolverCreatesTenantProfileHome(t *testing.T) {
	root := t.TempDir()
	resolver := credentials.NewResolver(root)

	got, err := resolver.ResolveProfile("team-a", "team-default")
	if err != nil {
		t.Fatalf("resolve profile: %v", err)
	}

	want := filepath.Join(root, "team-a", "homes", "team-default")
	if got != want {
		t.Fatalf("profile path mismatch\nwant: %s\n got: %s", want, got)
	}
	if info, err := os.Stat(got); err != nil || !info.IsDir() {
		t.Fatalf("expected profile dir to exist, info=%#v err=%v", info, err)
	}
}

func TestResolverRejectsUnsafeProfileID(t *testing.T) {
	resolver := credentials.NewResolver(t.TempDir())

	if _, err := resolver.ResolveProfile("team-a", "../personal"); err == nil {
		t.Fatal("expected unsafe profile id to be rejected")
	}
}
