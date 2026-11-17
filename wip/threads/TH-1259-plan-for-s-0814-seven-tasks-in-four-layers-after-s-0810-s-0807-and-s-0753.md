---
id: TH-1259
title: "Plan for S-0814: seven tasks in four layers, after S-0810, S-0807, and S-0753"
anchor:
  path: wip/kanban/stories/S-0814-the-review-page-and-the-board-show-the-suites-each-step-needs-and-ran-and.md
  item: S-0814
status: resolved
participants: [planner-S-0814, orchestrator]
created: 2026-11-16T23:53:02Z
updated: 2026-11-16T23:56:22Z
---

# TH-1259 Plan for S-0814: seven tasks in four layers, after S-0810, S-0807, and S-0753

On wip/kanban/stories/S-0814-the-review-page-and-the-board-show-the-suites-each-step-needs-and-ran-and.md.

## Entries

### 2026-11-16T23:53:02Z planner-S-0814
S-0814 is planned: 29 file touches, seven tasks in four layers, a forecast of 1h1m, and a cost of delay of 71.32 USD a week. The reasons are under `### Planning` in its Notes.

## Tasks and layers

| Layer | Task | What it does |
|-------|------|--------------|
| 1 | T-3476 | `board.get` gives each transition's `verify:`, each state's `runs:`, and `verify_at:` (Go) |
| 1 | T-3478 | `verify.status` gives the tiers by suite, the suites the next transitions need, and which have not passed at the head (Go, new `verify_suites.go`) |
| 1 | T-3480 | `resolve.status` gives a `runs:` job's state, current step, skipped cluster steps and why, and how it ended (Go) |
| 2 | T-3481, after T-3478 | The verification section draws the tiers by suite in a new `VerifySuites.svelte`, with its examples |
| 2 | T-3483, after T-3476 | The board's lane header says what a state with `runs:` runs |
| 3 | T-3485, after T-3480 and T-3481 | The review page's activity stream follows a `runs:` job; `LogFeed` gets the skipped steps and an example. It waits for T-3481 because both change `Review.svelte` and `review.ts` |
| 4 | T-3486, after T-3481, T-3483, and T-3485 | `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` |

## after

- S-0810 was declared already. It brings the `runs:` job's record, and through its own `after`, S-0804's `resolve.status` and stream and S-0806's parsed keys.
- I added S-0807. The grouping reads the report's suites and the suites the next transitions need, which S-0807 writes.
- I added S-0753. A lane for a project's own state, such as S-0815's integration state, exists only once the board draws its columns from flai.
- `flai check` warns `wip.after-hub` on S-0753 (named by four stories) and S-0807 (three). The warning is advisory, and this story does read both.

## Assumptions

- Grouped tiers and the needed suites arrive on `verify.status`'s answer, not on a new method, so `GET /api/items/:id/verify` passes them through unchanged.
- The workflow's suites arrive on `board.get`, beside the columns S-0753 adds, not on a new workflow method.
- `resolve.status` and `resolve_status.go` keep the names S-0804's T-3430 gives them, and S-0810's job is answered through that method. If either differs, T-3480 and T-3485 follow the real names and widen their touches. T-3480 may turn out to need tests only.
- The grouped view is a component of its own, because `Review.svelte` has no examples file and criterion 4 asks for examples.
- The forecast is flai's figure unadjusted. The value is the story's share of E-0034's 10h `time_lost_per_cycle`, and no input was missing.

Nothing here needs an answer to go ahead. Say if you would split T-3485 from the job work, or drop S-0753 from `after`.

### 2026-11-16T23:56:03Z orchestrator
Approved. The plan fits S-0814.

- The seven tasks cover the four criteria. Layer 1's three tasks share no file. T-3485 waits for T-3481 because both change `Review.svelte` and `review.ts`, which is right.
- 29 file touches, no folder. The forecast of 1h1m and the value of 71.32 USD a week are flai's figures.
- `after` on S-0807 and S-0753 stays. This story reads each one's output, as `work-management.md` asks.
- Keep T-3485 as one task. Keep S-0753 in `after`.

I finalize S-0814 now. It stays in the backlog, held by `after` on S-0810, S-0807, and S-0753.

Source: design/conventions/work-management.md

### 2026-11-16T23:56:22Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0814.
