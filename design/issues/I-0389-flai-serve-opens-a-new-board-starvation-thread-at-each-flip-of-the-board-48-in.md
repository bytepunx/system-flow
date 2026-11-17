---
id: I-0389
title: flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator
class: efficiency
status: open
count: 1
cost: 1m
first_reported: 2026-11-16T06:16:26Z
last_reported: 2026-11-16T06:16:26Z
updated: 2026-11-17T08:32:58Z
usage:
  strategic:
    - kind: analyzer
      seconds: 99
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 47
          output: 6584
          cache_read: 349360
          cache_write: 34306
          cost: 1.1031
        - model: claude-haiku-5-5
          input: 36
          output: 14843
          cache_read: 2398931
          cache_write: 47318
          cost: 0.1557
---

# I-0389 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator

## Description
flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator

## Instances

### 2026-11-16T06:16:26Z
Report: design/analysis/2026-11-16-all.md.
Analysis 2026-11-16: the board flips between starved and not while drafts are finalized and promoted one at a time, and each flip opens a thread the operator must read and a ready-empty event.

## Impact
48 threads titled "The board is starved: nothing in ready can start and no backlog story can be promoted" by flai serve on wip/kanban/board.md from 2026-11-12T00:41:16Z to 2026-11-16T06:01:11Z, 28 on 2026-11-13; TH-1135 lived 51 s, TH-1151 45 s, TH-1166 6 s; TH-1134 gained 12 entries in four minutes, TH-1196 17 in sixteen. ADR-0285 §5: "The next starving opens a new thread". digest.weeks[] W46 threads_to_operator 244 against 178 stories completed; 48 ready-empty events since 2026-11-09.

- time_lost_per_cycle: 30m

## Remediation

Found by the analysis in [2026-11-16-all.md](../analysis/2026-11-16-all.md).

Story S-0847 remediates this issue, created from it at 2026-11-17T08:32:58Z.
