---
id: T-2273
type: task
nature: improvement
title: The planner's, the orchestrator's, the analyzer's, and the curator's runs append a started, ended, or failed agent event naming their role and item
status: backlog
parent: S-0532
owner: alex
created: 2026-10-28T21:09:04Z
updated: 2026-10-28T21:09:04Z
transitions: []
stream: S-0532
tags: [flai, serve, events]
touches: [flai/internal/serve/plan.go, flai/internal/serve/plan_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/analyze.go, flai/internal/serve/analyze_test.go, flai/internal/serve/curate.go, flai/internal/serve/curate_test.go]
after: [T-2271]
---
# T-2273 The planner's, the orchestrator's, the analyzer's, and the curator's runs append a started, ended, or failed agent event naming their role and item

## Work

Call the serve helper where each strategic role's run starts and ends, beside its journal entry:

- the planner in `plan.go` (its start, and `planEnded` after `judgePlan`), with the epic or story it planned as the item;
- the orchestrator in `orchestrate.go` (its start, and `orchestrateEnded`), including each retry `orchestrateRetry` makes, so that a restart storm shows as one `failed` and one `started` per minute;
- the analyzer in `analyze.go` (its start, and `analyzeEnded`);
- the curator in `curate.go` (its start, and `curateEnded`).

`to` is `started`, `ended` for a run that did its work or asked, and `failed` otherwise, with the role's own judged reason as `why`. A start that could not be made is `failed` with the error. The operator's stop is written by the stop path, which the story's agent task changes; this task does not touch it.

A project's own role (S-0417) starts through `flai/internal/serve/role.go` once that story lands. If it is on main when this task starts, call the helper there too and add its file to this task's touches; if not, the built-in four are the scope, as the plan's thread says.

This task waits for T-2271 alone, and runs beside the story's agent task, since the two share no file.

## Done when

- [ ] A test per role reads back one `started` and one `ended` or `failed` event for a run, with the role, the item where it has one, and the `why` its journal entry gives
- [ ] An orchestrator run that fails and is retried writes one `failed` and one `started` per attempt
- [ ] `flai test` passes on the changed paths

## Notes

Planned by planner-S-0532. S-0514 changes `curate.go` and `curate_test.go` too; S-0532 names it in `after`, so its change is on main first.
