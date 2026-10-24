---
id: T-1834
type: task
nature: feature
title: flai serve starts the curator after each acceptance, coalesced, and when digest.schedule comes round, while the curate action is on
status: done
parent: S-0432
owner: alex
created: 2026-10-16T06:07:04Z
updated: 2026-10-24T14:25:35Z
transitions:
  - to: ready
    at: 2026-10-24T14:14:01Z
    by: agent-S-0432
  - to: in-progress
    at: 2026-10-24T14:14:02Z
    by: agent-S-0432
  - to: done
    at: 2026-10-24T14:25:35Z
    by: agent-S-0432
stream: S-0432
tags: [cli]
touches: [flai/internal/serve/curate_schedule.go, flai/internal/serve/curate_schedule_test.go, flai/internal/serve/curate.go, flai/internal/serve/curate_test.go, flai/internal/serve/serve.go]
after: [T-1832, T-1833]
usage:
  source: log
  seconds: 693
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 50
      output: 20017
      cache_read: 3303036
      cache_write: 136930
      cost: 1.8298
---
# T-1834 flai serve starts the curator after each acceptance, coalesced, and when digest.schedule comes round, while the curate action is on

## Work

Add `flai/internal/serve/curate_schedule.go`, a curator scheduler per served project, wired into the launcher loop in `serve.go` beside `newAnalysisScheduler`:

- After an acceptance: when the scheduler's look sees a story newly `done` since its last look and `curate` is on, it starts the curator with the trigger `accepted S-nnnn`; acceptances seen while a curator runs are coalesced into one start after it ends, with the trigger naming them.
- On `digest.schedule`: as the analysis scheduler does (`cadence`), with the trigger `schedule <spec>`; missed times come round once; a refused start is logged at info and not retried until the next time.
- While `curate` is off it starts nothing and keeps up with both, so turning the action on does not act on the past.

Waits for T-1832, which gives the schedule setting.

## Done when

`curate_schedule_test.go` covers one acceptance, two during a run coalesced into one start, an acceptance while the action is off, the schedule coming round, and a refused start; `flai test flai/internal/serve` passes.

## Notes

The replanner's look on work completing (`replan.go`, ADR-0084) is the model for seeing a story newly done.
