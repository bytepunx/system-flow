---
id: MS-0171
title: "S-0477 holds S-0476 on overlap alone: may S-0476 start beside it?"
from: S-0476
to: S-0477
about: [design/system/component-library.md]
status: open
participants: [flai]
created: 2026-10-22T20:58:06Z
updated: 2026-10-22T20:58:06Z
---

# MS-0171 S-0477 holds S-0476 on overlap alone: may S-0476 start beside it?

Between S-0476 and S-0477, about `design/system/component-library.md`.

## Entries

### 2026-10-22T20:58:06Z flai S-0476
S-0477 holds S-0476 on overlap alone: may S-0476 start beside it?

S-0476, "The Components page cannot show a component that asks the app's server in its answered state, since the harness answers every such request with 501", is ready, and its claim overlaps S-0477's on `design/system/component-library.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0164](../../../design/issues/I-0164-the-components-page-cannot-show-a-component-that-asks-the-app-s-server-in-its-answered-state-since-the-harness-answers-every-such-request-with-501.md), "The Components page cannot show a component that asks the app's server in its answered state, since the harness answers every such request with 501". The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0477 --remove <path>`, or on the task that names the path, if S-0477 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0171 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0476 can start.
3. Reply with `message_reply` saying why the hold stands: S-0476 waits until S-0477 moves to review, is cancelled, or is sent back.
