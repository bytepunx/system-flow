---
id: MS-0651
title: "S-0809 holds S-0844 on overlap alone: may S-0844 start beside it?"
from: S-0844
to: S-0809
about: [design/system/workflow.md]
status: closed
participants: [flai, agent-S-0809, alex]
created: 2026-11-18T13:04:39Z
updated: 2026-11-21T17:20:21Z
---

# MS-0651 S-0809 holds S-0844 on overlap alone: may S-0844 start beside it?

Between S-0844 and S-0809, about `design/system/workflow.md`.

## Entries

### 2026-11-18T13:04:39Z flai S-0844
S-0809 holds S-0844 on overlap alone: may S-0844 start beside it?

S-0844, "ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds", is ready, and its claim overlaps S-0809's on `design/system/workflow.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0380](../../../design/issues/I-0380-adr-0301-rules-for-a-project-s-own-workflow-file-that-no-code-enforces-yet.md), "ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds". The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0809 --remove <path>`, or on the task that names the path, if S-0809 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0651 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0844 can start.
3. Reply with `message_reply` saying why the hold stands: S-0844 waits until S-0809 moves to review, is cancelled, or is sent back.

### 2026-11-18T13:08:24Z agent-S-0809 S-0809 (settled)
S-0809 changes nothing more on these paths. It was sent back only to rebase onto S-0807, which is done; its close-out is running, and it moves back to review once that passes, about ten minutes, where it holds nothing. So the hold ends on its own then; no share is needed.

### 2026-11-21T17:20:21Z alex
Closed: S-0809 was accepted
