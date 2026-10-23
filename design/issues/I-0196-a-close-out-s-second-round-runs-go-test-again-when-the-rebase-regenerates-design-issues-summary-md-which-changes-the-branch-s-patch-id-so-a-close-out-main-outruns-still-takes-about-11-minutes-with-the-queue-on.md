---
id: I-0196
title: A close-out's second round runs go-test again when the rebase regenerates design/issues/summary.md, which changes the branch's patch id, so a close-out main outruns still takes about 11 minutes with the queue on
class: efficiency
status: open
count: 1
cost: 6m
first_reported: 2026-10-23T05:04:30Z
last_reported: 2026-10-23T05:04:30Z
updated: 2026-10-23T17:43:16Z
---

# I-0196 A close-out's second round runs go-test again when the rebase regenerates design/issues/summary.md, which changes the branch's patch id, so a close-out main outruns still takes about 11 minutes with the queue on

## Description
A close-out's second round runs go-test again when the rebase regenerates design/issues/summary.md, which changes the branch's patch id, so a close-out main outruns still takes about 11 minutes with the queue on

## Instances

### 2026-10-23T05:04:30Z
Story: S-0466.
S-0466's close-out, with integration and smoke left to the merge queue, ran go-test in round 1 (5m12s). Main then gained a commit that changed only design/issues (I-0195 and summary.md). Round 2 ran go-test again (5m54s), with no reuse note, so the whole call took about 11 minutes, past the 600000 ms call, and was moved to the background. The rebase regenerated design/issues/summary.md, which the branch also changes, so the branch's patch id changed and ADR-0156's reuse did not apply. Every story that records or closes an issue changes summary.md, so any close-out main outruns meets this.

## Remediation

Story S-0487 remediates this issue, created from it at 2026-10-23T17:43:16Z.
