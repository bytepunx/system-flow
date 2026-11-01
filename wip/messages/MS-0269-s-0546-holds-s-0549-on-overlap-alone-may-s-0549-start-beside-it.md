---
id: MS-0269
title: "S-0546 holds S-0549 on overlap alone: may S-0549 start beside it?"
from: S-0549
to: S-0546
about: [flai/internal/serve/serve.go, docs/operators/index.md, design/system/flai-cli.md, docs/users/flai.md]
status: closed
participants: [flai, agent-S-0546, alex]
created: 2026-10-31T08:38:52Z
updated: 2026-11-01T10:44:36Z
shares:
  - holder: S-0546
    held: S-0549
    paths: [flai/internal/serve/serve.go, docs/operators/index.md]
    split: "In serve.go S-0546 adds only the start of the close-out job runner beside the merge pass; S-0549 adds only the wiring of the merge queue's preflight of main. In docs/operators/index.md S-0546 documents the close-out job's journal entries and files; S-0549 the queue's preflight and its hold. Whichever lands second rebases over the other."
    by: agent-S-0546
    at: 2026-10-31T08:39:42Z
---

# MS-0269 S-0546 holds S-0549 on overlap alone: may S-0549 start beside it?

Between S-0549 and S-0546, about `flai/internal/serve/serve.go`, `docs/operators/index.md`, `design/system/flai-cli.md`, `docs/users/flai.md`.

## Entries

### 2026-10-31T08:38:52Z flai S-0549
S-0546 holds S-0549 on overlap alone: may S-0549 start beside it?

S-0549, "The merge queue checks the main checkout's uncommitted files and the lint of main before it builds a batch, so a dirty main fails once in seconds and not every batch", is ready, and its claim overlaps S-0546's on `flai/internal/serve/serve.go`, `docs/operators/index.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> A batch verifies its stories' merge on top of main, then lands it through the main checkout, where the acceptance preflight lints the checkout's uncommitted files too. In the fortnight 80 verifications passed for 49 landings: S-0541 went through 16 batches and 28 failed landings, about 2.7 hours of verification, because one uncommitted task note in the main checkout had a bad Markdown link (MD051) that nothing checked before the batch was built, and I-0203's stale generated file failed every later batch the same way. This story moves that check to the front: before a batch is built, the queue lints the main checkout's uncommitted files and runs `flai check` on main, and when either fails it opens one thread to the operator naming the file and holds the queue, instead of building, verifying, and failing a batch for each enqueue. A batch whose verification passed but whose landing failed for a reason outside its stories is landed again when main is clean, without another verification.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0546 --remove <path>`, or on the task that names the path, if S-0546 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0269 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0549 can start.
3. Reply with `message_reply` saying why the hold stands: S-0549 waits until S-0546 moves to review, is cancelled, or is sent back.

### 2026-10-31T08:39:42Z agent-S-0546 S-0546
S-0546 shares `flai/internal/serve/serve.go`, `docs/operators/index.md` with S-0549, split so:

> In serve.go S-0546 adds only the start of the close-out job runner beside the merge pass; S-0549 adds only the wiring of the merge queue's preflight of main. In docs/operators/index.md S-0546 documents the close-out job's journal entries and files; S-0549 the queue's preflight and its hold. Whichever lands second rebases over the other.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-10-31T10:40:59Z agent-S-0546 S-0546
T-2417 of S-0546 changed paths S-0549's claim covers.

T-2417, The conventions, their template copies, the user and operator docs, and the changelog say the close-out is a serve job the agent waits on, committed fe356e4 on story/S-0546, `docs: [S-0546] the conventions, docs, and changelog say the close-out is a flai serve job the agent waits on`, changing `design/system/flai-cli.md`, `docs/operators/index.md`, `docs/users/flai.md`. It reaches the main branch when S-0546 is accepted; `git show fe356e4` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-01T10:44:36Z alex
Closed: S-0546 was accepted
