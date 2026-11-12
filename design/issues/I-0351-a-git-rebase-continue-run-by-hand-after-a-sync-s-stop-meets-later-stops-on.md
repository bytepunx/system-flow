---
id: I-0351
title: A git rebase --continue run by hand after a sync's stop meets later stops on issue files and summary.md that flai's sync would have merged itself
class: efficiency
status: closed
count: 1
cost: 5m
first_reported: 2026-11-11T18:18:50Z
last_reported: 2026-11-11T18:18:50Z
updated: 2026-11-12T19:26:36Z
---

# I-0351 A git rebase --continue run by hand after a sync's stop meets later stops on issue files and summary.md that flai's sync would have merged itself

## Description
A git rebase --continue run by hand after a sync's stop meets later stops on issue files and summary.md that flai's sync would have merged itself

## Instances

### 2026-11-11T18:18:50Z
Story: S-0689.
S-0689, sent back on TH-1082. flai stream sync S-0689 stopped on design/system/flai-cli.md, which its section merge could not settle, and said to resolve it, git add it, and run git rebase --continue. Done by hand, that continue replayed the next commit, S-0689's close-out, which stopped again on design/issues/I-0347 (both sides added instances) and design/issues/summary.md. flai's sync merges issue instances (ADR-0126) and regenerates summary.md (ADR-0098) only in the rebase it continues itself, so the agent had to settle those by hand, here by taking main's side, as the close-out commit only recorded a note outside the story. A flai command that continues a rebase stopped on a path the agent resolved, with the same settling of issue and generated files, would avoid it.

## Remediation

Story S-0726 remediates this issue, created from it at 2026-11-11T21:55:30Z.
Closed 2026-11-12T19:26:36Z: S-0726 (ADR-0294): flai stream sync, and flai task done called again, now continue a rebase whose stop the agent resolved and staged, settling the later stops on issue files and design/issues/summary.md as the sync's own rebase does and running the steps after a clean rebase; no advice tells an agent to run git rebase --continue any more. TestSyncContinuesARebaseWhoseStopTheAgentResolvedAndSettlesTheStopsAfterIt reproduces I-0351's case.
