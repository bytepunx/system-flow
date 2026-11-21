---
id: MS-0657
title: "S-0810 holds S-0844 on overlap alone: may S-0844 start beside it?"
from: S-0844
to: S-0810
about: [docs/operators/runbooks/change-the-board.md, design/system/workflow.md]
status: open
participants: [flai]
created: 2026-11-21T17:20:30Z
updated: 2026-11-21T17:20:30Z
---

# MS-0657 S-0810 holds S-0844 on overlap alone: may S-0844 start beside it?

Between S-0844 and S-0810, about `docs/operators/runbooks/change-the-board.md`, `design/system/workflow.md`.

## Entries

### 2026-11-21T17:20:30Z flai S-0844
S-0810 holds S-0844 on overlap alone: may S-0844 start beside it?

S-0844, "ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds", is ready, and its claim overlaps S-0810's on `docs/operators/runbooks/change-the-board.md`, `design/system/workflow.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0380](../../../design/issues/I-0380-adr-0301-rules-for-a-project-s-own-workflow-file-that-no-code-enforces-yet.md), "ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds". The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0810 --remove <path>`, or on the task that names the path, if S-0810 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0657 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0844 can start.
3. Reply with `message_reply` saying why the hold stands: S-0844 waits until S-0810 moves to review, is cancelled, or is sent back.
