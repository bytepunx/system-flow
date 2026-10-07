---
id: T-1176
type: task
nature: feature
title: plan_backlog_stories is a permission, flai plan --candidates lists unplanned backlog stories, and flai holds the orchestrator's story plans and cost of delay edits to it
status: done
parent: S-0328
owner: alex
created: 2026-10-07T19:40:18Z
updated: 2026-10-07T20:00:10Z
transitions:
  - to: ready
    at: 2026-10-07T19:40:48Z
    by: agent-S-0328
  - to: in-progress
    at: 2026-10-07T19:40:48Z
    by: agent-S-0328
  - to: done
    at: 2026-10-07T20:00:10Z
    by: agent-S-0328
stream: S-0328
tags: []
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/workitem/plancandidates.go, flai/internal/workitem/plancandidates_test.go, flai/internal/workitem/promotable.go, flai/cmd/plan.go, flai/cmd/plan_candidates_test.go, flai/cmd/plan_orchestrate_test.go, flai/cmd/orchestrate_permissions_test.go, flai/cmd/edit.go, flai/cmd/edit_test.go, flai/internal/mcpserver/items_write.go, flai/internal/mcpserver/items_write_test.go, flai/internal/mcpserver/plan.go, flai/internal/mcpserver/plan_test.go, docs/operators/settings.md, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 1162
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 133
      output: 57980
      cache_read: 9835476
      cache_write: 247240
      cost: 4.673
---
# T-1176 plan_backlog_stories is a permission, flai plan --candidates lists unplanned backlog stories, and flai holds the orchestrator's story plans and cost of delay edits to it

## Work

- `manifest.Permissions` gains `PlanBacklogStories` (`plan_backlog_stories`), `PermitPlanBacklogStories`, a place in `PermissionNames` after `plan_backlog_epics`, and a case in `Allows`. `settings.go` gives it a meaning and a risk for the Orchestrator page.
- `PlanCandidatesOf` also judges stories: a story in the backlog, not archived, that lacks touches, a forecast duration, a cost of delay value, or a task that is not cancelled. Its reason names what it lacks. Stories follow the epics, in ID order. Each candidate's JSON gains `type`, `epic` or `story`. The `leave` map works for stories as for epics.
- `planHeld` in `cmd/plan.go` covers stories too. A story is also left out while a planner runs for its epic, and when its newest planner run ended and the story's `updated` is not after the run's end.
- `orchestratorPlans` takes a story under `plan_backlog_stories` as it takes an epic under `plan_backlog_epics`. The help text and the MCP `plan` description say so.
- Under `FLAI_ROLE=orchestrate`, an `item_edit` or `flai edit` that gives cost of delay inputs is held by flai: it needs `plan_backlog_stories`, a story in the backlog, no cost of delay inputs on the story or its epic, and no change but the inputs (`revenue_per_week`, `penalty_per_week`, `time_lost_per_cycle`). A value, a clear, or any other field is refused. The planner works out the value on its next run.
- `docs/operators/settings.md` gets the setting's row, and `docs/users/flai-reference.md` is regenerated, so the settings test passes.
- Tests show each case on and off. `cmd/orchestrate_permissions_test.go` changes only where flai's own refusals change; T-1177 changes it for the guard after this task.

This task waits on none and runs in layer 1 beside T-1174, sharing no path with it.

## Done when

- [ ] `flai test` passes on the paths changed.
- [ ] With the permission off, the orchestrator's story plan and cost of delay inputs are refused, naming `orchestration.permissions.plan_backlog_stories`. With it on, a listed story plans, and a backlog story without inputs takes them.

## Notes
