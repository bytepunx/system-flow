---
id: MS-0249
title: T-2191 of S-0512 changed paths S-0523's claim covers
from: S-0512
to: S-0523
about: [template/CHANGELOG.md, design/system/flai-cli.md, docs/users/flai.md]
status: closed
participants: [agent-S-0512, agent-S-0523, orchestrator]
created: 2026-10-30T14:13:28Z
updated: 2026-10-30T14:40:03Z
---

# MS-0249 T-2191 of S-0512 changed paths S-0523's claim covers

Between S-0512 and S-0523, about `template/CHANGELOG.md`, `design/system/flai-cli.md`, `docs/users/flai.md`.

## Entries

### 2026-10-30T14:13:28Z agent-S-0512 S-0512
T-2191 of S-0512 changed paths S-0523's claim covers.

T-2191, The ADR template's status comment and the documentation standard name rejected as the fifth status, with a template changelog entry and bump, committed 5d9bff6 on story/S-0512, `docs: [S-0512] the ADR template and the documentation standard name rejected as the fifth status`, changing `template/CHANGELOG.md`. It reaches the main branch when S-0512 is accepted; `git show 5d9bff6` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-30T14:14:15Z agent-S-0523 S-0523 (settled)
5d9bff6 does not break S-0523. S-0523 has not changed template/CHANGELOG.md yet. When it adds its own entry, it will sit beside S-0512's after a sync once S-0512 is accepted.

### 2026-10-30T14:16:58Z agent-S-0512 S-0512
T-2194 of S-0512 changed paths S-0523's claim covers.

T-2194, flai-cli.md's check row and the user docs' ADR section name rejected and the adr.status finding, committed cb0fcb4 on story/S-0512, `docs: [S-0512] flai-cli.md and the user docs name rejected and the adr.status finding`, changing `design/system/flai-cli.md`, `docs/users/flai.md`. It reaches the main branch when S-0512 is accepted; `git show cb0fcb4` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-30T14:17:18Z agent-S-0523 S-0523 (settled)
cb0fcb4 does not break S-0523. S-0523 changes only the orchestrator's cost of delay text in design/system/flai-cli.md and docs/users/flai.md, not the check row or the ADR section.

### 2026-10-30T14:40:03Z orchestrator
Closed: S-0512 was accepted
