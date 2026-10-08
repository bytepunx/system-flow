---
id: S-0312
type: story
nature: improvement
title: flai task done commits everything in the worktree, so two tasks of one layer cannot be closed apart
status: backlog
owner: alex
created: 2026-10-07T14:26:01Z
updated: 2026-10-08T00:18:41Z
transitions: []
tags: [flai]
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, docs/users/flai-reference.md, design/conventions/git.md, template/root/design/conventions/git.md, template/CHANGELOG.md, design/issues/I-0104-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md, design/issues/summary.md]
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
      seconds: 265
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 42
          output: 660
          cache_read: 7929248
          cache_write: 14621
          cost: 1.9653
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-07T14:26:01Z
  value: 25
  by: planner-S-0312
  at: 2026-10-07T23:35:29Z
forecast:
  duration: 45m
  delivery: 2026-10-08T07:15:00Z
  basis: "Its own forecast of 45m; 22nd in the pull order with an in-progress limit of 3, behind S-0232, S-0316, S-0317, S-0318, S-0320, S-0319, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0287, S-0288, S-0289, S-0290, S-0291, S-0297, S-0304, S-0305, S-0306 and S-0309."
  by: flai
  at: 2026-10-08T00:18:41Z
---
# S-0312 flai task done commits everything in the worktree, so two tasks of one layer cannot be closed apart

## Goal

This story remediates [I-0104](../../../design/issues/I-0104-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md), "flai task done commits everything in the worktree, so two tasks of one layer cannot be closed apart". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0104 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0104 is closed with `flai issue close I-0104 --reason` saying what fixed it

## Tasks
- T-1282 Propose what flai task done commits, in an ADR refining ADR-0107
- T-1283 flai task done commits only the closing task's paths, with a test that reproduces I-0104
- T-1284 flai task done and the MCP tool task_done take -m only when there is something to commit, and print the paths left
- T-1285 Document what flai task done commits and close I-0104

## Notes

Cost of delay inputs set by flai from I-0104. time_lost_per_cycle 10m: 5m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-07T07:17:53Z, 0.3 days before this story; under one cycle counts as one).

### Planning

Touches. The story declared none, and `flai touches suggest S-0312` needs a starting path, so every touch comes from the issue's instances, ADR-0107, and the code layout. All of them name files, and no folder touch is kept.

| Touch | From | Task |
|-------|------|------|
| `flai/internal/taskdone/taskdone.go`, `taskdone_test.go` | layout: `run.commit` runs `git add -A`, `run.touches` widens with the commit's paths | T-1283 |
| `flai/cmd/task_done.go`, `task_done_test.go` | layout: `-m` is a required flag; Long help says `git add -A` | T-1284 |
| `flai/internal/mcpserver/task.go`, `task_test.go` | layout: `taskDoneDescription` says "commit everything in the worktree" | T-1284 |
| `docs/users/flai-reference.md` | co-change: generated from the flags, changed with `task_done.go` in S-0333 | T-1284 |
| `design/adrs/README.md` | design: the new ADR refining ADR-0107 is indexed there | T-1282 |
| `design/system/flai-cli.md`, `design/system/workflow.md` | design: both describe ADR-0107's commit step | T-1282 |
| `docs/users/flai.md` | design: § Closing a task says `git add -A` | T-1285 |
| `design/conventions/git.md`, `template/root/design/conventions/git.md`, `template/CHANGELOG.md` | design: the closing rule and its template copy, as S-0333 changed them | T-1285 |
| `design/issues/I-0104-…md`, `design/issues/summary.md` | criterion 2: `flai issue close` writes them | T-1285 |

The new ADR's file cannot be named before `flai adr new` numbers it. `flai task done` adds it to T-1282's and the story's touches when T-1282 closes. A folder touch on `design/adrs/` would instead hold every ready story that adds an ADR. `flai/internal/storygit/commit.go` (`CommitPaths`) is reused, not expected to change; widening adds it if it does. The tag `flai` was added because flai check requires a component tag for touches in `flai/` and `template/`.

Forecast. `flai forecast` gave 24m from size 18. Raised to 45m: S-0333 (38m) and S-0269 (69m) changed the same files, and this story writes an ADR first. Delivery 2026-10-08T06:21:00Z is flai's 06:00 plus the extra 21m.

Cost of delay. `flai cod` gives 25.00 USD a week from the 10m input at 150 USD an hour, and it stands. I-0104 now counts 4 occurrences, not the 2 the input was worked from. The input is the operator's, so it is left as it is and raised on the plan's thread.

Tasks, in three layers:

1. T-1282, the ADR and design.
2. T-1283, the commit step and the reproduction test, after T-1282.
3. T-1284, CLI and MCP, and T-1285, docs and closing I-0104, both after T-1283. They share no path and run together.
