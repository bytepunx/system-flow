---
id: S-0742
type: story
nature: feature
title: flai serve starts, restarts, and stops agents, plans again, and finds its reasons by a state's kind and a transition's effect instead of a state's name
status: ready
parent: E-0034
owner: alex
created: 2026-11-12T15:59:30Z
updated: 2026-11-14T23:48:54Z
transitions:
  - to: ready
    at: 2026-11-14T17:47:40Z
    by: alex
tags: [cli]
topics: [cli, orchestration, planning]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_launch.go, flai/internal/serve/agents_start.go, flai/internal/serve/agents_restart.go, flai/internal/serve/agents_activity.go, flai/internal/serve/restart.go, flai/internal/serve/start.go, flai/internal/serve/stop.go, flai/internal/serve/agents_test.go, flai/internal/serve/restart_test.go, flai/internal/serve/start_test.go, flai/internal/serve/stop_test.go, flai/internal/serve/readyempty.go, flai/internal/serve/readyempty_test.go, flai/internal/serve/replan.go, flai/internal/serve/replan_test.go, flai/internal/serve/plan.go, flai/internal/serve/plan_test.go, flai/internal/serve/role.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/orchestrate_reasons_test.go, flai/internal/serve/orchestrate_hold.go, flai/internal/serve/analyze.go, flai/internal/serve/curate.go, flai/internal/serve/curate_schedule.go, flai/internal/serve/formula_schedule.go, flai/internal/serve/usage.go, flai/internal/serve/budget.go, flai/internal/serve/confirm.go, flai/internal/serve/remind.go, flai/internal/serve/remind_test.go, flai/internal/serve/commit.go, flai/internal/serve/claudecheck.go, flai/internal/serve/components.go, flai/internal/serve/checks.go, flai/internal/serve/verifyjob.go, flai/internal/serve/upgrade.go]
after: [S-0737]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 29
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 7
          output: 43
          cache_read: 330386
          cache_write: 51750
          cost: 0.1122
cost_of_delay:
  value: 142.52
  by: system-flow
  at: 2026-11-14T17:44:04Z
forecast:
  duration: 1h20m
  delivery: 2026-11-15T02:33:00Z
  basis: "Its own forecast of 1h20m; 7th in the pull order with an in-progress limit of 3, behind S-0738, S-0741, S-0746, S-0739, S-0743 and S-0744."
  by: flai
  at: 2026-11-14T23:48:54Z
finalized:
  by: orchestrator
  at: 2026-11-14T17:45:53Z
---
# S-0742 flai serve starts, restarts, and stops agents, plans again, and finds its reasons by a state's kind and a transition's effect instead of a state's name

## Goal

flai serve runs a story's agent when the story makes the transition whose effect is the pull, not when it enters `ready`, and keeps its other runs on the kind of the state each item is in. An operator's workflow with a state of its own between `in-progress` and `review` keeps the agent running there, starts nothing for it that the baseline does not, and is starved, reminded, and replanned as the baseline is.

## Acceptance criteria

- [ ] The files in `flai/internal/serve` that start, restart, stop, and measure a story's agent, start the strategic roles, find the orchestrator's reasons, open the starved-board thread, plan again, and remind name no state: each reads a state's kind or a transition's effect from the workflow
- [ ] A story's agent starts on the transition with the pull effect, restarts while its story is in a state of the active-work kind, and is stopped or left to end as today when it leaves that kind
- [ ] The starved-board thread (ADR-0285) and the `ready-empty` event count the stories in the queue state the pull takes from and the WIP limit of the active-work states
- [ ] Tests run the starts, the restarts, the ready-empty count, and the orchestrator's reasons against a workflow with a state renamed and one added between `in-progress` and `review`, with the outcome the kinds give
- [ ] With the baseline, every existing `serve` test passes unchanged

## Tasks

- T-3213 A story's agent starts on the pull transition, restarts while its story is in active work, and stops as today when it leaves that kind
- T-3215 The ready-empty event, the starved-board thread, the replanner, and the planner's starts count by kind and WIP limit
- T-3216 The strategic roles' starts and the orchestrator's reasons and hold read a state's kind and a transition's effect
- T-3217 flai serve's usage, budget, confirmations, reminders, commits, checks, and upgrades read a state's kind

## Notes

- `serve` holds about 100 references to a state's name in 32 files. This story takes the agents' and the strategic roles' files; the kluster, forge, and merge queue files are a story of their own, and `replan.go`'s and `orchestrate_reasons.go`'s use of the role conditions changes with the role condition story.

### Planning

- Touches, all files:
  - Layout: every file of `flai/internal/serve` with a reference to a state's name but the kluster, forge, and merge queue files (S-0743's), with the test files beside the ones that change behaviour.
  - Co-change: `flai/internal/serve/agents.go` (47% with the agents' files in `flai touches suggest`), the file the `agents_*.go` files were split from.
  - `orchestrate_reasons.go` and `replan.go` are touched by S-0748 as well, for the role condition; the trial merge judges that, and the tasks leave the condition to S-0748.
  - Left out from `flai touches suggest`: `cmd/serve_actions.go`, `serve.go`, `harness`, and `hostapi/writes.go` (12% to 22%). None names a state.
- `after`: S-0737. Every change here reads a state's kind, the pull transition, and the WIP limits from the workflow `Repo` carries from S-0737 on, and the tests use its test hook.
- Tasks: T-3213 (agents), T-3215 (ready-empty and replanning), T-3216 (strategic roles), and T-3217 (the rest), all in one layer: no two share a path.
- Forecast 1h20m, `flai forecast`'s 109 s per unit of size for size 43, up from its 1h4m for the first draft's 35 once the test files and `agents.go` were added. The work is wide but mechanical, and the four tasks can run at once.
- Cost of delay: 142.52 USD a week, as `flai cod` gives it: this story's share of E-0034's 1500.00 USD a week, the operator's `time_lost_per_cycle: 10h` (TH-1115) at 150 USD an hour, by its forecast of 1h20m over the 14h2m of the epic's 19 open stories. Kept: the epic's outcome needs every story, so a share by size is the split that fits.
