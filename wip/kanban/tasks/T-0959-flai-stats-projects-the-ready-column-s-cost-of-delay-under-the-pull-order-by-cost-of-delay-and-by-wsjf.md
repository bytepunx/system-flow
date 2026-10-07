---
id: T-0959
type: task
nature: feature
title: flai stats projects the ready column's cost of delay under the pull order, by cost of delay, and by WSJF
status: done
parent: S-0213
owner: alex
created: 2026-10-05T05:46:37Z
updated: 2026-10-07T07:50:53Z
transitions:
  - to: ready
    at: 2026-10-07T07:40:00Z
    by: agent-S-0213
  - to: in-progress
    at: 2026-10-07T07:40:00Z
    by: agent-S-0213
  - to: done
    at: 2026-10-07T07:50:52Z
    by: agent-S-0213
stream: S-0213
tags: [flai]
touches: [flai/internal/metrics/costorder.go, flai/internal/metrics/costorder_test.go, flai/internal/metrics/costofdelay.go, flai/internal/metrics/metrics.go, flai/internal/planning/forecast.go, flai/internal/statsread, flai/internal/metrics/costofdelay_test.go]
after: [T-0953, T-0955]
usage:
  source: log
  seconds: 652
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 72
      output: 28647
      cache_read: 4851561
      cache_write: 120877
      cost: 2.2687
---
# T-0959 flai stats projects the ready column's cost of delay under the pull order, by cost of delay, and by WSJF

## Work

Compute `cost_of_delay.order` as T-0953 defines it, in a new `flai/internal/metrics/costorder.go`:

- Take the ready stories in the board's pull order (`workitem.PullSequence`). Order them by `cod` and `wsjf` with S-0217's policy function in `workitem/policy.go`, which this story waits for, rather than a second copy of it.
- Lay each order onto the in-progress limit as `planning.playOut` lays the pull order. Export from `planning/forecast.go` the smallest piece that replays a given order: lanes freed by the stories in progress, and each story's forecast duration. `flai forecast`'s output must not change.
- Emit for each order its points and `total`, then `saving` and `left_out`, under `CostOfDelay` in `costofdelay.go`.

`metrics.Options` gains the board's order and the planning settings the replay needs. `flai/internal/statsread`, which builds the options for `flai stats` and for the dashboard's `stats.get`, passes them, so both carry the projection.

It waits for T-0953, whose definition it implements, and for T-0955, which changes the same types in `costofdelay.go`.

## Done when

- fixture tests in `costorder_test.go` pin the three orders' points, totals, and saving to the cent on a ready column where the pull order is not the cheapest; a story without a value and one without a forecast are in `left_out`; an empty ready column gives empty series
- `flai forecast` tests pass unchanged
- `flai stats --json` and `stats.get` carry `cost_of_delay.order` (the hostapi contract test passes)
- `go test ./internal/metrics/ ./internal/planning/ ./internal/statsread/ ./cmd/` passes

## Notes
