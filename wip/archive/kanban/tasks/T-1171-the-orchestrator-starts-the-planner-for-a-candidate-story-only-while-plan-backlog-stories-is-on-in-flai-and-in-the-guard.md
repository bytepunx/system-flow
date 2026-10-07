---
id: T-1171
type: task
nature: feature
title: The orchestrator starts the planner for a candidate story only while plan_backlog_stories is on, in flai and in the guard
status: cancelled
parent: S-0328
owner: alex
created: 2026-10-07T19:39:42Z
updated: 2026-10-07T19:42:15Z
transitions:
  - to: cancelled
    at: 2026-10-07T19:42:15Z
    by: agent-S-0328
stream: S-0328
tags: [flai, guard, orchestrator, planner]
touches: [flai/cmd/plan.go, flai/cmd/plan_orchestrate_test.go, flai/cmd/orchestrate_permissions_test.go, flai/cmd/guard.go, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/internal/guard/orchestrate_test.go, flai/internal/mcpserver/plan.go, flai/internal/mcpserver/plan_test.go]
after: [T-1169, T-1170]
---
# T-1171 The orchestrator starts the planner for a candidate story only while plan_backlog_stories is on, in flai and in the guard

## Work

Let the orchestrator ask for the planner on a story, as it may on an epic:

- `orchestratorPlans` in `flai/cmd/plan.go` takes a story `flai plan --candidates` lists while `plan_backlog_stories` is on, and an epic while `plan_backlog_epics` is on. Its refusals name the permission that is off, or why the story is not a candidate. Update the `--candidates` flag's help and the command's long text to say it lists stories too.
- `flai/internal/guard/guard.go` gives `flai plan` and the MCP tool `plan` on a story the need `PermitPlanBacklogStories` in place of the `plansEpics` refusal, for the orchestrator only. The planner, the analyzer, and sub-agents are still refused the tool.
- The MCP tool's description in `flai/internal/mcpserver/plan.go`, and the guard's help in `flai/cmd/guard.go`, say so.

It waits for T-1169, which declares the permission, and T-1170, which lists the story candidates it checks against.

## Done when

- With `plan_backlog_stories` on, the orchestrator starts the planner for a candidate story through `plan` and `flai plan`, and the run records the trigger `orchestrator`.
- With it off, or for a story the candidates do not list, both the guard and flai refuse, naming why; an epic is still governed by `plan_backlog_epics` alone.
- `TestTheOrchestratorNeverPlansAStoryEditsAnItemOrPlacesOneByHand` is rewritten for the new rule, and `plan_orchestrate_test.go`, `guard_test.go`, and `plan_test.go` cover both permissions on and off.
- `flai test` passes on the paths changed.

## Notes
- 2026-10-07T19:42:15Z: moved to cancelled: duplicates the story agent's plan (T-1174, T-1176 to T-1179), written at start while the planner ran; its refinements are folded into T-1176 to T-1178
