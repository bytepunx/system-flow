---
id: T-2375
type: task
nature: improvement
title: An ADR refining ADR-0174 and ADR-0221 says a queue tier whose when matches the story runs in the close-out and again in the batch
status: backlog
parent: S-0547
owner: alex
created: 2026-10-31T07:39:16Z
updated: 2026-10-31T07:39:16Z
transitions: []
stream: S-0547
tags: [flai]
touches: [design/adrs/README.md]
---
# T-2375 An ADR refining ADR-0174 and ADR-0221 says a queue tier whose when matches the story runs in the close-out and again in the batch

## Work

ADR-0174 decided that while the merge queue can batch, `flai verify` and the close-out leave every tier marked `queue` to the batch. This story reverses that for a tier whose `when` matches the story, so the decision needs an ADR of its own before the code and the docs describe it.

Write it with `flai adr new` (the number and slug are flai's, which is why the story keeps the folder touch `design/adrs`), refining ADR-0174 and ADR-0221, and accept it:

- A tier with `queue: true` may carry `when`, a list of path patterns in the language `paths` uses.
- While the queue can batch, `flai verify` and the close-out run such a tier when a pattern of its `when` matches one of the story's touches (its claim: the story's and its open tasks' touches) or a path its branch changes, and leave it to the batch otherwise, as today.
- The batch's verification never leaves a tier out, so it runs the tier again as the landing check.
- A tier run in the close-out at the story's head counts for ADR-0221's hand acceptance, and `flai verify --all` at the same patch reuses it.
- `when` on a tier without `queue: true` is a `manifest.tests` error.
- Context: the fortnight's 24 failed queue verifications and 14 send-backs (E-0031), I-0221, and I-0203.

Waits for nothing: the ADR states the decision the other tasks build to.

## Done when

- The ADR is accepted, refines ADR-0174 and ADR-0221, and is listed in `design/adrs/README.md`.
- `flai check --strict` and the markdown lint pass on it.

## Notes

Drafted by planner-S-0547. The decision's points are the plan's assumptions, named on the plan's thread; change them there or in the story's `## Decisions` if the agent decides otherwise.
