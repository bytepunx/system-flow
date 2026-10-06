---
id: T-1008
type: task
nature: improvement
title: The design, an ADR, and the user and operator docs describe the guard rule and how to wait
status: done
parent: S-0285
owner: alex
created: 2026-10-06T10:38:24Z
updated: 2026-10-06T10:54:18Z
transitions:
  - to: ready
    at: 2026-10-06T10:51:06Z
    by: agent-S-0285
  - to: in-progress
    at: 2026-10-06T10:51:06Z
    by: agent-S-0285
  - to: done
    at: 2026-10-06T10:54:18Z
    by: agent-S-0285
stream: S-0285
tags: []
touches: [design/system/flai-cli.md, design/system/agent-context.md, docs/users/flai.md, docs/operators/settings.md, docs/operators/index.md, design/adrs]
after: [T-1004, T-1005, T-1006]
usage:
  source: log
  seconds: 192
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 44
      output: 17149
      cache_read: 2326569
      cache_write: 68263
      cost: 1.2356
---
# T-1008 The design, an ADR, and the user and operator docs describe the guard rule and how to wait

## Work

The design and the docs describe the change: `design/system/flai-cli.md` (the `flai guard` row, the MCP tools), `design/system/agent-context.md` (the sub-agents section), `docs/users/flai.md` (`wait_for_events`, the guard and `.claude/settings.json` sections), `docs/operators/settings.md` (the hooks rows) and `docs/operators/index.md` where it describes how flai serve's agent waits. An ADR, written by the story's agent with `flai adr new`, records the guard rule and the way to wait; the design links it. Waits for T-1004, T-1005, and T-1006, whose behaviour it describes.

## Done when

- [ ] each document above says what the guard refuses, the hooks it needs, how a story's agent waits for a sub-agent, and what `wait_for_events` is for
- [ ] the ADR exists and the design links it

## Notes
