---
id: T-2650
type: task
nature: improvement
title: metrics.Compute runs its sections at once, with claims after replayHolds
status: done
parent: S-0607
owner: alex
created: 2026-11-07T23:49:17Z
updated: 2026-11-08T02:31:08Z
transitions:
  - to: ready
    at: 2026-11-08T02:18:04Z
    by: agent-S-0607
  - to: in-progress
    at: 2026-11-08T02:18:05Z
    by: agent-S-0607
  - to: done
    at: 2026-11-08T02:31:08Z
    by: agent-S-0607
stream: S-0607
tags: [cli]
touches: [flai/internal/metrics/metrics.go, flai/internal/metrics/metrics_test.go]
after: [T-2637]
usage:
  source: log
  seconds: 783
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 51
      output: 13478
      cache_read: 6364488
      cache_write: 392052
      cost: 3.6935
---
# T-2650 metrics.Compute runs its sections at once, with claims after replayHolds

## Work

Cause 17: `Compute` in `flai/internal/metrics/metrics.go` computes its sections one after another, about 450 ms at 30 days. Run them at once, to about 280 ms.

Sections that depend on others keep that order:

| Section | Waits for |
|---------|-----------|
| the per-item loop (`Derive`, `perItem`, `rep.Items`) | `threadWaits`, `forgeWaits`, `replayHolds`, `mergeQueue`, `CostPerAgentHour` |
| `summarise`, and the groups when `By` is set | the per-item loop |
| `aging`, `readyWait` | `summarise`, for `CycleTime.P85` and `QueueTime` |
| `forecasts`, `waiting`, `strategicDays` | the per-item loop |
| `claims` | `replayHolds` |
| `strategic` | `strategicIssues`, for `onIssues` |

These wait for nothing beyond the inputs: `throughput`, the burn-ups, `cfd`, `spendReport` with `spendOverTime` and `strategicTotals`, `costOfDelay`, `costOrder`, `concurrency`, `turns`, `coordination`, `digestSection`, and `closeOuts`.

- Each goroutine writes only its own field of `rep`. Leave `rep.Burnup` to one goroutine, since it is one map. Leave the nil-to-empty fixes at the end, after all have ended.
- Check that no section changes what another reads: `items`, `all`, `opt`'s slices, `perItem`, `waits`, and `replay`. Copy a slice before a section sorts it, as `strategic` already does with `slices.Clone`. Where a section changes shared input, fix it in that section's file and add that file to the task's touches.
- Use no new module.

It waits for T-2637, whose golden test shows the report did not change.

## Done when

- `go test -race ./internal/metrics/ ./internal/statsread/` passes, the golden test of T-2637 among it.
- `Compute` runs the sections that wait for nothing at once, and each dependent section after the sections the table names.

## Notes
