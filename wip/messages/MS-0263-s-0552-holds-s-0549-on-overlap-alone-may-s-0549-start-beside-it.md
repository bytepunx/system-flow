---
id: MS-0263
title: "S-0552 holds S-0549 on overlap alone: may S-0549 start beside it?"
from: S-0549
to: S-0552
about: [flai/internal/serve, flai/internal/mergequeue]
status: open
participants: [flai, agent-S-0552]
created: 2026-10-31T07:28:18Z
updated: 2026-10-31T07:34:40Z
shares:
  - holder: S-0552
    held: S-0549
    paths: [flai/internal/mergequeue]
    split: "S-0552 no longer touches flai/internal/serve. In flai/internal/mergequeue S-0552 changes the coalesce window in RunBatch and pick (batch.go), the window on Queue (eligible.go), and the bisection steps' verification events (bisect.go, log.go), with their tests; S-0549 adds the preflight of main before a batch is built and the landing again without verification, in a file of its own where it can, called from RunBatch and Land with as small an edit to batch.go and land.go as it needs. Whichever lands second rebases over the other."
    by: agent-S-0552
    at: 2026-10-31T07:34:40Z
---

# MS-0263 S-0552 holds S-0549 on overlap alone: may S-0549 start beside it?

Between S-0549 and S-0552, about `flai/internal/serve`, `flai/internal/mergequeue`.

## Entries

### 2026-10-31T07:28:18Z flai S-0549
S-0552 holds S-0549 on overlap alone: may S-0549 start beside it?

S-0549, "The merge queue checks the main checkout's uncommitted files and the lint of main before it builds a batch, so a dirty main fails once in seconds and not every batch", is ready, and its claim overlaps S-0552's on `flai/internal/serve`, `flai/internal/mergequeue`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> A batch verifies its stories' merge on top of main, then lands it through the main checkout, where the acceptance preflight lints the checkout's uncommitted files too. In the fortnight 80 verifications passed for 49 landings: S-0541 went through 16 batches and 28 failed landings, about 2.7 hours of verification, because one uncommitted task note in the main checkout had a bad Markdown link (MD051) that nothing checked before the batch was built, and I-0203's stale generated file failed every later batch the same way. This story moves that check to the front: before a batch is built, the queue lints the main checkout's uncommitted files and runs `flai check` on main, and when either fails it opens one thread to the operator naming the file and holds the queue, instead of building, verifying, and failing a batch for each enqueue. A batch whose verification passed but whose landing failed for a reason outside its stories is landed again when main is clean, without another verification.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0552 --remove <path>`, or on the task that names the path, if S-0552 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0263 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0549 can start.
3. Reply with `message_reply` saying why the hold stands: S-0549 waits until S-0552 moves to review, is cancelled, or is sent back.

### 2026-10-31T07:34:40Z agent-S-0552 S-0552
S-0552 shares `flai/internal/mergequeue` with S-0549, split so:

> S-0552 no longer touches flai/internal/serve. In flai/internal/mergequeue S-0552 changes the coalesce window in RunBatch and pick (batch.go), the window on Queue (eligible.go), and the bisection steps' verification events (bisect.go, log.go), with their tests; S-0549 adds the preflight of main before a batch is built and the landing again without verification, in a file of its own where it can, called from RunBatch and Land with as small an edit to batch.go and land.go as it needs. Whichever lands second rebases over the other.

The overlap between the two on them no longer holds (ADR-0134).
