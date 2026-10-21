---
id: I-0153
title: Ready stories waited up to 3.5 days unheld with in-progress lanes free and review under its limit
class: impression
status: open
count: 1
cost: 75h
first_reported: 2026-10-15T22:48:48Z
last_reported: 2026-10-15T22:48:48Z
updated: 2026-10-21T11:58:17Z
usage:
  strategic:
    - kind: analyzer
      seconds: 113
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 97
          output: 6305
          cache_read: 370754
          cache_write: 41687
          cost: 1.2427
        - model: claude-haiku-5-5
          input: 8
          output: 6718
          cache_read: 204809
          cache_write: 23258
          cost: 0.0083
---

# I-0153 Ready stories waited up to 3.5 days unheld with in-progress lanes free and review under its limit

## Description
Ready stories waited up to 3.5 days unheld with in-progress lanes free and review under its limit

## Instances

### 2026-10-15T22:48:48Z
Report: design/analysis/2026-10-15-bottlenecks.md.
Five E-0016 chart stories sat in ready about 75 h each beyond their holds while lanes were free; the cause is not measured, so this is an impression for the operator to confirm or dismiss.

## Impact
flai stats items queue_time_seconds S-0216 305,070, S-0215 303,610, S-0214 301,459, S-0213 298,706, S-0212 295,977 against held_seconds 35,693, 34,246, 32,110, 29,392, 26,665; claims.days in_progress 1, 1, 2 of 5 and cfd review 1, 0, 0 of 5 over 2026-10-04 to 10-06. The metrics do not record why flai serve did not start them; the orchestrator ran only from 2026-10-07.

- time_lost_per_cycle: 75h

## Remediation

Found by the analysis in [2026-10-15-bottlenecks.md](../analysis/2026-10-15-bottlenecks.md).

Story S-0469 remediates this issue, created from it at 2026-10-21T11:58:17Z.
