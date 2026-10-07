---
id: I-0104
title: flai task done commits everything in the worktree, so two tasks of one layer cannot be closed apart
class: efficiency
status: open
count: 1
cost: 5m
first_reported: 2026-10-07T07:17:53Z
last_reported: 2026-10-07T07:17:53Z
updated: 2026-10-07T07:17:53Z
---

# I-0104 flai task done commits everything in the worktree, so two tasks of one layer cannot be closed apart

## Description
flai task done commits everything in the worktree, so two tasks of one layer cannot be closed apart

## Instances

### 2026-10-07T07:17:53Z
Story: S-0212.
S-0212 ran T-0958 (chart page) and T-0962 (docs) together in one worktree. To close T-0962 apart, its docs were committed by hand, then `flai task done T-0962 -m ...` was run for the sync and move; its `git add -A` committed T-0958's uncommitted page and table under T-0962's placeholder message, and widened T-0962's touches with T-0958's paths. Fixed by amending the message and resetting T-0962's touches with `flai touches`. `task done` could commit only the task's touches (or take paths), and not require `-m` when there is nothing to commit.

## Remediation
