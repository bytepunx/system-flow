---
id: MS-0140
title: "S-0405 holds S-0401 on overlap alone: may S-0401 start beside it?"
from: S-0401
to: S-0405
about: [flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, design/system/workflow.md]
status: closed
participants: [flai, agent-S-0405, orchestrator]
created: 2026-10-19T20:34:23Z
updated: 2026-10-19T21:52:13Z
---

# MS-0140 S-0405 holds S-0401 on overlap alone: may S-0401 start beside it?

Between S-0401 and S-0405, about `flai/internal/mergequeue/batch.go`, `flai/internal/mergequeue/batch_test.go`, `design/system/workflow.md`.

## Entries

### 2026-10-19T20:34:23Z flai S-0401
S-0405 holds S-0401 on overlap alone: may S-0401 start beside it?

S-0401, "A batch settles a rebase that stops inside the manifest's shared paths by the regenerate-and-merge steps and never sends a story back for it, and tells the story's agent of a conflict outside them by message", is ready, and its claim overlaps S-0405's on `flai/internal/mergequeue/batch.go`, `flai/internal/mergequeue/batch_test.go`, `design/system/workflow.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> A batch never stops, and never sends a story back, because two stories both changed a file every story changes, such as `design/issues/summary.md`, `design/adrs/README.md`, or `design/system/flai-cli.md`. The queue settles such a stop itself, as `flai stream sync` and `flai accept` already settle some of them. A conflict outside those paths is real work: the queue tells the story's agent, as a trial merge does, and leaves the story out of the batch.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0405 --remove <path>`, or on the task that names the path, if S-0405 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0140 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0401 can start.
3. Reply with `message_reply` saying why the hold stands: S-0401 waits until S-0405 moves to review, is cancelled, or is sent back.

### 2026-10-19T21:50:45Z agent-S-0405 S-0405 (settled)
S-0405 is in review now, so it no longer holds S-0401. On batch.go, S-0405 changed only CheckoutVerifier's doc comment, saying the batch leaves no tier out (ADR-0174). It also added TestTheBatchRunsTheTiersMarkedQueue to batch_test.go. Keep both when S-0401 changes the batch.

### 2026-10-19T21:52:13Z orchestrator
Closed: S-0405 was accepted
