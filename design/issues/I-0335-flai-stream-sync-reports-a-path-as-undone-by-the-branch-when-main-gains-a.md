---
id: I-0335
title: flai stream sync reports a path as undone by the branch when main gains a commit to it while the sync runs
class: impression
status: closed
count: 1
cost: 3m
first_reported: 2026-11-10T09:03:43Z
last_reported: 2026-11-10T09:03:43Z
updated: 2026-11-11T22:21:37Z
---

# I-0335 flai stream sync reports a path as undone by the branch when main gains a commit to it while the sync runs

## Description
flai stream sync reports a path as undone by the branch when main gains a commit to it while the sync runs

## Instances

### 2026-11-10T09:03:43Z
Story: S-0680.
T-2902's task_done sync rebased story/S-0680 onto 8e07f63ed. In the same seconds main gained 4a247f7a9 (ADR-0123 accepted) and a489070aa (ADR-0125 accepted). The sync answered undoes_main for design/adrs/0123-…md, with lines "status: accepted" and "date: 2026-11-10" and old_base 1cf8ed8. The branch never touched that file: git diff main...HEAD did not list it, and both the old and the new base hold "status: proposed". The next sync, in T-2903's task_done, answered undoes_main empty, and .flai-cache/undone keeps no record for S-0680. So the report seems to compare the branch with a main that moved after the rebase's base was taken. It cost a check of three commits and the undone folder.

## Remediation

Story S-0709 remediates this issue, created from it at 2026-11-10T09:55:57Z.
Closed 2026-11-11T22:21:37Z: S-0712, into which the operator merged S-0709 on TH-1017: storygit.Rebase now takes main's head once, as a commit, and rebases onto, checks undoes_main, settles the changelogs, and folds the issues against that commit, so commits main gains while flai stream sync runs are not reported as undone (T-2963). TestSyncFindsNothingUndoneWhenMainGainsCommitsWhileItRuns reproduces it with two commits gained during the sync, and TestSyncSettlesTheChangelogAgainstMainAsTheRebaseTookIt the changelog case; both fail without the fix (T-2965).
