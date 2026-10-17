---
id: MS-0133
title: T-1628 of S-0396 changed paths S-0395's claim covers
from: S-0396
to: S-0395
about: [design/adrs/0168-under-pull-request-acceptance-flai-stats-reports-a-checks-wait-and-a-reviewers.md, design/adrs/0169-under-pull-request-acceptance-every-acceptance-the-operator-s-and-the.md, design/adrs/README.md, design/system/flaiover-dashboard.md, design/system/flai-cli.md, docs/users/flai.md]
status: closed
participants: [agent-S-0396, agent-S-0395, orchestrator]
created: 2026-10-17T17:37:40Z
updated: 2026-10-17T18:25:19Z
---

# MS-0133 T-1628 of S-0396 changed paths S-0395's claim covers

Between S-0396 and S-0395, about `design/adrs/0168-under-pull-request-acceptance-flai-stats-reports-a-checks-wait-and-a-reviewers.md`, `design/adrs/0169-under-pull-request-acceptance-every-acceptance-the-operator-s-and-the.md`, `design/adrs/README.md`, `design/system/flaiover-dashboard.md`, `design/system/flai-cli.md`, `docs/users/flai.md`.

## Entries

### 2026-10-17T17:37:40Z agent-S-0396 S-0396
T-1628 of S-0396 changed paths S-0395's claim covers.

T-1628, An ADR defines the checks wait and the reviewers wait, and metrics.md carries their definitions and fields, committed 1f93c16 on story/S-0396, `docs: [S-0396] metrics.md defines the checks wait and the reviewers wait and their fields, under ADR-0168`, changing `design/adrs/0168-under-pull-request-acceptance-flai-stats-reports-a-checks-wait-and-a-reviewers.md`. It reaches the main branch when S-0396 is accepted; `git show 1f93c16` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-17T17:41:55Z agent-S-0395 S-0395
T-1624 of S-0395 changed paths S-0396's claim covers.

T-1624, An ADR records that under pull-request acceptance every acceptance, the orchestrator's included, waits for the required checks and reviews, committed 1086fcc on story/S-0395, `docs: [S-0395] an ADR records that under pull-request acceptance every acceptance, the orchestrator's included, waits for the required checks and reviews`, changing `design/adrs/0169-under-pull-request-acceptance-every-acceptance-the-operator-s-and-the.md`, `design/adrs/README.md`. It reaches the main branch when S-0395 is accepted; `git show 1086fcc` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-17T17:46:58Z agent-S-0396 S-0396
T-1629 of S-0396 changed paths S-0395's claim covers.

T-1629, The metrics package computes the checks and reviewers waits per item, per week, and among the longest waits, from the kept forge events statsread loads, committed 902d739 on story/S-0396, `feat: [S-0396] the metrics package computes the checks and reviewers waits per item, per week, and among the longest waits, from the kept forge events statsread loads`, changing `design/adrs/0168-under-pull-request-acceptance-flai-stats-reports-a-checks-wait-and-a-reviewers.md`. It reaches the main branch when S-0396 is accepted; `git show 902d739` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-17T17:57:33Z agent-S-0395 S-0395
The claims of S-0395 and S-0396 now overlap.

A write grew the claim of S-0395 to cover `design/system/flaiover-dashboard.md`, which the claim of S-0396 covers too. Which of the two stories changes them first? Reply here to agree. If you do not agree, either agent may ask the operator with `flai message escalate` or the MCP tool `message_escalate`.

### 2026-10-17T18:01:26Z agent-S-0396 S-0396
T-1632 of S-0396 changed paths S-0395's claim covers.

T-1632, The dashboard and CLI design and the user docs describe the checks and reviewers waits, committed 9601bbd on story/S-0396, `docs: [S-0396] the dashboard and CLI design and the user docs describe the checks and reviewers waits`, changing `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, `docs/users/flai.md`. It reaches the main branch when S-0396 is accepted; `git show 9601bbd` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-17T18:01:51Z agent-S-0395 S-0395
T-1627 of S-0395 changed paths S-0396's claim covers.

T-1627, The workflow, the strategic agents' design, and the CLI docs describe the forge gate on acceptance, committed 74e2d7c on story/S-0395, `docs: [S-0395] the workflow, the strategic agents' design and convention, and the CLI and dashboard docs describe the forge gate on acceptance`, changing `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, `docs/users/flai.md`. It reaches the main branch when S-0395 is accepted; `git show 74e2d7c` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-17T18:01:54Z agent-S-0395 S-0395
The claims of S-0395 and S-0396 now overlap.

A write grew the claim of S-0395 to cover `design/system/flaiover-dashboard.md`, which the claim of S-0396 covers too. Which of the two stories changes them first? Reply here to agree. If you do not agree, either agent may ask the operator with `flai message escalate` or the MCP tool `message_escalate`.

### 2026-10-17T18:22:20Z agent-S-0395 S-0395 (settled)
S-0395's work is finished and it is going to review now, so it changes these paths no further.

In `design/system/flaiover-dashboard.md` it edited only two rows: `/review/:id` (the confirmation's forge-blocker box) and `GET /api/items/:id/acceptance` (`forge_blockers` and `pull_request`). It also added ADR-0169 and its index row, a `### The acceptance gate` section in `flai-cli.md`, and a section in `docs/users/flai.md` under pull-request acceptance. `story/S-0395` merged cleanly with `story/S-0396` at the last sync. Whichever story is accepted second rebases onto the other.

### 2026-10-17T18:25:19Z orchestrator
Closed: S-0396 was accepted
