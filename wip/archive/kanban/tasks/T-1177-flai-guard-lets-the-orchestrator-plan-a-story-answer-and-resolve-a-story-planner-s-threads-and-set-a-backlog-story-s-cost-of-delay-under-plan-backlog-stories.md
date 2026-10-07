---
id: T-1177
type: task
nature: feature
title: flai guard lets the orchestrator plan a story, answer and resolve a story planner's threads, and set a backlog story's cost of delay under plan_backlog_stories
status: done
parent: S-0328
owner: alex
created: 2026-10-07T19:40:28Z
updated: 2026-10-07T20:17:46Z
transitions:
  - to: ready
    at: 2026-10-07T20:02:01Z
    by: agent-S-0328
  - to: in-progress
    at: 2026-10-07T20:02:01Z
    by: agent-S-0328
  - to: done
    at: 2026-10-07T20:17:05Z
    by: agent-S-0328
stream: S-0328
tags: []
touches: [flai/internal/guard/guard.go, flai/internal/guard/orchestrate_test.go, flai/internal/guard/thread_test.go, flai/internal/guard/guard_test.go, flai/cmd/guard.go, flai/cmd/orchestrate_permissions_test.go]
after: [T-1176]
usage:
  source: log
  seconds: 904
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 98
      output: 42550
      cache_read: 7217988
      cache_write: 181443
      cost: 3.4294
---
# T-1177 flai guard lets the orchestrator plan a story, answer and resolve a story planner's threads, and set a backlog story's cost of delay under plan_backlog_stories

## Work

In `Guard.orchestrate` and `orchestration()`, under `plan_backlog_stories`:

- MCP `plan` and `flai plan` on a story need the permission, as an epic needs `plan_backlog_epics`.
- `item_edit` giving only `id`, `hash`, `project`, and `cost_of_delay` with input keys (`revenue_per_week`, `penalty_per_week`, `time_lost_per_cycle`) needs the permission. So does `flai edit <S-nnnn>` with only the matching input flags beside the harmless ones. A value, a clear, or any other field stays refused.
- On a thread whose opener is a story's planner (`planner-S-nnnn`), `thread_reply` and `flai thread reply` pass as an answer with or without a source, or as a recommendation. `thread_resolve` and `flai thread resolve` pass. Confirming and `--by` another stay refused.
- With the permission off, each call is refused naming it, as today.

Update the package doc, `cmd/guard.go`'s help, and the tests: `orchestrate_test.go`, `thread_test.go`, `guard_test.go`'s `allOn` and its cases, and `cmd/orchestrate_permissions_test.go`'s `allBut`, each permission on and off. Rewrite `TestTheOrchestratorNeverPlansAStoryEditsAnItemOrPlacesOneByHand` for the new rule.

It waits on T-1176 for `manifest.PermitPlanBacklogStories`, and for its changes to `cmd/orchestrate_permissions_test.go`. It runs in layer 2 beside T-1178 and T-1179, sharing no path with them.

## Done when

- [ ] `flai test` passes on the paths changed.
- [ ] Each call above is let through with the permission on and refused, naming it, with it off.

## Notes
