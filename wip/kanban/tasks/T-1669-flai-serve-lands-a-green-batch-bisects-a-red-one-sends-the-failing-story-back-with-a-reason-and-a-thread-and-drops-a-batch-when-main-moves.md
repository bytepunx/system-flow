---
id: T-1669
type: task
nature: feature
title: flai serve lands a green batch, bisects a red one, sends the failing story back with a reason and a thread, and drops a batch when main moves
status: backlog
parent: S-0400
owner: alex
created: 2026-10-16T05:02:54Z
updated: 2026-10-16T05:02:54Z
transitions: []
stream: S-0400
tags: [cli]
touches: [flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go]
after: [T-1659, T-1644]
---
# T-1669 flai serve lands a green batch, bisects a red one, sends the failing story back with a reason and a thread, and drops a batch when main moves

## Work

- After a batch is verified:
  - when it passed, land it with T-1659's landing
  - when it failed, bisect it with T-1644, land the stories before the failing one, and leave the ones after it enqueued.
- Send the failing story back from review to in-progress as `flai serve`, with a reason naming the failing tier, its first findings, and the commit. Open a thread on the story for its agent saying the same.
- When main's head moves while a batch is under verification, drop the batch and build it again.
- Test each path with a fake verify.
- Waits for T-1659, which lands, and T-1644, which bisects.

## Done when

`flai serve` ends each batch in landings or a send-back as above, and the tests pass.

## Notes
