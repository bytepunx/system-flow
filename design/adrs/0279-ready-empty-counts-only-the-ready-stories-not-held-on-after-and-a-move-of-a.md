---
id: ADR-0279
title: "ready-empty counts only the ready stories not held on after, and a move of a story to ready warns when a story it names in after is not done"
status: accepted
date: 2026-11-10
supersedes: []
superseded_by: []
refines: [ADR-0211, ADR-0230]
---

# ADR-0279 ready-empty counts only the ready stories not held on after, and a move of a story to ready warns when a story it names in after is not done

## Context

[I-0285](../issues/I-0285-after-chains-on-planned-stories-hold-ready-stories-for-days-while-in-progress.md), found by the analysis of 2026-11-07, records that `after` holds kept ready stories waiting for days while in progress stayed under two of a limit of five. W44 held 704,200 s under `after` and W45 203,155 s; `concurrency.mean_busy` was 1.81 of 5.

Its instances show three things.

- The longest holds were on stories planned before [ADR-0230](0230-a-ready-story-whose-claim-overlaps-a-story-in-progress-is-held-only-when-a.md) made `after` name only a story whose output it reads: S-0425 (326,792 s), S-0424 (133,349 s), S-0423, S-0422, and S-0420, all made on 2026-10-16 in E-0024's chain, and S-0554 (74,959 s), made on 2026-10-30. Of the stories with `after` started since 2026-11-03, the longest hold is S-0579's 6,259 s, and W46 has none. The planner's rule needs no change.
- Each long hold sat in ready because it was moved there by hand while a story it names in `after` was not done: S-0425 on 2026-10-30, S-0424 and S-0554 on 2026-11-01. `flai promote --candidates` leaves such a story out, but `flai move` to ready says nothing of the hold; only the move to in-progress warns.
- `ready-empty` ([ADR-0211](0211-the-orchestrator-s-trigger-vocabulary-gains-the-events-moved-thread-planned.md) §1) counts the ready stories, held ones counted. With as many held stories in ready as the in-progress limit leaves room for, no `ready-empty` is told, the orchestrator is not started to promote, and lanes stand empty. This is how in progress stayed at 1 to 4 of 5 while the stories held reached 6 on 2026-11-01.

## Decision

**`ready-empty` counts only the ready stories not held on `after`, and a move of a story to ready warns when a story it names in `after` is not done.** This refines ADR-0211 §1 and ADR-0230; the rest of each stands.

1. **What ready-empty counts.** A ready story held on `after` (hold code `after`, alone or with an overlap) is not counted against the in-progress room. A story held on overlap alone is still counted: a trial merge may clear it at the next look ([ADR-0230](0230-a-ready-story-whose-claim-overlaps-a-story-in-progress-is-held-only-when-a.md)), and the room it would fill is not empty. A `no-touches` hold is counted too, since its fix is a touch, not an acceptance. A blocked story is left out, as before. The event's `to` is that count, and its words and the orchestrator's reason say the ready stories not held on `after`.
2. **The move to ready warns.** `flai move`, `item_move`, and the dashboard's move of a story to ready warn, as the move to in-progress warns of a hold, when a story it names in `after` is not done: the warning gives the hold's reason, which names each such story and what clears it. It is a warning, not a refusal: the move is the operator's, and a story may wait in ready on purpose.

## Consequences

- A ready column full of stories held on `after` drains in the event's eyes, so the orchestrator is started to promote stories that can start, and the lanes fill.
- The operator who moves a story to ready by hand learns at once that it will wait, and why.
- `ready-empty` can be told while the ready column holds stories, which is what the room needs. A later acceptance that clears an `after` hold can lift the count to the room again, and the next draining is told again, once.
- `flai stats` and its replay of holds are unchanged.

## Alternatives considered

- **Refuse a move to ready of a story held on `after`.** The operator may want a story at the head of ready, waiting for its predecessor; a refusal takes that choice away.
- **Leave out every held story, overlap too.** A trial merge clears most overlap holds at the next look; counting them out would start the orchestrator to promote stories that the next look starts anyway.
- **Change the planner's rule.** ADR-0230 already did; since then the `after` holds are short.
