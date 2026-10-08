---
id: S-0326
type: story
nature: remediation
title: flai verify --record-issues opens an issue on the story branch that another branch opened under the same title meanwhile
status: in-progress
owner: alex
created: 2026-10-07T18:59:58Z
updated: 2026-10-08T07:08:09Z
transitions:
  - to: ready
    at: 2026-10-08T04:21:11Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-08T06:23:31Z
    by: agent-S-0326
tags: [cli]
topics: [cli, git, continuous-improvement]
touches: [design/adrs/README.md, flai/internal/issues/fold.go, flai/internal/issues/fold_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/stream.go, flai/cmd/branch.go, flai/internal/taskdone/taskdone.go, flai/cmd/stream_sync_test.go, design/system/continuous-improvement.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0112-flai-verify-record-issues-opens-an-issue-on-the-story-branch-that-another-branch-opened-under-the-same-title-meanwhile.md, design/issues/summary.md, design/adrs/0126-a-rebase-of-a-story-branch-merges-an-issue-file-both-sides-changed-by-its.md, flai/internal/issues/merge.go, flai/internal/issues/merge_test.go, flai/internal/issues/generated.go, flai/cmd/stream_sync.go, docs/users/flai-reference.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2789
  estimated: true
  turns:
    - day: 2026-10-08
      ceremony: 4
      hand_edits: 1
      work: 44
  models:
    - model: claude-opus-5-5
      input: 296
      output: 1872
      cache_read: 16147709
      cache_write: 888512
      cost: 7.6016
  strategic:
    - kind: orchestrator
      seconds: 249
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 62
          output: 1032
          cache_read: 21618198
          cache_write: 40327
          cost: 5.3363
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-07T18:59:58Z
  value: 25
  by: planner-S-0326
  at: 2026-10-08T00:34:49Z
forecast:
  duration: 45m
  delivery: 2026-10-08T07:05:00Z
  basis: "Its own forecast of 45m; 1st in the pull order with an in-progress limit of 3, behind S-0232, S-0323 and S-0336."
  by: flai
  at: 2026-10-08T06:13:49Z
finalized:
  by: orchestrator
  at: 2026-10-08T04:21:07Z
---
# S-0326 flai verify --record-issues opens an issue on the story branch that another branch opened under the same title meanwhile

## Goal

This story remediates [I-0112](../../../design/issues/I-0112-flai-verify-record-issues-opens-an-issue-on-the-story-branch-that-another-branch-opened-under-the-same-title-meanwhile.md), "flai verify --record-issues opens an issue on the story branch that another branch opened under the same title meanwhile". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0112 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0112 is closed with `flai issue close I-0112 --reason` saying what fixed it

## Tasks
- T-1319 An ADR refining ADR-0085 and ADR-0098 records the remedy for I-0112, proposed from its instance
- T-1320 issues.Fold folds each open issue a branch added into the main branch's open issue of the same title
- T-1321 A sync, a task's close, and an acceptance fold the branch's duplicate issue after the rebase, with a test that reproduces I-0112
- T-1322 The design and the users' guide say that a sync folds a branch's issue into the main branch's issue of the same title
- T-1323 I-0112 is closed with flai issue close, saying that a sync folds a branch's duplicate issue
- T-1332 issues.Merge merges an issue file both sides of a stopped rebase changed, keeping the instances of both

## Notes

Cost of delay inputs set by flai from I-0112. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T09:09:31Z, 0.4 days before this story; under one cycle counts as one).

### Planning

The cause: `issues.RecordOnce` (`flai/internal/issues/record.go`) matches titles in the story's own checkout only, so a close-out cannot see an issue another story branch opened under the same title. Once that branch is accepted and this one syncs, both files stand on the branch, and the next close-out bumps the lower ID (`issues.FindOpenByTitle`), leaving the other a duplicate. The proposed remedy, for T-1319's ADR to weigh: fold the branch's issue into the main branch's issue of the same title once a rebase completes, in the three places that rebase a story branch.

Touches, by where each came from:

| Touch | Source |
|-------|--------|
| `design/adrs/`, `design/adrs/README.md` | design: the remedy is recorded as an ADR refining ADR-0085 and ADR-0098, as S-0315's was |
| `flai/internal/issues/fold.go`, `fold_test.go` | layout: a new file beside `record.go` and `generated.go` |
| `flai/internal/storygit/sync.go` | co-change (named as a start); holds `SyncOptions` and the rebase |
| `flai/internal/storygit/sync_test.go` | co-change: 31% with the start paths |
| `flai/cmd/stream.go`, `flai/cmd/branch.go` | co-change (50%, 25%) and layout: they pass `issues.Generated` to the rebase |
| `flai/internal/taskdone/taskdone.go` | layout: the third caller passing `issues.Generated` |
| `flai/cmd/stream_sync_test.go` | co-change (named as a start); holds the summary sync tests whose helpers the reproduction reuses |
| `design/system/continuous-improvement.md`, `design/system/flai-cli.md`, `docs/users/flai.md` | design: where the summary's regeneration on a sync is documented today |
| `design/issues/I-0112-…md`, `design/issues/summary.md` | criteria: criterion 2 closes I-0112 |

`flai touches suggest` also listed `flai/cmd/stream_sync.go` (62%), `docs/users/flai-reference.md`, `docs/operators/settings.md`, `flai/cmd/check.go`, and the issues package's other tests. They are left out: the remedy adds no flag or setting and does not change what the close-out records.

One folder touch is kept: `design/adrs/`, on T-1319, because `flai adr new` names the ADR's file only once it numbers it. Narrow it to that file when T-1319 is done.

Forecast: flai gave 20m (73 s per unit of size over 17 done large remediation stories, size 16). Raised to 45m: five tasks with an ADR, a new function, a hook through three rebase callers, and a two-branch reproduction test, where S-0315, two tasks in the same area, took 17m. Delivery is flai's 2026-10-08T08:19:00Z plus the 25m added.

Cost of delay: 25 USD a week, as `flai cod` gives it from the inputs flai set from I-0112 (10m lost per 168h cycle at 150 USD an hour). It stands: the inputs are the operator's, and one occurrence in 0.4 days gives no reason to adjust.
