---
id: I-0388
title: flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch
class: efficiency
status: open
count: 4
cost: 41m
first_reported: 2026-11-16T06:16:21Z
last_reported: 2026-11-27T22:47:15Z
updated: 2026-11-27T22:47:15Z
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

### 2026-11-17T09:43:10Z
Story: S-0831.
S-0831's close-out job S-0831-20261117T091529Z, asked for at 09:15:29Z, ended interrupted at 09:33:37Z in its integration tier, while flai 1.104.1 and then 1.105.0 were published on main; the job asked for again at 09:33:42Z reused gofmt, vet, golangci-lint, and markdown, ran integration again from scratch (9m13s), and passed at 09:42:56Z.

### 2026-11-17T09:43:23Z
Story: S-0806.
S-0806's close-out job S-0806-20261117T091319Z was terminated at 09:33:30Z in round 2, 9m36s into the integration tier, after flai 1.105.0 was published. The log reads "Terminated", the integration tier "stopped: context canceled", and the verdict "interrupted". Asked again, it passed in 9m8s from scratch. It cost about 20 minutes of wall clock.

### 2026-11-27T22:47:15Z
Report: design/analysis/2026-11-27-all.md.
Analysis of 2026-11-27: 13 close-outs interrupted by an auto-upgrade restart since 2026-11-16, 21 installs, each acceptance still publishing a release under threshold count 1; W47 close-out failures 8 of 76 at p50 602 s. S-0846 is in progress on this issue today.

## Impact
closeouts.weeks[] W46 failed 11; the serve journal names each as "interrupted, flai serve stopped while it ran" (S-0603, S-0606, S-0689, S-0706, S-0715, S-0730, S-0742, S-0743, S-0748, S-0750, S-0759), 2 more in W45 and 2 in W47. The journal holds 16 "it installed it and restarts" auto-upgrade entries in W46 (59 over the window) and 48 flai/v* tags in W46 under the threshold release policy. closeouts p50 474 s in W46. ADR-0158 waits for writes before the restart and says nothing of the close-out jobs ADR-0229 added.

Serve journal since 2026-11-16: 21 "asked flai host to upgrade flai ... it installed it and restarts", 46 auto-upgrade failures for a tag not yet published, 13 close-outs "interrupted, flai serve stopped while it ran" (S-0739 twice, S-0747, S-0806, S-0831, S-0833 among them). flai stats closeouts.weeks W47: failed 8 of 76, p50 601.6 s, p85 1,168.5 s, the highest of the window. git tag: 12 flai/v* tags on 2026-11-17, 4 on 2026-11-18, 4 on 2026-11-21; orchestration.release is threshold value 1000 count 1. The host's flai is 1.111.1 behind the 1.111.2 pushed at 2026-11-27T22:34Z.

- time_lost_per_cycle: 2h

## Remediation

Found by the analysis in [2026-11-16-all.md](../analysis/2026-11-16-all.md).

Story S-0846 remediates this issue, created from it at 2026-11-17T08:32:57Z.

Found by the analysis in [2026-11-27-all.md](../analysis/2026-11-27-all.md).
