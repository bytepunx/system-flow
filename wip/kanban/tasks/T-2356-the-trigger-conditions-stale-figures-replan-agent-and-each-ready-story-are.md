---
id: T-2356
type: task
nature: feature
title: The trigger conditions stale-figures, replan-agent, and each-ready-story are implemented in flai/internal/roles/conditions.go
status: backlog
parent: S-0554
owner: alex
created: 2026-10-30T20:48:48Z
updated: 2026-10-30T20:48:48Z
transitions: []
stream: S-0554
tags: [cli]
touches: [flai/internal/roles/conditions.go, flai/internal/roles/conditions_test.go]
after: [T-2354]
---
# T-2356 The trigger conditions stale-figures, replan-agent, and each-ready-story are implemented in flai/internal/roles/conditions.go

## Work

`Matches` in `flai/internal/roles/conditions.go` reports the three trigger conditions not implemented. Implement them as ADR-0207 § 5 and the new ADR say, with the rules the replanner holds today in `flai/internal/serve/replan.go`. `stale-figures` passes an edit of a backlog or ready story's goal, criteria, or touches by someone other than a planner, when its forecast or cost of delay value is older than the edit (`editTrigger`). `replan-agent` passes, while `planning.replan` is `agent`, each story whose delivery moved (`replan`). `each-ready-story` gives one start per ready story on a schedule. Take what each needs from the event or the board as arguments, so the conditions stay free of `serve`.

Second layer: it waits for the ADR. It shares no path with the start prompt task.

## Done when

`conditions_test.go` covers each condition passing and refusing on the cases `replan_test.go` covers for the replanner: an edit by a planner, a fresh forecast, a field outside the three, `planning.replan` other than `agent`, and an empty ready column.

## Notes

Drafted by the planner of E-0024 on 2026-10-30. In S-0419's worktree that day, the conditions are named at `conditions.go:126` and `Matches` at `:89`; the replanner's rules are `editTrigger` (`replan.go:378`), `moveTriggers` (`:427`), and `replan` (`:195`).
