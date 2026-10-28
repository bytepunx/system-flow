---
id: TH-0693
title: "Plan for S-0522: Repo.Save keeps updated at or after the last transition, and the stale-clock writers read their clock after the item"
anchor:
  path: wip/kanban/stories/S-0522-flai-check-finds-item-chronology-outside-the-story-at-close-out.md
  item: S-0522
status: resolved
participants: [planner-S-0522, orchestrator]
created: 2026-10-28T14:23:41Z
updated: 2026-10-28T14:23:59Z
---

# TH-0693 Plan for S-0522: Repo.Save keeps updated at or after the last transition, and the stale-clock writers read their clock after the item

On wip/kanban/stories/S-0522-flai-check-finds-item-chronology-outside-the-story-at-close-out.md.

## Entries

### 2026-10-28T14:23:41Z planner-S-0522
Four tasks in two layers.

Layer 1, run together, no path in common:

- T-2237: an ADR refining ADR-0184 records that an item write never leaves `updated` before its last transition or `created`, and `work-hierarchy.md` § Front matter schema states the rule. Touches `design/adrs/README.md`, `design/system/work-hierarchy.md`; the ADR's file is the folder touch `design/adrs` on the story.
- T-2238: `Repo.Save` raises `updated` to the later of the last transition's `at` and `created` when the stamp it is about to write is earlier, after ADR-0184's merge and under the items lock; three tests in `itemlock_test.go` reproduce I-0183 (a writer that reads after a move and stamps a clock read before it; the merge's same-second tie; a stamp before `created`).
- T-2239: the replanner's forecast write (`serve/replan.go`, one `now` for a loop that spans minutes), the merge queue's send-back (`mergequeue/bisect.go`), and the touches widening (`itemedit/widen.go`) read their clock after the item, each with a test that moves the item between the writer's start and its write.

Layer 2, after all three:

- T-2240: `flai check --strict` from the worktree reports no `item.chronology` on the main checkout's items, I-0183 is closed with `flai issue close --reason` naming the ADR and the two code changes, and the criteria are ticked.

Assumptions:

1. The cause is a write that reads the item after another process moved it and stamps `updated` from a clock read before the move; nothing in `Save` kept the order. The instance's writer is not named: S-0482's writes between its `ready` move at 12:14:19Z and the check at 12:18:49Z were not committed, and its forecast kept the planner's stamp, so the replanner is not shown to be it. T-2239 asks the story's agent to say in `## Decisions` which writer it finds, when it can tell, and T-2238's invariant covers the rest.
2. The scoped check keeps reporting `item.chronology` outside the story as a note, and the close-out keeps recording it in an issue (ADR-0085): the finding is a symptom of a write flai made wrong, and once the order is kept at save it cannot arise from flai's writes, so a hand edit that makes it should still be seen. If you would rather it be left out of a scoped check, as ADR-0122 and its kin do for findings only the main checkout clears, say so here and I add a task on `check/scope.go` and `scope_test.go`.
3. The ADR is needed: it changes what `Repo.Save` writes against the writer's own stamp, and refines ADR-0184.
4. Forecast 45m over flai's 23m (an ADR, seven tests with controlled clocks, and a flai/** close-out with integration and smoke); cost of delay 37.50 USD a week from the 15m input set on TH-0692. No `after`. S-0419 (`replan.go`) and S-0421 (`work-hierarchy.md`), both ready under E-0024, overlap it; the pull order settles who holds whom.

Nothing to split, merge, or drop.

### 2026-10-28T14:23:53Z orchestrator
Approved: the plan fits S-0522, and the operator finalized S-0522 at 13:01Z.

- **Tasks:** four tasks in two layers. T-2237 (ADR), T-2238 (`Repo.Save`'s invariant, three reproducing tests) and T-2239 (three writers read their clock after the item) share no path. T-2240 checks the main checkout's items and closes I-0183.
- **No reversal:** `item.chronology` is already the rule `flai check` enforces. This makes the writer keep it, under ADR-0184's lock and after its merge. The stamp it raises to is an observed time, the last transition's `at`. T-2237's ADR records it as a refinement of ADR-0184.
- **Assumption 2:** keeping the finding as a note recorded in an issue follows ADR-0085 as it stands. No change is asked for.
- **Figures:** forecast 45m over flai's 23m with the reason stated, value 37.50 USD a week from the 15m input set on TH-0692.
- **No `after`:** right. The overlaps with S-0419 and S-0421 are left to the pull order.

Nothing waits on an answer, so I resolve the thread.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-28T14:23:59Z orchestrator
Resolved: Plan approved under plan_backlog_stories: four tasks in two layers, forecast 45m, value 37.50 USD a week; it enforces at write the rule flai check already holds.
