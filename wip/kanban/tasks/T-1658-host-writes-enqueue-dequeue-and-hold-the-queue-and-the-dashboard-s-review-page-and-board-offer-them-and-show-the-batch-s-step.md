---
id: T-1658
type: task
nature: feature
title: Host writes enqueue, dequeue, and hold the queue, and the dashboard's review page and board offer them and show the batch's step
status: backlog
parent: S-0399
owner: alex
created: 2026-10-16T05:02:03Z
updated: 2026-10-16T05:02:03Z
transitions: []
stream: S-0399
tags: [cli, dashboard]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, flaiover/src/routes/board/+page.svelte, "flaiover/src/routes/api/items/[id]/queue/+server.ts", flaiover/src/routes/api/queue/+server.ts]
after: [T-1640, T-1641]
---
# T-1658 Host writes enqueue, dequeue, and hold the queue, and the dashboard's review page and board offer them and show the batch's step

## Work

- Host writes `queue.add`, `queue.remove`, and `queue.hold`, gated by the `merge-queue` host action, running the `flai queue` commands (ADR-0016).
- The review page gets a queue button: enqueue, or dequeue when enqueued. Show a refusal's reasons verbatim.
- The board gets a hold switch beside the queue's panel, and the panel shows the batch's step as it changes.
- Test the writes' gating, the button, the switch, and the step's display.
- Waits for T-1640, which adds the host action to `writes.go` first, and for T-1641, whose commands the writes run.

## Done when

The operator can enqueue, dequeue, and hold from the dashboard, the step shows, and the Go and vitest tests pass.

## Notes
