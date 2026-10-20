---
id: T-1671
type: task
nature: feature
title: flai queue, flai serve's log, and the dashboard name what the main checkout is behind and the command that catches it up
status: done
parent: S-0402
owner: alex
created: 2026-10-16T05:03:03Z
updated: 2026-10-20T04:44:22Z
transitions:
  - to: ready
    at: 2026-10-20T04:41:01Z
    by: agent-S-0402
  - to: in-progress
    at: 2026-10-20T04:41:01Z
    by: agent-S-0402
  - to: done
    at: 2026-10-20T04:44:22Z
    by: agent-S-0402
stream: S-0402
tags: [cli, dashboard]
touches: [flai/cmd/queue.go, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, flai/cmd/board.go, flai/cmd/queue_held_test.go, design/components/flaiover/src/lib/components/MergeQueue.examples.json]
after: [T-1662]
usage:
  source: log
  seconds: 201
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 34
      output: 11541
      cache_read: 3163372
      cache_write: 64364
      cost: 1.3392
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
