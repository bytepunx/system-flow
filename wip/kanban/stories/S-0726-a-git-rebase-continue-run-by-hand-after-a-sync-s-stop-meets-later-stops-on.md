---
id: S-0726
type: story
nature: improvement
title: A git rebase --continue run by hand after a sync's stop meets later stops on issue files and summary.md that flai's sync would have merged itself
status: review
owner: alex
created: 2026-11-11T21:55:30Z
updated: 2026-11-12T19:36:54Z
transitions:
  - to: ready
    at: 2026-11-12T15:59:56Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-12T18:39:04Z
    by: agent-S-0726
  - to: review
    at: 2026-11-12T19:36:54Z
    by: agent-S-0726
tags: [flai, template]
topics: [sync, git]
touches: [flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/stream.go, flai/cmd/branch.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/task_done.go, flai/cmd/storycommit.go, flai/cmd/stream_restore.go, flai/internal/storygit/commit.go, flai/internal/taskdone/taskdone.go, flai/internal/verify/story.go, flai/internal/harness/harness.go, flai/internal/mcpserver/task.go, scripts/close-out.sh, template/root/scripts/close-out.sh, docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md, design/system/workflow.md, design/conventions/git.md, template/root/design/conventions/git.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/adrs/README.md, design/issues/I-0351-a-git-rebase-continue-run-by-hand-after-a-sync-s-stop-meets-later-stops-on.md, design/issues/summary.md, design/adrs/0294-flai-stream-sync-continues-a-rebase-whose-stop-the-agent-resolved-and-staged.md, flai/cmd/accept.go, flai/cmd/accept_conflict_test.go, flai/cmd/branch_test.go, flai/cmd/task_done_test.go, flai/internal/harness/harness_test.go, flai/internal/mcpserver/task_test.go, flai/internal/taskdone/taskdone_test.go, flai/internal/verify/story_test.go, design/issues/I-0347-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3484
  priced_by: harness
  turns:
    - day: 2026-11-12
      ceremony: 3
      hand_edits: 1
      work: 44
  models:
    - model: claude-opus-5-5
      input: 346
      output: 122298
      cache_read: 23708717
      cache_write: 1140945
      cost: 13.8582
  strategic:
    - kind: planner
      seconds: 175
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 12
          output: 38
          cache_read: 800461
          cache_write: 29412
          cost: 0.0567
        - model: claude-opus-5-5
          input: 164
          output: 17192
          cache_read: 13644317
          cache_write: 222012
          cost: 4.6377
    - kind: orchestrator
      seconds: 139
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 37
          output: 233
          cache_read: 1515597
          cache_write: 199154
          cost: 0.4895
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-11T21:55:30Z
  value: 12.5
  by: planner-S-0726
  at: 2026-11-12T15:55:52Z
forecast:
  duration: 2h
  delivery: 2026-11-12T20:10:00Z
  basis: "Its own forecast of 2h; 3rd in the pull order with an in-progress limit of 3, behind S-0659, S-0683, S-0709, S-0723 and S-0732."
  by: flai
  at: 2026-11-12T16:04:53Z
finalized:
  by: orchestrator
  at: 2026-11-12T15:59:18Z
---
# S-0726 A git rebase --continue run by hand after a sync's stop meets later stops on issue files and summary.md that flai's sync would have merged itself

## Goal

This story remediates [I-0351](../../../design/issues/I-0351-a-git-rebase-continue-run-by-hand-after-a-sync-s-stop-meets-later-stops-on.md), "A git rebase --continue run by hand after a sync's stop meets later stops on issue files and summary.md that flai's sync would have merged itself". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0351 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0351 is closed with `flai issue close I-0351 --reason` saying what fixed it

## Tasks
- T-3154 Record in an ADR and the design that flai stream sync continues a rebase whose stop the agent resolved
- T-3155 storygit's sync continues a rebase whose stop the agent resolved and settles the stops after it, with a real-git test of I-0351
- T-3156 flai stream sync's help and report, and every other git rebase --continue advice after a sync's stop, say to git add and run flai stream sync again
- T-3157 The user docs and the git and work-management conventions, here and in the template, say to git add a resolved stop and run flai stream sync again
- T-3158 Close I-0351 with the reason the sync's continue of a resolved stop fixed it

## Notes

Cost of delay inputs set by flai from I-0351. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-11T18:18:50Z, 0.2 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from I-0351's one instance (S-0689): `flai stream sync`, run again on a worktree whose rebase waits on a stop the agent resolved and staged, continues that rebase itself. It uses the loop its own rebase uses, `ContinueOverResolved` in `flai/internal/storygit/sync.go`. Later stops on issue files (ADR-0126), `design/issues/summary.md` (ADR-0098), the ADR index, and markdown merged by sections (ADR-0175, ADR-0276) are then settled as in a sync that never stopped. It finishes with the steps after a clean rebase: `regenerateAtHead`, the changelogs, the undone check, fold, and fix. It still refuses a rebase whose stop has unmerged paths or conflict markers, naming them. Every message that says to run `git rebase --continue` after a sync's stop says instead to `git add` the paths and run `flai stream sync` again. Today `Rebase` refuses any rebase in progress (`StopRebaseInProgress`), and `flai stream sync`'s help says a rebase in progress is never continued for you. This story reverses that, so it records the decision in an ADR first.

Touches, by where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/storygit/sync.go`, `sync_test.go` | layout: `Rebase`, `ContinueOverResolved`, the stop messages, and their real-git tests |
| `flai/cmd/stream.go`, `stream_sync.go`, `branch.go`, `stream_sync_test.go` | layout: the command's help, its stop report, and its real-git tests |
| `flai/cmd/task_done.go`, `storycommit.go`, `stream_restore.go`, `flai/internal/storygit/commit.go`, `flai/internal/taskdone/taskdone.go`, `flai/internal/verify/story.go`, `flai/internal/harness/harness.go`, `flai/internal/mcpserver/task.go` | layout: each prints the `git rebase --continue` advice |
| `scripts/close-out.sh`, `template/root/scripts/close-out.sh` | layout: the close-out's stop hint |
| `docs/users/flai.md`, `docs/users/flai-reference.md` | design: the user docs, and the reference generated from the help |
| `design/system/flai-cli.md`, `design/system/workflow.md` | design: the sync's bullet and § Branches and collisions |
| `design/conventions/git.md`, `design/conventions/work-management.md` and their `template/root/` copies | design: the conventions that give the advice |
| `template/CHANGELOG.md` | co-change: 6% of the commits that change these paths; the template's conventions change |
| `design/adrs/0291-…md`, `design/adrs/README.md` | design: the ADR the fix needs; the number is predicted, so the task that writes it narrows it |
| I-0351 and `design/issues/summary.md` | declared by criterion 2 |

The co-change list's other paths, such as `docs/operators/index.md` and `design/system/flaiover-dashboard.md`, are hub files with no part in this change, so they are left out. No folder touch is kept. The story names no `after`. S-0709 also changes `flai stream sync`, but this story reads nothing it writes, so a trial merge judges the two.

Forecast: `flai forecast` gave 54m, at 107 s per unit of size over 175 large-band improvement stories, times size 30. It is raised to 2h. Most of the 29 touches are one-line changes to advice text. The depth is in the rebase loop: recovering the branch's old base and head from the rebase's state for the undone check, and a test across two stops. The size measures neither.

Cost of delay: `flai cod` gives 12.50 USD a week from flai's input, 5m lost per 168h cycle at 150 USD an hour. It stands as given.
