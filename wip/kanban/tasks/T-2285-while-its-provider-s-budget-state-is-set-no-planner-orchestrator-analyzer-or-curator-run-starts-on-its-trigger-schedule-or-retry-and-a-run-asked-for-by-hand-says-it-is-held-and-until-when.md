---
id: T-2285
type: task
nature: improvement
title: While its provider's budget state is set, no planner, orchestrator, analyzer, or curator run starts, on its trigger, schedule, or retry, and a run asked for by hand says it is held and until when
status: backlog
parent: S-0534
owner: alex
created: 2026-10-28T21:26:16Z
updated: 2026-10-28T21:26:16Z
transitions: []
stream: S-0534
tags: [flai, serve, budget]
touches: [flai/internal/serve/plan.go, flai/internal/serve/plan_test.go, flai/internal/serve/replan.go, flai/internal/serve/replan_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/orchestrate_hold.go, flai/internal/serve/orchestrate_hold_test.go, flai/internal/serve/analyze.go, flai/internal/serve/analyze_test.go, flai/internal/serve/analysis_schedule.go, flai/internal/serve/analysis_schedule_test.go, flai/internal/serve/curate.go, flai/internal/serve/curate_test.go, flai/internal/serve/curate_schedule.go, flai/internal/serve/curate_schedule_test.go]
after: [T-2283]
---
# T-2285 While its provider's budget state is set, no planner, orchestrator, analyzer, or curator run starts, on its trigger, schedule, or retry, and a run asked for by hand says it is held and until when

## Work

Check the budget state (T-2282) before each role's run starts, with the provider of the role's agent (`planning.agent`, `orchestration.agent`, `analysis.agent`, and the curator's) merged over the project's. Check before any run is recorded. Check at every start path:

- the planner: `planCheck` and `l.plan` in `plan.go`, and the replanner's queue in `replan.go`
- the orchestrator: `look` and `orchestrate` in `orchestrate.go`
- the analyzer: `analyzeCheck` and `l.analyze` in `analyze.go`, and its schedule in `analysis_schedule.go`
- the curator: `curateCheck` and `l.curate` in `curate.go`, and its schedule in `curate_schedule.go`

While the provider is held, nothing is lost:

- The replanner's queue keeps its items.
- The analysis and curate schedules keep their due run.
- The curator's coalesced acceptances stay queued.
- The orchestrator is not started; its `orchestrateRetry` clock does not run, and it starts at the first look after the state clears.

Reasons queued for the orchestrator under S-0536 wait the same way if that story is on main first. Once `until` has passed, a role's start may be the provider's one probe (`Probe`), as a story's agent's may.

A run the operator asks for by hand is refused with nothing recorded and the provider and `until` named. That covers `flai plan`, `flai analyze`, `flai curate`, and `flai serve orchestrate start` (`OrchestratorStart` in `orchestrate_hold.go`), and their host API and MCP forms. Project-defined roles start through the path S-0417 adds. If that path is on main when this task runs, hold it the same way and add its file to this task's touches.

This task waits for T-2283, which changes the roles' end paths first. It runs beside the story agents' hold.

## Done when

- [ ] With the state set, no role starts on that provider from its trigger, schedule, queue, or retry, and a role on another provider starts as before.
- [ ] Once the state clears, the queued planner items, the due schedules, the curator's acceptances, and the orchestrator start at the next look.
- [ ] `flai plan`, `flai analyze`, `flai curate`, and `flai serve orchestrate start` refuse with the provider and `until`, and record no run.
- [ ] `flai test` passes on the changed paths.

## Notes

Planned by planner-S-0534.
