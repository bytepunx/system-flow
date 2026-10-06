---
id: T-1018
type: task
nature: remediation
title: Fix the cause so the stopped orchestrator's activity is always logged, and the test passes under load
status: cancelled
parent: S-0292
owner: alex
created: 2026-10-06T11:34:24Z
updated: 2026-10-06T11:35:51Z
transitions:
  - to: cancelled
    at: 2026-10-06T11:35:51Z
    by: agent-S-0292
stream: S-0292
tags: [flai, serve]
touches: [flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/activity.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/plan.go]
after: [T-1017]
---
# T-1018 Fix the cause so the stopped orchestrator's activity is always logged, and the test passes under load

## Work

Build the fix T-1017 recorded in the narrative's `## Decisions`. It waits for T-1017 because what it changes depends on the cause found there, and both change `orchestrate_test.go`.

- If the activity entry is dropped or raced, fix it in flai serve: in `orchestrateEnded`, `stop`, `await`, or `logRunEnd`. Do not lengthen the test's wait to hide it. A run that ends without its activity logged must be warned of, not passed over silently.
- If the planner's run end (`planEnded` in `plan.go`) takes the same path, fix it the same way.
- If the cause is only the test's fixed 5 s deadline under load, give `waitFor` room for a loaded machine, and say in `## Decisions` why that is the cause and not a hidden defect.
- If what flai serve does changes, update `design/system/strategic-agents.md` and add it to this task's and the story's touches.

## Done when

- The test from T-1017 passes.
- `go test -race -count=20 -run TestTheOrchestratorIsStopped ./internal/serve` passes while the rest of the module's tests run in parallel.
- `scripts/flai-test.sh` passes.
- The first criterion of S-0292 is ticked.

## Notes

Planned by planner-S-0292. `plan.go` is predicted from its shared use of `logRunEnd`; drop it from the touches if the planner's path is not affected.
- 2026-10-06T11:35:51Z: moved to cancelled: T-1016, written before the planner's draft arrived, already reproduced, fixed, and closes I-0090 in one task
