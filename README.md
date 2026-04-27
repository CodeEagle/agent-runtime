# Agent Runtime

Agent Runtime is a portable, multi-tenant runtime for shared CLI agents such as Claude Code, Codex, Gemini CLI, and OpenCode.

The goal is not to replace apps like cc-connect. The goal is to give many apps one shared place for CLI installation, login state, workspaces, job execution, and later interactive terminals. LazyCat should be a packaging target, not a core runtime assumption.

## Current Slice

This repo currently contains the first runtime slice:

- Built-in Web UI at `/` for terminal login, CLI management, and tenant policy visibility.
- Official-source CLI install cards for Claude Code, Codex, Gemini CLI, OpenCode, iFlow, Kimi, and Qoder.
- JSON config loader.
- Tenant-scoped token policies.
- Persistent tool registry with add/update/delete HTTP APIs.
- Credential profile home resolver.
- Workspace resolver with symlink escape protection.
- PTY-backed WebSocket terminal for human login flows.
- Asynchronous job manager.
- Local process executor.
- HTTP API for health, readiness, status, tools, tenants, jobs, job events, and terminal sessions.

Not implemented yet:

- Background CLI installer/updater jobs.
- Persistent audit database.
- CLI login-state health probes.

## API

```text
GET  /api/health
GET  /api/ready
GET  /api/status
GET  /api/tools
POST /api/tools
DELETE /api/tools/{name}
GET  /api/tenants
WS   /api/terminal
POST /api/jobs
GET  /api/jobs/{id}
GET  /api/jobs/{id}/events
```

Example job:

```json
{
  "tenant": "team-a",
  "tool": "codex",
  "args": ["exec", "fix tests"],
  "workspace": "repo-main",
  "credential_profile": "team-default",
  "timeout_seconds": 900
}
```

## Local Development

This machine does not currently have `go` installed globally. If Nix is available, run commands through a temporary Go shell:

```bash
nix shell nixpkgs#go -c go test ./...
nix shell nixpkgs#go -c go run ./cmd/agent-runtime -config configs/local.json
```

Build and run the container:

```bash
docker build -t agent-runtime:local .
docker run --rm -p 8080:8080 -v agent-runtime-data:/data agent-runtime:local
```

Create a job:

```bash
curl -sS \
  -H 'Authorization: Bearer dev-token' \
  -H 'Content-Type: application/json' \
  -d '{"tenant":"team-a","tool":"codex","args":["exec","fix tests"],"workspace":"repo-main","credential_profile":"team-default","timeout_seconds":60}' \
  http://127.0.0.1:8080/api/jobs
```

`configs/local.json` registers the known CLI command names up front. Install them from the Web UI's official-source cards, then use the terminal or job API against the shared tenant credential profile.

Open the Web UI at `/`, enter a token such as `dev-token`, connect the terminal with tenant `team-a`, workspace `repo-main`, and credential profile `team-default`, then use the login shortcut buttons for CLI auth flows.

## Storage Layout

The core runtime only needs a configurable tenants directory:

```text
<data_dir>/
  tools/
    registry.json
  tenants/<tenant_id>/
    homes/<credential_profile>/
    workspaces/<workspace_id>/
```

For LazyCat, packaging can point `data_dir` at `/data`. For local development, it can point at `/tmp/agent-runtime` or another absolute path.

## Design Docs

- [Design spec](docs/superpowers/specs/2026-04-27-agent-runtime-design.md)
- [Implementation plan](docs/superpowers/plans/2026-04-27-agent-runtime-first-slice.md)
