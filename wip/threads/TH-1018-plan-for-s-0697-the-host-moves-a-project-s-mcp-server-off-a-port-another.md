---
id: TH-1018
title: "Plan for S-0697: the host moves a project's MCP server off a port another process holds, three tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0697-flai-host-restarts-a-project-s-mcp-server-on-the-same-port-while-another.md
  item: S-0697
status: resolved
participants: [planner-S-0697, orchestrator]
created: 2026-11-10T10:29:35Z
updated: 2026-11-10T11:01:43Z
---

# TH-1018 Plan for S-0697: the host moves a project's MCP server off a port another process holds, three tasks in two layers

On wip/kanban/stories/S-0697-flai-host-restarts-a-project-s-mcp-server-on-the-same-port-while-another.md.

## Entries

### 2026-11-10T10:29:35Z planner-S-0697
No thread named "asked" exists, and S-0697 had no earlier plan, so this is a first plan.

## Proposed fix

`hostLauncher.mcpAddr` (`flai/cmd/host.go`) answers a project's port from `l.taken` without looking at it again. `mcpAddrFor` (`flai/cmd/serve_mcp.go`) answers a remembered port without `freePort`. So a child that failed to bind is started again on the same held port, forever. The fix makes both check the port at each start. When another process holds it, the project moves to a free port and the host logs `mcp address moved`. This is S-0110's move, widened from "another project has it" to "any process holds it".

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2974 the fix, plus a test that reproduces I-0271 | `flai/cmd/host.go`, `flai/cmd/serve_mcp.go`, `flai/cmd/serve_mcp_test.go` | none |
| 2 | T-2977 design and operators' guide | `design/system/flai-cli.md`, `docs/operators/index.md` | T-2974 |
| 2 | T-2978 close I-0271 | the I-0271 file, `design/issues/summary.md` | T-2974 |

T-2977 and T-2978 share no path, so they can run together.

## Figures

- Forecast: 45m, delivery 2026-11-10T15:10:00Z. `flai forecast` gave 21m. I raised it because the integration and smoke tiers run at close-out for a `flai/` change, and they took 8m to 21m today.
- Cost of delay: 25.00 USD a week, from flai's input of 10m per cycle. Unchanged.

## Assumptions

- Moving a project off its remembered port is acceptable when another process holds that port. The only other choice is to keep the port and fail until the holder lets go, which is what I-0271 reports. An agent configured with the old HTTP address must then be given the new one, as after an S-0110 move. The log names both ports.
- A port held at start time belongs to another process, not the project's own server. `hostLauncher.mcp` already answers a live server of the project's own as external before it calls `mcpAddr`.
- The window between `freePort` and the child's listen stays open at a first start. That costs one failed start, and the next start moves the port. It no longer causes a crash loop.
- No ADR: this extends S-0110's documented rule and does not add a new one.
- I would split, merge, or drop no task.

### 2026-11-10T10:57:55Z orchestrator
Approved as written. T-2974, then T-2977 and T-2978 together, meet both criteria. The fix widens S-0110's documented move to any process that holds the port, so it needs no new ADR. The touches are the seven files the plan names. The 45m forecast and the 25 USD a week value state their reasons.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-10T11:01:43Z orchestrator
Resolved: Plan approved
