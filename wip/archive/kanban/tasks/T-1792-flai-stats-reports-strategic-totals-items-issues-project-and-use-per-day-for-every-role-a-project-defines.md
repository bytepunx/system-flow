---
id: T-1792
type: task
nature: feature
title: flai stats reports strategic totals, items, issues, project, and use per day for every role a project defines
status: done
parent: S-0425
owner: alex
created: 2026-10-16T05:49:27Z
updated: 2026-11-02T19:34:17Z
transitions:
  - to: ready
    at: 2026-11-02T19:23:43Z
    by: agent-S-0425
  - to: in-progress
    at: 2026-11-02T19:23:44Z
    by: agent-S-0425
  - to: done
    at: 2026-11-02T19:34:17Z
    by: agent-S-0425
stream: S-0425
tags: [cli]
touches: [flai/internal/metrics/strategic.go, flai/internal/metrics/strategic_test.go, flai/internal/metrics/usage.go, flai/internal/metrics/usage_test.go, flai/internal/statsread/statsread.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go, flai/internal/workitem/activity.go, flai/internal/metrics/metrics.go]
after: [T-1791]
usage:
  source: log
  seconds: 633
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 45
      output: 14915
      cache_read: 2676185
      cache_write: 70552
      cost: 1.2538
---
# T-1792 flai stats reports strategic totals, items, issues, project, and use per day for every role a project defines

## Work

Make `strategic`, `strategicIssues`, `StrategicAgent`, and `StrategicKindSpend` in `flai/internal/metrics` iterate the defined roles instead of the three kinds. Keep the JSON shape for the three, so that the dashboard of today reads it, and add each defined role beside them, built-in kinds first and the rest by name. `statsread` reads the project's roles as well as the activity documents, so that a role with no activity document yet is listed with zeros, as criterion 3 asks of every role a project defines. `strategic_days[].agents` carries each role by name. `cmd/stats.go` prints the same lines per role. Test a fixture role's totals and per-day use; the existing stats tests stay unchanged.

Waits for the ADR. Runs beside T-2073, since they share no path.

## Done when

`go test ./internal/metrics/... ./internal/statsread/... ./cmd/...` passes.

## Notes

Rewritten on TH-0609 with the orchestrator's approval, 2026-10-23.
