---
id: T-1791
type: task
nature: feature
title: A defined role's activity is charged by its file's rule under its name in usage.strategic, on the items or issues it reaches and every item above them
status: backlog
parent: S-0425
owner: alex
created: 2026-10-16T05:49:23Z
updated: 2026-10-16T05:49:31Z
transitions: []
stream: S-0425
tags: [cli]
touches: [flai/internal/serve/activity.go, flai/internal/serve/activity_test.go, flai/internal/workitem/usage.go, flai/internal/workitem/usage_test.go, flai/internal/usage/usage.go, flai/internal/usage/usage_test.go]
after: [T-1790]
---
# T-1791 A defined role's activity is charged by its file's rule under its name in usage.strategic, on the items or issues it reaches and every item above them

## Work

`StrategicKinds` in `flai/internal/usage/usage.go` and the charge code in `flai/internal/serve/activity.go` (`chargePlanned`, `chargeNamed`, and `chargeIssues`) know three kinds.

- Make the kinds the defined roles.
- Pick the charge by the role's `charge` field: `item` charges the item its run was started for, `named` the work items its entry names, and `issues` the issues it names.
- `ChargeStrategic` and `StrategicErrors` in `workitem/usage.go` take any role's name.
- Test a fixture role under each rule. The existing charge tests stay unchanged.

Waits for the ADR, which fixes the rules.

## Done when

`go test ./internal/serve/... ./internal/workitem/... ./internal/usage/...` passes.

## Notes
