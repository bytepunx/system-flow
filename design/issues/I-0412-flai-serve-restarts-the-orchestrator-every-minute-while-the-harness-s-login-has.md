---
id: I-0412
title: flai serve restarts the orchestrator every minute while the harness's login has expired, 666 runs in 14 hours, with no hold and no thread to the operator
class: defect
status: open
count: 1
cost: 14h
first_reported: 2026-11-27T22:46:27Z
last_reported: 2026-11-27T22:46:27Z
updated: 2026-11-27T22:46:27Z
usage:
  strategic:
    - kind: analyzer
      seconds: 85
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 55
          output: 5532
          cache_read: 391176
          cache_write: 29514
          cost: 0.9652
        - model: claude-opus-5-5
          input: 29
          output: 4664
          cache_read: 1264948
          cache_write: 28114
          cost: 0.487
---

# I-0412 flai serve restarts the orchestrator every minute while the harness's login has expired, 666 runs in 14 hours, with no hold and no thread to the operator

## Description
flai serve restarts the orchestrator every minute while the harness's login has expired, 666 runs in 14 hours, with no hold and no thread to the operator

## Instances

### 2026-11-27T22:46:27Z
Report: design/analysis/2026-11-27-all.md.
From 2026-11-17T17:24Z to 2026-11-18T07:07Z the claude-code login had expired; every orchestrator start failed in 90 ms with "Failed to authenticate: OAuth session expired and could not be refreshed" and flai serve started the next on its following look, 666 times, until the operator's session renewed the login. The orchestrator accepted, ordered, and answered nothing for 13.7 hours and the operator learned of it from the journal alone; the activity document shows none of it since a run that fails before its first tool call logs nothing.

## Impact
Serve journal: 666 orchestrator runs between 2026-11-17T17:24:45Z and 2026-11-18T07:07:39Z "ended, failed: ended (exit 1)", one every 73 s, every one for the reason "planned S-0822, wrote T-3520, T-3521, T-3522"; each log under ~/.flai/serve/agents/sf-orchestrator-20261117T17…Z.log to …20261118T07…Z.log ends with "Failed to authenticate: OAuth session expired and could not be refreshed", terminal_reason api_error, duration_ms 90, total_cost_usd 0. 697 exit-1 endings in the 30-day window; the same loop ran 2026-10-26T12Z to 2026-10-28T17Z at about 270 an hour. No hold, thread, or notice in the journal for those hours. ADR-0209 holds starts on budget exhaustion, ADR-0214 on the curator's usage limit, ADR-0108 bounds a story agent's restarts; an authentication failure matches none.

- time_lost_per_cycle: 14h

## Remediation

Found by the analysis in [2026-11-27-all.md](../analysis/2026-11-27-all.md).
