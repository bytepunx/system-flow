---
id: I-0292
title: No analysis.schedule or analysis.formula is set, so the analyzer ran once in thirty days
class: impression
status: open
count: 1
first_reported: 2026-11-07T08:17:36Z
last_reported: 2026-11-07T08:17:36Z
updated: 2026-11-10T08:49:03Z
usage:
  strategic:
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3102
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 14
          output: 2840
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
---

# I-0292 No analysis.schedule or analysis.formula is set, so the analyzer ran once in thirty days

## Description
No analysis.schedule or analysis.formula is set, so the analyzer ran once in thirty days

## Instances

### 2026-11-07T08:17:36Z
Report: design/analysis/2026-11-07-all.md.
Analysis 2026-11-07: this run was started by hand through S-0564's analysis-cadence instance; strategic-agents.md § The analysis schedule says unset means none.

## Impact
flai stats strategic[] analyzer: activities 1, last_run 2026-10-15T22:49:08Z; system-flow.yaml analysis block names an agent only; I-0150 to I-0153 closed and the waits and holds they named returned in W44

## Remediation

Found by the analysis in [2026-11-07-all.md](../analysis/2026-11-07-all.md).

Story S-0692 remediates this issue, created from it at 2026-11-10T08:49:03Z.
