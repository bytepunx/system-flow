---
id: MS-0660
title: "S-0810 holds S-0845 on overlap alone: may S-0845 start beside it?"
from: S-0845
to: S-0810
about: [design/system/workflow.md]
status: open
participants: [flai]
created: 2026-11-21T17:20:30Z
updated: 2026-11-21T17:20:30Z
---

# MS-0660 S-0810 holds S-0845 on overlap alone: may S-0845 start beside it?

Between S-0845 and S-0810, about `design/system/workflow.md`.

## Entries

### 2026-11-21T17:20:30Z flai S-0845
S-0810 holds S-0845 on overlap alone: may S-0845 start beside it?

S-0845, "A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves", is ready, and its claim overlaps S-0810's on `design/system/workflow.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0386](../../../design/issues/I-0386-a-story-in-review-whose-close-out-fails-after-main-moved-waits-for-the-operator.md), "A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves", which the analysis in [2026-11-16-all.md](../../../design/analysis/2026-11-16-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0810 --remove <path>`, or on the task that names the path, if S-0810 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0660 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0845 can start.
3. Reply with `message_reply` saying why the hold stands: S-0845 waits until S-0810 moves to review, is cancelled, or is sent back.
