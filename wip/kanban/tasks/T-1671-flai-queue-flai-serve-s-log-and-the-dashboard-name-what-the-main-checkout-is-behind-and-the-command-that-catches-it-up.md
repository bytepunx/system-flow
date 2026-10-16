---
id: T-1671
type: task
nature: feature
title: flai queue, flai serve's log, and the dashboard name what the main checkout is behind and the command that catches it up
status: backlog
parent: S-0402
owner: alex
created: 2026-10-16T05:03:03Z
updated: 2026-10-16T05:03:03Z
transitions: []
stream: S-0402
tags: [cli, dashboard]
touches: [flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/cmd/queue.go, flai/cmd/queue_test.go, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts]
after: [T-1662]
---
# T-1671 flai queue, flai serve's log, and the dashboard name what the main checkout is behind and the command that catches it up

## Work

- From the state T-1662 records, `flai queue` prints the commits and paths the checkout is behind and the command that catches it up, as the ADR names it, until it has caught up. `flai serve` logs it once per landing.
- The queue panel shows the same, with the command to copy.
- Clear it once the checkout's head and files match main.
- Test the message, its clearing, and the panel.
- Waits for T-1662, which records it.

## Done when

The operator sees what their checkout is behind wherever the queue is shown, it clears when caught up, and the tests pass.

## Notes
