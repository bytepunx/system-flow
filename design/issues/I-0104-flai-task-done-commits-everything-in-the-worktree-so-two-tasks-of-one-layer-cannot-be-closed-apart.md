---
id: I-0104
title: flai task done commits everything in the worktree, so two tasks of one layer cannot be closed apart
class: efficiency
status: open
count: 3
cost: 4m
first_reported: 2026-10-07T07:17:53Z
last_reported: 2026-10-07T14:32:41Z
updated: 2026-10-07T14:32:41Z
---

# I-0104 flai task done commits everything in the worktree, so two tasks of one layer cannot be closed apart

## Description
flai task done commits everything in the worktree, so two tasks of one layer cannot be closed apart

## Instances

### 2026-10-07T07:17:53Z
Story: S-0212.
S-0212 ran T-0958 (chart page) and T-0962 (docs) together in one worktree. To close T-0962 apart, its docs were committed by hand, then `flai task done T-0962 -m ...` was run for the sync and move; its `git add -A` committed T-0958's uncommitted page and table under T-0962's placeholder message, and widened T-0962's touches with T-0958's paths. Fixed by amending the message and resetting T-0962's touches with `flai touches`. `task done` could commit only the task's touches (or take paths), and not require `-m` when there is nothing to commit.

### 2026-10-07T09:05:53Z
Story: S-0215.
S-0215: closing T-0927 committed T-0923's flai/internal/metrics files under the chart's message; I amended the message and committed T-0932 and T-0934 by hand before closing each with nothing to commit.

### 2026-10-07T14:32:41Z
Story: S-0298.
S-0298, layer 1: T-1038 (design docs) and T-1039 (selfupgrade) ran together in one worktree. I committed T-1039's files by hand and called task_done T-1039, which committed T-1038's uncommitted design files (c211addb) under T-1039's subject. Closing two tasks of one layer apart needs each closed before the next's files are written, or a hand commit of every task's files first.

## Remediation

Story S-0312 remediates this issue, created from it at 2026-10-07T14:26:01Z.
