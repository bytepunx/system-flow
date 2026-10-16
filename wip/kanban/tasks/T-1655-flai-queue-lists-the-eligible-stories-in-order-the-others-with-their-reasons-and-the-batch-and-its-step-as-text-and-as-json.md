---
id: T-1655
type: task
nature: feature
title: flai queue lists the eligible stories in order, the others with their reasons, and the batch and its step, as text and as JSON
status: backlog
parent: S-0398
owner: alex
created: 2026-10-16T05:01:49Z
updated: 2026-10-16T05:01:49Z
transitions: []
stream: S-0398
tags: [cli]
touches: [flai/cmd/queue.go, flai/cmd/queue_test.go, flai/cmd/root.go]
after: [T-1639]
---
# T-1655 flai queue lists the eligible stories in order, the others with their reasons, and the batch and its step, as text and as JSON

## Work

- New command `flai queue`, registered in `root.go`, reading through `mergequeue`.
- Print the eligible stories in queue order, the other stories in review with their reasons, and the batch and its step, or none.
- `--json` gives the same as data, with each reason's `kind`.
- Test the output on a fixture board with eligible and ineligible stories and a state file.
- Waits for T-1639, whose package it reads.

## Done when

`flai queue` and `flai queue --json` print the queue as above, and the tests pass.

## Notes
