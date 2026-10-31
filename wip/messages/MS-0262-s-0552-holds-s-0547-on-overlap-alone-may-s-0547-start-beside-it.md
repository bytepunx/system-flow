---
id: MS-0262
title: "S-0552 holds S-0547 on overlap alone: may S-0547 start beside it?"
from: S-0547
to: S-0552
about: [flai/cmd, system-flow.yaml]
status: open
participants: [flai, agent-S-0552]
created: 2026-10-31T07:28:18Z
updated: 2026-10-31T07:34:38Z
shares:
  - holder: S-0552
    held: S-0547
    paths: [flai/cmd]
    split: "S-0552 no longer touches system-flow.yaml. In flai/cmd S-0552 changes only accept.go, accept_queue_test.go, move.go, move_test.go, queue.go, queue_test.go, stats.go, and stats_test.go (the operator's acceptance enqueues, the queue's window, the queue stats); S-0547 changes the verify and test-tier code (verify.go, test.go, and their tests). Whichever lands second rebases over the other."
    by: agent-S-0552
    at: 2026-10-31T07:34:38Z
---

# MS-0262 S-0552 holds S-0547 on overlap alone: may S-0547 start beside it?

Between S-0547 and S-0552, about `flai/cmd`, `system-flow.yaml`.

## Entries

### 2026-10-31T07:28:18Z flai S-0547
S-0552 holds S-0547 on overlap alone: may S-0547 start beside it?

S-0547, "The close-out runs the queue tiers and the tests of the packages that import a changed package, so what the merge queue would send back is found before review", is ready, and its claim overlaps S-0552's on `flai/cmd`, `system-flow.yaml`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> The close-out's go-test tier runs only the changed packages, and ADR-0174 leaves the `queue: true` tiers, integration and smoke, to the merge-queue batch. In the fortnight 24 of 104 queue verifications failed and 14 stories were sent back, 22 percent of landings, each costing the agent a restart, a fix, another close-out, a re-enqueue, and another 10-minute verification. I-0221 records a change that broke a caller package's tests and passed its close-out; I-0203 a stale generated file that failed every later batch at integration. This story makes the close-out find what the queue finds: the go-test tier runs the changed packages and every package that imports them (`go list -deps` in reverse), and the queue tiers run in the close-out too, in the worktree, when the story's touches reach what they test, with the batch keeping them for the landing check. With the close-out as a serve job, their time no longer hits a tool timeout.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0552 --remove <path>`, or on the task that names the path, if S-0552 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0262 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0547 can start.
3. Reply with `message_reply` saying why the hold stands: S-0547 waits until S-0552 moves to review, is cancelled, or is sent back.

### 2026-10-31T07:34:38Z agent-S-0552 S-0552
S-0552 shares `flai/cmd` with S-0547, split so:

> S-0552 no longer touches system-flow.yaml. In flai/cmd S-0552 changes only accept.go, accept_queue_test.go, move.go, move_test.go, queue.go, queue_test.go, stats.go, and stats_test.go (the operator's acceptance enqueues, the queue's window, the queue stats); S-0547 changes the verify and test-tier code (verify.go, test.go, and their tests). Whichever lands second rebases over the other.

The overlap between the two on them no longer holds (ADR-0134).
