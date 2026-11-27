---
id: I-0411
title: The acceptance's lint fix refuses every acceptance while the curator's uncommitted digest entries in the main checkout fail markdownlint, so stories that passed review wait for the operator
class: efficiency
status: open
count: 1
cost: 76h
first_reported: 2026-11-27T22:46:18Z
last_reported: 2026-11-27T22:46:18Z
updated: 2026-11-27T22:46:18Z
---

# I-0411 The acceptance's lint fix refuses every acceptance while the curator's uncommitted digest entries in the main checkout fail markdownlint, so stories that passed review wait for the operator

## Description
The acceptance's lint fix refuses every acceptance while the curator's uncommitted digest entries in the main checkout fail markdownlint, so stories that passed review wait for the operator

## Instances

### 2026-11-27T22:46:18Z
Report: design/analysis/2026-11-27-all.md.
On 2026-11-18 five stories that had passed their close-outs (S-0812, S-0809, S-0839, S-0843, S-0810) were refused within three hours on design/digest/D-0671-…md:20 MD028, which the curator wrote in the main checkout and nothing committed; the orchestrator opened TH-1309, TH-1310, TH-1311, and TH-1314 asking the operator to fix the line, and the stories waited 76 to 82 hours until the operator's session on 2026-11-21. 20 such refusals since 2026-11-15 over 16 stories. The hazard stands today: 26 digest files uncommitted in the main checkout, one failing the lint (D-0653:20).

## Impact
Serve journal: 20 orchestrator refusals 2026-11-15T00:52Z to 2026-11-18T13:22Z reading "cannot be accepted: the main checkout's uncommitted files hold lint the fix of the markdown tier could not fix", naming 16 stories and the curator's digest entries D-0671:20 (MD028), D-0523, D-0635, D-0653, D-0662, D-0669. flai stats waiting.longest: S-0809 273,530 s, S-0810 280,637 s, S-0812 291,495 s, S-0811 291,953 s, S-0834 296,845 s in review from 2026-11-18 to 2026-11-21. W47 review.total_seconds 1,808,331 over 50 stories. Today 26 digest files are uncommitted on main and markdownlint-cli2 0.20.0 over them fails on D-0653:20 MD028. workflow.md says uncommitted main files no longer wait for the operator (ADR-0189, ADR-0248); ADR-0140's preflight check refuses on them; ADR-0318 stops only new entries.

- time_lost_per_cycle: 76h

## Remediation

Found by the analysis in [2026-11-27-all.md](../analysis/2026-11-27-all.md).
