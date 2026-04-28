# Agent Runtime

Agent Runtime is a portable, multi-tenant runtime for shared CLI agents such as Claude Code, Codex, Gemini CLI, and OpenCode.

The goal is not to replace apps like cc-connect. The goal is to give many apps one shared place for CLI installation, login state, workspaces, job execution, and later interactive terminals. LazyCat should be a packaging target, not a core runtime assumption.

## Current Slice

This repo currently contains the first runtime slice:

- Built-in Web UI at `/` for a dark single-page CLI control plane with silent user session bootstrap, UI-only CLI install/auth actions, runtime user count, and an API reference view.
- Official-source CLI install cards for Claude Code, Codex, Gemini CLI, OpenCode, iFlow, Kimi, and Qoder.
- JSON config loader.
- Tenant-scoped user and token policies with `admin` and `tenant` roles.
- Persistent tool registry with add/update/delete HTTP APIs.
- Credential profile home resolver.
- Workspace resolver with symlink escape protection.
- PTY-backed WebSocket terminal API for integrations and hidden UI actions.
- Tenant-aware file APIs for `tenants/<tenant>/workspaces` and `tenants/<tenant>/homes`, including bounded file preview.
- CLI executable health probes based on the real runtime `PATH`; registered tools are not reported as available until the command is actually present.
- Asynchronous job manager.
- Local process executor.
- HTTP API for health, readiness, status, login, tools, users, tenants, tenant tokens, files, jobs, job events, terminal sessions, and `/openapi.json`.

Not implemented yet:

- Background CLI installer/updater jobs.
- Persistent audit database.
- CLI login-state health probes beyond executable/version detection.

## API

```text
GET  /api/health
GET  /api/ready
GET  /api/status
GET  /openapi.json
POST /api/login
GET  /api/session
GET  /api/tools
POST /api/tools
DELETE /api/tools/{name}
GET  /api/tenants
GET  /api/users
POST /api/users
DELETE /api/users/{id}
GET  /api/tokens
POST /api/tokens
DELETE /api/tokens/{id}
GET  /api/files?tenant=team-a&space=workspaces&path=/
GET  /api/files/raw?tenant=team-a&space=workspaces&path=/repo-main/README.md
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

`configs/local.json` registers the known CLI command names up front. They show as unavailable until the command is found in the runtime `PATH`. Install and authorize them from the Web UI's official-source cards, then use the job API against the shared tenant credential profile.

Open the Web UI at `/`. The UI silently bootstraps the default `admin` / `admin` human session from the sample config, keeps terminal details hidden, and exposes a one-click API reference view similar to Swagger docs.

`dev-token` remains an admin bearer token for service/API calls in the sample config. `team-a-token` remains a tenant bearer token that can only see and browse `team-a`; human Web UI login is handled through the `users` config.

## Storage Layout

The core runtime only needs a configurable tenants directory:

```text
<data_dir>/
  tools/
    registry.json
  tenants/
    registry.json
    <tenant_id>/
      homes/<credential_profile>/
      workspaces/<workspace_id>/
```

For LazyCat, packaging can point `data_dir` at `/data`. For local development, it can point at `/tmp/agent-runtime` or another absolute path.

## Design Docs

- [Design spec](docs/superpowers/specs/2026-04-27-agent-runtime-design.md)
- [Implementation plan](docs/superpowers/plans/2026-04-27-agent-runtime-first-slice.md)
