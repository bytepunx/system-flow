---
id: T-2679
type: task
nature: improvement
title: Each item's transitions are parsed once a call, and costOfDelay and claims sweep the days once from the first item's day
status: backlog
parent: S-0608
owner: alex
created: 2026-11-07T23:50:26Z
updated: 2026-11-07T23:50:39Z
transitions: []
stream: S-0608
tags: [cli]
touches: [flai/internal/metrics/metrics.go, flai/internal/metrics/costofdelay.go, flai/internal/metrics/claims.go, flai/internal/metrics/coordination.go, flai/internal/metrics/costofdelay_test.go, flai/internal/metrics/claims_test.go]
after: [T-2662]
---
# T-2679 Each item's transitions are parsed once a call, and costOfDelay and claims sweep the days once from the first item's day

## Work

`stateAt` (`flai/internal/metrics/metrics.go`) parses every transition's time again on every call, and `costOfDelay` (`flai/internal/metrics/costofdelay.go`, with `waited`) and `claims` (`flai/internal/metrics/claims.go`) call it for every item on every day of the window, from the window's start even before the first item existed. In `Compute`, parse each item's transition times once into a timeline that the burnup, the cfd, `costOfDelay`, `claims`, and `coordination.go` read. Make `costOfDelay` and `claims` sweep the days once, moving each item's state forward as the day passes its next transition, rather than asking every item on every day. Answer the days before the first item was created as empty without a loop, as `dayRange` does for the burnup.

Waits for T-2662, whose golden test shows the report is the same at every window.

## Done when

- `costofdelay_test.go` and `claims_test.go` show a window that starts before the first item gives the same days as before, and that the sweep agrees with `stateAt` on every day.
- `go test ./internal/metrics` passes, the golden test of T-2662 at all five windows among them.

## Notes
