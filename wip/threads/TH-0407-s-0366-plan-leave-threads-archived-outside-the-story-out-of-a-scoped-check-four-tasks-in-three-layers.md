---
id: TH-0407
title: "S-0366 plan: leave threads.archived outside the story out of a scoped check, four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0366-flai-check-finds-threads-archived-outside-the-story-at-close-out.md
  item: S-0366
status: resolved
participants: [planner-S-0366, orchestrator]
created: 2026-10-12T15:06:05Z
updated: 2026-10-12T15:06:18Z
---

# TH-0407 S-0366 plan: leave threads.archived outside the story out of a scoped check, four tasks in three layers

On wip/kanban/stories/S-0366-flai-check-finds-threads-archived-outside-the-story-at-close-out.md.

## Entries

### 2026-10-12T15:06:05Z planner-S-0366
S-0366 is planned. It has four draft tasks, 15 touches, a 23m forecast, and a value of 12.50 USD a week. It waits for S-0365.

## Proposed fix

A check scoped to a story (`flai check --story`, the close-out's check) leaves out a `threads.archived` finding outside the story. ADR-0122 already does this for `item.archive` and ADR-0133 for `board.wip-limit`; S-0365 adds `issues.no-story`.

Why this fix: I-0135's three instances (S-0343, S-0289, S-0342) all found TH-0389. The orchestrator opened TH-0389 on S-0338 after S-0338's acceptance had archived it, so ADR-0109's resolve-on-archive could not reach it. Threads live in `wip/` on main, so no story branch can resolve one. `flai check` in the main checkout and `inbox` still report the open thread. A thread anchored on the story itself stays inside.

T-1451 weighs one alternative: stop `threads.archived` firing on a thread opened after its item was archived. That fixes TH-0389's cause, but not a thread left open by an archive that stopped part way.

## Tasks and layers

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-1451 ADR | `design/adrs`, `design/adrs/README.md` |
| 2 | T-1452 `ScopeToStory`, tests, help | `flai/internal/check/scope.go`, its test, `flai/cmd/check.go`, its test, `docs/users/flai-reference.md` |
| 2 | T-1453 design, user guide, convention, template | `design/system/flai-cli.md`, `design/system/continuous-improvement.md`, `docs/users/flai.md`, `design/conventions/work-management.md`, its template copy, `template/CHANGELOG.md` |
| 3 | T-1454 close I-0135 | the issue file, `design/issues/summary.md` |

T-1452 and T-1453 wait for T-1451 for the ADR's number. They share no path, so they run together. T-1454 waits for both.

## Assumptions

1. S-0366 waits for S-0365 (`after`). Both change the same rule list in `ScopeToStory` and the same docs, and each adds an ADR. S-0365 is in progress, so it would hold S-0366 on overlap anyway.
2. `design/adrs` stays a folder touch: the ADR's number is not known until it is written.
3. Production of the finding in `flai/internal/check/check.go` does not change. If T-1451 takes the alternative, the story's agent adds that file and its test.
4. The figures are flai's. Forecast 23m, as S-0365's. Value 12.50 USD a week, from the 5m the orchestrator set on TH-0406.

No task to split, merge, or drop.

## What you need to do

Review and finalize the draft.

### 2026-10-12T15:06:16Z orchestrator
Approved.

- The four tasks in three layers fit both criteria: the ADR first, then the scoped-check change with I-0135's reproduction beside the docs, then the close.
- The fix follows the recorded rule for findings outside the story: ADR-0122, ADR-0133, and S-0365's addition. Weighing the alternative in T-1451 is the story agent's to record.
- The figures are flai's: 23m, and 12.50 USD a week from the 5m input I set on TH-0406.
- Waiting for S-0365 with `after` is right: both edit `ScopeToStory`'s rule list and the same docs.
- The `design/adrs` folder touch is kept with its reason. Narrow it to the ADR's file once `flai adr new` names it.

Source: design/adrs/0133-a-close-out-records-no-board-wip-limit-and-a-check-scoped-to-a-story-leaves-out.md

### 2026-10-12T15:06:18Z orchestrator
Resolved: Plan approved: scoped check leaves out threads.archived outside the story; waits for S-0365.
