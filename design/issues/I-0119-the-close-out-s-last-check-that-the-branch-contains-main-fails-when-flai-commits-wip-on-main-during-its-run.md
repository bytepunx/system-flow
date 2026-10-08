---
id: I-0119
title: The close-out's last check that the branch contains main fails when flai commits wip on main during its run
class: efficiency
status: open
count: 2
cost: 10m
first_reported: 2026-10-08T04:29:28Z
last_reported: 2026-10-08T04:35:55Z
updated: 2026-10-08T04:35:55Z
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

## Remediation
