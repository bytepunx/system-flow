---
id: S-0709
type: story
nature: improvement
title: flai stream sync reports a path as undone by the branch when main gains a commit to it while the sync runs
status: backlog
owner: alex
created: 2026-11-10T09:55:57Z
updated: 2026-11-12T03:38:48Z
transitions: []
tags: [cli]
topics: [git]
touches: [flai/internal/storygit/sync.go, flai/internal/storygit/undone.go, flai/internal/storygit/sync_test.go, flai/internal/storygit/undone_test.go, docs/users/flai.md, design/issues/I-0335-flai-stream-sync-reports-a-path-as-undone-by-the-branch-when-main-gains-a.md, design/issues/summary.md]
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
      seconds: 202
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 290
          output: 1074
          cache_read: 10507352
          cache_write: 964959
          cost: 0.7722
        - model: claude-opus-5-5
          input: 749
          output: 28856
          cache_read: 34098653
          cache_write: 3123026
          cost: 11.0893
    - kind: orchestrator
      seconds: 14
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 37
          cache_read: 165562
          cache_write: 2974
          cost: 0.0473
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 3m
    by: flai
    at: 2026-11-10T09:55:57Z
  value: 7.5
  by: planner-S-0709
  at: 2026-11-10T10:27:00Z
forecast:
  duration: 1h30m
  delivery: 2026-11-12T12:32:00Z
  basis: "Its own forecast of 1h30m; 8th in the pull order with an in-progress limit of 3, behind S-0241, S-0659, S-0660, S-0661, S-0662, S-0663 and S-0683."
  by: flai
  at: 2026-11-12T03:38:48Z
---
# S-0709 flai stream sync reports a path as undone by the branch when main gains a commit to it while the sync runs

## Goal

This story remediates [I-0335](../../../design/issues/I-0335-flai-stream-sync-reports-a-path-as-undone-by-the-branch-when-main-gains-a.md), "flai stream sync reports a path as undone by the branch when main gains a commit to it while the sync runs". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0335 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0335 is closed with `flai issue close I-0335 --reason` saying what fixed it

## Tasks
- T-2959 storygit.Rebase compares the rebased branch with main's commit as the rebase took it, not with main's name read again after it
- T-2960 Say in docs/users/flai.md that sync compares with main as its rebase took it, and close I-0335

## Notes

Cost of delay inputs set by flai from I-0335. time_lost_per_cycle 3m: 3m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T09:03:43Z, 0 days before this story; under one cycle counts as one).

### Planning

Proposed solution. `storygit.Rebase` in `flai/internal/storygit/sync.go` names main by its branch name throughout: `git rebase <main>` resolves it when the rebase starts, and `checkUndone` then calls `FindUndone` in `flai/internal/storygit/undone.go` with the same name, which `git rev-parse` and the diffs resolve again after the rebase. A commit that reaches main in between, such as an ADR's acceptance, is then in `oldBase..main` and missing from `HEAD`, so `main..HEAD` shows the branch removing its lines and the path is reported as undone, though the branch never changed it. The fix reads main's commit once with `git rev-parse` before the rebase, rebases onto that commit, and gives that commit, not the name, to `FindUndone` and to `SettleChangelogs`, which reads main's changelog after the rebase in the same way. The three-dot comparisons (`fix`, `checkSync`) go through the merge base and are not affected. `res.Base` stays the branch name, so the sync's output does not change. I-0338 (S-0712), seen in `flai task done`'s sync on S-0682, has the same cause; the plan's thread proposes merging S-0712 into this story.

Touches:

- `flai/internal/storygit/sync.go`: design, the rebase and `checkUndone`'s call, where the fix goes.
- `flai/internal/storygit/undone.go`: declared by the issue's cause, `FindUndone`'s comment and its `main` argument, now a commit; co-change with `sync.go`.
- `flai/internal/storygit/sync_test.go`: co-change (76%), the real-git test that reproduces the race.
- `flai/internal/storygit/undone_test.go`: layout, `FindUndone`'s scripted tests if its argument's meaning changes.
- `docs/users/flai.md`: design, the "Main's change undone" paragraph under Story branches, which says what `old_base` and main are compared at.
- `design/issues/I-0335-…md` and `design/issues/summary.md`: the second criterion, `flai issue close I-0335`.

Left out: `flai/cmd/stream_sync.go`, `flai/cmd/stream.go`, `flai/cmd/branch.go`, and their tests, which co-change at 35 to 59% but print or call `Rebase` and `Sync` without comparing with main themselves, so they get the fix with no change. No folder touch.

Figures:

- Forecast: flai's 2m rests on size 2, two criteria and no touches, and is too small for a change that needs a real-git test with a commit landing between the rebase and the check and a close-out that runs the integration and smoke tiers for `flai/`. Raised to 1h30m; delivery moved by the same 1h28m, from 2026-11-10T14:55Z to 16:30Z.
- Cost of delay: 7.50 USD a week, as `flai cod` gives it from the 3m a cycle flai set from I-0335. It stands. Merged with S-0712, I-0338's 3m a cycle would double it, which is the operator's input to give.

No `after`: the story reads no other open story's output.
