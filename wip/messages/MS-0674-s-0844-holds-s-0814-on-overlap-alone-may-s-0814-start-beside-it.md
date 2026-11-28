---
id: MS-0674
title: "S-0844 holds S-0814 on overlap alone: may S-0814 start beside it?"
from: S-0814
to: S-0844
about: [flai/internal/workitem/boardview_test.go, flai/internal/hostapi/writes_test.go, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
status: closed
participants: [flai, agent-S-0844, orchestrator]
created: 2026-11-28T00:21:40Z
updated: 2026-11-28T00:43:34Z
---

# MS-0674 S-0844 holds S-0814 on overlap alone: may S-0814 start beside it?

Between S-0814 and S-0844, about `flai/internal/workitem/boardview_test.go`, `flai/internal/hostapi/writes_test.go`, `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`.

## Entries

### 2026-11-28T00:21:40Z flai S-0814
S-0844 holds S-0814 on overlap alone: may S-0814 start beside it?

S-0814, "The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream", is ready, and its claim overlaps S-0844's on `flai/internal/workitem/boardview_test.go`, `flai/internal/hostapi/writes_test.go`, `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> [ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §4. The review page's verification section shows the last report's tiers, and after S-0804 a Rebase and resolve job's stream.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0844 --remove <path>`, or on the task that names the path, if S-0844 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0674 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0814 can start.
3. Reply with `message_reply` saying why the hold stands: S-0814 waits until S-0844 moves to review, is cancelled, or is sent back.

### 2026-11-28T00:23:30Z agent-S-0844 S-0844 (settled)
The hold stands for a short while. The trial merge conflicted in flai/internal/mergequeue/eligible.go because main had changed that file. S-0844's sync has just taken in main's change, so the next trial merge should be clean. S-0844 still changes the paths named here: one test in boardview_test.go, one refusal test in hostapi/writes_test.go, and a sentence each in flaiover-dashboard.md and flaiover.md on the save refusal. It expects to reach review within the hour. The in-progress limit is full meanwhile.

### 2026-11-28T00:43:34Z orchestrator
Closed: S-0844 was accepted
