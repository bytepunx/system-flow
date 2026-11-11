---
id: S-0712
type: story
nature: remediation
title: flai task done's sync reports undoes_main for paths the story never changed when main gains a commit during the sync
status: in-progress
owner: alex
created: 2026-11-10T10:10:59Z
updated: 2026-11-11T21:55:01Z
transitions:
  - to: ready
    at: 2026-11-11T16:15:04Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-11T21:55:01Z
    by: agent-S-0712
tags: [flai, git]
topics: [git]
touches: [flai/internal/storygit/sync.go, flai/internal/storygit/undone.go, flai/internal/storygit/undone_test.go, flai/internal/storygit/sync_test.go, design/system/workflow.md, docs/users/flai.md, design/issues/I-0338-flai-task-done-s-sync-reports-undoes-main-for-paths-the-story-never-changed.md, design/issues/I-0335-flai-stream-sync-reports-a-path-as-undone-by-the-branch-when-main-gains-a.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 711
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-11
      work: 22
  models:
    - model: claude-opus-5-5
      input: 100
      output: 601
      cache_read: 5096452
      cache_write: 382334
      cost: 2.0336
  strategic:
    - kind: planner
      seconds: 127
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 118
          output: 22224
          cache_read: 6439546
          cache_write: 191153
          cost: 3.2214
    - kind: orchestrator
      seconds: 121
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 27
          output: 262
          cache_read: 1819158
          cache_write: 106577
          cost: 0.5311
cost_of_delay:
  inputs:
    time_lost_per_cycle: 6m
    by: orchestrator
    at: 2026-11-11T09:04:58Z
  value: 15
  by: planner-S-0712
  at: 2026-11-11T09:06:14Z
forecast:
  duration: 1h
  delivery: 2026-11-11T20:25:00Z
  basis: "Its own forecast of 1h; 2nd in the pull order with an in-progress limit of 3, behind S-0689, S-0707, S-0714, S-0715 and S-0711."
  by: flai
  at: 2026-11-11T18:17:17Z
finalized:
  by: alex
  at: 2026-11-11T09:03:17Z
---
# S-0712 flai task done's sync reports undoes_main for paths the story never changed when main gains a commit during the sync

## Goal

This story remediates [I-0338](../../../design/issues/I-0338-flai-task-done-s-sync-reports-undoes-main-for-paths-the-story-never-changed.md), "flai task done's sync reports undoes_main for paths the story never changed when main gains a commit during the sync". The issue recommends no solution yet: propose one from its instances before building it.

It also remediates [I-0335](../../../design/issues/I-0335-flai-stream-sync-reports-a-path-as-undone-by-the-branch-when-main-gains-a.md), "flai stream sync reports a path as undone by the branch when main gains a commit to it while the sync runs", which has the same cause. The operator merged S-0709 into this story on TH-1017.

## Acceptance criteria
- [ ] The cause I-0338 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0338 is closed with `flai issue close I-0338 --reason` saying what fixed it
- [ ] I-0335 is closed with `flai issue close I-0335 --reason` saying what fixed it

## Tasks
- T-2963 Take main's head once as a commit in storygit.Rebase and pass that commit to the undone check and the changelog settle
- T-2965 Reproduce I-0338 and I-0335 in real-git sync tests where main gains a commit between the rebase and the checks after it
- T-2967 Say in workflow.md and the users guide that the undone check compares against the commit the rebase took
- T-2970 Close I-0338 and I-0335 with the reason that the sync now compares against the commit its rebase took

## Notes

Cost of delay inputs set by flai from I-0338. time_lost_per_cycle 3m: 3m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T09:13:03Z, 0 days before this story; under one cycle counts as one).

### Planning

Replanned on 2026-11-11 for the operator's answer on TH-1017, which confirmed the orchestrator's recommendation: merge S-0709 (I-0335) into this story, add the criterion that closes I-0335, have T-2963 also pass the pinned commit to `SettleChangelogs`, close both issues in T-2970, and set `time_lost_per_cycle` to 6m. Cancelling S-0709 is the operator's or the orchestrator's move, not this plan's.

Proposed solution. The cause is a race in `storygit.Rebase` (`flai/internal/storygit/sync.go`). It reads the merge base with main by the branch name, rebases onto `main` by name, and then `FindUndone` (`flai/internal/storygit/undone.go`) runs `git rev-parse main` again and diffs against it. A commit main gains after the rebase is missing from the branch, so `diff main HEAD` lists its lines as removed, and `diff oldBase main` lists them as added. In I-0338 this happened with S-0706's creation commit; in I-0335, with two commits main gained during S-0680's sync. The fix: resolve main to one commit before the merge-base read, and use that commit for the merge base, the rebase, the undone check, and `settleChangelogs`, which passes `res.Base`, the name, to `SettleChangelogs` and so reads main's changelog after the rebase too (found by S-0709's plan, TH-1014). `flai accept` rebases through the same function and takes the fix with it. The merge queue's `build.settleChangelogs` (`flai/internal/mergequeue/batch.go`) passes `onto`, the batch head it rebased onto, so it is left alone.

Touches, file by file, no folder:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/storygit/sync.go` | layout, co-change | `Rebase` takes main's commit once and passes it to the check and `settleChangelogs` (T-2963) |
| `flai/internal/storygit/undone.go` | layout | `FindUndone` and `checkUndone` take that commit (T-2963) |
| `flai/internal/storygit/undone_test.go` | layout | the scripted case (T-2963) |
| `flai/internal/storygit/sync_test.go` | co-change (76%) | the real-git reproductions (T-2965) |
| `design/system/workflow.md` | design | its S-0371 paragraph describes the check (T-2967) |
| `docs/users/flai.md` | design | its "Main's change undone." bullet (T-2967) |
| `design/issues/I-0338-…md`, `design/issues/I-0335-…md`, `design/issues/summary.md` | criteria | `flai issue close` writes them (T-2970) |

`flai touches suggest` lists, as before, `design/system/flai-cli.md`, `docs/users/flai-reference.md`, `docs/operators/index.md`, and other documents that change with most commits. None is taken: `res.Base` keeps the branch's name, so the JSON answer, the help text, and the generated reference stay as they are, and accept needs no change of its own.

Forecast: 1h, raised from flai's 27m. flai's figure is 135 s per unit of size over 59 done remediation stories in the medium band, times size 12 (3 criteria, 9 touches). It is raised for three reasons: the reproduction needs a runner that commits on main in the middle of a sync, with a second case for the changelog; `settleChangelogs` takes the pinned commit too; and a story that changes `flai/` runs the integration and smoke tiers at close-out. Delivery moves by the same 33m, from flai's 2026-11-11T22:25:00Z to 22:58:00Z, at 23rd in the pull order.

Cost of delay value: 15 USD a week, as `flai cod` gives it from the input of 6m lost per 168h cycle at 150 USD an hour. The orchestrator set the input to 6m on 2026-11-11, I-0335's 3m added to I-0338's 3m, as the operator confirmed on TH-1017. It stands.

No `after`: the story reads no output of another open story. S-0709 is not named, since it is to be cancelled, not built on.
