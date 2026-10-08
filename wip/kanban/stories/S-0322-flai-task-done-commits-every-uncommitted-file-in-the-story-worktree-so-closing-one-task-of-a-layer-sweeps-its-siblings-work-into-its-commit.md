---
id: S-0322
type: story
nature: improvement
title: flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit
status: in-progress
owner: alex
created: 2026-10-07T18:59:53Z
updated: 2026-10-08T08:53:52Z
transitions:
  - to: ready
    at: 2026-10-08T04:36:45Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-08T08:13:36Z
    by: agent-S-0322
tags: [flai]
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, docs/users/flai-reference.md, design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, design/conventions/git.md, template/root/design/conventions/git.md, template/CHANGELOG.md, design/issues/I-0104-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md, design/issues/I-0108-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md, design/issues/summary.md, design/adrs/0128-flai-task-done-commits-the-paths-the-closing-task-covers-and-those-no-other.md, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1136
  estimated: true
  turns:
    - day: 2026-10-08
      ceremony: 2
      hand_edits: 2
      work: 32
  models:
    - model: claude-opus-5-5
      input: 258
      output: 1499
      cache_read: 11678857
      cache_write: 530477
      cost: 5.4886
  strategic:
    - kind: planner
      seconds: 369
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 159
          output: 42403
          cache_read: 7915816
          cache_write: 411856
          cost: 5.3736
    - kind: orchestrator
      seconds: 663
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 211
          output: 3347
          cache_read: 87842484
          cache_write: 79081
          cost: 21.662
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-07T18:59:53Z
  value: 12.5
  by: planner-S-0322
  at: 2026-10-08T04:33:13Z
forecast:
  duration: 55m
  delivery: 2026-10-08T09:42:00Z
  basis: "Its own forecast of 55m; 2nd in the pull order with an in-progress limit of 3, behind S-0232, S-0321, S-0339 and S-0342."
  by: flai
  at: 2026-10-08T08:07:59Z
finalized:
  by: orchestrator
  at: 2026-10-08T04:36:42Z
---
# S-0322 flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit

## Goal

This story remediates [I-0108](../../../design/issues/I-0108-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md), "flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0108 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0104 is closed with `flai issue close I-0104 --reason` saying what fixed it
- [x] I-0108 is closed with `flai issue close I-0108 --reason` saying what fixed it

## Tasks
- T-1328 flai task done commits only the closing task's paths, with tests that reproduce I-0104's two tasks and I-0108's three
- T-1329 Document what flai task done commits, and close I-0104 and I-0108
- T-1330 Record what flai task done commits, in an ADR refining ADR-0107
- T-1331 flai task done, the MCP tool task_done, and the host method task.done take a message only when there is something to commit, and print the paths left

## Notes

Cost of delay inputs set by flai from I-0108. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T08:06:26Z, 0.5 days before this story; under one cycle counts as one).

### Planning

The operator cancelled S-0312 as a duplicate of this story and added criterion 2, closing I-0104, to it. This story therefore builds the fix itself: the remedy the operator confirmed on TH-0337 for S-0312. `flai task done` commits the changed paths the closing task's touches cover, plus changed paths no other open task covers. It leaves and lists paths only another open task covers, and needs `-m` only when there is something to commit. The `after` on S-0312 is removed, since a cancelled story is never done. S-0312's cancelled tasks T-1282 to T-1285 are carried over as T-1330, T-1328, T-1331, and T-1329.

The second planning run, after the story was finalized, found one gap: the host method `task.done` in `flai/internal/hostapi/writes.go` refuses a blank message, so the dashboard could not close a task with nothing to commit. T-1331 now changes it too, and T-1329 and T-1330 name its sentences in `docs/users/flai.md` and `design/system/flai-cli.md`.

Touches. The story declared three; all are kept. `flai touches suggest S-0322` lists, from them, the busiest design and issue files at 13% co-change or less. Of those, `design/system/flai-cli.md`, `docs/users/flai.md`, `design/adrs/README.md`, and `docs/users/flai-reference.md` are taken because a task changes them; the rest, such as `design/system/flaiover-dashboard.md` and other issues' files, are not. All touches name files; no folder touch is kept.

| Touch | From | Task |
|-------|------|------|
| `flai/internal/taskdone/taskdone_test.go` | declared | T-1328 |
| `flai/internal/taskdone/taskdone.go` | layout: `run.commit` runs `git add -A`; `run.touches` widens with the commit's paths | T-1328 |
| `flai/cmd/task_done.go`, `task_done_test.go` | layout: `-m` is a required flag; the Long help says `git add -A` | T-1331 |
| `flai/internal/mcpserver/task.go`, `task_test.go` | layout: `taskDoneDescription` says "commit everything in the worktree" | T-1331 |
| `flai/internal/hostapi/writes.go`, `writes_test.go` | layout: `task.done` refuses a blank `message`; its refusal cases list one | T-1331 |
| `docs/users/flai-reference.md` | co-change: generated from the flags | T-1331 |
| `design/adrs/README.md` | co-change and design: the new ADR refining ADR-0107 is indexed there | T-1330 |
| `design/system/flai-cli.md`, `design/system/workflow.md` | co-change and design: both describe ADR-0107's commit step; `flai-cli.md` also the `task.done` host method | T-1330 |
| `docs/users/flai.md` | co-change and design: § Closing a task says `git add -A` and that `task.done` takes a message | T-1329 |
| `design/conventions/git.md`, `template/root/design/conventions/git.md`, `template/CHANGELOG.md` | design: the closing rule and its template copy | T-1329 |
| `design/issues/I-0104-…md`, `design/issues/I-0108-…md` (declared), `design/issues/summary.md` (declared) | criteria 2 and 3: `flai issue close` writes them | T-1329 |

The new ADR's file cannot be named before `flai adr new` numbers it; `flai task done` adds it when T-1330 closes. A folder touch on `design/adrs/` would hold every ready story that adds an ADR. `flai/internal/storygit/commit.go` (`CommitPaths`) is reused, not expected to change; widening adds it if it does. `flaiover/src/lib/server/agent.ts` only lists `task.done` among the required methods and does not change.

S-0326, ahead in the pull order, touches `flai/internal/taskdone/taskdone.go` too, so whichever starts second is held until the other reaches review.

Forecast. `flai forecast` gave 26m (size 20 at 78 s per unit), 7th in the pull order. Raised to 55m: S-0312's identical plan was raised from 24m to 45m because it writes an ADR first and S-0333 (38m) and S-0269 (69m) changed the same files; this story adds I-0108's three-task test and a second issue close; and T-1331 now changes the host method and its tests as well, about 5m more. Delivery 2026-10-08T06:43:00Z is flai's 06:38 for 50m plus the extra 5m; it does not count a hold behind S-0326.

Cost of delay. `flai cod` gives 12.50 USD a week from the 5m input at 150 USD an hour, and it stands. The input was worked from I-0108's one occurrence only. I-0104 now counts 5, and S-0312's 25 USD a week went with its cancellation. The input is the operator's, so it is left as it is; TH-0358 named the figure a 30m input would give.

Tasks, in three layers:

1. T-1330, the ADR and design.
2. T-1328, the commit step and the I-0104 and I-0108 reproduction tests, after T-1330.
3. T-1331, CLI, MCP, and the host method, and T-1329, docs and closing both issues, both after T-1328. They share no path and run together.
