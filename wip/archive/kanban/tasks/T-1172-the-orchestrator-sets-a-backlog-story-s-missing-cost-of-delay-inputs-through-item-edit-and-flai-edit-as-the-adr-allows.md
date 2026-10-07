---
id: T-1172
type: task
nature: feature
title: The orchestrator sets a backlog story's missing cost of delay inputs through item_edit and flai edit, as the ADR allows
status: cancelled
parent: S-0328
owner: alex
created: 2026-10-07T19:39:56Z
updated: 2026-10-07T19:42:16Z
transitions:
  - to: cancelled
    at: 2026-10-07T19:42:16Z
    by: agent-S-0328
stream: S-0328
tags: [flai, guard, orchestrator, cost-of-delay]
touches: [flai/internal/workitem/promotable.go, flai/internal/mcpserver/items_write.go, flai/internal/mcpserver/items_write_test.go, flai/cmd/edit.go, flai/cmd/edit_test.go, flai/cmd/orchestrate_permissions_test.go, flai/internal/guard/guard.go, flai/internal/guard/orchestrate_test.go]
after: [T-1171]
---
# T-1172 The orchestrator sets a backlog story's missing cost of delay inputs through item_edit and flai edit, as the ADR allows

## Work

Let the orchestrator choose a cost of delay, or take the planner's recommended one, on a story it had planned, within the permissions T-1168's ADR names (recommended: `plan_backlog_stories` with `answer_threads` at `autonomous`):

- `flai/internal/guard/guard.go` lets the orchestrator's `item_edit` through when it gives only `id`, `hash`, `project`, and `cost_of_delay` with input keys (`revenue_per_week`, `penalty_per_week`, `time_lost_per_cycle`), and `flai edit` with only the matching input flags, while those permissions are on. Every other edit stays refused as before.
- `Repo.OrchestratorPermits` in `flai/internal/workitem/promotable.go`, called from `flai/internal/mcpserver/items_write.go` and `flai/cmd/edit.go` under `FLAI_ROLE=orchestrate`, refuses the write unless the item is a backlog story with no cost of delay inputs on it or its epic, so that it never overwrites the operator's inputs. The inputs it writes are stamped `by: orchestrator`.

It waits for T-1171, which changes `guard.go`, `orchestrate_test.go`, and `orchestrate_permissions_test.go` before it; the two share those files, so they run one after the other.

## Done when

- With the permissions on, the orchestrator sets a backlog story's missing inputs, `flai cod` then gives its value, and the inputs record `by: orchestrator`.
- With a permission off, on a story that already has inputs or whose epic has them, on a story not in the backlog, or with any other field in the edit, the guard or flai refuses, naming why.
- `orchestrate_test.go`, `orchestrate_permissions_test.go`, `items_write_test.go`, and `edit_test.go` cover each case.
- `flai test` passes on the paths changed.

## Notes

If the operator's answer on the plan thread keeps cost of delay inputs the operator's alone, cancel this task and record why in `## Decisions`: the orchestrator then only recommends, which `answer_threads` allows already.
- 2026-10-07T19:42:16Z: moved to cancelled: duplicates the story agent's plan (T-1174, T-1176 to T-1179), written at start while the planner ran; its refinements are folded into T-1176 to T-1178
