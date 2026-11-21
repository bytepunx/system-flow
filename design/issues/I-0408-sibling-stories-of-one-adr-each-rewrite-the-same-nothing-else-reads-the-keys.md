---
id: I-0408
title: Sibling stories of one ADR each rewrite the same 'nothing else reads the keys yet' sentence in a hub design doc and runbook, so each acceptance stops the next on a rebase conflict
class: efficiency
status: open
count: 2
cost: 38m
first_reported: 2026-11-21T18:15:40Z
last_reported: 2026-11-21T18:37:51Z
updated: 2026-11-21T18:37:51Z
---

# I-0408 Sibling stories of one ADR each rewrite the same 'nothing else reads the keys yet' sentence in a hub design doc and runbook, so each acceptance stops the next on a rebase conflict

## Description
Sibling stories of one ADR each rewrite the same 'nothing else reads the keys yet' sentence in a hub design doc and runbook, so each acceptance stops the next on a rebase conflict

## Instances

### 2026-11-21T18:15:40Z
Story: S-0811.
ADR-0312's stories S-0807, S-0809, S-0810, and S-0811 each rewrote the one sentence in design/system/workflow.md § Suites a step needs and runs, and its twin in docs/operators/runbooks/change-the-board.md, that lists which steps read verify, runs, and verify_at yet and which stories are still to come. S-0811 stopped on that sentence twice: on the rebase after S-0807's acceptance (2026-11-20, TH-1317) and again after S-0809's and S-0810's (2026-11-21, TH-1319, which sent it back from review). flai could not settle it by sections because both sides change the same lines. Each time an agent had to be started, resolve, sync, and close out again. The same sentence also went stale on main: S-0809 left the runbook saying the batch runs as it did. A list of who has landed, kept in prose that every sibling edits, is a hub line. Writing each step's behaviour in its own paragraph or section, without a running 'still to come' list, would let the siblings merge apart.

### 2026-11-21T18:37:51Z
Story: S-0811.
A third stop, minutes after S-0811 went back to review: S-0812's acceptance put its own clause into the same sentence in design/system/workflow.md and the board runbook, so S-0811's acceptance stopped its rebase there again and sent it back (TH-1322). Both stories had agreed the split in MS-0626, and each did its part. The sentence still conflicts every time, because each story rewrites the whole line.

## Remediation
