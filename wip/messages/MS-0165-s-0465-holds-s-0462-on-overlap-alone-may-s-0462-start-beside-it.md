---
id: MS-0165
title: "S-0465 holds S-0462 on overlap alone: may S-0462 start beside it?"
from: S-0462
to: S-0465
about: [flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
status: open
participants: [flai]
created: 2026-10-21T08:32:26Z
updated: 2026-10-21T08:32:26Z
---

# MS-0165 S-0465 holds S-0462 on overlap alone: may S-0462 start beside it?

Between S-0462 and S-0465, about `flaiover/src/routes/board/+page.svelte`, `flaiover/src/routes/board/board.svelte.test.ts`, `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`.

## Entries

### 2026-10-21T08:32:26Z flai S-0462
S-0465 holds S-0462 on overlap alone: may S-0462 start beside it?

S-0462, "The board's lanes sit side by side and scroll sideways on a wide screen, and collapse to their headings, toggled by a tap, on a narrow one", is ready, and its claim overlaps S-0465's on `flaiover/src/routes/board/+page.svelte`, `flaiover/src/routes/board/board.svelte.test.ts`, `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> The board reads as a board at every width. Today its six lanes stack in one column below 640px and wrap into two columns up to 1024px (`grid-cols-1 sm:grid-cols-2 lg:grid-cols-6`), all inside the layout's `max-w-6xl`. On a wide screen the lanes never stack: they sit in one row as wide as the main area, and the operator scrolls sideways between them. On a phone each lane collapses to its heading, which stays visible, and tapping a heading opens or closes its lane.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0465 --remove <path>`, or on the task that names the path, if S-0465 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0165 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0462 can start.
3. Reply with `message_reply` saying why the hold stands: S-0462 waits until S-0465 moves to review, is cancelled, or is sent back.
