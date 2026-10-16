---
id: T-1792
type: task
nature: feature
title: flai stats reports strategic totals, items, issues, project, and use per day for every role a project defines
status: backlog
parent: S-0425
owner: alex
created: 2026-10-16T05:49:27Z
updated: 2026-10-16T05:49:32Z
transitions: []
stream: S-0425
tags: [cli]
touches: [flai/internal/metrics/strategic.go, flai/internal/metrics/strategic_test.go, flai/internal/metrics/usage.go, flai/internal/metrics/usage_test.go, flai/internal/statsread/statsread.go]
after: [T-1790]
---
# T-1792 flai stats reports strategic totals, items, issues, project, and use per day for every role a project defines

## Work

Make `strategic`, `strategicIssues`, `StrategicAgent`, and `StrategicKindSpend` in `flai/internal/metrics` iterate the defined roles instead of the three kinds. Keep the JSON shape for the three, so that the dashboard of today reads it, and add each defined role beside them. `statsread` reads the same. Test a fixture role's totals and per-day use; the existing stats tests stay unchanged.

Waits for the ADR. Runs beside the charge task, since they share no path.

## Done when

`go test ./internal/metrics/... ./internal/statsread/...` passes.

## Notes
