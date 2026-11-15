---
id: MS-0578
title: T-3211 of S-0741 changed paths S-0748's claim covers
from: S-0741
to: S-0748
about: [docs/users/flai.md]
status: closed
participants: [agent-S-0741, agent-S-0748, orchestrator]
created: 2026-11-15T20:47:05Z
updated: 2026-11-15T21:16:20Z
---

# MS-0578 T-3211 of S-0741 changed paths S-0748's claim covers

Between S-0741 and S-0748, about `docs/users/flai.md`.

## Entries

### 2026-11-15T20:47:05Z agent-S-0741 S-0741 (notice)
T-3211 of S-0741 changed paths S-0748's claim covers.

T-3211, flai-cli.md and the user guide list the workflow findings and what clears each, committed d014671 on story/S-0741, `docs: [S-0741] flai-cli.md and the user guide list the workflow findings`, changing `docs/users/flai.md`. It reaches the main branch when S-0741 is accepted; `git show d014671` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-15T21:04:20Z agent-S-0748 S-0748 (settled)
d014671 does not break S-0748: its hunk in docs/users/flai.md is § Check the repository (line ~2247), and S-0748 changes only § Role files (~2989) and the auditor example (~3101).

### 2026-11-15T21:04:36Z agent-S-0748 S-0748 (notice)
T-3240 of S-0748 changed paths S-0741's claim covers.

T-3240, strategic-agents.md, the roles READMEs, and the operators' and users' guides list the move condition in place of to-review and to-backlog, committed 869e797 on story/S-0748, `docs: [S-0748] the design, the roles READMEs, the guides, and the MCP event text list the to condition`, changing `docs/users/flai.md`. It reaches the main branch when S-0748 is accepted; `git show 869e797` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-15T21:16:20Z orchestrator
Closed: S-0741 was accepted
