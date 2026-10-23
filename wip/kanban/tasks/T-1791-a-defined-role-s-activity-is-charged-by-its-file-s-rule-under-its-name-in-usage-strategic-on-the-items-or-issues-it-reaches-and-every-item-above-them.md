---
id: T-1791
type: task
nature: feature
title: A defined role's activity is charged by its file's rule under its name in usage.strategic, on the items or issues it reaches and every item above them
status: backlog
parent: S-0425
owner: alex
created: 2026-10-16T05:49:23Z
updated: 2026-10-23T15:36:40Z
transitions: []
stream: S-0425
tags: [cli]
touches: [flai/internal/serve/activity.go, flai/internal/serve/activity_test.go, flai/internal/workitem/usage.go, flai/internal/workitem/usage_test.go, flai/internal/usage/usage.go, flai/internal/usage/usage_test.go, flai/internal/issues/usage.go, flai/internal/issues/usage_test.go, flai/cmd/activity.go, flai/cmd/activity_test.go, flai/internal/mcpserver/activity.go, flai/internal/mcpserver/activity_test.go]
after: [T-1790, T-2073]
---
# T-1791 A defined role's activity is charged by its file's rule under its name in usage.strategic, on the items or issues it reaches and every item above them

## Work

`StrategicKinds` in `flai/internal/usage/usage.go` and the charge code in `flai/internal/serve/activity.go` (`charge`, `chargePlanned`, `chargeNamed`, and `chargeIssues`) know three kinds.

- Make the kinds the defined roles. `StrategicKinds` stays the list of the three built-in kinds, ordered first; defined roles follow by name.
- Pick the charge by the role's `charge` field, read from its role file (T-2073): `item` charges the item its run was started for, found from the run recorded under `roles` in `serve/agents.json` as `planned()` reads `plans`; `named` the work items its entry names; `issues` the issues it names. A role without `charge` charges nothing. A built-in kind with no role file keeps its rule by name.
- `ChargeStrategic` and `StrategicErrors` in `workitem/usage.go`, and `ChargeStrategic` in `issues/usage.go`, take any role's name.
- `chargedTo` in `cmd/activity.go` says what an activity charged by its rule, not its kind, and `activity_log`'s description in `mcpserver/activity.go` states the three rules as the role file's.
- Test a fixture role under each rule. The existing charge tests stay unchanged.

Waits for the ADR, which fixes the rules, and for T-2073, which gives it the parsed `charge` field. Runs beside T-1793; they share no path.

## Done when

`go test ./internal/serve/... ./internal/workitem/... ./internal/usage/... ./internal/issues/... ./internal/mcpserver/... ./cmd/...` passes.

## Notes

Rewritten on TH-0609 with the orchestrator's approval, 2026-10-23.
