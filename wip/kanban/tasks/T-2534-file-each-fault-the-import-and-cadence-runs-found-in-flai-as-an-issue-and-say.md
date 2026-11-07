---
id: T-2534
type: task
nature: improvement
title: File each fault the import and cadence runs found in flai as an issue, and say in this task's Notes when they found none
status: done
parent: S-0564
owner: alex
created: 2026-11-07T03:48:11Z
updated: 2026-11-07T08:52:08Z
transitions:
  - to: ready
    at: 2026-11-07T08:51:40Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T08:51:43Z
    by: system-flow
  - to: done
    at: 2026-11-07T08:52:08Z
    by: agent-S-0564
stream: S-0564
tags: [template]
touches: [design/issues/summary.md]
after: [T-2523, T-2526]
usage:
  source: log
  seconds: 25
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 633
      cache_read: 302320
      cache_write: 3051
      cost: 0.0975
---
# T-2534 File each fault the import and cadence runs found in flai as an issue, and say in this task's Notes when they found none

## Work

Read what the two runs did, from T-2514's and T-2520's Notes, both instances' Notes, their threads, and the serve journal. A fault is flai or a formula falling short. Examples:

- the board's Import did not start the formula
- a step's command refused or reported the wrong thing
- `role_run` did not start the analyzer or did not end the step
- the listing step found no report
- a run needed a hand step the formula does not name

For each fault, record it with `flai issue new`, or bump the open issue that records it with `flai issue bump`. Name the run and the instance in the occurrence. A finding of the analyzer's report is not a fault of the run: the report files it already.

When the runs found no fault, say so in this task's Notes with `flai edit`.

This task waits for T-2523 and T-2526: both runs must be over. It runs beside the task that records the runs in S-0564's Notes, with which it shares no path.

## Done when

- Each fault the runs found is recorded or bumped as an issue, and `design/issues/summary.md` is current.
- When the runs found none, this task's Notes say so.

## Notes
