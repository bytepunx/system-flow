---
id: I-0119
title: The close-out's last check that the branch contains main fails when flai commits wip on main during its run
class: efficiency
status: open
count: 6
cost: 9m
first_reported: 2026-10-08T04:29:28Z
last_reported: 2026-10-08T09:55:20Z
updated: 2026-10-08T09:55:20Z
---

# I-0119 The close-out's last check that the branch contains main fails when flai commits wip on main during its run

## Description
The close-out's last check that the branch contains main fails when flai commits wip on main during its run

## Instances

### 2026-10-08T04:29:28Z
Story: S-0324.
S-0324's close-out passed every verify step twice in a row (about 4m30s each, mostly integration and smoke). Both times it stopped at its final sync check, because flai had committed `wip/` on main during the run: replans after a reorder or a cancel (99293995, fe393c71) and an edit of S-0323's criteria (1b80ca80). Those commits touch only `wip/`, never what the branch changes, yet each needs another sync and another full run, which the next replan can race again.

### 2026-10-08T04:35:55Z
Story: S-0316.
S-0316's close-out passed every step three times (00:11Z, 00:21Z, 04:24Z), but each run stopped at the last sync check: main moved during the ten-minute run, once with an acceptance, once with a forecast replan, and once with S-0324's acceptance, publication, and replan. Each meant another sync and another full run.

### 2026-10-08T04:55:54Z
Story: S-0320.
S-0320's first close-out passed every verify step, then stopped at its last check that the branch contains main: S-0316 was accepted and flai 1.39.5 published on main during the run. The sync that followed conflicted on I-0118, which S-0316 and this close-out had each bumped. Resolving that by hand and running the close-out again cost about ten minutes.

### 2026-10-08T08:55:09Z
Story: S-0321.
S-0321's close-out passed every step of flai verify twice, once in 13 minutes and once in 11, and both times stopped at its last check that the branch contains main. During the second run main gained 17 wip-only commits, the planner's "edit draft" and "replan forecasts" among them.

### 2026-10-08T09:16:38Z
Story: S-0322.
S-0322's sixth close-out passed every step of flai verify, then stopped at the last check that the branch contains main: flai had committed `chore: replan forecasts after accepted S-0288` (wip only) on main during the run. Its first and fifth close-outs stopped at the same check after stories were accepted meanwhile.

### 2026-10-08T09:55:20Z
Story: S-0297.
S-0297's close-out passed every verify step, then stopped at the last sync check: main moved while its integration tier ran.

## Remediation

Story S-0347 remediates this issue, created from it at 2026-10-08T08:08:22Z.
