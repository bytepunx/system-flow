---
id: T-2410
type: task
nature: improvement
title: flai stats reports the share of the window with n stories in progress at once, and the ready wait at p50 and p85
status: backlog
parent: S-0551
owner: alex
created: 2026-10-31T07:44:26Z
updated: 2026-10-31T07:44:26Z
transitions: []
stream: S-0551
tags: [flai, go, metrics]
touches: [flai/internal/metrics/concurrency.go, flai/internal/metrics/concurrency_test.go, flai/internal/metrics/metrics.go, flai/cmd/stats.go, flai/cmd/stats_test.go]
after: [T-2388]
---
# T-2410 flai stats reports the share of the window with n stories in progress at once, and the ready wait at p50 and p85

## Work

Compute, in a new `flai/internal/metrics/concurrency.go`, the share of the stats window during which n stories were in progress at once, for each n from 0 up, from the stories' state intervals, and the mean n while one or more ran. Add it to the report's summary in `metrics.go` under the JSON key `metrics.md` names. Report the ready wait as the stories' queue time at p50 and p85, which `Summary` already computes ([ADR-0190](../../../design/adrs/0190-an-item-s-queue-time-is-the-seconds-it-spent-in-ready-before-it-started-so-a.md)); print both in `flai/cmd/stats.go`'s flow group, the queue time under the name ready wait.

It waits for the ADR task, which defines the metric in `metrics.md`.

## Done when

- A test in `concurrency_test.go` shows the shares for overlapping and idle intervals summing to 1 over the window.
- A test in `stats_test.go` shows `flai stats` and `flai stats --json` giving the shares and the ready wait at p50 and p85.
- `flai test flai/internal/metrics flai/cmd` passes.

## Notes

S-0552 in progress changes `flai/cmd/stats.go` and `stats_test.go` for the queue's stats; agree the order in the two stories' conversation.
