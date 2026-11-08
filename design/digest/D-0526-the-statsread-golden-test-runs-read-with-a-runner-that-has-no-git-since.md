---
id: D-0526
title: The statsread golden test runs Read with a runner that has no git, since testdata/good sits in this checkout and the real runner would read the monorepo's history into the claims drift
kind: decision
sources: [S-0607, T-2637]
topics: [server-side, back-end, go]
paths: [flai/internal/statsread/report_test.go, flai/internal/statsread/testdata/good-report.json, flai/internal/metrics/testdata/good]
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

Sources: S-0607, T-2637.
