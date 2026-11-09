---
id: T-2647
type: task
nature: improvement
title: The replanner writes its moved forecasts through the batch entry
status: done
parent: S-0599
owner: alex
created: 2026-11-07T23:49:08Z
updated: 2026-11-08T00:43:14Z
transitions:
  - to: ready
    at: 2026-11-08T00:19:57Z
    by: system-flow
  - to: in-progress
    at: 2026-11-08T00:20:06Z
    by: system-flow
  - to: done
    at: 2026-11-08T00:43:14Z
    by: agent-S-0599
stream: S-0599
tags: [cli]
touches: [flai/internal/serve/replan.go, flai/internal/serve/replan_test.go, flai/internal/itemedit/forecasts.go, flai/internal/itemedit/forecasts_test.go]
after: [T-2641]
usage:
  source: log
  seconds: 1388
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 25
      output: 3910
      cache_read: 4711909
      cache_write: 80237
      cost: 1.528
---
# T-2647 The replanner writes its moved forecasts through the batch entry

## Work

In `flai/internal/serve/replan.go`, make `replan` build one change per move from `replays` and write them with one call to `itemedit.ApplyForecasts`, in place of the loop over `itemedit.Apply`.

- `written` is the items the call wrote. Each item it did not write is warned of as `forecast not replanned`, with its item and error, as now.
- The files it returns are committed in one commit as now, `chore: replan forecasts after <triggers>`, unless autocommit is off.
- The log line `forecasts replanned` keeps its fields.

Waits for T-2641: it calls `ApplyForecasts`.

## Done when

- The replan tests in `flai/internal/serve/replan_test.go` pass unchanged: the same forecasts move, and `replan-agent` queues the planner for the same stories.
- A new test shows a replan that moves several forecasts writes `wip/agents/index.md` once and makes one commit.
- A new test shows a forecast another writer changed between the replay and the write is warned of and left as that writer saved it, and the other forecasts are written.
- `flai test flai/internal/serve` passes.

## Notes
