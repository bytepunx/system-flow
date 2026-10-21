---
id: MS-0164
title: "S-0465 holds S-0459 on overlap alone: may S-0459 start beside it?"
from: S-0459
to: S-0465
about: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
status: open
participants: [flai]
created: 2026-10-21T08:32:26Z
updated: 2026-10-21T08:32:26Z
---

# MS-0164 S-0465 holds S-0459 on overlap alone: may S-0459 start beside it?

Between S-0459 and S-0465, about `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`.

## Entries

### 2026-10-21T08:32:26Z flai S-0459
S-0465 holds S-0459 on overlap alone: may S-0459 start beside it?

S-0459, "The settings page becomes five Configuration pages, Features, Agents, File System, Projects, and Harnesses, each reached from the rail", is ready, and its claim overlaps S-0465's on `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> The operator finds a setting under the Configuration page it belongs to, instead of scrolling one page of ten sections. The rail's Configuration group links five pages:
>
> | Page | Sections |
> |------|----------|
> | Features | Host actions, Checks |
> | Agents | Default agent for new stories, Providers, Planning again |
> | File System | Shared paths, Import folders |
> | Projects | Projects, Tokens |
> | Harnesses | Starting agents, one block per harness |
>
> The split is the operator's answer to TH-0553 question 4. Each page shows only its own sections. What each section reads, saves, and gates is unchanged.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0465 --remove <path>`, or on the task that names the path, if S-0465 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0164 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0459 can start.
3. Reply with `message_reply` saying why the hold stands: S-0459 waits until S-0465 moves to review, is cancelled, or is sent back.
