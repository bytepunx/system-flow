---
id: T-2358
type: task
nature: feature
title: The planner starts from its role file's triggers through the event and schedule launchers, which take an item-taking built-in role, in place of the replanner's own triggers
status: backlog
parent: S-0554
owner: alex
created: 2026-10-30T20:49:07Z
updated: 2026-10-30T20:49:07Z
transitions: []
stream: S-0554
tags: [cli]
touches: [flai/internal/serve/replan.go, flai/internal/serve/role_events.go, flai/internal/serve/role_events_test.go, flai/internal/serve/role_schedule.go, flai/internal/serve/role_schedule_test.go]
after: [T-2356]
---
# T-2358 The planner starts from its role file's triggers through the event and schedule launchers, which take an item-taking built-in role, in place of the replanner's own triggers

## Work

`eventRoles` in `flai/internal/serve/role_events.go` and `scheduledRoles` in `role_schedule.go` skip the built-in roles, and `lookAt` skips a role that takes an item. Let both take the planner, read its `trigger` entries from its role file, and evaluate their conditions with the conditions task's functions. A schedule on an item-taking role under `each-ready-story` queues one run per ready story. Starts still go through the replanner's queue and `planCheck`, so the planner's coalescing, records under `plans`, and journal entries stay as they are. Remove the triggers `replan.go` decides by itself (`editTrigger`, `moveTriggers`, the reorder and schedule loops), keeping its queue and `drain`.

Third layer: it waits for the conditions task. It shares no path with the start prompt rendering task. The analyzer's and the curator's triggers wait for it, since they use the launchers it changes.

## Done when

`replan_test.go` passes unchanged. `role_events_test.go` and `role_schedule_test.go` gain cases for the planner started on an edit, an acceptance under `planning.replan: agent`, a reorder, and its schedule, and for a project's planner role file whose trigger list drops the schedule, which then queues nothing on it.

## Notes

Drafted by the planner of E-0024 on 2026-10-30. In S-0419's worktree that day, the skips are at `role_events.go:145`, `role_schedule.go:201`, and `:142`.
