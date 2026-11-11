---
id: S-0678
type: story
nature: improvement
title: agent.status with 13 projects at once is not under 200 ms
status: ready
owner: alex
created: 2026-11-10T08:34:55Z
updated: 2026-11-11T02:28:16Z
transitions:
  - to: ready
    at: 2026-11-10T10:57:56Z
    by: alex
tags: [cli]
topics: [server-side]
touches: [flai/internal/serve/agents_activity.go, flai/internal/serve/activity_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/storygit/trial.go, flai/internal/storygit/trial_test.go, design/system/server-performance.md, design/issues/I-0313-agent-status-with-13-projects-at-once-is-not-under-200-ms.md, design/issues/summary.md]
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
      seconds: 15
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 15
          cache_read: 156045
          cache_write: 9701
          cost: 0.0456
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1s
    by: flai
    at: 2026-11-10T08:34:55Z
  value: 0.04
  by: planner-S-0678
  at: 2026-11-10T10:32:11Z
forecast:
  duration: 1h30m
  delivery: 2026-11-11T04:42:00Z
  basis: "Its own forecast of 1h30m; 2nd in the pull order with an in-progress limit of 3, behind S-0690, S-0692, S-0693 and S-0691."
  by: flai
  at: 2026-11-11T02:28:16Z
finalized:
  by: alex
  at: 2026-11-10T10:57:53Z
---
# S-0678 agent.status with 13 projects at once is not under 200 ms

## Goal

This story remediates [I-0313](../../../design/issues/I-0313-agent-status-with-13-projects-at-once-is-not-under-200-ms.md), "agent.status with 13 projects at once is not under 200 ms". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0313 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0313 is closed with `flai issue close I-0313 --reason` saying what fixed it

## Tasks
- T-3004 agent.status times serve.Activity's parts as phases of the request: the runs, the holds with their trial merges, and elsewhere
- T-3008 The trial judge keeps its verdicts per root and commit pair across calls, and reads the branch heads with one git process a call
- T-3013 Measure 13 agent.status calls at once before and after, record them in server-performance.md, and close I-0313

## Notes

Cost of delay inputs set by flai from I-0313. time_lost_per_cycle 1s: 1s per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-08T02:32:22Z, 2.3 days before this story; under one cycle counts as one).

### Planning

The proposed solution, from I-0313's one instance and the code as it is on main. S-0603 left this repository's `agent.status` at about 1.1 s with 13 at once; alone, its `serve.Activity` took 540 to 600 ms. Of that, the threads (180 to 210 ms) and the conversations (70 to 180 ms) are cause 13's, which S-0602 has since kept parsed: a warm thread read is now about 7 ms. What is left is the holds with their trial merges and `elsewhere`. `Activity` makes a new trial judge on every call (`withTrial`, `storygit.NewTrialJudge`), so each call runs `git --version`, a `git rev-parse` per branch judged, and a `git merge-tree` per overlap, though no branch moved. The fix keeps the trial verdicts across calls by their commits (T-3008), the phases say where the rest goes (T-3004), and the measurement shows whether the target is met (T-3013).

Touches, by where each came from:

| Touch | From |
|-------|------|
| `flai/internal/serve/agents_activity.go` | design: `serve.Activity`, named in I-0313 and in server-performance.md § Cause 12 |
| `flai/internal/serve/activity_test.go` | layout: the tests of `Activity` |
| `flai/cmd/serve_actions.go` | design: `Host.Agent`, where the `agent.*` phases are tracked |
| `flai/cmd/serve_actions_test.go` | co-change: 69% of the commits that changed the three code paths above |
| `flai/internal/storygit/trial.go` | layout: `NewTrialJudge`, which `withTrial` calls on every `Activity` |
| `flai/internal/storygit/trial_test.go` | layout: its tests |
| `design/system/server-performance.md` | design: where causes 12 and 13 are measured |
| `design/issues/I-0313-agent-status-with-13-projects-at-once-is-not-under-200-ms.md` | criteria: the issue closed |
| `design/issues/summary.md` | layout: regenerated when the issue closes |

The story declared no touches. `flai touches suggest` from the three code paths listed `flai/internal/hostapi/writes.go`, `flai/internal/serve/agents.go`, the docs reference, and the configuration; none is changed by this fix, since `serve_actions.go` is a hub that co-changes with every host method. No folder touch. No `after`: S-0602, whose kept threads and conversations this story builds on, is done.

Forecast: `flai forecast` gave 20m, 107 s per unit of size times 11. Raised to 1h30m: S-0603, the first pass at the same call with the same measurement, took 2h35m of agent work, S-0604 1h04m, and S-0602 47m. Most of the time is the measurement program and the runs on a loaded host.

Cost of delay: 0.04 USD a week, as `flai cod` works it out from the input flai set from I-0313, 1s lost per 168h cycle. Not adjusted: the input is the operator's, and the issue's cost is its one occurrence.
