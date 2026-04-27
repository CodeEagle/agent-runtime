package workspaces_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"agent-runtime/internal/workspaces"
)

func TestResolveWithinAllowsNormalRelativePath(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "repo", "src")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := workspaces.ResolveWithin(root, "repo/../repo/src")
	if err != nil {
		t.Fatalf("expected path to resolve, got %v", err)
	}

	want, err := filepath.EvalSymlinks(nested)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("resolved path mismatch\nwant: %s\n got: %s", want, got)
	}
}

func TestResolveWithinRejectsRelativeEscape(t *testing.T) {
	root := t.TempDir()

	_, err := workspaces.ResolveWithin(root, "../outside")
	if err == nil {
		t.Fatal("expected escape to be rejected")
	}
	if !strings.Contains(err.Error(), "escapes") {
		t.Fatalf("expected escape error, got %v", err)
	}
}

func TestResolveWithinRejectsSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("secret"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "link")); err != nil {
		t.Fatal(err)
	}

	_, err := workspaces.ResolveWithin(root, "link/secret.txt")
	if err == nil {
		t.Fatal("expected symlink escape to be rejected")
	}
	if !strings.Contains(err.Error(), "escapes") {
		t.Fatalf("expected escape error, got %v", err)
	}
}
