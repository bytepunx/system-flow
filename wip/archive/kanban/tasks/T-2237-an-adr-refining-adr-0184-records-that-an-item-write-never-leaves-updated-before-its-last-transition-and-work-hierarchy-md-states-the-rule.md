---
id: T-2237
type: task
nature: improvement
title: An ADR refining ADR-0184 records that an item write never leaves updated before its last transition, and work-hierarchy.md states the rule
status: done
parent: S-0522
owner: alex
created: 2026-10-28T14:20:59Z
updated: 2026-10-28T14:38:24Z
transitions:
  - to: ready
    at: 2026-10-28T14:25:03Z
    by: agent-S-0522
  - to: in-progress
    at: 2026-10-28T14:25:04Z
    by: agent-S-0522
  - to: done
    at: 2026-10-28T14:38:24Z
    by: agent-S-0522
stream: S-0522
tags: [cli]
touches: [design/adrs/README.md, design/system/work-hierarchy.md, design/adrs/0208-repo-save-never-leaves-a-work-item-s-updated-before-its-last-transition-or-its.md]
usage:
  source: log
  seconds: 800
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 19
      output: 7506
      cache_read: 800609
      cache_write: 42641
      cost: 0.5875
---
# T-2237 An ADR refining ADR-0184 records that an item write never leaves updated before its last transition, and work-hierarchy.md states the rule

## Work

Write the ADR with `flai adr new`, refining ADR-0184, from I-0183's instance. Context: `flai check`'s `item.chronology` (`flai/internal/check/check.go`, `history`) is an error when `updated` is before the last transition's `at` or before `created`; a move stamps the transition and `updated` from one clock (`flai/internal/workitem/rules.go`), but nothing in `Repo.Save` (`flai/internal/workitem/repo.go`, `itemlock.go`) keeps the order, and several writers take their clock before they read the item, so a write that reads the item after another process moved it, and stamps `updated` from a clock read before the move, leaves `updated` earlier than the transition. The instance: at S-0473's close-out on 2026-10-21T12:18:49Z, S-0482's file in the main checkout had its `ready` transition of 12:14:19Z by the orchestrator and an `updated` before it; ADR-0184's merge landed at 12:40Z, after the occurrence, when a save was still the last writer's; the merge of ADR-0184 keeps the order no better, since `updated` takes the later value only when both sides changed it, and the later of two writes may carry the earlier clock. Decision, in two parts. First, `Repo.Save` raises `updated` to the later of the last transition's `at` and `created` when the stamp it is about to write is earlier, under the items lock and after the merge, so that no write flai makes leaves an item that `item.chronology` reports. Second, a writer takes its clock after it has read the item it stamps, so that `updated` is the write's own time: the replanner's forecast write, the merge queue's send-back, and the touches widening. Alternatives considered: refusing the save (the writer cannot do better than write again with a later clock, which is what raising does); leaving `item.chronology` out of a check scoped to a story, as ADR-0122 and its kin do for findings only the main checkout clears (the finding is a symptom of a write flai made wrong, and hiding it hides the defect; with the order kept at save it cannot arise from flai's writes, and a hand edit that makes it should still be seen). Consequences: `updated` may read later than the clock a writer passed, by the gap between its read and the move; `flai check` on the main checkout still reports a hand-edited file. Give it the topics of ADR-0184 with `flai adr topics` and accept it; `flai adr new` writes the index row in `design/adrs/README.md`.

In `design/system/work-hierarchy.md` § Front matter schema, under Rules, add after the line on `transitions` as the source of truth that `updated` is never before the last transition's `at` or before `created`: `Repo.Save` raises it when a writer's clock is earlier, and `flai check` reports `item.chronology` on a file a hand edit left otherwise, linking the ADR. Bump the document's `updated`.

This task waits for nothing: it records the decision while the task on `Repo.Save` implements it.

## Done when

The ADR is under `design/adrs/` with status `accepted`, `refines` names ADR-0184, it names I-0183 and the S-0482 instance, its row is in `design/adrs/README.md`, `design/system/work-hierarchy.md` states the rule with the ADR linked, and `flai check --strict` and the markdown lint are clean on the files changed.

## Notes

The ADR's number is the next free one when the task runs: S-0502 (in progress) has taken 0206, and S-0518's T-2223 takes the next when it runs.
