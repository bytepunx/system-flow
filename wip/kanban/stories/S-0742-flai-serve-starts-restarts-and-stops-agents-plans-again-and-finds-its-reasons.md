---
id: S-0742
type: story
nature: feature
title: flai serve starts, restarts, and stops agents, plans again, and finds its reasons by a state's kind and a transition's effect instead of a state's name
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T15:59:30Z
updated: 2026-11-13T07:23:59Z
transitions: []
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
      seconds: 11
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 13
          cache_read: 108314
          cache_write: 18020
          cost: 0.0361
draft: true
forecast:
  duration: 1h20m
  delivery: 2026-11-13T16:54:00Z
  basis: "Its own forecast of 1h20m; 16th in the pull order with an in-progress limit of 3, behind S-0659, S-0241, S-0660, S-0661, S-0662, S-0663, S-0730, S-0733, S-0734, S-0735, S-0736, S-0737, S-0738, S-0739, S-0740 and S-0741."
  by: flai
  at: 2026-11-13T07:23:59Z
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
- Cost of delay: not set yet, pending TH-1097, as on S-0736.
