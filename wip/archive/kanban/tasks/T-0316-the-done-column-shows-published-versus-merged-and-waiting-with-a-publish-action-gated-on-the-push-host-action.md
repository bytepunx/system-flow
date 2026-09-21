---
id: T-0316
type: task
nature: feature
title: The done column shows published versus merged-and-waiting, with a Publish action gated on the push host action
status: done
parent: S-0087
owner: alex
created: 2026-09-21T03:32:55Z
updated: 2026-09-21T04:09:00Z
transitions:
  - to: ready
    at: 2026-09-21T03:58:38Z
    by: system-flow
  - to: in-progress
    at: 2026-09-21T03:58:39Z
    by: system-flow
  - to: done
    at: 2026-09-21T04:09:00Z
    by: system-flow
stream: S-0087
tags: []
---
# T-0316 The done column shows published versus merged-and-waiting, with a Publish action gated on the push host action

## Work
The board's done column shows, per card, whether it is published or merged and waiting (T2's plan, or the absence of one, decides which), and a Publish action at the top of the column, present only when something is pending, that calls T3's batch publish as a host action, gated the same way the automatic push was before this story (`flai serve enable push`, ADR-0031, S-0078). The action streams or reports progress like acceptance already does (NDJSON, or the same pattern `UnpushedNotice` and the review page's accept flow use), and shows the plan (which components, which stories, which bump) before or as it runs.

## Done when
- A done card visibly says published or waiting; the distinction is correct after T1 (merge, no tag) and after a publish (T3)
- The Publish action is absent or disabled with nothing pending, and disabled with an explanation when the push host action is off
- Clicking it shows what will be released (components, bumps, story IDs) and reports what happened, success or flai's error verbatim, the same as accept's confirmation and review page do today
- Component and route tests for all three states: nothing pending, pending with the action enabled, pending with it disabled

## Notes
Depends on T2 and T3 existing to ask of flai. No per-card selection (decided with the operator): the action always covers everything pending.
