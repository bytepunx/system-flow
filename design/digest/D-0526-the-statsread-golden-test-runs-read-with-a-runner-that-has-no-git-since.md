---
id: D-0526
title: The statsread golden test runs Read with a runner that has no git, since testdata/good sits in this checkout and the real runner would read the monorepo's history into the claims drift
kind: decision
sources: [S-0607, T-2637, S-0608, T-2662]
topics: [server-side, back-end, go]
paths: [flai/internal/statsread/report_test.go, flai/internal/statsread/testdata/good-report.json, flai/internal/metrics/testdata/good, flai/internal/metrics/windows_test.go, flai/internal/metrics/testdata/windows.golden.json]
date: 2026-11-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0526 The statsread golden test runs Read with a runner that has no git, since testdata/good sits in this checkout and the real runner would read the monorepo's history into the claims drift

S-0607 added `TestReportOverTheGoodFixtureIsTheGolden` (`flai/internal/statsread/report_test.go`), which pins the report over `flai/internal/metrics/testdata/good` in `flai/internal/statsread/testdata/good-report.json`, before it made `statsread.Read` and `metrics.Compute` run at once. From S-0607's `## Decisions` (`wip/archive/agents/S-0607.md`):

```text
T-2637's golden runs `Read` with a runner that has no git: `testdata/good` sits in this checkout, so the real runner would read the monorepo's history into the claims drift, which changes with every commit. One compact line per case keeps the 356 KB file diffable by case.
```

A story that changes the golden or adds a case keeps the git-less runner and one line per case. The design records the test in `design/system/server-performance.md` § Cause 17, but not why it has no git.

Because that golden has no git, it leaves `claims.drift` out and has no 90-day case. So S-0608 pinned the touches drift in a second golden beside it, rather than extending it: `TestReportAtEveryWindowIsTheGolden` (`flai/internal/metrics/windows_test.go`, `flai/internal/metrics/testdata/windows.golden.json`) calls `metrics.Compute` with commits and projects given, at 1, 7, 30, 90, and 365 days. From S-0608's `## Decisions` (`wip/archive/agents/S-0608.md`):

```text
T-2662's golden stays beside S-0607's `TestReportOverTheGoodFixtureIsTheGolden` (`flai/internal/statsread/report_test.go`) rather than extending it: that one reads through `statsread.Read` without git, so `claims.drift` is left out and there is no 90-day case; T-2662's computes with commits and projects given, so the touch drift T-2683 keeps is pinned too.
```

So a story that needs the drift or a window pinned uses the metrics golden, and one that changes the read uses the statsread golden. The design names both tests (§ Cause 17 and § Cause 18) but not why there are two.

Sources: S-0607, T-2637, S-0608, T-2662.
