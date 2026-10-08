---
id: S-0322
type: story
nature: improvement
title: flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit
status: backlog
owner: alex
created: 2026-10-07T18:59:53Z
updated: 2026-10-08T04:34:12Z
transitions: []
tags: [flai]
touches: [flai/internal/taskdone/taskdone_test.go, design/issues/I-0108-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md, design/issues/summary.md]
after: [S-0312]
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
    - kind: planner
      seconds: 94
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 60
          output: 9856
          cache_read: 2951519
          cache_write: 95061
          cost: 1.8226
    - kind: orchestrator
      seconds: 273
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 95
          output: 1575
          cache_read: 38865911
          cache_write: 31111
          cost: 9.5834
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-07T18:59:53Z
  value: 12.5
  by: planner-S-0322
  at: 2026-10-08T04:33:13Z
forecast:
  duration: 20m
  delivery: 2026-10-08T11:57:00Z
  basis: "Its own forecast of 20m; 25th in the pull order with an in-progress limit of 3, behind S-0232, S-0316, S-0324, S-0318, S-0320, S-0319, S-0309, S-0326, S-0287, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0288, S-0289, S-0290, S-0291, S-0297, S-0304, S-0305, S-0306, S-0313 and S-0321."
  by: flai
  at: 2026-10-08T04:33:38Z
---
# S-0322 flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit

## Goal

This story remediates [I-0108](../../../design/issues/I-0108-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md), "flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0108 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0104 is closed with flai issue close I-0104 --reason saying what fixed it
- [ ] I-0108 is closed with `flai issue close I-0108 --reason` saying what fixed it

## Tasks
- T-1328 Test that three tasks of one layer closed in turn with flai task done each get a commit of their own, as I-0108 reports
- T-1329 Close I-0108, naming S-0312's commit step and T-1328's test

## Notes

Cost of delay inputs set by flai from I-0108. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T08:06:26Z, 0.5 days before this story; under one cycle counts as one).

### Planning

I-0108 is the fault I-0104 records, and S-0312 (ready, planned) remediates I-0104: its T-1283 replaces `git add -A` in `flai/internal/taskdone/taskdone.go` with a commit of the closing task's paths. This plan therefore builds no fix of its own. It waits for S-0312 (`after`), tests I-0108's own instance (three tasks of one layer, from S-0274) against S-0312's commit step, and closes I-0108. TH-0306 asks the operator whether to cancel S-0322 as a duplicate instead; the plan stands if it is kept.

Touches. The story declared none. `flai touches suggest S-0322` needs a starting path; from `taskdone_test.go` and `summary.md` it lists only the busiest design and issue files (at most 12% co-change), none of which this story changes, so none is taken. All touches name files; no folder touch is kept.

| Touch | From | Task |
|-------|------|------|
| `flai/internal/taskdone/taskdone_test.go` | layout: the tests of the commit step S-0312 changes | T-1328 |
| `design/issues/I-0108-…md`, `design/issues/summary.md` | criterion 2: `flai issue close` writes them | T-1329 |

`flai/internal/taskdone/taskdone.go` is not claimed: S-0312 changes it, and this story only reads it. If T-1328 finds a fault there, widening adds it.

Forecast. `flai forecast` gave 9m (size 5 at 97 s per unit). Raised to 20m: the agent must read S-0312's new commit step before writing the test, and a change under `flai/` selects the Go tiers at close-out. Delivery 2026-10-08T12:13:00Z is flai's 12:02 plus the extra 11m.

Cost of delay. `flai cod` gives 12.50 USD a week from the 5m input at 150 USD an hour, and it stands. The input is the operator's and is left as it is.

Tasks, in two layers:

1. T-1328, the three-task test, or the name of an S-0312 test that already covers it.
2. T-1329, closing I-0108, after T-1328.
