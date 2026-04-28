package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"agent-runtime/internal/config"
)

func TestLoadAppliesDefaultsAndParsesTokens(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	raw := `{
		"server": {"address": "127.0.0.1:9090"},
		"storage": {"data_dir": "/tmp/agent-runtime-test"},
		"tools": [
			{"name": "codex", "path": "/usr/bin/codex", "version": "test", "credential_env": "CODEX_HOME", "credential_subdir": ".codex"}
		],
		"tokens": [
			{
				"token": "secret-token",
				"subject": "service-account:cc-connect",
				"tenant": "team-a",
				"allowed_tools": ["codex"],
				"allowed_workspaces": ["repo-*"],
				"allowed_credential_profiles": ["team-default"],
				"max_job_seconds": 600
			}
		],
		"users": [
			{"username": "admin", "password": "admin", "token": "secret-token"}
		]
	}`
	if err := os.WriteFile(path, []byte(raw), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.Server.Address != "127.0.0.1:9090" {
		t.Fatalf("unexpected server address %q", cfg.Server.Address)
	}
	if cfg.Storage.ToolsDir != "/tmp/agent-runtime-test/tools" {
		t.Fatalf("expected derived tools dir, got %q", cfg.Storage.ToolsDir)
	}
	if cfg.Storage.TenantsDir != "/tmp/agent-runtime-test/tenants" {
		t.Fatalf("expected derived tenants dir, got %q", cfg.Storage.TenantsDir)
	}
	if len(cfg.Tools) != 1 || cfg.Tools[0].Name != "codex" {
		t.Fatalf("unexpected tools: %#v", cfg.Tools)
	}
	if len(cfg.Tokens) != 1 {
		t.Fatalf("expected one token, got %#v", cfg.Tokens)
	}
	if cfg.Tokens[0].Policy.MaxJobDuration != 10*time.Minute {
		t.Fatalf("expected max duration 10m, got %s", cfg.Tokens[0].Policy.MaxJobDuration)
	}
	if len(cfg.UserStore()) != 1 || cfg.UserStore()[0].Username != "admin" {
		t.Fatalf("expected one user login, got %#v", cfg.UserStore())
	}
}

func TestLoadRejectsRelativeDataDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"storage":{"data_dir":"relative"}}`), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := config.Load(path); err == nil {
		t.Fatal("expected relative data_dir to be rejected")
	}
}
