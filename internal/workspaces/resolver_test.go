package workspaces_test

import (
	"os"
	"path/filepath"
	"testing"

	"agent-runtime/internal/workspaces"
)

func TestResolverCreatesTenantWorkspace(t *testing.T) {
	root := t.TempDir()
	resolver := workspaces.NewResolver(root)

	got, err := resolver.ResolveWorkspace("team-a", "repo-main")
	if err != nil {
		t.Fatalf("resolve workspace: %v", err)
	}

	realRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(realRoot, "team-a", "workspaces", "repo-main")
	if got != want {
		t.Fatalf("workspace path mismatch\nwant: %s\n got: %s", want, got)
	}
	if info, err := os.Stat(got); err != nil || !info.IsDir() {
		t.Fatalf("expected workspace dir to exist, info=%#v err=%v", info, err)
	}
}

func TestResolverRejectsUnsafeWorkspaceID(t *testing.T) {
	resolver := workspaces.NewResolver(t.TempDir())

	if _, err := resolver.ResolveWorkspace("team-a", "../repo"); err == nil {
		t.Fatal("expected unsafe workspace id to be rejected")
	}
}
