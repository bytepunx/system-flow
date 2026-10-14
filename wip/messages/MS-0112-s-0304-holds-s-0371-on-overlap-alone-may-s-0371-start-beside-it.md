---
id: MS-0112
title: "S-0304 holds S-0371 on overlap alone: may S-0371 start beside it?"
from: S-0371
to: S-0304
about: [flai/cmd/stream.go]
status: open
participants: [flai]
created: 2026-10-14T20:56:20Z
updated: 2026-10-14T20:56:20Z
---

# MS-0112 S-0304 holds S-0371 on overlap alone: may S-0371 start beside it?

Between S-0371 and S-0304, about `flai/cmd/stream.go`.

## Entries

### 2026-10-14T20:56:20Z flai S-0371
S-0304 holds S-0371 on overlap alone: may S-0371 start beside it?

S-0371, "flai stream sync replays a story's revert of its own change over an identical change on main, undoing main's", is ready, and its claim overlaps S-0304's on `flai/cmd/stream.go`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0124](../../../design/issues/I-0124-flai-stream-sync-replays-a-story-s-revert-of-its-own-change-over-an-identical-change-on-main-undoing-main-s.md), "flai stream sync replays a story's revert of its own change over an identical change on main, undoing main's". The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0304 --remove <path>`, or on the task that names the path, if S-0304 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0112 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0371 can start.
3. Reply with `message_reply` saying why the hold stands: S-0371 waits until S-0304 moves to review, is cancelled, or is sent back.
