---
id: T-3238
type: task
nature: feature
title: flai guard and flai serve's triggers read the new move condition, and a role starts on a move into a state of the workflow's own
status: in-progress
parent: S-0748
owner: alex
created: 2026-11-12T16:08:46Z
updated: 2026-11-15T20:40:51Z
transitions:
  - to: ready
    at: 2026-11-15T20:40:51Z
    by: agent-S-0748
  - to: in-progress
    at: 2026-11-15T20:40:51Z
    by: agent-S-0748
stream: S-0748
tags: [cli]
touches: [flai/internal/guard/roles.go, flai/internal/guard/roles_test.go, flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/replan.go, flai/internal/serve/orchestrate_reasons_test.go, flai/internal/serve/role_events_test.go, flai/cmd/guard_roles_test.go, flai/internal/hostapi/settings.go, flai/internal/hostapi/people_test.go]
after: [T-3231]
---
# T-3238 flai guard and flai serve's triggers read the new move condition, and a role starts on a move into a state of the workflow's own

## Work

Waits for T-3231: it reads the condition that task adds. Make `guard/roles.go`'s call checks and `serve`'s occurrences in `orchestrate_reasons.go` and `replan.go` give and read the state a move goes to through the new condition. Test that a role whose trigger names a state added between `in-progress` and `review` starts on a move into it, and that the planner's `item_move` stays held to backlog.

## Done when

The guard and the triggers read the new condition, and the tests pass on both workflows.

## Notes
