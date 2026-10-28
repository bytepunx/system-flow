---
id: D-0361
title: The chart data takes about 4 s warm, against 10 ms measured on 2026-09-29, almost all of it in metrics.replayHolds
kind: cause
sources: [I-0155]
topics: [metrics, performance, flai, cli, go, code, all]
paths: [flai/internal/metrics/claims.go, flai/internal/metrics/claims_test.go, flai/internal/metrics/claims_bench_test.go, flai/internal/workitem/hold.go, design/system/server-performance.md, design/issues/I-0155-the-chart-data-takes-about-4-s-warm-against-10-ms-measured-on-2026-09-29-almost-all-of-it-in-metrics-replayholds.md, design/issues/summary.md, flai/internal/workitem/hold_test.go, design/issues/I-0184-every-stats-get-and-flai-stats-runs-one-git-log-for-claims-drift-about-450-ms-of-a-900-ms-warm-answer.md]
date: 2026-10-21
status: absorbed
absorbed_by: "design/system/server-performance.md#a-regression-in-the-chart-data-the-hold-replay-s-0471"
asked_again: []
cited: []
---

# D-0361 The chart data takes about 4 s warm, against 10 ms measured on 2026-09-29, almost all of it in metrics.replayHolds

I-0155 was closed on 2026-10-21.

The cause:

```text
The chart data takes about 4 s warm, against 10 ms measured on 2026-09-29, almost all of it in metrics.replayHolds
```

The fix:

```text
S-0471: replayHolds judges the holds again only at an event time that can change one, over only the open stories, with a HoldsMemo of claims between judgements, and PathsOverlap no longer allocates. One replay at this repository's scale went from 2,726 ms to 71 ms with the same result. On this repository on 2026-10-21, stats.get warm went from 5,695 ms to 894 ms and replayHolds from 73 to 22 percent of the samples. Half of what remains is a git log per call, recorded as I-0184.
```

Source: I-0155.
