package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"agent-runtime/internal/policy"
	"agent-runtime/internal/tenants"
	"agent-runtime/internal/tools"
)

type Config struct {
	Server  ServerConfig  `json:"server"`
	Storage StorageConfig `json:"storage"`
	Tools   []tools.Tool  `json:"tools"`
	Tokens  []TokenConfig `json:"tokens"`
	Users   []UserConfig  `json:"users"`
}

type ServerConfig struct {
	Address string `json:"address"`
}

type StorageConfig struct {
	DataDir    string `json:"data_dir"`
	ToolsDir   string `json:"tools_dir"`
	TenantsDir string `json:"tenants_dir"`
}

type TokenConfig struct {
	Token                     string        `json:"token"`
	SubjectID                 string        `json:"subject"`
	TenantID                  string        `json:"tenant"`
	Role                      string        `json:"role"`
	AllowedTools              []string      `json:"allowed_tools"`
	AllowedWorkspaces         []string      `json:"allowed_workspaces"`
	AllowedCredentialProfiles []string      `json:"allowed_credential_profiles"`
	MaxJobSeconds             int           `json:"max_job_seconds"`
	Policy                    policy.Policy `json:"-"`
}

type UserConfig struct {
	Username                  string   `json:"username"`
	Password                  string   `json:"password"`
	Token                     string   `json:"token"`
	SubjectID                 string   `json:"subject"`
	TenantID                  string   `json:"tenant"`
	Role                      string   `json:"role"`
	AllowedTools              []string `json:"allowed_tools"`
	AllowedWorkspaces         []string `json:"allowed_workspaces"`
	AllowedCredentialProfiles []string `json:"allowed_credential_profiles"`
	MaxJobSeconds             int      `json:"max_job_seconds"`
}

func Default() Config {
	return Config{
		Server: ServerConfig{
			Address: "127.0.0.1:8080",
		},
		Storage: StorageConfig{
			DataDir: "/var/lib/agent-runtime",
		},
	}
}

func Load(path string) (Config, error) {
	cfg := Default()

	if path != "" {
		raw, err := os.ReadFile(path)
		if err != nil {
			return Config{}, fmt.Errorf("read config: %w", err)
		}
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return Config{}, fmt.Errorf("parse config: %w", err)
		}
	}

	if cfg.Server.Address == "" {
		cfg.Server.Address = Default().Server.Address
	}
	if cfg.Storage.DataDir == "" {
		cfg.Storage.DataDir = Default().Storage.DataDir
	}
	if !filepath.IsAbs(cfg.Storage.DataDir) {
		return Config{}, fmt.Errorf("storage.data_dir must be absolute")
	}
	if cfg.Storage.ToolsDir == "" {
		cfg.Storage.ToolsDir = filepath.Join(cfg.Storage.DataDir, "tools")
	}
	if cfg.Storage.TenantsDir == "" {
		cfg.Storage.TenantsDir = filepath.Join(cfg.Storage.DataDir, "tenants")
	}
	if !filepath.IsAbs(cfg.Storage.ToolsDir) {
		return Config{}, fmt.Errorf("storage.tools_dir must be absolute")
	}
	if !filepath.IsAbs(cfg.Storage.TenantsDir) {
		return Config{}, fmt.Errorf("storage.tenants_dir must be absolute")
	}

	for i := range cfg.Tokens {
		token := &cfg.Tokens[i]
		if token.Role == "" {
			token.Role = "tenant"
		}
		duration := time.Duration(token.MaxJobSeconds) * time.Second
		token.Policy = policy.Policy{
			SubjectID:                 token.SubjectID,
			TenantID:                  token.TenantID,
			Role:                      token.Role,
			AllowedTools:              append([]string(nil), token.AllowedTools...),
			AllowedWorkspaces:         append([]string(nil), token.AllowedWorkspaces...),
			AllowedCredentialProfiles: append([]string(nil), token.AllowedCredentialProfiles...),
			MaxJobDuration:            duration,
		}
	}

	return cfg, nil
}

func (cfg Config) PolicyStore() map[string]policy.Policy {
	out := make(map[string]policy.Policy, len(cfg.Tokens))
	for _, token := range cfg.Tokens {
		out[token.Token] = token.Policy
	}
	return out
}

func (cfg Config) UserStore() []tenants.UserRequest {
	out := make([]tenants.UserRequest, 0, len(cfg.Users))
	for _, user := range cfg.Users {
		out = append(out, tenants.UserRequest{
			Username:                  user.Username,
			Password:                  user.Password,
			Token:                     user.Token,
			SubjectID:                 user.SubjectID,
			TenantID:                  user.TenantID,
			Role:                      user.Role,
			AllowedTools:              append([]string(nil), user.AllowedTools...),
			AllowedWorkspaces:         append([]string(nil), user.AllowedWorkspaces...),
			AllowedCredentialProfiles: append([]string(nil), user.AllowedCredentialProfiles...),
			MaxJobSeconds:             user.MaxJobSeconds,
		})
	}
	return out
}
