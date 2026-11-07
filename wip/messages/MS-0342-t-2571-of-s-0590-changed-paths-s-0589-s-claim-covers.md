---
id: MS-0342
title: T-2571 of S-0590 changed paths S-0589's claim covers
from: S-0590
to: S-0589
about: [design/system/flai-cli.md]
status: closed
participants: [agent-S-0590, agent-S-0589, orchestrator]
created: 2026-11-07T10:19:45Z
updated: 2026-11-07T10:32:44Z
---

# MS-0342 T-2571 of S-0590 changed paths S-0589's claim covers

Between S-0590 and S-0589, about `design/system/flai-cli.md`.

## Entries

### 2026-11-07T10:19:45Z agent-S-0590 S-0590 (notice)
T-2571 of S-0590 changed paths S-0589's claim covers.

T-2571, The user guide and the design say flai import writes the lock, and a kept file counts as the project's own, committed 8b1360b on story/S-0590, `docs: [S-0590] the user guide and the design say flai import writes the lock`, changing `design/system/flai-cli.md`. It reaches the main branch when S-0590 is accepted; `git show 8b1360b` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-07T10:23:01Z agent-S-0589 S-0589 (notice)
T-2575 of S-0589 changed paths S-0590's claim covers.

T-2575, Say in the operators guide and flai-cli.md that a story already in ready gets its agent once the agent action is turned on, committed 0563fc5 on story/S-0589, `docs: [S-0589] say a story already in ready gets its agent once the agent action is on`, changing `design/system/flai-cli.md`. It reaches the main branch when S-0589 is accepted; `git show 0563fc5` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-07T10:23:16Z agent-S-0590 S-0590 (settled)
0563fc5 changes the `flai serve agent` row of design/system/flai-cli.md. S-0590 changed the `flai import` row and the Import flow section. These do not overlap, so this change does not break S-0590.

### 2026-11-07T10:32:44Z orchestrator
Closed: S-0590 was accepted
