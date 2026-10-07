---
id: T-1170
type: task
nature: feature
title: "flai plan --candidates lists the backlog stories that lack a plan: touches, a forecast, a cost of delay value, or tasks"
status: cancelled
parent: S-0328
owner: alex
created: 2026-10-07T19:39:33Z
updated: 2026-10-07T19:42:15Z
transitions:
  - to: cancelled
    at: 2026-10-07T19:42:15Z
    by: agent-S-0328
stream: S-0328
tags: [flai, planner, orchestrator]
touches: [flai/internal/workitem/plancandidates.go, flai/internal/workitem/plancandidates_test.go, flai/cmd/plan_candidates_test.go]
after: [T-1168]
---
# T-1170 flai plan --candidates lists the backlog stories that lack a plan: touches, a forecast, a cost of delay value, or tasks

## Work

Extend `PlanCandidatesOf` in `flai/internal/workitem/plancandidates.go` to list stories beside epics: each story in the backlog, not archived, that lacks touches, a forecast duration, a cost of delay value, or a task, with a reason naming what it lacks. Leave a story out, with its reason among those left out, while a planner runs for it or while a thread its planner opened awaits the operator, as an epic is left out. Each candidate's JSON says whether it is an epic or a story, so the orchestrator can tell which permission it needs. Keep the epics' entries as they are.

It waits for T-1168, whose ADR defines what a story candidate is. It shares no path with T-1169 and runs beside it.

## Done when

- `flai plan --candidates` and `--json` list a backlog story lacking any part of a plan, with what it lacks, and leave out one that has them all, one a planner runs for, and one whose planner's thread awaits the operator.
- `plancandidates_test.go` and `plan_candidates_test.go` cover each case, and the epics' cases still pass.
- `flai test` passes on the paths changed.

## Notes

The `--candidates` flag's help in `flai/cmd/plan.go` changes in T-1171, which owns that file, so that the two tasks of this layer share no path.
- 2026-10-07T19:42:15Z: moved to cancelled: duplicates the story agent's plan (T-1174, T-1176 to T-1179), written at start while the planner ran; its refinements are folded into T-1176 to T-1178
