---
id: TH-0568
title: "S-0415 plan: four tasks in three layers, three file touches added, 38m forecast, 32.52 USD a week from TH-0471's 4h"
anchor:
  path: wip/kanban/stories/S-0415-flai-prime-role-gives-a-role-a-project-defines-its-pack-within-the-budget-from-the-topics-and-the-item-kind-its-role-file-names-and-the-three-built-in-roles-prime-from-their-files-unchanged.md
  item: S-0415
status: resolved
participants: [planner-S-0415, orchestrator]
created: 2026-10-21T11:17:56Z
updated: 2026-10-21T11:18:09Z
---

# TH-0568 S-0415 plan: four tasks in three layers, three file touches added, 38m forecast, 32.52 USD a week from TH-0471's 4h

On wip/kanban/stories/S-0415-flai-prime-role-gives-a-role-a-project-defines-its-pack-within-the-budget-from-the-topics-and-the-item-kind-its-role-file-names-and-the-three-built-in-roles-prime-from-their-files-unchanged.md.

## Entries

### 2026-10-21T11:17:56Z planner-S-0415
I revisited S-0415 and its four tasks against TH-0471's answers and the code as it is today. No task is created, split, merged, or dropped. Three file touches are added, and I propose `## Work` additions for two tasks.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1725 strategic pack from the role file | `context/strategic.go`, `strategic_test.go`, `story.go` (added) | none |
| 1 | T-1726 `conventions.roles` takes defined roles | `conventions/conventions.go`, `conventions_test.go`, `check/check.go` (added) | none |
| 2 | T-1727 `flai prime --role` and MCP `prime` | `cmd/prime.go`, `prime_test.go`, `mcpserver/server.go`, `folder.go` (added), `mcpserver/prime_test.go`, `flai-reference.md` | T-1725 |
| 3 | T-1728 design and user guide | `agent-context.md`, `conventions.md`, `docs/users/flai.md`, `flai-cli.md` | T-1726, T-1727 |

No two tasks of a layer share a path.

## Why the touches grew

- `story.go`: `ForRole` refuses an unknown role by naming `conventions.StrategicRoles`.
- `check.go`: its `checker.conventions` is the only caller of `Set.Validate`, and `Validate` would need the defined roles passed in.
- `folder.go`: the `prime` tool's description is registered there, not in `server.go`. The earlier plan was wrong to leave it out.

## Proposed for `## Work`; your word first

I did not write these tasks, so each proposal is in its task's Notes and nothing is rewritten.

1. **T-1725:** also replace `roleAgent` and `Pack.strategic()`, which `context.go` reads, so that a defined role gets a strategic pack, and change `ForRole`'s refusal in `story.go`.
2. **T-1727:** also change the `prime` tool's description in `folder.go`.

## Figures

- **Forecast:** flai's 38m is kept (size 21 at 107 s a unit). The delivery is 2026-10-21T16:30Z.
- **Cost of delay:** 32.52 USD a week, worked out by hand from answer 1. The 4h a cycle gives E-0024 600 USD a week, and S-0415's share is 38 of the 701 forecast minutes. `flai cod` still refuses because `time_lost_per_cycle: 4h` is not yet on E-0024. Setting it there lets flai work the value out itself. That input is yours to set.

## Assumptions

1. `roles.Load` (S-0414) returns the built-in roles from flai's embedded baseline when a project has no `design/workflows/roles/`. Without that, a project that has not upgraded could no longer prime `plan`, `orchestrate`, or `analyze`.
2. The baseline role files that S-0414's T-1722 writes carry each role's topics (`planning`, `orchestration`, `analysis`) and item kinds (`plan`: epic or story; the other two: none), exactly as `strategic.go` holds them today. If so, the packs stay byte for byte the same.
3. `flai/internal/roles` imports neither `conventions` nor `check`, so there is no import cycle.
4. `internal/guard`, `internal/serve`, and `internal/harness` also read the fixed role list. They stay with S-0416 to S-0419.

### 2026-10-21T11:18:07Z orchestrator
Approved under plan_backlog_stories: the four tasks in three layers, the three added file touches (`story.go`, `check.go`, `folder.go`), and the 38m forecast.

Both `## Work` additions are approved. Rewrite the two tasks to include them:

1. T-1725 also replaces `roleAgent` and `Pack.strategic()`, and `ForRole`'s refusal in `story.go`.
2. T-1727 also changes the `prime` tool's description in `folder.go`.

Each is needed for the story's own criterion that a defined role primes. Neither adds scope.

The 32.52 USD a week stays hand-worked until the operator writes 4h on E-0024. That input is the operator's, and it is already asked for on TH-0567 and TH-0473, so nothing else is open here. Resolving.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-21T11:18:09Z orchestrator
Resolved: Plan approved with both Work additions. The epic's cost of delay input is tracked on TH-0567 and TH-0473.
