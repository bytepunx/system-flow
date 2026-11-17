---
id: S-0828
type: story
nature: improvement
title: board.get's undone.read runs one git diff per recorded path of each story's undone record, one after another, on every call
status: ready
owner: alex
created: 2026-11-17T08:17:19Z
updated: 2026-11-17T11:48:16Z
transitions:
  - to: ready
    at: 2026-11-17T11:48:10Z
    by: orchestrator
tags: [flai, git, performance]
topics: [performance]
touches: [flai/internal/storygit/undone.go, flai/internal/storygit/undone_test.go, flai/internal/hostapi/hostapi.go, design/system/server-performance.md, design/issues/I-0349-board-get-s-undone-read-runs-one-git-diff-per-recorded-path-of-each-story-s.md, design/issues/summary.md]
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
      seconds: 41
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 38
          output: 113
          cache_read: 2835276
          cache_write: 63253
          cost: 0.1979
        - model: claude-opus-5-5
          input: 161
          output: 7257
          cache_read: 10734875
          cache_write: 114626
          cost: 3.276
    - kind: orchestrator
      seconds: 69
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 16
          output: 90
          cache_read: 748303
          cache_write: 130458
          cost: 0.2726
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2m
    by: flai
    at: 2026-11-17T08:17:19Z
  value: 5
  by: planner-S-0828
  at: 2026-11-17T08:28:36Z
forecast:
  duration: 1h
  delivery: 2026-11-17T13:39:00Z
  basis: "Its own forecast of 1h; 4th in the pull order with an in-progress limit of 3, behind S-0806, S-0836, S-0811, S-0809 and S-0812."
  by: flai
  at: 2026-11-17T11:48:16Z
finalized:
  by: orchestrator
  at: 2026-11-17T08:33:56Z
---
# S-0828 board.get's undone.read runs one git diff per recorded path of each story's undone record, one after another, on every call

## Goal

This story remediates [I-0349](../../../design/issues/I-0349-board-get-s-undone-read-runs-one-git-diff-per-recorded-path-of-each-story-s.md), "board.get's undone.read runs one git diff per recorded path of each story's undone record, one after another, on every call". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0349 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0349 is closed with `flai issue close I-0349 --reason` saying what fixed it

## Tasks
- T-3525 RecordedUndone runs one git diff for all of a record's paths and one per old base, not one or two per path
- T-3527 Measure board.get's undone.read with the batched diffs, record it in server-performance.md, and close I-0349

## Notes

Cost of delay inputs set by flai from I-0349. time_lost_per_cycle 2m: 2m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-11T08:52:49Z, 6 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from I-0349's instance: `storygit.RecordedUndone` (`flai/internal/storygit/undone.go`) runs `git diff main story/S-nnnn -- <path>` for each recorded path, and `git diff <old_base> main -- <path>` again for each path the branch still removes lines of. Run one `git diff` of main against the branch for all the record's paths, and one of each distinct `old_base` against main for the paths still removing lines, and split each multi-file diff by file. A record of 20 paths with one old base then runs 2 diffs, not 20 to 40, for every caller: `board.get`, the merge queue's eligibility, the accept preview, `flai stream restore`, and the sync. The second candidate the issue names, keeping the answer while the branch's head and main's are unchanged, as S-0678 keeps the trial merges, is left out: it needs a `git rev-parse` per call to know the heads, which costs about what the two batched diffs do, and adds a cache to invalidate. T-3527's measurement says whether it is still needed.

Touches, each a file; no folder touch:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/storygit/undone.go` | layout | `RecordedUndone` and `readHunks`, the per-path diffs |
| `flai/internal/storygit/undone_test.go` | layout | the test that counts the git diffs a record runs |
| `flai/internal/hostapi/hostapi.go` | design | `boardUndone` and `undoneWorkers`, whose comment says a read is a git process per recorded path |
| `design/system/server-performance.md` | design | § The merge queue's git read once (S-0698) records `undone.read`'s 20 diffs; the fix is measured there |
| `design/issues/I-0349-…md`, `design/issues/summary.md` | criteria | `flai issue close I-0349` writes both |

`flai touches suggest S-0828` listed nothing: the story declared no touches. No other open story touches these files.

Figures:

- Forecast 1h, not flai's 2m. flai's figure is 48 s per unit of size times size 2, from 2 criteria and no touches. S-0698, S-0604, and S-0800, board.get and undone-record fixes of the same shape, each took about an hour (3363, 3842, and 3455 s). Delivery is flai's 14:56 UTC, from 17th in the pull order, moved an hour later for the longer duration.
- Cost of delay 5.00 USD a week, as `flai cod S-0828` computes it from the operator's input, 2m lost per 168h cycle at 150 USD an hour. It stands.

No `after`: the story reads nothing another open story writes.
