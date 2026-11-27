---
id: I-0415
title: The orchestrator's accept takes no --keep-removal, so a story that rewrites a line main added waits in review for the operator even after the story's agent and the orchestrator agree the removal is right
class: efficiency
status: open
count: 1
cost: 148h
first_reported: 2026-11-27T22:46:52Z
last_reported: 2026-11-27T22:46:52Z
updated: 2026-11-27T22:46:52Z
---

# I-0415 The orchestrator's accept takes no --keep-removal, so a story that rewrites a line main added waits in review for the operator even after the story's agent and the orchestrator agree the removal is right

## Description
The orchestrator's accept takes no --keep-removal, so a story that rewrites a line main added waits in review for the operator even after the story's agent and the orchestrator agree the removal is right

## Instances

### 2026-11-27T22:46:52Z
Report: design/analysis/2026-11-27-all.md.
S-0811 passed review and every close-out step, but its branch rewrites main's 'nothing else reads the keys yet' sentence in two files, which ADR-0155 records as undoing main's change. The orchestrator checked the two removals, found them deliberate rewrites, and could only recommend the operator's flai accept --keep-removal command on TH-1321; the story has waited 148 hours in review since, holding S-0813 to S-0815 behind it in E-0034.

## Impact
flai stats aging: S-0811 in review 829,391 s (9d14h), over_p85, the one story in review today; its current wait 531,995 s since 2026-11-21T18:49:20Z (waiting.longest). flai board --json queue.ineligible for S-0811: undone_main on design/system/workflow.md and docs/operators/runbooks/change-the-board.md "keep the removal by accepting it with flai accept S-0811 --keep-removal", thread_open TH-1321, verify_failed at sync before the head moved. TH-1321 (orchestrator, 2026-11-21T18:34Z): "My accept tool takes no --keep-removal, and the other way out, flai stream restore, would bring back main's stale sentence beside the story's new one"; agent-S-0811 agreed at 18:38Z. ADR-0155 makes --keep-removal the operator's; ADR-0093 gives the orchestrator's accept no such flag. S-0811's worktree records this as I-0410, not yet on main; this issue carries its title so the two fold at acceptance (ADR-0126).

- time_lost_per_cycle: 148h

## Remediation

Found by the analysis in [2026-11-27-all.md](../analysis/2026-11-27-all.md).
