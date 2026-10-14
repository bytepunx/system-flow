---
id: MS-0114
title: "S-0372 holds S-0374 on overlap alone: may S-0374 start beside it?"
from: S-0374
to: S-0372
about: [design/system/workflow.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md]
status: open
participants: [flai]
created: 2026-10-14T20:57:06Z
updated: 2026-10-14T20:57:06Z
---

# MS-0114 S-0372 holds S-0374 on overlap alone: may S-0374 start beside it?

Between S-0374 and S-0372, about `design/system/workflow.md`, `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`.

## Entries

### 2026-10-14T20:57:06Z flai S-0374
S-0372 holds S-0374 on overlap alone: may S-0374 start beside it?

S-0374, "An acceptance on main during a close-out's tiers fails its final sync check, so the whole close-out runs again", is ready, and its claim overlaps S-0372's on `design/system/workflow.md`, `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0136](../../../design/issues/I-0136-an-acceptance-on-main-during-a-close-out-s-tiers-fails-its-final-sync-check-so-the-whole-close-out-runs-again.md), "An acceptance on main during a close-out's tiers fails its final sync check, so the whole close-out runs again". The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0372 --remove <path>`, or on the task that names the path, if S-0372 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0114 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0374 can start.
3. Reply with `message_reply` saying why the hold stands: S-0374 waits until S-0372 moves to review, is cancelled, or is sent back.
