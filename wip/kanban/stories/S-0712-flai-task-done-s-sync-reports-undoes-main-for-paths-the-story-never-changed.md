---
id: S-0712
type: story
nature: remediation
title: flai task done's sync reports undoes_main for paths the story never changed when main gains a commit during the sync
status: backlog
owner: alex
created: 2026-11-10T10:10:59Z
updated: 2026-11-11T08:47:48Z
transitions: []
tags: [flai, git]
topics: [git]
touches: [flai/internal/storygit/sync.go, flai/internal/storygit/undone.go, flai/internal/storygit/undone_test.go, flai/internal/storygit/sync_test.go, design/system/workflow.md, docs/users/flai.md, design/issues/I-0338-flai-task-done-s-sync-reports-undoes-main-for-paths-the-story-never-changed.md, design/issues/summary.md]
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
      seconds: 47
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 8
          output: 82
          cache_read: 1167497
          cache_write: 34082
          cost: 0.3284
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 3m
    by: flai
    at: 2026-11-10T10:10:59Z
  value: 7.5
  by: planner-S-0712
  at: 2026-11-10T10:28:50Z
forecast:
  duration: 45m
  delivery: 2026-11-11T21:58:00Z
  basis: "Its own forecast of 45m; 21st in the pull order with an in-progress limit of 3, behind S-0693, S-0698, S-0702, S-0694, S-0705, S-0681, S-0718, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663, S-0683, S-0689, S-0699, S-0701, S-0703, S-0706, S-0707, S-0709, S-0710 and S-0711."
  by: flai
  at: 2026-11-11T08:47:48Z
---
# S-0712 flai task done's sync reports undoes_main for paths the story never changed when main gains a commit during the sync

## Goal

This story remediates [I-0338](../../../design/issues/I-0338-flai-task-done-s-sync-reports-undoes-main-for-paths-the-story-never-changed.md), "flai task done's sync reports undoes_main for paths the story never changed when main gains a commit during the sync". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0338 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0338 is closed with `flai issue close I-0338 --reason` saying what fixed it

## Tasks
- T-2963 Take main's head once as a commit in storygit.Rebase and compare the undone check against that commit
- T-2965 Reproduce I-0338 in a real-git sync test where main gains a commit between the rebase and the undone check
- T-2967 Say in workflow.md and the users guide that the undone check compares against the commit the rebase took
- T-2970 Close I-0338 with the reason that the sync now compares against the commit its rebase took

## Notes

Cost of delay inputs set by flai from I-0338. time_lost_per_cycle 3m: 3m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T09:13:03Z, 0 days before this story; under one cycle counts as one).

### Planning

Proposed solution. The cause is a race in `storygit.Rebase` (`flai/internal/storygit/sync.go`). It reads the merge base with main by the branch name, rebases onto `main` by name, and then `FindUndone` (`flai/internal/storygit/undone.go`) runs `git rev-parse main` again and diffs against it. A commit main gains after the rebase is missing from the branch, so `diff main HEAD` lists its lines as removed, and `diff oldBase main` lists them as added. In I-0338 this happened with S-0706's creation commit. The fix: resolve main to one commit before the merge-base read, and use that commit for the merge base, the rebase, and the check. `flai accept` rebases through the same function and takes the fix with it. I-0335, which S-0709 remediates, has the same cause; see the plan's thread.

Touches, file by file, no folder:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/storygit/sync.go` | layout, co-change | `Rebase` takes main's commit once (T-2963) |
| `flai/internal/storygit/undone.go` | layout | `FindUndone` and `checkUndone` take that commit (T-2963) |
| `flai/internal/storygit/undone_test.go` | layout | the scripted case (T-2963) |
| `flai/internal/storygit/sync_test.go` | co-change (76%) | the real-git reproduction (T-2965) |
| `design/system/workflow.md` | design | its S-0371 paragraph describes the check (T-2967) |
| `docs/users/flai.md` | design | its "Main's change undone." bullet (T-2967) |
| `design/issues/I-0338-…md`, `design/issues/summary.md` | criteria | `flai issue close` writes both (T-2970) |

`flai touches suggest` from `sync.go` and `undone.go` also lists `flai/cmd/stream.go`, `flai/cmd/stream_sync.go`, their tests, `flai/cmd/branch.go`, `flai/cmd/accept.go`, and `docs/users/flai-reference.md`. None is taken: `res.Base` keeps the branch's name, so the JSON answer, the help text, and the generated reference stay as they are, and accept needs no change of its own.

Forecast: 45m, adjusted from flai's 23m. flai's figure is 134 s per unit of size over 54 done remediation stories in the medium band, times size 10 (2 criteria, 8 touches). It is raised because the reproduction needs a runner that commits on main in the middle of a sync. That test runs only without `-short`, and a story that changes `flai/` runs the integration and smoke tiers at close-out. Delivery moves by the same 22m, from flai's 2026-11-10T16:05:00Z to 16:27:00Z, at 43rd in the pull order.

Cost of delay value: 7.50 USD a week, as `flai cod` gives it from flai's input of 3m lost per 168h cycle at 150 USD an hour. It stands: the issue has one instance. The value would double if I-0335's time lost were counted here. That input is the operator's, so it is not set.

No `after`: the story reads no output of another open story.
