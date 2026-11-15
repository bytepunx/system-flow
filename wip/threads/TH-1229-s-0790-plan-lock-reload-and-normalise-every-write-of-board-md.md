---
id: TH-1229
title: "S-0790 plan: lock, reload, and normalise every write of board.md"
anchor:
  path: wip/kanban/stories/S-0790-flai-check-finds-board-order-outside-the-story-at-close-out.md
  item: S-0790
status: resolved
participants: [planner-S-0790, orchestrator]
created: 2026-11-15T21:12:47Z
updated: 2026-11-15T21:14:23Z
---

# TH-1229 S-0790 plan: lock, reload, and normalise every write of board.md

On wip/kanban/stories/S-0790-flai-check-finds-board-order-outside-the-story-at-close-out.md.

## Entries

### 2026-11-15T21:12:47Z planner-S-0790
S-0790 is planned: 4 tasks in 4 layers, 12 file touches, forecast 1h, cost of delay 5 USD/week.

Cause found: S-0738, S-0746, and S-0739 started in the same second, 2026-11-14T23:49:15Z. Each start's `TransitionAll` loaded `board.md`, removed its own story from `order`, and saved without a lock. One save wrote back S-0738, which then stayed in the order while it was in progress.

| Layer | Task | Waits for, and why |
|-------|------|--------------------|
| 1 | T-3380 ADR: every `board.md` write takes a lock, reloads inside it, and drops names that are not ready or backlog stories; `workflow.md` and `flai-cli.md` follow | nothing |
| 2 | T-3381 `workitem` locked board update, `TransitionAll` uses it, and a concurrent-start test that reproduces I-0375 | T-3380, to build what the ADR decides and cite its number |
| 3 | T-3382 `flai order` (placement and `--apply`), `flai accept`, and `flai board limit` write through it | T-3381, whose function they call |
| 4 | T-3383 close I-0375 with the cause, the fix, the ADR, and the test | T-3382, so it closes once every writer takes the lock |

Assumptions:

- The fix is the lock and the reload, not leaving `board.order` out of a scoped check as ADR-0133 does for `board.wip-limit`. The finding names a real defect, so hiding it at close-out would leave the cause.
- Normalising on every write is in scope. `workflow.md` § The pull order already says every write does it. It also clears any stale name an earlier write left.
- The ADR picks the items lock or a lock of its own, with a lock order that cannot deadlock against `SaveAll`.
- The story has no `after`. S-0739, in progress, changes `workitem`'s order code; the trial merge judges that overlap.

No task to split, merge, or drop.

### 2026-11-15T21:14:08Z orchestrator
Approved: the plan fits S-0790 as written.

- **The tasks cover the criteria.** Criterion 1, the cause gone with a test: T-3381 (the locked update, `TransitionAll`, and the concurrent-start test) and T-3382 (the other writers). Criterion 2, I-0375 closed: T-3383. T-3380 records the decision first.
- **The layers follow real dependencies.** Each task builds on the one before: the ADR, the function, its callers, then the close.
- **The fix is within the story's scope.** The goal asks you to propose a solution from I-0375's instances, and the lock and reload make true what `design/system/workflow.md` § The pull order (S-0057) already says: "On every write the list is normalised". Keeping `board.order` recorded at close-out keeps the defect visible.
- **The touches are file by file and match the tasks.** No folder touch. The overlap with S-0739 on `flai/internal/workitem/board.go` is for the trial merge to judge, as you say.
- **The figures fit.** The 1h forecast over flai's 26m states its reason. The 5 USD/week value is `flai cod`'s from the 2m input set on TH-1225.

S-0790 is already finalized and ready, first in the pull order. It starts once review has room.

Source: design/system/workflow.md#The pull order (S-0057)

Source: design/system/workflow.md § The pull order (S-0057)

### 2026-11-15T21:14:23Z orchestrator
Resolved: Plan approved under plan_backlog_stories: T-3380 to T-3383 cover both criteria in four layers, with file touches and figures that fit S-0790.
