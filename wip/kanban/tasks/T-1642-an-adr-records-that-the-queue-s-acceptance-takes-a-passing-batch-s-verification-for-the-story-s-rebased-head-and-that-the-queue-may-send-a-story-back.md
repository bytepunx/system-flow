---
id: T-1642
type: task
nature: feature
title: An ADR records that the queue's acceptance takes a passing batch's verification for the story's rebased head, and that the queue may send a story back
status: backlog
parent: S-0400
owner: alex
created: 2026-10-16T05:00:55Z
updated: 2026-10-16T05:00:55Z
transitions: []
stream: S-0400
tags: [cli]
touches: [design/adrs]
---
# T-1642 An ADR records that the queue's acceptance takes a passing batch's verification for the story's rebased head, and that the queue may send a story back

## Work

- Write the ADR with `flai adr new`, refining ADR-0093 and ADR-0032. It records:
  - that a green batch lands each story through `flai accept`'s flow, as whoever enqueued it
  - that the verified commit is the story's rebased head in a batch that passed
  - how a red batch is bisected
  - that `flai serve` moves the failing story from review to in-progress with a reason and a thread, a second maker of that move besides a person.
- First in its layer: it waits for no task.

## Done when

The ADR is accepted, indexed, and `flai check --strict` passes.

## Notes
