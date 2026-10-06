---
id: T-0949
type: task
nature: feature
title: system-flow.yaml takes analysis.agent and analysis.schedule, and flai check reports a bad one
status: done
parent: S-0223
owner: alex
created: 2026-10-05T05:46:03Z
updated: 2026-10-06T20:14:59Z
transitions:
  - to: ready
    at: 2026-10-06T20:08:34Z
    by: agent-S-0223
  - to: in-progress
    at: 2026-10-06T20:08:34Z
    by: agent-S-0223
  - to: review
    at: 2026-10-06T20:14:59Z
    by: agent-S-0223
  - to: done
    at: 2026-10-06T20:14:59Z
    by: agent-S-0223
stream: S-0223
tags: [flai]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/check/check.go, flai/internal/check/analysis_test.go]
usage:
  source: log
  seconds: 308
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 22
      output: 129
      cache_read: 519455
      cache_write: 68030
      cost: 0.2557
---
# T-0949 system-flow.yaml takes analysis.agent and analysis.schedule, and flai check reports a bad one

## Work

Add an `Analysis` block to `manifest.Manifest` (`flai/internal/manifest/manifest.go`, beside `Planning`):

- `analysis.agent`, an `Agent` as `planning.agent` is, and `Manifest.AnalysisAgent()`, which merges it over the project's `agent` field by field, config key by config key, and role by role, as `PlanningAgent()` does.
- `analysis.schedule`, a five-field cron expression in UTC or `daily`, parsed with `internal/cron` as `planning.schedule` is, and `Analysis.AnalysisSchedule()` returning the parsed schedule, nil when unset. An expression that never comes round is refused, as the planner's is.

`flai check` reports an `analysis.schedule` it cannot parse and an `analysis.agent` that is not valid as `manifest.analysis` (`flai/internal/check/check.go`), each saying how to fix it.

The manifest's documentation is the last task's, so that one task writes `project-manifest.md`.

This task waits for nothing. It runs with the prompt task, whose paths it does not share.

## Done when

- the manifest parses `analysis.agent` and `analysis.schedule`, and a test pins both as unset by default
- `AnalysisAgent()` merges over the project's agent, and a test pins the merge as `PlanningAgent()`'s test does
- a bad schedule and a bad agent are each reported as `manifest.analysis` on a fixture
- `go test ./internal/manifest/ ./internal/check/` passes

## Notes
