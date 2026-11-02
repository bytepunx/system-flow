---
id: T-2014
type: task
nature: improvement
title: release_publish refuses, naming flai formula run publish, while the project defines a publish formula, and stays for a project with none
status: done
parent: S-0423
owner: alex
created: 2026-10-21T13:10:20Z
updated: 2026-11-02T04:22:49Z
transitions:
  - to: ready
    at: 2026-11-02T04:13:22Z
    by: agent-S-0423
  - to: in-progress
    at: 2026-11-02T04:13:24Z
    by: agent-S-0423
  - to: done
    at: 2026-11-02T04:22:49Z
    by: agent-S-0423
stream: S-0423
tags: [cli]
touches: [flai/internal/mcpserver/orchestrate.go, flai/internal/mcpserver/orchestrate_test.go]
usage:
  source: log
  seconds: 565
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 57
      output: 289
      cache_read: 2984065
      cache_write: 123034
      cost: 1.0935
---
# T-2014 release_publish refuses, naming flai formula run publish, while the project defines a publish formula, and stays for a project with none

## Work

In `flai/internal/mcpserver/orchestrate.go`, `releasePublish` reads the project's formulas with `flai/internal/formulas` before `mayPublish`. When one named `publish` exists, it refuses, changing nothing, saying that publishing is the `publish` formula and naming `flai formula run publish --ready`. With no such formula it publishes as today, so a project that has not taken the baseline formula keeps the tool.

Waits for the ADR (T-2015), which records the decision it enforces. Runs beside T-1780, since they share no path.

## Done when

`go test ./internal/mcpserver/...` passes, with a test for each case. The three existing `release_publish` tests are unchanged.

## Notes

Criterion 3: the orchestrator reaches the remote through the formula and no other way, and `release_publish` stays only where there is no `publish` formula.
