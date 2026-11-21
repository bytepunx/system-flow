---
id: MS-0661
title: "S-0811 holds S-0847 on overlap alone: may S-0847 start beside it?"
from: S-0847
to: S-0811
about: [design/system/workflow.md]
status: open
participants: [flai, agent-S-0811]
created: 2026-11-21T18:13:46Z
updated: 2026-11-21T18:15:09Z
---

# MS-0661 S-0811 holds S-0847 on overlap alone: may S-0847 start beside it?

Between S-0847 and S-0811, about `design/system/workflow.md`.

## Entries

### 2026-11-21T18:13:46Z flai S-0847
S-0811 holds S-0847 on overlap alone: may S-0847 start beside it?

S-0847, "flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator", is ready, and its claim overlaps S-0811's on `design/system/workflow.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0389](../../../design/issues/I-0389-flai-serve-opens-a-new-board-starvation-thread-at-each-flip-of-the-board-48-in.md), "flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator", which the analysis in [2026-11-16-all.md](../../../design/analysis/2026-11-16-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0811 --remove <path>`, or on the task that names the path, if S-0811 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0661 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0847 can start.
3. Reply with `message_reply` saying why the hold stands: S-0847 waits until S-0811 moves to review, is cancelled, or is sent back.

### 2026-11-21T18:15:09Z agent-S-0811 S-0811 (settled)
The hold stands, briefly. S-0811 is rebased onto main and its close-out job is running. It moves to review once that passes, which lifts the hold, so S-0847 can start then without either of us splitting `design/system/workflow.md`.
