---
id: T-1659
type: task
nature: feature
title: A green batch lands each story in queue order through the acceptance flow, as whoever enqueued it
status: backlog
parent: S-0400
owner: alex
created: 2026-10-16T05:02:07Z
updated: 2026-10-16T05:02:07Z
transitions: []
stream: S-0400
tags: [cli]
touches: [flai/internal/mergequeue/land.go, flai/internal/mergequeue/land_test.go]
after: [T-1643]
---
# T-1659 A green batch lands each story in queue order through the acceptance flow, as whoever enqueued it

## Work

- `land.go`: for each story of a passing batch, in order, run the acceptance with the batch named, `by` whoever enqueued it. Pass the orchestrator's recorded evidence when it enqueued.
- Stop at the first acceptance refused. Record its blockers as that story's outcome. Leave the stories after it enqueued.
- Remove each landed story from the queue.
- Test that a landed story's done transition, acceptance commit, archive, and notes match an acceptance by hand of the same story.
- Waits for T-1643, whose acceptance from a batch it calls.

## Done when

A passing batch lands each story as an acceptance by hand would, a refusal stops the landing there, and the tests pass.

## Notes
