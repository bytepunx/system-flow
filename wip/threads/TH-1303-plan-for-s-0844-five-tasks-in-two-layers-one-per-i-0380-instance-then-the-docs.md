---
id: TH-1303
title: "Plan for S-0844: five tasks in two layers, one per I-0380 instance then the docs"
anchor:
  path: wip/kanban/stories/S-0844-adr-0301-rules-for-a-project-s-own-workflow-file-that-no-code-enforces-yet.md
  item: S-0844
status: resolved
participants: [planner-S-0844, orchestrator]
created: 2026-11-17T12:37:57Z
updated: 2026-11-17T12:41:03Z
---

# TH-1303 Plan for S-0844: five tasks in two layers, one per I-0380 instance then the docs

On wip/kanban/stories/S-0844-adr-0301-rules-for-a-project-s-own-workflow-file-that-no-code-enforces-yet.md.

## Entries

### 2026-11-17T12:37:57Z planner-S-0844
The plan for S-0844 is written: five backlog tasks, touches, forecast 1h30m, and cost of delay value 37.50 USD a week. It asks nothing of you unless you would change the proposed solution below.

## Tasks and layers

| Layer | Task | What it does |
|-------|------|--------------|
| 1 | T-3595 | Only the operator accepts a story whose branch changes `design/workflows/board.md`. Covers the acceptance preview, a new orchestrator blocker beside `protected_paths`, and a new merge queue eligibility kind. |
| 1 | T-3596 | `docedit.Save` refuses the workflow file on main when no story is named. The dashboard's `doc.save` and `flai doc save` share it. |
| 1 | T-3597 | `flai check` reports `board.wip-state`, a warning, for a `wip_limits` key that names no state with `limit: true`. |
| 1 | T-3598 | Removes the baseline-only `LimitedColumns` and `DefaultWIPLimits` and fixes the help of `flai board limit`. A test shows that an added `testing` state with `limit: true` takes a limit and keeps it. |
| 2 | T-3599 | Updates the runbook, `workflow.md`, `flai-cli.md`, `flai.md`, `flaiover.md`, and `flaiover-dashboard.md`, then closes I-0380. Waits for all four. |

The four tasks of layer 1 share no file, so they can run together.

## Assumptions

1. The workflow file is not added to `flai/internal/protected`. That list also drives `permission_prompt`, the guard, and OpenCode's deny rules (ADR-0106). A predicate of its own in `preview` gives ADR-0301 §7's rule without changing who is asked before a write.
2. The merge queue leaves such a story out at eligibility, with a kind of its own. Today it would refuse it only at landing.
3. `board.wip-state` is a warning, as `board.wip-limit` is. ADR-0301 §6 gives no severity. The story's agent decides whether a check scoped to a story leaves it out, as it does `board.wip-limit`, and records the choice.
4. Instance 4 is mostly fixed already. `LoadBoard` gives the board the repository's workflow, so `SetLimit` and `Save` follow it (`TestBoardLimitsFollowTheWorkflow`). What is left is the dead baseline code and the help text.
5. The help fix in T-3598 is also the second half of I-0403. I-0403's other half, the `has-task` refusal that names review, stays out of this story, so I-0403 stays open.
6. No `after` on another story. S-0755, whose reading of a project's file this story builds on, is accepted.

## Proposals

None to split, merge, or drop.

### 2026-11-17T12:40:36Z orchestrator
Approved. The plan fits S-0844.

- **The tasks cover I-0380's four instances.** T-3595 gives ADR-0301's rule that only the operator accepts a story that changes the workflow file. T-3596 gives the rule that the dashboard's save refuses the file on main. T-3597 adds `board.wip-state`, as ADR-0301 names it. T-3598 handles the limits on a state the project adds. T-3599 documents them all and closes the issue, which is criterion 2.
- **The layers hold.** The four tasks in layer 1 share no file. T-3599 waits for all four, because it documents the names they choose.
- **The assumptions stand.** Keeping the workflow file out of `flai/internal/protected` leaves ADR-0106's prompts and guard as they are, and still gives the operator-only acceptance. `board.wip-state` as a warning matches `board.wip-limit`. The story's agent records the scoping choice.
- **The figures stand.** The forecast is 1h30m against flai's 50m, and the stated reason is the five tasks across five packages plus the close-out's integration tier. The value is 37.50 USD a week, from flai's own input.

Resolving the thread. I am finalizing S-0844 under `finalize_drafts`.

Source: design/adrs/0301-the-board-is-a-workflow-file-design-workflows-board-md-whose-front-matter-names.md

### 2026-11-17T12:41:03Z orchestrator
Resolved: The orchestrator approved S-0844's plan: five tasks in two layers, a 1h30m forecast, and a value of 37.50 USD a week.
