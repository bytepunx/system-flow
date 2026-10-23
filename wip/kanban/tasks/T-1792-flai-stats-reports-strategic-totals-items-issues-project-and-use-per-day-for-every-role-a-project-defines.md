---
id: T-1792
type: task
nature: feature
title: flai stats reports strategic totals, items, issues, project, and use per day for every role a project defines
status: backlog
parent: S-0425
owner: alex
created: 2026-10-16T05:49:27Z
updated: 2026-10-23T15:36:47Z
transitions: []
stream: S-0425
tags: [cli]
touches: [flai/internal/metrics/strategic.go, flai/internal/metrics/strategic_test.go, flai/internal/metrics/usage.go, flai/internal/metrics/usage_test.go, flai/internal/statsread/statsread.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go]
after: [T-1790]
---
# T-1792 flai stats reports strategic totals, items, issues, project, and use per day for every role a project defines

## Work

Make `strategic`, `strategicIssues`, `StrategicAgent`, and `StrategicKindSpend` in `flai/internal/metrics` iterate the defined roles instead of the three kinds. Keep the JSON shape for the three, so that the dashboard of today reads it, and add each defined role beside them, built-in kinds first and the rest by name. `statsread` reads the project's roles as well as the activity documents, so that a role with no activity document yet is listed with zeros, as criterion 3 asks of every role a project defines. `strategic_days[].agents` carries each role by name. `cmd/stats.go` prints the same lines per role. Test a fixture role's totals and per-day use; the existing stats tests stay unchanged.

Waits for the ADR. Runs beside T-2073, since they share no path.

## Done when

`go test ./internal/metrics/... ./internal/statsread/... ./cmd/...` passes.

## Notes

Rewritten on TH-0609 with the orchestrator's approval, 2026-10-23.
