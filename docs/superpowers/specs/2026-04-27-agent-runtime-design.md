# Agent Runtime Design

## Goal

Build a portable, multi-tenant CLI agent runtime that can run as a normal development service and later be packaged for LazyCat without binding the core runtime to LazyCat-specific paths or assumptions.

## Scope For The First Slice

The first slice provides a working backend foundation:

- JSON configuration with portable storage roots.
- Tenant-aware policy checks for tools, workspaces, credential profiles, terminal access, and job duration.
- Workspace path resolution that rejects symlink and relative-path escape attempts.
- Shared tool registry with per-job policy enforcement.
- Credential profile home resolution.
- Asynchronous job manager backed by a local process executor.
- HTTP API for health, status, tool listing, job creation, job lookup, and job event streaming.

The first slice intentionally does not implement the full xterm.js web terminal or CLI installer. Those depend on the stable execution, tenant, and credential boundaries built here.

## Architecture

The core repo is deployment-neutral. LazyCat support should live in packaging and adapter directories, while the runtime packages stay usable by a local binary, Docker image, or future Kubernetes executor.

Core packages:

- `internal/config`: Load and validate runtime configuration.
- `internal/policy`: Authorize tenant-scoped requests.
- `internal/workspaces`: Resolve workspace IDs to safe filesystem paths.
- `internal/credentials`: Resolve tenant credential profile homes.
- `internal/tools`: Register and expose CLI tools.
- `internal/jobs`: Validate, store, execute, and stream job state.
- `internal/execution`: Run jobs through the local process executor.
- `internal/api`: Serve the HTTP API.
- `cmd/agent-runtime`: Wire configuration, services, and HTTP server.

## Multi-Tenant Model

The runtime treats tenants, service accounts, credential profiles, and workspaces as first-class concepts. The first version can run in one process with filesystem isolation and token policies, but API and storage names must not assume a single user.

Recommended storage layout:

```text
<data_dir>/
  tools/
  tenants/<tenant_id>/
    homes/<credential_profile>/
    workspaces/<workspace_id>/
    jobs/<job_id>/
```

Tool binaries may be shared across tenants. Credential homes and workspaces are tenant-scoped.

## API

Initial API:

```text
GET  /api/health
GET  /api/ready
GET  /api/status
GET  /api/tools
POST /api/jobs
GET  /api/jobs/{id}
GET  /api/jobs/{id}/events
```

The event endpoint uses server-sent events in the first slice because it is standard-library friendly and adequate for non-interactive jobs. Interactive terminal support should use WebSocket in the next slice:

```text
WS /api/terminal?tenant=<id>&workspace=<id>&credential_profile=<id>
```

## Job Request Shape

```json
{
  "tenant": "team-a",
  "tool": "codex",
  "args": ["exec", "fix tests"],
  "workspace": "repo-main",
  "credential_profile": "team-default",
  "env": {
    "CODEX_HOME": "/custom/home"
  },
  "timeout_seconds": 900
}
```

Callers should use workspace IDs and credential profile IDs instead of raw paths. The runtime maps IDs to paths after authorization.

## Security Boundaries

- Service calls never receive a raw shell API.
- Each request is authenticated by bearer token.
- Tokens map to policies; policies decide tools, workspaces, profiles, terminal access, and max duration.
- Workspace paths are resolved through `filepath.Clean`, absolute roots, and symlink evaluation.
- Audit records store caller ID, tenant, tool, args, workspace, profile, duration, exit code, and event summary, but not bearer tokens or secret environment values.

## Testing Strategy

Tests lock the critical boundaries first:

- Policy allows valid requests and rejects tenant/tool/workspace/profile/terminal/duration violations.
- Workspace resolution rejects symlink escape attempts.
- Job manager enforces policy before execution and records terminal job state.
- API rejects missing credentials and accepts authorized job creation.

