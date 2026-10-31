---
id: MS-0266
title: "S-0550 holds S-0549 on overlap alone: may S-0549 start beside it?"
from: S-0549
to: S-0550
about: [flai/internal/check/check.go, flai/internal/check/fix.go, flai/internal/check/fix_test.go, flai/internal/check/issues_test.go]
status: open
participants: [flai, agent-S-0550]
created: 2026-10-31T08:16:11Z
updated: 2026-10-31T08:16:55Z
shares:
  - holder: S-0550
    held: S-0549
    paths: [flai/internal/check/check.go, flai/internal/check/fix.go, flai/internal/check/fix_test.go, flai/internal/check/issues_test.go]
    split: "S-0550 adds the generated-index checks (issues summary and ADR index behind their folders) and flai check --fix: new code in check.go's rule list, the new fix.go and fix_test.go, and cases in issues_test.go. S-0549 adds the merge queue's preflight of main, calling check from mergequeue, and does not edit fix.go, fix_test.go, or issues_test.go; any check.go edit it needs is outside the index rules. Whichever lands second rebases over the other."
    by: agent-S-0550
    at: 2026-10-31T08:16:55Z
---

# MS-0266 S-0550 holds S-0549 on overlap alone: may S-0549 start beside it?

Between S-0549 and S-0550, about `flai/internal/check/check.go`, `flai/internal/check/fix.go`, `flai/internal/check/fix_test.go`, `flai/internal/check/issues_test.go`.

## Entries

### 2026-10-31T08:16:11Z flai S-0549
S-0550 holds S-0549 on overlap alone: may S-0549 start beside it?

S-0549, "The merge queue checks the main checkout's uncommitted files and the lint of main before it builds a batch, so a dirty main fails once in seconds and not every batch", is ready, and its claim overlaps S-0550's on `flai/internal/check/check.go`, `flai/internal/check/fix.go`, `flai/internal/check/fix_test.go`, `flai/internal/check/issues_test.go`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> A batch verifies its stories' merge on top of main, then lands it through the main checkout, where the acceptance preflight lints the checkout's uncommitted files too. In the fortnight 80 verifications passed for 49 landings: S-0541 went through 16 batches and 28 failed landings, about 2.7 hours of verification, because one uncommitted task note in the main checkout had a bad Markdown link (MD051) that nothing checked before the batch was built, and I-0203's stale generated file failed every later batch the same way. This story moves that check to the front: before a batch is built, the queue lints the main checkout's uncommitted files and runs `flai check` on main, and when either fails it opens one thread to the operator naming the file and holds the queue, instead of building, verifying, and failing a batch for each enqueue. A batch whose verification passed but whose landing failed for a reason outside its stories is landed again when main is clean, without another verification.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0550 --remove <path>`, or on the task that names the path, if S-0550 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0266 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0549 can start.
3. Reply with `message_reply` saying why the hold stands: S-0549 waits until S-0550 moves to review, is cancelled, or is sent back.

### 2026-10-31T08:16:55Z agent-S-0550 S-0550
S-0550 shares `flai/internal/check/check.go`, `flai/internal/check/fix.go`, `flai/internal/check/fix_test.go`, `flai/internal/check/issues_test.go` with S-0549, split so:

> S-0550 adds the generated-index checks (issues summary and ADR index behind their folders) and flai check --fix: new code in check.go's rule list, the new fix.go and fix_test.go, and cases in issues_test.go. S-0549 adds the merge queue's preflight of main, calling check from mergequeue, and does not edit fix.go, fix_test.go, or issues_test.go; any check.go edit it needs is outside the index rules. Whichever lands second rebases over the other.

The overlap between the two on them no longer holds (ADR-0134).
