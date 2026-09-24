---
id: T-0375
type: task
nature: feature
title: flai host runs, starts, stops, and reports the host, and flai serve runs under it
status: done
parent: S-0106
owner: alex
created: 2026-09-24T01:30:20Z
updated: 2026-09-24T01:42:34Z
transitions:
  - to: ready
    at: 2026-09-24T01:35:48Z
    by: system-flow
  - to: in-progress
    at: 2026-09-24T01:35:48Z
    by: system-flow
  - to: done
    at: 2026-09-24T01:42:34Z
    by: system-flow
stream: S-0106
tags: []
touches: [flai/cmd]
---
# T-0375 flai host runs, starts, stops, and reports the host, and flai serve runs under it

## Work

- `flai host` (foreground), `start` (detached, waits for its state), `stop`, `status`, `restart|start|stop <serve|mcp|all>`, `check`, `upgrade`, each with `--json`.
- The host runs `flai serve --exit-with <host pid>` and, when serve asks, `flai mcp http --exit-with <host pid>` per project; children get `FLAI_HOST_URL` and `FLAI_HOST_TOKEN`.
- `flai serve` gets a hidden `--exit-with`; under a host its MCP launcher asks the host (`/mcp/ensure`, `/mcp/release`) and starts no process itself; without a host it keeps no MCP server.
- `flai serve start` and `stop` go through the host; `upgrade` runs `flai self-upgrade`, then the host replaces itself with the new binary and starts its children again.

## Done when

- A test runs a real host with this binary's commands in a temp config and sees serve start, restart, and stop with it; behaviour tests pass; lint clean.

## Notes
