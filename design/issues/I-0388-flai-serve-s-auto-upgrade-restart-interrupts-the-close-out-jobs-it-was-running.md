---
id: I-0388
title: flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch
class: efficiency
status: open
count: 1
cost: 8m
first_reported: 2026-11-16T06:16:21Z
last_reported: 2026-11-16T06:16:21Z
updated: 2026-11-16T06:16:21Z
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

# I-0388 flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch

## Description
flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch

## Instances

### 2026-11-16T06:16:21Z
Report: design/analysis/2026-11-16-all.md.
Analysis 2026-11-16: 12 close-outs in W46 ended "interrupted, flai serve stopped while it ran" at an auto-upgrade restart and were not resumed; each was asked again from the start.

## Impact
closeouts.weeks[] W46 failed 11; the serve journal names each as "interrupted, flai serve stopped while it ran" (S-0603, S-0606, S-0689, S-0706, S-0715, S-0730, S-0742, S-0743, S-0748, S-0750, S-0759), 2 more in W45 and 2 in W47. The journal holds 16 "it installed it and restarts" auto-upgrade entries in W46 (59 over the window) and 48 flai/v* tags in W46 under the threshold release policy. closeouts p50 474 s in W46. ADR-0158 waits for writes before the restart and says nothing of the close-out jobs ADR-0229 added.

- time_lost_per_cycle: 2h

## Remediation

Found by the analysis in [2026-11-16-all.md](../analysis/2026-11-16-all.md).
