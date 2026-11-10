---
id: TH-0968
title: "Plan for S-0648: Save refuses a moved item, Archive locks, Measure follows the item, the stray T-2760 copy goes"
anchor:
  path: wip/kanban/stories/S-0648-flai-check-finds-item-duplicate-id-outside-the-story-at-close-out.md
  item: S-0648
status: resolved
participants: [planner-S-0648, orchestrator, agent-S-0648]
created: 2026-11-10T06:37:54Z
updated: 2026-11-10T08:10:58Z
---

# TH-0968 Plan for S-0648: Save refuses a moved item, Archive locks, Measure follows the item, the stray T-2760 copy goes

On wip/kanban/stories/S-0648-flai-check-finds-item-duplicate-id-outside-the-story-at-close-out.md.

## Entries

### 2026-11-10T06:37:54Z planner-S-0648
S-0648 is planned: four tasks in three layers, touches predicted file by file, a forecast of 1h, and a cost of delay of 5.00 USD a week.

## The cause

Every I-0323 instance is the open copy `wip/kanban/tasks/T-2760-fetch-the-remote-s-branch-and-tags.md`. It sits beside its archived copy from S-0642's acceptance, commit 436e8f612. The two copies differ only in `usage`, so a usage measurement wrote the open copy. The explorer read this from the code; it is not reproduced yet:

- `serve.Measure`'s `put` reads a task and then saves it.
- `Archive` moves the file without the items lock.
- `mergeSaved` takes a missing file as "nothing saved since" and writes the whole item back at its old path.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2827 `Repo.Save` refuses to recreate an item file moved since it was read, and `Archive` moves under the items lock | `flai/internal/workitem/itemlock.go`, `itemlock_test.go`, `archive.go`, `archive_test.go`, `design/system/work-hierarchy.md` | none |
| 1 | T-2829 Remove the stray open copy of T-2760 | `wip/kanban/tasks/T-2760-fetch-the-remote-s-branch-and-tags.md` | none |
| 2 | T-2828 `Measure` writes a usage to the item where it now is | `flai/internal/serve/usage.go`, `usage_test.go` | T-2827, whose refusal it handles |
| 3 | T-2830 Close I-0323 | `design/issues/I-0323-…md`, `design/issues/summary.md` | T-2827, T-2828, T-2829 |

T-2827 and T-2829 share no path and can run together.

## Assumptions

- A `Save` of an item read from a file that is now gone is refused, not re-homed into `wip/archive`. Only `Measure`, which knows the write is a usage, reads again and writes to the archived file. An archived item may so gain its measured usage after acceptance.
- `Archive` can take the items lock with no caller deadlocking. T-2827 checks `flai/cmd/accept.go` and `flai/cmd/archive.go`, and changes a caller only if one holds the lock.
- T-2760's archived copy stays as the acceptance wrote it, with the estimate, cost 0.2198. The harness figure, 0.3414, goes with the deleted copy and is recorded in T-2829's Notes. If you would rather fold the harness figure into the archived copy, say so and T-2829 changes.
- No ADR: the rule is written into `design/system/work-hierarchy.md`'s ADR-0184 bullet. If the story's agent judges it a new decision, it adds one, and `flai task done` widens the touches.
- The cost of delay input, `time_lost_per_cycle: 2m`, is the one the orchestrator set on TH-0967.

Nothing to split, merge, or drop.

### 2026-11-10T06:39:40Z orchestrator (recommendation)
**Recommended:** take the plan as written, with one change: T-2827 records its rule in a new ADR refining ADR-0184. Its touches, and the story's, gain `design/adrs/<new>.md` and `design/adrs/README.md`. The planner adds them on its next run.

**Why I do not approve it outright.**

- ADR-0184 § Decision says "A missing file, a file that does not parse, or an item made in memory is written as before." T-2827 makes `Repo.Save` refuse to recreate an item file that was there when it was read. That reverses part of a tool-enforced rule. `design/conventions/decisions.md` asks for an ADR for it, not only the `design/system/work-hierarchy.md` bullet the plan names.
- S-0648's goal asks for the remedy to be proposed before it is built. Confirming it is your decision, not mine.

**The plan otherwise fits S-0648.**

- Criterion 1: T-2827 (the reproducing test and the `Save`/`Archive` fix) and T-2828 (`Measure` writes to the archived copy). T-2829 removes the stray T-2760 copy.
- Criterion 2: T-2830.
- The layers follow the real dependencies. T-2828 handles T-2827's refusal. T-2827 and T-2829 share no path and run together.
- The 1h forecast states why it is over flai's 22m. The 5 USD a week value comes from the `time_lost_per_cycle: 2m` set on TH-0967.
- The planner's assumptions are sound. Keeping T-2760's archived copy as written, with the harness figure recorded in T-2829's Notes, leaves `wip/archive` unedited by hand.

This is the same change TH-0969 asks for before I finalize the draft. Confirming this recommendation settles both.

Source: design/adrs/0184-every-work-item-write-takes-one-items-lock-and-merges-a-copy-read-before.md § Decision

### 2026-11-10T08:09:45Z agent-S-0648
S-0648 is in progress with the plan as written and the one change you recommend. T-2827 records its rule in a new ADR that refines ADR-0184, written with `flai adr new` and committed with the task, and `flai task done` widens the touches to it. `design/conventions/decisions.md` asks for an ADR when a tool-enforced rule changes, so I am following it without waiting for a confirmation. If the operator would rather have no ADR, say so here and I will drop it before review.

Source: design/conventions/decisions.md § Rules

### 2026-11-10T08:10:34Z orchestrator
Approved. The operator finalized S-0648, and its agent works the plan with the one change I recommended: T-2827 records the new rule in an ADR refining ADR-0184, and `flai task done` widens the touches to it. That meets `design/conventions/decisions.md` for a change to a tool-enforced rule. Nothing more is owed on this plan, so I resolve this thread.

Source: design/conventions/decisions.md § Rules

### 2026-11-10T08:10:58Z orchestrator
Resolved: Plan approved: the operator finalized S-0648 and its agent works it with T-2827 adding an ADR refining ADR-0184.
