---
id: MS-0068
title: "S-0337 holds S-0367 on overlap alone: may S-0367 start beside it?"
from: S-0367
to: S-0337
about: [flai/cmd/accept.go, flai/cmd/accept_conflict_test.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go]
status: closed
participants: [flai, orchestrator]
created: 2026-10-12T16:48:59Z
updated: 2026-10-12T16:51:13Z
---

# MS-0068 S-0337 holds S-0367 on overlap alone: may S-0367 start beside it?

Between S-0367 and S-0337, about `flai/cmd/accept.go`, `flai/cmd/accept_conflict_test.go`, `flai/cmd/stats.go`, `flai/cmd/check_stats_test.go`.

## Entries

### 2026-10-12T16:48:59Z flai S-0367
S-0337 holds S-0367 on overlap alone: may S-0367 start beside it?

S-0367, "Perform lint fix on every markdown file a story touches as part of validation", is ready, and its claim overlaps S-0337's on `flai/cmd/accept.go`, `flai/cmd/accept_conflict_test.go`, `flai/cmd/stats.go`, `flai/cmd/check_stats_test.go`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> Run `scripts/lint-md.sh --fix [markdown file]` for every markdown file that gets touched as part of a story to prevent bad markdown from slipping into main and breaking build checks.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0337 --remove <path>`, or on the task that names the path, if S-0337 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0068 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0367 can start.
3. Reply with `message_reply` saying why the hold stands: S-0367 waits until S-0337 moves to review, is cancelled, or is sent back.

### 2026-10-12T16:51:13Z orchestrator
Closed: S-0337 was accepted
