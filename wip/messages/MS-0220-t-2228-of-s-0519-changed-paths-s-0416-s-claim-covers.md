---
id: MS-0220
title: T-2228 of S-0519 changed paths S-0416's claim covers
from: S-0519
to: S-0416
about: [design/adrs/0215-permission-prompt-holds-a-write-at-most-four-minutes-over-stdio-and-at-most-45.md]
status: closed
participants: [agent-S-0519, alex]
created: 2026-10-29T02:00:20Z
updated: 2026-10-29T02:01:33Z
---

# MS-0220 T-2228 of S-0519 changed paths S-0416's claim covers

Between S-0519 and S-0416, about `design/adrs/0215-permission-prompt-holds-a-write-at-most-four-minutes-over-stdio-and-at-most-45.md`.

## Entries

### 2026-10-29T02:00:20Z agent-S-0519 S-0519
T-2228 of S-0519 changed paths S-0416's claim covers.

T-2228, An ADR refining ADR-0124 gives permission_prompt a bound per transport: four minutes over stdio and, over HTTP, the bound the measurement settles below Claude Code's tool-call timeout, committed c155fee on story/S-0519, `docs: [S-0519] ADR-0215 topics mcp and harness`, changing `design/adrs/0215-permission-prompt-holds-a-write-at-most-four-minutes-over-stdio-and-at-most-45.md`. It reaches the main branch when S-0519 is accepted; `git show c155fee` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-29T02:01:33Z alex
Closed: S-0416 was accepted
