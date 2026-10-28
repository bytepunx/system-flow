---
id: ADR-0208
title: "Repo.Save never leaves a work item's updated before its last transition or its created, and a writer takes its clock after it reads the item it stamps"
status: accepted
date: 2026-10-28
supersedes: []
superseded_by: []
refines: [ADR-0184]
topics: [cli]
---

# ADR-0208 Repo.Save never leaves a work item's updated before its last transition or its created, and a writer takes its clock after it reads the item it stamps

## Context

`flai check` reports `item.chronology` as an error when an item's `updated` is before the last transition's `at` or before `created` (`flai/internal/check/check.go`, `history`). A move stamps the transition and `updated` from one clock (`flai/internal/workitem/rules.go`), so a move alone keeps the order. Nothing in `Repo.Save` (`flai/internal/workitem/repo.go`, `itemlock.go`) keeps it, though. Several writers take their clock before they read the item: the replanner's forecast write (`flai/internal/serve/replan.go`), the merge queue's send-back (`flai/internal/mergequeue/bisect.go`), and the touches widening (`flai/internal/itemedit/widen.go`). Such a writer can read the item after another process moved it and stamp `updated` from a clock read before the move. The item then has an `updated` earlier than its last transition.

I-0183 records one instance. At S-0473's close-out on 2026-10-21T12:18:49Z, S-0482's file in the main checkout had its `ready` transition of 12:14:19Z by the orchestrator and an `updated` before it. ADR-0184's merge landed at 12:40Z, after the occurrence, when a save was still the last writer's. The merge keeps the order no better: `mergeParts` in `flai/internal/workitem/itemlock.go` takes the later `updated` only when both sides changed it, and the later of two writes may carry the earlier clock.

## Decision

No write flai makes leaves an item that `item.chronology` reports, in two parts.

- `Repo.Save` raises `updated` to the later of the last transition's `at` and `created` when the stamp it is about to write is earlier. It does so under the items lock and after the merge (ADR-0184), so it sees the transitions another writer saved.
- A writer takes its clock after it has read the item it stamps, so that `updated` is the write's own time. This covers the replanner's forecast write, the merge queue's send-back, and the touches widening.

## Consequences

- `flai check` no longer finds `item.chronology` on a file only flai wrote, in a story's check or at a close-out.
- `updated` may read later than the clock a writer passed, by the gap between its read and the move it raced.
- `flai check` on the main checkout still reports a file a hand edit left out of order.

## Alternatives considered

- Refuse the save. The writer cannot do better than write again with a later clock, which is what raising does.
- Leave `item.chronology` out of a check scoped to a story, as ADR-0122 and its kin do for findings only the main checkout clears. The finding is a symptom of a write flai made wrong, and hiding it hides the defect. With the order kept at save it cannot arise from flai's writes, and a hand edit that makes it should still be seen.
