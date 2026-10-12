---
id: I-0136
title: An acceptance on main during a close-out's tiers fails its final sync check, so the whole close-out runs again
class: efficiency
status: open
count: 6
cost: 6m
first_reported: 2026-10-12T15:19:01Z
last_reported: 2026-10-12T21:11:37Z
updated: 2026-10-12T21:11:37Z
---

# I-0136 An acceptance on main during a close-out's tiers fails its final sync check, so the whole close-out runs again

## Description
An acceptance on main during a close-out's tiers fails its final sync check, so the whole close-out runs again

## Instances

### 2026-10-12T15:19:01Z
Story: S-0232.
S-0232's close-out passed every step twice on 2026-10-12 (about 6 minutes each: flaiover, integration, and smoke). Both times the final `flai verify --sync-only` failed because other stories had been accepted meanwhile, putting the branch 5 and then 7 commits behind main, with paths outside wip/. Each failure means another sync and another full run. I-0119, now closed, covered only flai's wip-only commits on main. Possible remedies: let the final check accept a main that merges cleanly with the branch and touches none of its paths, or rerun only the tiers the new main commits select.

### 2026-10-12T16:28:21Z
Story: S-0337.
S-0337's second close-out passed every tier, including about 4 minutes of integration and 1.5 minutes of the flaiover tier. Its final `flai verify --sync-only` then failed, because 15 commits had reached main during the run with changes outside wip/, such as .github/workflows/release-flaiover.yml and issue files. The branch had to be synced and the whole close-out run a third time.

### 2026-10-12T16:43:24Z
Story: S-0337.
On S-0337's fourth close-out, every tier passed again. S-0235 was then accepted during the run, with 9 commits that changed .github/workflows, design/system/flai-cli.md, and other paths. The final sync check failed again. That is twice in one story, about 12 minutes of tiers run again, each time with the branch rebasing onto main cleanly.

### 2026-10-12T18:23:27Z
Story: S-0373.
S-0373's close-out passed every tier (integration 4m43s, smoke 28s), then its --sync-only check failed: S-0368's acceptance had landed on main during the tiers, 7 commits changing paths outside wip/ that S-0373 does not touch. Synced cleanly and ran the whole close-out again.

### 2026-10-12T20:57:39Z
Story: S-0360.
S-0360's close-out passed every tier (integration 5m24s, smoke 1m48s), then its final sync check failed because S-0370 was accepted on main meanwhile; the close-out runs again in full.

### 2026-10-12T21:11:37Z
Story: S-0352.
S-0352's close-out passed every step and tier (integration 5m25s, smoke 1m56s), then its final sync check found main 7 commits ahead, outside wip/, from acceptances made during the tiers, so the whole close-out runs again.

## Remediation

Story S-0374 remediates this issue, created from it at 2026-10-12T16:52:54Z.
