---
id: S-0311
type: story
nature: feature
title: The Orchestrator needs a role it can run tests in review under
status: review
owner: alex
created: 2026-10-07T14:24:35Z
updated: 2026-10-07T14:49:55Z
transitions:
  - to: ready
    at: 2026-10-07T14:24:36Z
    by: alex
  - to: in-progress
    at: 2026-10-07T14:29:36Z
    by: agent-S-0311
  - to: review
    at: 2026-10-07T14:49:55Z
    by: agent-S-0311
tags: [cli]
topics: [orchestrator]
touches: [flai/internal/verify/run.go, flai/internal/verify/run_test.go, flai/internal/verify/proc_test.go, flai/internal/mcpserver/verify_test.go, design/system/flai-cli.md, design/system/strategic-agents.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/verify.go, flai/cmd/test.go, flai/cmd/verify_test.go, flai/cmd/test_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1231
  models:
    - model: claude-opus-5-5
      input: 204
      output: 57158
      cache_read: 8522925
      cache_write: 360674
      cost: 5.1732
  strategic:
    - kind: orchestrator
      seconds: 1277
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 48
          output: 672
          cache_read: 8969916
          cache_write: 19104
          cost: 2.3453
        - model: claude-sonnet-5-5
          input: 8
          output: 32
          cache_read: 70202
          cache_write: 37758
          cost: 0.0724
forecast:
  duration: 30m
  delivery: 2026-10-07T16:00:00Z
  basis: "flai forecast gave 13m (74 s per unit of size over 3 medium feature stories, size 10); raised to 30m for three tasks with real-process and MCP tests and a regenerated reference, as comparable small stories S-0308, S-0265, and S-0296 took 5 to 12 agent minutes each with one task's worth of change"
  by: planner-S-0311
  at: 2026-10-07T14:30:14Z
---
# S-0311 The Orchestrator needs a role it can run tests in review under

## Goal

The orchestrator needs the ability to execute test runs without hitting flai guard's rules about orchestrator permissions in test scenarios. This likely means using whatever FLAI_ROLE identity will avoid tripping the guard's check.

## Acceptance criteria
- [x] Orchestrator is able to complete test runs as a role that does not break flai's guard checks against the `orchestrator` role

## Tasks
- T-1161 Run every test tier under FLAI_ROLE=verify, whatever role started flai verify or flai test
- T-1162 Prove flai verify, flai test, and the MCP tool verify pass a flai-writing tier under the orchestrator's role
- T-1163 Document that test tiers run under FLAI_ROLE=verify, in the commands' help, the user guide, and the design

## Notes

### Planning

Cause, from TH-0260: `OS.Run` in `flai/internal/verify/proc.go` gives each tier flai's own environment. Under the orchestrator that holds `FLAI_ROLE=orchestrate`, and flai's own checks in `flai/cmd/move.go`, `accept.go`, `edit.go`, `order.go`, `plan.go`, and `flai/internal/mcpserver/server.go` refuse the in-process fixture writes of the Go tests. The `flai guard` hook is not involved. The plan runs every tier under `FLAI_ROLE=verify`, the existing role of a run that only checks, which none of those checks match.

Touches:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/cmd` | declared | Kept, as the story declared it. It is a folder touch: T-1162 and T-1163 name the files in it (`verify.go`, `test.go`, `verify_test.go`, `test_test.go`), which narrow it in the claim (ADR-0096) |
| `flai/internal/verify/run.go` | layout | `runTier` builds each tier's environment |
| `flai/internal/verify/run_test.go`, `flai/internal/verify/proc_test.go` | layout | The runner's tests |
| `flai/internal/mcpserver/verify_test.go` | co-change, layout | The MCP tool `verify` runs the same tiers; `mcpserver` co-changes with `flai/cmd` in 13% of commits |
| `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flai-reference.md` | co-change, design | The top three co-changes of `flai/cmd` (23 to 25%); they describe `flai verify` and `flai test` |
| `design/system/strategic-agents.md` | design | Accepting a story (S-0221), step 1, is the orchestrator's verify |

Left out: `flai/internal/hostapi/writes.go`, `flai/internal/serve/agents.go`, `flai/internal/guard/guard.go`, and the rest of `flai touches suggest`'s list. Each co-changes with `flai/cmd` in general, but the fix does not reach it.

Forecast: `flai forecast` gave 13m, from 74 s per unit of size over 3 medium feature stories, times size 10. Raised to 30m: the plan has three tasks, a real-process test, a test through the MCP server, and a regenerated reference. S-0308, S-0265, and S-0296, each about one task of change, took 5 to 12 agent minutes. The delivery, 16:00Z, keeps flai's offset from the start and adds the 17 minutes.

Cost of delay: no value is written. The story has no inputs and no epic, and `flai cod S-0311` cannot work one out. The inputs are asked for on TH-0283.
