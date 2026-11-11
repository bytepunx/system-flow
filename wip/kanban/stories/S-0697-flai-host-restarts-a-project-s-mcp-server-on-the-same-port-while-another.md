---
id: S-0697
type: story
nature: remediation
title: flai host restarts a project's MCP server on the same port while another process holds it, so the project has no MCP server until that process lets go
status: ready
owner: alex
created: 2026-11-10T08:54:33Z
updated: 2026-11-11T05:20:51Z
transitions:
  - to: ready
    at: 2026-11-11T02:31:09Z
    by: alex
tags: [flai, host, mcp]
topics: [mcp, host]
touches: [flai/cmd/host.go, flai/cmd/serve_mcp.go, flai/cmd/serve_mcp_test.go, design/system/flai-cli.md, docs/operators/index.md, design/issues/I-0271-flai-host-restarts-a-project-s-mcp-server-on-the-same-port-while-another.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 22
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 100
          output: 371
          cache_read: 5264790
          cache_write: 210092
          cost: 0.3686
        - model: claude-opus-5-5
          input: 130
          output: 10624
          cache_read: 9188806
          cache_write: 175332
          cost: 3.193
    - kind: orchestrator
      seconds: 7
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 2
          cache_read: 102139
          cache_write: 491
          cost: 0.0281
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-10T08:54:33Z
  value: 25
  by: planner-S-0697
  at: 2026-11-10T10:28:31Z
forecast:
  duration: 45m
  delivery: 2026-11-11T07:48:00Z
  basis: "Its own forecast of 45m; 5th in the pull order with an in-progress limit of 3, behind S-0678, S-0686, S-0691, S-0694, S-0681, S-0718 and S-0696."
  by: flai
  at: 2026-11-11T05:20:51Z
finalized:
  by: orchestrator
  at: 2026-11-10T11:04:41Z
---
# S-0697 flai host restarts a project's MCP server on the same port while another process holds it, so the project has no MCP server until that process lets go

## Goal

This story remediates [I-0271](../../../design/issues/I-0271-flai-host-restarts-a-project-s-mcp-server-on-the-same-port-while-another.md), "flai host restarts a project's MCP server on the same port while another process holds it, so the project has no MCP server until that process lets go". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0271 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0271 is closed with `flai issue close I-0271 --reason` saying what fixed it

## Tasks
- T-2974 flai host moves a project's MCP server off a port another process holds, and a test reproduces I-0271
- T-2977 The design and the operators' guide say a project's MCP server moves off a port another process holds
- T-2978 Close I-0271 with the reason that the host now moves a held port

## Notes

Cost of delay inputs set by flai from I-0271. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T04:41:15Z, 3.2 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from I-0271's one instance: `hostLauncher.mcpAddr` in `flai/cmd/host.go` answers a project's port from `l.taken` without looking at it again, and `mcpAddrFor` in `flai/cmd/serve_mcp.go` answers a remembered port without `freePort`. So a child that failed to bind is started again on the same held port. The fix makes both look: when the port the project has, from `l.taken` or remembered, cannot be listened on at a start, and `readMCPState` found no live server of the project's own (which `mcp` already answers as external), the project moves to the next free port as a new one would, `l.taken` is updated, and the host logs `mcp address moved` naming the holder's port and the new one. This is S-0110's move, which `docs/operators/index.md` already documents for two projects on one port, widened to any other process. The bind race the issue names stays at a first start, but costs one failed start rather than a crash loop.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/cmd/host.go` | design: I-0271 names `hostLauncher.mcpAddr`; co-change |
| `flai/cmd/serve_mcp.go` | design: `mcpAddrFor` answers the remembered port unchecked; co-change |
| `flai/cmd/serve_mcp_test.go` | co-change (4 of 14 commits); the mcpAddr tests live there |
| `design/system/flai-cli.md` | design: its `flai mcp http` row says how the host picks a port; co-change; S-0110 changed it |
| `docs/operators/index.md` | design: says when a project moves port; S-0110 changed it |
| `design/issues/I-0271-…md`, `design/issues/summary.md` | criterion 2, `flai issue close` |

`flai touches suggest` also listed `flai/cmd/serve.go`, `flai/cmd/host_test.go`, and `flai/internal/host/host.go` (21 to 43 percent). They are left out: the fix stays in the launcher, and the host's child loop already calls `spec` again at each start. No folder touch. No `after`: the story reads no other open story's output.

Figures:

- Forecast: `flai forecast` gave 21m. Raised to 45m, because a `flai/` change runs the integration and smoke tiers at close-out, and those took 8m to 21m on today's close-outs (S-0672, S-0669). Delivery is moved back by the same 24m, to 2026-11-10T15:10:00Z.
- Cost of delay: 25.00 USD a week, as `flai cod` gave it from flai's input (10m lost per 168h cycle at 150 USD an hour). Kept: one occurrence, no penalty, no revenue.
