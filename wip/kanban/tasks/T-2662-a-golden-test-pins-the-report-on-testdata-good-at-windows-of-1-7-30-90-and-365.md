---
id: T-2662
type: task
nature: improvement
title: A golden test pins the report on testdata/good at windows of 1, 7, 30, 90, and 365 days before anything changes
status: in-progress
parent: S-0608
owner: alex
created: 2026-11-07T23:49:46Z
updated: 2026-11-09T09:12:41Z
transitions:
  - to: ready
    at: 2026-11-09T09:12:40Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T09:12:41Z
    by: system-flow
stream: S-0608
tags: [cli]
touches: [flai/internal/metrics/windows_test.go, flai/internal/metrics/testdata/windows.golden.json]
---
# T-2662 A golden test pins the report on testdata/good at windows of 1, 7, 30, 90, and 365 days before anything changes

## Work

Write `flai/internal/metrics/windows_test.go`. It computes the report with `metrics.Compute` over `testdata/good` at a fixed `Now`, for windows of 1, 7, 30, 90, and 365 days, with commits and projects given so that `claims`, `touchDrift`, and the cost order are filled. It compares each report, as JSON, with `testdata/windows.golden.json`, and an `-update` flag writes that file. Write the golden file from the code as it is now, before any other task of S-0608 changes it. Every later task runs this test to show the report did not change.

Waits for nothing.

## Done when

- `go test ./internal/metrics -run Windows` passes on the code as it is, and the golden file holds all five windows.
- Each window's `claims`, `cost_of_delay`, and cost order sections are not empty in the golden file.

## Notes
