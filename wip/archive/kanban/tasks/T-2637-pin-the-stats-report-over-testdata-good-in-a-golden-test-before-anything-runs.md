---
id: T-2637
type: task
nature: improvement
title: Pin the stats report over testdata/good in a golden test before anything runs at once
status: done
parent: S-0607
owner: alex
created: 2026-11-07T23:48:37Z
updated: 2026-11-08T02:07:28Z
transitions:
  - to: ready
    at: 2026-11-08T01:46:05Z
    by: agent-S-0607
  - to: in-progress
    at: 2026-11-08T01:46:07Z
    by: agent-S-0607
  - to: done
    at: 2026-11-08T02:07:28Z
    by: agent-S-0607
stream: S-0607
tags: [cli]
touches: [flai/internal/statsread/report_test.go, flai/internal/statsread/testdata/good-report.json]
usage:
  source: log
  seconds: 1281
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 36
      output: 9610
      cache_read: 4538007
      cache_write: 279541
      cost: 2.6336
---
# T-2637 Pin the stats report over testdata/good in a golden test before anything runs at once

## Work

Criterion 1 asks that the report stay the same as before. Pin "before" while the code is still serial, so the later tasks are checked against it.

- Add `flai/internal/statsread/report_test.go`. It opens `../metrics/testdata/good` with `workitem.Open`, calls `statsread.Read`, sets a fixed `Now`, and calls `metrics.Compute`.
- Cover the windows and options the dashboard asks for: 1 day, 7 days, 30 days, and 365 days, each by story; one by `nature`; one with `bucket: hour`.
- Compare the JSON of each report with `flai/internal/statsread/testdata/good-report.json`. Write that file from the unchanged code, behind an `-update` flag.
- The test lives in `statsread`, not `metrics`: `statsread` imports `metrics`, so a test in `metrics` that calls `Read` would be an import cycle.
- `testdata/good` is not a git checkout, so `storygit.Committed` fails and `claims.drift` stays out. That holds before and after, so the golden is stable.

This task waits for none.

## Done when

- `go test -race ./internal/statsread/ -run Report` passes on the unchanged `Read` and `Compute`.
- `good-report.json` is committed, written by the unchanged code.

## Notes
