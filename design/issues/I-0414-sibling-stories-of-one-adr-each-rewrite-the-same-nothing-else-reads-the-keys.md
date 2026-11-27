---
id: I-0414
title: Sibling stories of one ADR each rewrite the same 'nothing else reads the keys yet' sentence in a hub design doc and runbook, so each acceptance stops the next on a rebase conflict
class: efficiency
status: open
count: 1
cost: 40m
first_reported: 2026-11-27T22:46:44Z
last_reported: 2026-11-27T22:46:44Z
updated: 2026-11-27T22:46:44Z
---

# I-0414 Sibling stories of one ADR each rewrite the same 'nothing else reads the keys yet' sentence in a hub design doc and runbook, so each acceptance stops the next on a rebase conflict

## Description
Sibling stories of one ADR each rewrite the same 'nothing else reads the keys yet' sentence in a hub design doc and runbook, so each acceptance stops the next on a rebase conflict

## Instances

### 2026-11-27T22:46:44Z
Report: design/analysis/2026-11-27-all.md.
ADR-0312's stories S-0807, S-0809, S-0810, S-0811, and S-0812 each rewrote whole the one sentence in workflow.md § Suites a step needs and runs and its twin in the board runbook that lists which steps read the suite keys yet and which stories are still to come; a section merge cannot settle two sides that change the same line, so each sibling's acceptance stopped the next's rebase. S-0811 took three send-backs, each a sync, a close-out (W47 p50 602 s), and an acceptance again. The last report's I-0387 was the same shape on server-performance.md.

## Impact
S-0811 was sent back from review three times on design/system/workflow.md and docs/operators/runbooks/change-the-board.md: TH-1317 after S-0807's acceptance, TH-1319 (2026-11-21T18:13Z) after S-0809's and S-0810's, TH-1322 (18:36Z) after S-0812's, each "stopped rebasing onto main on conflicts". flai stats coordination.weeks W47: conflicts.met_sync 15 over 36 paths. Serve journal since 2026-11-16: 9 close-outs stopped at rebase, more than at any tier. S-0811's worktree records this as I-0408 (count 2, cost 38m), not yet on main; this issue carries its title so the two fold at acceptance (ADR-0126).

- time_lost_per_cycle: 2h

## Remediation

Found by the analysis in [2026-11-27-all.md](../analysis/2026-11-27-all.md).
