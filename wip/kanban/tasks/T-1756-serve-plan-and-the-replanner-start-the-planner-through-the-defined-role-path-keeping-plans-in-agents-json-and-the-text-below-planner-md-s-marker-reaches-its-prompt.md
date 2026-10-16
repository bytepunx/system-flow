---
id: T-1756
type: task
nature: improvement
title: serve.Plan and the replanner start the planner through the defined-role path, keeping plans in agents.json, and the text below planner.md's marker reaches its prompt
status: backlog
parent: S-0419
owner: alex
created: 2026-10-16T05:46:13Z
updated: 2026-10-16T05:46:13Z
transitions: []
stream: S-0419
tags: [cli]
touches: [flai/internal/serve/plan.go, flai/internal/serve/replan.go, flai/internal/serve/role.go, flai/internal/serve/role_test.go]
after: [T-1753]
---
# T-1756 serve.Plan and the replanner start the planner through the defined-role path, keeping plans in agents.json, and the text below planner.md's marker reaches its prompt

## Work

- Make `serve.Plan`, `planCheck`, `planStart`, and `planEnded` call the defined-role path for `planner` with the item and the trigger.
- Keep the run's record under `plans` and its refusal texts. Keep the replanner's queue as it is: it starts its runs through the same path.
- Add a test in `role_test.go` that text a project writes below the marker of `planner.md` reaches the planner's prompt.

Waits for the prompt task. The analyzer's and the orchestrator's tasks wait for this one, which shapes `role.go` for the built-in roles first.

## Done when

`plan_test.go`, `replan_test.go`, and the `flai/cmd` plan tests pass unchanged, and the marker test passes.

## Notes
