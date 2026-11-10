---
id: I-0285
title: after chains on planned stories hold ready stories for days while in-progress stays under two of five
class: efficiency
status: open
count: 1
first_reported: 2026-11-07T08:17:35Z
last_reported: 2026-11-07T08:17:35Z
updated: 2026-11-10T08:34:51Z
usage:
  strategic:
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 33
          output: 3103
          cache_read: 238317
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2841
          cache_read: 759050
          cache_write: 13855
          cost: 0.278
---

# I-0285 after chains on planned stories hold ready stories for days while in-progress stays under two of five

## Description
after chains on planned stories hold ready stories for days while in-progress stays under two of five

## Instances

### 2026-11-07T08:17:35Z
Report: design/analysis/2026-11-07-all.md.
Analysis 2026-11-07: E-0024's S-0420 to S-0425 each named the stories before it in after and ran one at a time; 46 of 93 stories completed since W44 carry held_seconds, 1,112,249 s together.

## Impact
claims.weeks[].held_seconds.after: W44 704,200 s, W45 203,155 s; S-0425 held 326,792 s (after S-0423, S-0424), S-0424 133,349 s; concurrency.mean_busy 1.81 of limit 5, share with 0 in progress 0.357

- time_lost_per_cycle: 56h

## Remediation

Found by the analysis in [2026-11-07-all.md](../analysis/2026-11-07-all.md).

Story S-0673 remediates this issue, created from it at 2026-11-10T08:34:51Z.
