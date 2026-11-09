---
id: T-2675
type: task
nature: improvement
title: costOrder builds its forecaster once and plays the three orders out on it
status: in-progress
parent: S-0608
owner: alex
created: 2026-11-07T23:50:18Z
updated: 2026-11-09T09:43:22Z
transitions:
  - to: ready
    at: 2026-11-09T09:43:21Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T09:43:22Z
    by: system-flow
stream: S-0608
tags: [cli]
touches: [flai/internal/planning/forecast.go, flai/internal/planning/forecast_test.go, flai/internal/metrics/costorder.go, flai/internal/metrics/costorder_test.go]
after: [T-2662]
---
# T-2675 costOrder builds its forecaster once and plays the three orders out on it

## Work

`costOrder` (`flai/internal/metrics/costorder.go`) calls `planning.PullTimes` once for each of its three orders, and each call builds `newForecaster` again (`flai/internal/planning/forecast.go`). Let `planning` build the forecaster once for the items, the order, the limit, the fallback, and now, and answer the pull times of each sequence from it; keep `PullTimes` as it answers today for its other callers. `costOrder` builds it once and plays the three orders out on it.

Waits for T-2662, whose golden test shows the cost order is the same.

## Done when

- `forecast_test.go` shows the pull times from one forecaster equal `PullTimes` for the same sequences.
- `costorder_test.go` shows `costOrder` builds the forecaster once.
- `go test ./internal/planning ./internal/metrics` passes, the golden test of T-2662 among them.

## Notes
