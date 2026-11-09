---
id: T-2683
type: task
nature: improvement
title: The hold replay and the touch drift are kept between requests until the item generation or the commits' key moves, and each request only slices them to its window
status: done
parent: S-0608
owner: alex
created: 2026-11-07T23:50:47Z
updated: 2026-11-09T10:13:53Z
transitions:
  - to: ready
    at: 2026-11-09T09:56:25Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T09:56:27Z
    by: system-flow
  - to: done
    at: 2026-11-09T10:13:53Z
    by: agent-S-0608
stream: S-0608
tags: [cli]
touches: [flai/internal/metrics/claims.go, flai/internal/metrics/metrics.go, flai/internal/metrics/claims_test.go, flai/internal/metrics/claims_bench_test.go, flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go, flai/internal/storygit/committed.go, flai/internal/storygit/committed_test.go]
after: [T-2668, T-2679]
usage:
  source: log
  seconds: 1045
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 62
      output: 437
      cache_read: 7085653
      cache_write: 90016
      cost: 2.6098
---
# T-2683 The hold replay and the touch drift are kept between requests until the item generation or the commits' key moves, and each request only slices them to its window

## Work

`replayHolds` (`flai/internal/metrics/claims.go`) replays all history on every call, and `touchDrift` reads every story's commits on every call. Neither depends on the window but through its slots: the replay's day and week slots are indexed from the window's start, and its open holds run to now. Keep, per repository root, the replay's timeline of hold intervals, independent of the window and of now, and the drift, keyed by the item store's generation from T-2668 and the commits' key `storygit.Committed` already keys its cache by (`flai/internal/storygit/committed.go`), which it answers alongside the commits. `statsread.Read` (`flai/internal/statsread/statsread.go`) hands both keys to `Compute` in its options. Each request slices the kept timeline to its window and closes the open holds at its now; a moved key replays again. Without keys, as in `flai stats`, nothing is kept and the answer is as today.

Waits for T-2668, whose generation is half the key, and for T-2679, which changes `claims.go` and `metrics.go` first.

## Done when

- `claims_test.go` shows a second request with the same keys does not replay, a moved generation or commits' key replays again, and the sliced replay equals `replayHoldsReference` in `claims_bench_test.go` at every window of T-2662.
- `BenchmarkReplayHolds` records a kept request's cost against a full replay.
- `go test ./internal/metrics ./internal/statsread ./internal/storygit` passes under `-race`, the golden test of T-2662 among them.

## Notes
