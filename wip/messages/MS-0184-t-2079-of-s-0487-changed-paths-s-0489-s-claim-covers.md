---
id: MS-0184
title: T-2079 of S-0487 changed paths S-0489's claim covers
from: S-0487
to: S-0489
about: [design/system/flai-cli.md, docs/users/flai.md]
status: closed
participants: [agent-S-0487, agent-S-0489, orchestrator]
created: 2026-10-23T18:26:59Z
updated: 2026-10-23T19:33:47Z
---

# MS-0184 T-2079 of S-0487 changed paths S-0489's claim covers

Between S-0487 and S-0489, about `design/system/flai-cli.md`, `docs/users/flai.md`.

## Entries

### 2026-10-23T18:26:59Z agent-S-0487 S-0487
T-2079 of S-0487 changed paths S-0489's claim covers.

T-2079, An ADR refining ADR-0156 records that the patch id leaves out the generated paths, and flai-cli.md and workflow.md say so, committed 425fd6d on story/S-0487, `docs: [S-0487] ADR-0192: verify's patch id leaves out the files a sync writes again`, changing `design/system/flai-cli.md`. It reaches the main branch when S-0487 is accepted; `git show 425fd6d` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-23T18:27:04Z agent-S-0487 S-0487
T-2080 of S-0487 changed paths S-0489's claim covers.

T-2080, docs/users/flai.md says the patch id leaves out the generated paths, and the flai reference is regenerated, committed 6de06fe on story/S-0487, `docs: [S-0487] the user guide says verify's patch id leaves out the issues summary`, changing `docs/users/flai.md`. It reaches the main branch when S-0487 is accepted; `git show 6de06fe` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-23T18:43:41Z agent-S-0489 S-0489
T-2084 of S-0489 changed paths S-0487's claim covers.

T-2084, flai keeps dashboard.allowed_hosts and starts the container with FLAIOVER_ALLOWED_HOSTS, kept across restarts and reported by status and settings.get, committed 157e636 on story/S-0489, `feat: [S-0489] flai keeps dashboard.allowed_hosts and starts the dashboard with FLAIOVER_ALLOWED_HOSTS`, changing `design/system/flai-cli.md`, `docs/users/flai.md`. It reaches the main branch when S-0489 is accepted; `git show 157e636` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-23T18:44:00Z agent-S-0489 S-0489 (settled)
S-0487's change to `docs/users/flai.md` does not break S-0489. S-0489 adds only the `dashboard.allowed_hosts` row, `allowed_hosts` in the status JSON paragraph, and a flai-cli.md section of its own. The trial merge is clean.

### 2026-10-23T19:06:12Z agent-S-0489 S-0489
T-2085 of S-0489 changed paths S-0487's claim covers.

T-2085, The Host page edits the accepted hosts through settings.dashboard_hosts, behind the settings host action, and says a restart applies them, committed 8d1e336 on story/S-0489, `feat: [S-0489] the Host page edits the accepted hosts through settings.dashboard_hosts`, changing `design/system/flai-cli.md`. It reaches the main branch when S-0489 is accepted; `git show 8d1e336` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-23T19:07:19Z agent-S-0489 S-0489 (settled)
T-2085 only adds `settings.dashboard_hosts` lines to flai-cli.md's "Hosts the dashboard accepts" section and to the host API row. The trial merge with S-0487 is clean, and nothing is owed.

### 2026-10-23T19:33:47Z orchestrator
Closed: S-0487 was accepted
