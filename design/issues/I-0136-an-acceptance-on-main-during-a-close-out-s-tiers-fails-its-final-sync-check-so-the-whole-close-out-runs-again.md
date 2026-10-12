---
id: I-0136
title: An acceptance on main during a close-out's tiers fails its final sync check, so the whole close-out runs again
class: efficiency
status: open
count: 1
cost: 6m
first_reported: 2026-10-12T15:19:01Z
last_reported: 2026-10-12T15:19:01Z
updated: 2026-10-12T15:19:01Z
---

# I-0136 An acceptance on main during a close-out's tiers fails its final sync check, so the whole close-out runs again

## Description
An acceptance on main during a close-out's tiers fails its final sync check, so the whole close-out runs again

## Instances

### 2026-10-12T15:19:01Z
Story: S-0232.
S-0232's close-out passed every step twice on 2026-10-12 (about 6 minutes each: flaiover, integration, and smoke). Both times the final `flai verify --sync-only` failed because other stories had been accepted meanwhile, putting the branch 5 and then 7 commits behind main, with paths outside wip/. Each failure means another sync and another full run. I-0119, now closed, covered only flai's wip-only commits on main. Possible remedies: let the final check accept a main that merges cleanly with the branch and touches none of its paths, or rerun only the tiers the new main commits select.

## Remediation
