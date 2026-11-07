---
id: I-0284
title: orchestration.permissions.set_cost_of_delay is off, so a cost of delay input the operator confirmed waits a day for the operator to write it
class: efficiency
status: open
count: 1
first_reported: 2026-11-07T08:17:35Z
last_reported: 2026-11-07T08:17:35Z
updated: 2026-11-07T08:17:35Z
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

# I-0284 orchestration.permissions.set_cost_of_delay is off, so a cost of delay input the operator confirmed waits a day for the operator to write it

## Description
orchestration.permissions.set_cost_of_delay is off, so a cost of delay input the operator confirmed waits a day for the operator to write it

## Instances

### 2026-11-07T08:17:35Z
Report: design/analysis/2026-11-07-all.md.
Analysis 2026-11-07: W44 threads.confirmed 12 waits, 261,253 s; two of them were inputs already decided that the orchestrator could not write.

## Impact
waiting.longest: S-0549 TH-0784 135,693 s and S-0547 TH-0789 134,328 s, both awaited alex, on an input confirmed on TH-0775; system-flow.yaml orchestration.permissions has no set_cost_of_delay (ADR-0224, built by S-0523)

- time_lost_per_cycle: 75h

## Remediation

Found by the analysis in [2026-11-07-all.md](../analysis/2026-11-07-all.md).
