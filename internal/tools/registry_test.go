package tools_test

import (
	"os"
	"testing"

	"agent-runtime/internal/tools"
)

func TestRegistryUpsertAndDelete(t *testing.T) {
	registry := tools.NewRegistry(nil)

	if err := registry.Upsert(tools.Tool{Name: "codex", Path: "/usr/bin/codex"}); err != nil {
		t.Fatalf("upsert tool: %v", err)
	}
	if _, ok := registry.Resolve("codex"); !ok {
		t.Fatalf("expected tool to resolve after upsert")
	}
	if err := registry.Delete("codex"); err != nil {
		t.Fatalf("delete tool: %v", err)
	}
	if _, ok := registry.Resolve("codex"); ok {
		t.Fatalf("expected tool to be removed")
	}
}

func TestRegistryRejectsUnsafeTool(t *testing.T) {
	registry := tools.NewRegistry(nil)

	if err := registry.Upsert(tools.Tool{Name: "../codex", Path: "/usr/bin/codex"}); err == nil {
		t.Fatalf("expected unsafe name to be rejected")
	}
	if err := registry.Upsert(tools.Tool{Name: "codex"}); err == nil {
		t.Fatalf("expected empty path to be rejected")
	}
}

func TestPersistentRegistryLoadsSavedTools(t *testing.T) {
	storePath := t.TempDir() + "/registry.json"

	registry, err := tools.NewPersistentRegistry(nil, storePath)
	if err != nil {
		t.Fatalf("create persistent registry: %v", err)
	}
	if err := registry.Upsert(tools.Tool{Name: "codex", Path: "/usr/bin/codex", Version: "test"}); err != nil {
		t.Fatalf("upsert tool: %v", err)
	}

	reloaded, err := tools.NewPersistentRegistry(nil, storePath)
	if err != nil {
		t.Fatalf("reload persistent registry: %v", err)
	}
	tool, ok := reloaded.Resolve("codex")
	if !ok {
		t.Fatalf("expected persisted tool to resolve")
	}
	if tool.Version != "test" {
		t.Fatalf("unexpected persisted tool: %#v", tool)
	}
}

func TestPersistentRegistryIgnoresLegacyEnvPlaceholder(t *testing.T) {
	storePath := t.TempDir() + "/registry.json"
	if err := os.WriteFile(storePath, []byte(`{
		"tools": [
			{"name": "codex", "path": "/usr/bin/env", "version": "placeholder", "credential_env": "CODEX_HOME", "credential_subdir": ".codex"}
		]
	}`), 0o600); err != nil {
		t.Fatal(err)
	}

	reloaded, err := tools.NewPersistentRegistry([]tools.Tool{
		{Name: "codex", Path: "codex", Version: "official", CredentialEnv: "CODEX_HOME", CredentialSubdir: ".codex"},
	}, storePath)
	if err != nil {
		t.Fatalf("reload persistent registry: %v", err)
	}
	tool, ok := reloaded.Resolve("codex")
	if !ok {
		t.Fatalf("expected tool to resolve")
	}
	if tool.Path != "codex" || tool.Version != "official" {
		t.Fatalf("expected default tool to replace legacy placeholder, got %#v", tool)
	}
}
