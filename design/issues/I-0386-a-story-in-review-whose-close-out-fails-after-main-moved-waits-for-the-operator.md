---
id: I-0386
title: A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves
class: efficiency
status: open
count: 1
cost: 30h
first_reported: 2026-11-16T06:16:10Z
last_reported: 2026-11-16T06:16:10Z
updated: 2026-11-17T08:32:56Z
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

# I-0386 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves

## Description
A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves

## Instances

### 2026-11-16T06:16:10Z
Report: design/analysis/2026-11-16-all.md.
Analysis 2026-11-16: three E-0034 stories in review for 87 h together whose re-close-outs stop after main moved (vet, sync, undone_main); nothing but the operator can send them back, while five backlog stories are held on after behind them and the board is starved.

## Impact
aging[] 2026-11-16: S-0739 in review 108,804 s, S-0742 103,608 s, S-0747 100,939 s, each over_p85. flai board --json queue.ineligible: S-0742 verify_failed at vet (84ee2851d320), S-0747 verify_failed at sync (a4d66ac3507f), S-0739 undone_main on flai/internal/workitem/board.go; threads TH-1233, TH-1245, TH-1246 are the orchestrator's recommendations that the operator send each back. ADR-0302 makes flai accept --resolve the operator's alone, and ADR-0248/ADR-0176 send back only on a rebase conflict or a bisected batch. flai check warns wip.after-hub on S-0753; TH-1247 lists S-0752 to S-0758 held by after behind them; 48 ready-empty events since 2026-11-09; concurrency.share[0] 0.404.

- time_lost_per_cycle: 30h

## Remediation

Found by the analysis in [2026-11-16-all.md](../analysis/2026-11-16-all.md).

Story S-0845 remediates this issue, created from it at 2026-11-17T08:32:56Z.
