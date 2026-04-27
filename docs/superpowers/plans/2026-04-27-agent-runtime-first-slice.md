# Agent Runtime First Slice Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the first runnable backend slice for a portable, multi-tenant CLI agent runtime.

**Architecture:** The service is a single Go binary with focused internal packages for config, policy, workspaces, credentials, tools, jobs, local execution, and HTTP API. LazyCat-specific packaging stays out of core runtime code.

**Tech Stack:** Go standard library, net/http, encoding/json, os/exec, filesystem-backed storage roots.

---

### Task 1: Policy And Path Safety

**Files:**
- Create: `internal/policy/policy_test.go`
- Create: `internal/policy/policy.go`
- Create: `internal/workspaces/path_test.go`
- Create: `internal/workspaces/path.go`

- [ ] Write failing tests for policy allow/deny behavior and workspace symlink escape rejection.
- [ ] Run `go test ./internal/policy ./internal/workspaces` and confirm the tests fail because implementation is missing.
- [ ] Implement the minimum policy and path resolver code.
- [ ] Run `go test ./internal/policy ./internal/workspaces` and confirm both packages pass.

### Task 2: Core Registry And Job Manager

**Files:**
- Create: `internal/tools/registry.go`
- Create: `internal/credentials/resolver.go`
- Create: `internal/jobs/manager_test.go`
- Create: `internal/jobs/manager.go`

- [ ] Write failing tests proving the job manager enforces policy, resolves workspace/profile/tool data, executes allowed jobs, and stores job events.
- [ ] Run `go test ./internal/jobs` and confirm failure.
- [ ] Implement the minimum registry, credential resolver, job manager, and in-memory event store.
- [ ] Run `go test ./internal/jobs` and confirm success.

### Task 3: Local Executor And HTTP API

**Files:**
- Create: `internal/execution/local.go`
- Create: `internal/api/server_test.go`
- Create: `internal/api/server.go`

- [ ] Write failing API tests for unauthorized requests, tool listing, job creation, and job lookup.
- [ ] Run `go test ./internal/api` and confirm failure.
- [ ] Implement HTTP handlers and local process executor.
- [ ] Run `go test ./internal/api` and confirm success.

### Task 4: Binary Wiring And Documentation

**Files:**
- Create: `cmd/agent-runtime/main.go`
- Create: `internal/config/config.go`
- Create: `configs/local.json`
- Create: `README.md`

- [ ] Add config defaults and JSON loading.
- [ ] Wire the binary from config to API server.
- [ ] Document local usage, tenant policy, LazyCat packaging boundary, and next slices.
- [ ] Run `go test ./...` and `go build ./cmd/agent-runtime`.

