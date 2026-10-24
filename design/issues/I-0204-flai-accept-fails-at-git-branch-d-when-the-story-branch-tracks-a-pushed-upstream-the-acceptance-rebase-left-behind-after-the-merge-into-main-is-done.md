---
id: I-0204
title: flai accept fails at git branch -d when the story branch tracks a pushed upstream the acceptance rebase left behind, after the merge into main is done
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-10-24T13:15:26Z
last_reported: 2026-10-24T13:15:26Z
updated: 2026-10-24T13:15:26Z
---

# I-0204 flai accept fails at git branch -d when the story branch tracks a pushed upstream the acceptance rebase left behind, after the merge into main is done

## Description
flai accept fails at git branch -d when the story branch tracks a pushed upstream the acceptance rebase left behind, after the merge into main is done

## Instances

### 2026-10-24T13:15:26Z
In gnosiac, flai accept S-0001 (local acceptance, forge action on, so flai serve had pushed story/S-0001 to origin) rebased the branch, fast-forwarded it into main, and removed the worktree, then failed at git branch -d story/S-0001: 'not deleting branch that is not yet merged to refs/remotes/origin/story/S-0001, even though it is merged to HEAD'. The acceptance rebase rewrote the commits, so origin's branch was no longer an ancestor and -d checked the upstream. The story stayed in review with its work already on main; accepting again fails the same way until the operator runs git branch -D by hand. mergeStoryBranch (flai/cmd/branch.go) should delete with -D after a successful --ff-only merge, as removeStoryCheckout in accept_remote.go already does for the same reason, and an acceptance should go on from a merge it already made.

## Impact
Every local acceptance of a story whose branch flai serve pushed for a pull request and whose acceptance rebases it; one occurrence in gnosiac on 2026-10-24, about ten minutes to diagnose and delete by hand

- time_lost_per_cycle: 10m

## Remediation
