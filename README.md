# Agent Runtime

Agent Runtime is a portable, multi-tenant runtime for shared CLI agents such as Claude Code, Codex, Gemini CLI, and OpenCode.

The goal is not to replace apps like cc-connect. The goal is to give many apps one shared place for CLI installation, login state, workspaces, job execution, and later interactive terminals. LazyCat should be a packaging target, not a core runtime assumption.

## Current Slice

This repo currently contains the first backend slice:

- Built-in Web UI at `/` for status, tools, and simple job submission.
- JSON config loader.
- Tenant-scoped token policies.
- Tool registry.
- Credential profile home resolver.
- Workspace resolver with symlink escape protection.
- Asynchronous job manager.
- Local process executor.
- HTTP API for health, readiness, status, tools, jobs, and job events.

Not implemented yet:

- CLI installer/updater.
- Web terminal and PTY WebSocket API.
- Persistent audit database.
- Web UI.
- Docker/LazyCat packaging.

## API

```text
GET  /api/health
GET  /api/ready
GET  /api/status
GET  /api/tools
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

`configs/local.json` maps `codex` and `claude` to `/usr/bin/env` so the job path is testable even before real CLIs are installed.

## Storage Layout

The core runtime only needs a configurable tenants directory:

```text
<data_dir>/
  tools/
  tenants/<tenant_id>/
    homes/<credential_profile>/
    workspaces/<workspace_id>/
```

For LazyCat, packaging can point `data_dir` at `/data`. For local development, it can point at `/tmp/agent-runtime` or another absolute path.

## Design Docs

- [Design spec](docs/superpowers/specs/2026-04-27-agent-runtime-design.md)
- [Implementation plan](docs/superpowers/plans/2026-04-27-agent-runtime-first-slice.md)
