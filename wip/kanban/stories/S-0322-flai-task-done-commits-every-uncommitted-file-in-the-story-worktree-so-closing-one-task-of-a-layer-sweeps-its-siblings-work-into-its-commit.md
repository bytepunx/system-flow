---
id: S-0322
type: story
nature: improvement
title: flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit
status: ready
owner: alex
created: 2026-10-07T18:59:53Z
updated: 2026-10-08T04:37:18Z
transitions:
  - to: ready
    at: 2026-10-08T04:36:45Z
    by: orchestrator
tags: [flai]
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, docs/users/flai-reference.md, design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, design/conventions/git.md, template/root/design/conventions/git.md, template/CHANGELOG.md, design/issues/I-0104-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md, design/issues/I-0108-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md, design/issues/summary.md]
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
      seconds: 229
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 119
          output: 26322
          cache_read: 6111238
          cache_write: 313004
          cost: 3.9001
    - kind: orchestrator
      seconds: 368
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 127
          output: 1990
          cache_read: 53151638
          cache_write: 60382
          cost: 13.1103
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-07T18:59:53Z
  value: 12.5
  by: planner-S-0322
  at: 2026-10-08T04:33:13Z
forecast:
  duration: 50m
  delivery: 2026-10-08T12:42:00Z
  basis: "flai forecast gave 26m (size 20 at 78 s per unit); raised to 50m because S-0312's same plan was raised to 45m for an ADR first and files S-0333 (38m) and S-0269 (69m) changed, and this story adds I-0108's three-task test and a second issue close."
  by: planner-S-0322
  at: 2026-10-08T04:36:08Z
finalized:
  by: orchestrator
  at: 2026-10-08T04:36:42Z
---
# S-0322 flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit

## Goal

This story remediates [I-0108](../../../design/issues/I-0108-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md), "flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0108 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0104 is closed with `flai issue close I-0104 --reason` saying what fixed it
- [ ] I-0108 is closed with `flai issue close I-0108 --reason` saying what fixed it

## Tasks
- T-1328 flai task done commits only the closing task's paths, with tests that reproduce I-0104's two tasks and I-0108's three
- T-1329 Document what flai task done commits, and close I-0104 and I-0108
- T-1330 Record what flai task done commits, in an ADR refining ADR-0107
- T-1331 flai task done and the MCP tool task_done take -m only when there is something to commit, and print the paths left

## Notes

Cost of delay inputs set by flai from I-0108. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T08:06:26Z, 0.5 days before this story; under one cycle counts as one).

### Planning

The operator cancelled S-0312 as a duplicate of this story and added criterion 2, closing I-0104, to it. This story therefore builds the fix itself: the remedy the operator confirmed on TH-0337 for S-0312. `flai task done` commits the changed paths the closing task's touches cover, plus changed paths no other open task covers. It leaves and lists paths only another open task covers, and needs `-m` only when there is something to commit. The `after` on S-0312 is removed, since a cancelled story is never done. S-0312's cancelled tasks T-1282 to T-1285 are carried over as T-1330, T-1328, T-1331, and T-1329.

Touches. The story declared three; all are kept. `flai touches suggest S-0322` lists, from them, the busiest design and issue files at 12% co-change or less. Of those, `design/system/flai-cli.md`, `docs/users/flai.md`, `design/adrs/README.md`, and `docs/users/flai-reference.md` are taken because a task changes them; the rest are other issues' files and are not. All touches name files; no folder touch is kept.

| Touch | From | Task |
|-------|------|------|
| `flai/internal/taskdone/taskdone_test.go` | declared | T-1328 |
| `flai/internal/taskdone/taskdone.go` | layout: `run.commit` runs `git add -A`; `run.touches` widens with the commit's paths | T-1328 |
| `flai/cmd/task_done.go`, `task_done_test.go` | layout: `-m` is a required flag; the Long help says `git add -A` | T-1331 |
| `flai/internal/mcpserver/task.go`, `task_test.go` | layout: `taskDoneDescription` says "commit everything in the worktree" | T-1331 |
| `docs/users/flai-reference.md` | co-change: generated from the flags | T-1331 |
| `design/adrs/README.md` | co-change and design: the new ADR refining ADR-0107 is indexed there | T-1330 |
| `design/system/flai-cli.md`, `design/system/workflow.md` | co-change and design: both describe ADR-0107's commit step | T-1330 |
| `docs/users/flai.md` | co-change and design: § Closing a task says `git add -A` | T-1329 |
| `design/conventions/git.md`, `template/root/design/conventions/git.md`, `template/CHANGELOG.md` | design: the closing rule and its template copy | T-1329 |
| `design/issues/I-0104-…md`, `design/issues/I-0108-…md` (declared), `design/issues/summary.md` (declared) | criteria 2 and 3: `flai issue close` writes them | T-1329 |

The new ADR's file cannot be named before `flai adr new` numbers it; `flai task done` adds it when T-1330 closes. A folder touch on `design/adrs/` would hold every ready story that adds an ADR. `flai/internal/storygit/commit.go` (`CommitPaths`) is reused, not expected to change; widening adds it if it does.

Forecast. `flai forecast` gave 26m (size 20 at 78 s per unit). Raised to 50m: S-0312's identical plan was raised from 24m to 45m because it writes an ADR first and S-0333 (38m) and S-0269 (69m) changed the same files. This story adds I-0108's three-task test and a second issue close. Delivery 2026-10-08T12:42:00Z is flai's 12:18 plus the extra 24m.

Cost of delay. `flai cod` gives 12.50 USD a week from the 5m input at 150 USD an hour, and it stands. The input was worked from I-0108's one occurrence only. I-0104 now counts 5, and S-0312's 25 USD a week went with its cancellation. The input is the operator's, so it is left as it is and raised on the plan's thread.

Tasks, in three layers:

1. T-1330, the ADR and design.
2. T-1328, the commit step and the I-0104 and I-0108 reproduction tests, after T-1330.
3. T-1331, CLI and MCP, and T-1329, docs and closing both issues, both after T-1328. They share no path and run together.
