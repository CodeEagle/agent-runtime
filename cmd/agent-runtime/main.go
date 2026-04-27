package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"agent-runtime/internal/api"
	"agent-runtime/internal/config"
	"agent-runtime/internal/credentials"
	"agent-runtime/internal/execution"
	"agent-runtime/internal/jobs"
	"agent-runtime/internal/tenants"
	"agent-runtime/internal/terminal"
	"agent-runtime/internal/tools"
	"agent-runtime/internal/workspaces"
)

func main() {
	configPath := flag.String("config", os.Getenv("AGENT_RUNTIME_CONFIG"), "path to JSON config")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	toolRegistry, err := tools.NewPersistentRegistry(cfg.Tools, filepath.Join(cfg.Storage.ToolsDir, "registry.json"))
	if err != nil {
		log.Fatalf("load tool registry: %v", err)
	}
	tenantStore := tenants.NewStore(cfg.PolicyStore())
	workspaceResolver := workspaces.NewResolver(cfg.Storage.TenantsDir)
	credentialResolver := credentials.NewResolver(cfg.Storage.TenantsDir)
	manager := jobs.NewManager(jobs.Options{
		Policies: tenantStore,
		Tools:    toolRegistry,
		ResolveWorkspace: func(tenantID string, workspaceID string) (string, error) {
			return workspaceResolver.ResolveWorkspace(tenantID, workspaceID)
		},
		ResolveCredentialProfile: func(tenantID string, profileID string) (string, error) {
			return credentialResolver.ResolveProfile(tenantID, profileID)
		},
		Executor: execution.LocalExecutor{},
	})
	terminalHandler := terminal.NewHandler(terminal.Options{
		Policies: tenantStore,
		Tools:    toolRegistry,
		ResolveWorkspace: func(tenantID string, workspaceID string) (string, error) {
			return workspaceResolver.ResolveWorkspace(tenantID, workspaceID)
		},
		ResolveCredentialProfile: func(tenantID string, profileID string) (string, error) {
			return credentialResolver.ResolveProfile(tenantID, profileID)
		},
	})

	handler := api.NewServer(api.Options{
		Jobs:     manager,
		Tools:    toolRegistry,
		Tenants:  tenantStore,
		Terminal: terminalHandler,
	})

	log.Printf("agent-runtime listening on %s", cfg.Server.Address)
	if err := http.ListenAndServe(cfg.Server.Address, handler); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
