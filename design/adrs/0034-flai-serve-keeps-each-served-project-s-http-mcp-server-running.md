---
id: ADR-0034
title: "flai serve keeps each served project's HTTP MCP server running"
status: accepted
date: 2026-09-23
supersedes: []
superseded_by: [ADR-0040]
refines: [ADR-0030]
---

# ADR-0034 flai serve keeps each served project's HTTP MCP server running

## Context

[ADR-0030](0030-mcp-is-served-by-flai-on-the-host-over-stdio-and-http-and-the-dashboard-s-api.md) put MCP over HTTP in a process of its own, one per project, started by hand with `flai mcp start`, and said it "is not part of `flai serve`: that process exists to dial dashboards, and an agent's server must work with no dashboard at all."

Since then `flai serve` has become the one process on the host that stands for the operator's projects: one per user, serving every project registered with it (S-0080), starting an agent when a story becomes ready (S-0079), and acting as the operator on the host (S-0078). The operator asked (S-0096) that the MCP be running whenever `flai serve` runs, so that agents can always reach it, and that it stop when `flai serve` stops. Nothing kept the HTTP server running: after a reboot, or a crash, agents configured for it found nothing until someone ran `flai mcp start` again in each project.

## Decision

**`flai serve` keeps each project it serves supplied with its HTTP MCP server.** This reverses one clause of ADR-0030, that the server "is not part of `flai serve`"; everything else ADR-0030 decided stands: the server is still `flai mcp http`, one process per project, with its state, token, log, and address in that project's `.flai-cache`, and `flai mcp` on stdio is unchanged.

- For each project it serves, `flai serve` starts `flai mcp http` as a child, unless one is already running for that project, which it uses and never stops, since `flai serve` did not start it.
- It looks again every 15 seconds and starts one again if the one it started has stopped.
- It stops the one it started when the project is no longer served, and when `flai serve` itself stops.
- The child is started with `--exit-with <flai serve's pid>` and ends by itself once that process is gone, so a `flai serve` that was killed, and never got to stop its children, leaves none behind.
- A project with no remembered address takes the first free port from 4243 (up to twenty), one start at a time, and keeps it, as `flai mcp start` keeps the address it used.

## Consequences

- An agent configured for a project's HTTP MCP finds it whenever `flai serve` runs; nothing needs starting by hand after a reboot.
- `flai mcp start`, `stop`, and `status` still work, and a server started with them is left alone. `flai mcp stop` on a server `flai serve` started ends it until the next look, when it is started again; to stop it for good, stop serving the project.
- Every served project listens on a loopback port whether or not an agent uses it: a process and a port per project, and a token file created in each project's `.flai-cache`.
- An agent's server now depends on `flai serve` for being started, though not for working: a server that runs answers with or without a dashboard, and `flai mcp start` remains for a project no `flai serve` serves.
- The ports a project gets depend on what was free when it was first served, so the address to configure is read from `flai mcp status` (or `flai serve status`), not assumed.

## Alternatives considered

- **One listener inside `flai serve` for every project**, the alternative ADR-0030 named. One port and one token, but the server would have to serve several projects at once, and an agent's address would carry a project key; the per-project process already exists and works.
- **Keep starting it by hand.** What ADR-0030 had; the operator asked for it to run whenever `flai serve` does.
- **Stop the child only from `flai serve`'s own shutdown.** Simpler, but a `flai serve` killed with SIGKILL, or by a crash, would leave servers behind; `--exit-with` covers that on every platform, where a Linux-only parent-death signal would not.
