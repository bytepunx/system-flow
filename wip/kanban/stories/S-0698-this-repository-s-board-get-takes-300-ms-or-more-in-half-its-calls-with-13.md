---
id: S-0698
type: story
nature: improvement
title: This repository's board.get takes 300 ms or more in half its calls with 13 boards at once, from the merge queue's git while stories are in review
status: backlog
owner: alex
created: 2026-11-10T08:54:37Z
updated: 2026-11-11T01:41:10Z
transitions: []
tags: [cli]
topics: [server-side, back-end]
touches: [flai/internal/storygit/storygit.go, flai/internal/storygit/storygit_test.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, design/system/server-performance.md, design/system/flai-cli.md, design/issues/I-0333-this-repository-s-board-get-takes-300-ms-or-more-in-half-its-calls-with-13.md, design/issues/summary.md]
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
      seconds: 6
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 2
          cache_read: 102139
          cache_write: 491
          cost: 0.0281
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-10T08:54:37Z
  value: 12.5
  by: planner-S-0698
  at: 2026-11-10T10:30:42Z
forecast:
  duration: 1h
  delivery: 2026-11-11T14:33:00Z
  basis: "Its own forecast of 1h; 23rd in the pull order with an in-progress limit of 3, behind S-0676, S-0678, S-0674, S-0704, S-0681, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663, S-0683, S-0686, S-0689, S-0690, S-0691, S-0692, S-0693, S-0694, S-0695, S-0696 and S-0697."
  by: flai
  at: 2026-11-11T01:41:10Z
finalized:
  by: orchestrator
  at: 2026-11-10T11:04:47Z
---
# S-0698 This repository's board.get takes 300 ms or more in half its calls with 13 boards at once, from the merge queue's git while stories are in review

## Goal

This story remediates [I-0333](../../../design/issues/I-0333-this-repository-s-board-get-takes-300-ms-or-more-in-half-its-calls-with-13.md), "This repository's board.get takes 300 ms or more in half its calls with 13 boards at once, from the merge queue's git while stories are in review". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0333 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0333 is closed with `flai issue close I-0333 --reason` saying what fixed it

## Tasks
- T-3000 storygit reads every local branch's head and the checked-out branch in one git for-each-ref
- T-3003 merge.queue reads the heads once and runs git status only when a story in review passed verify, at the same time as the heads
- T-3007 Measure this repository's board.get with stories in review before and after, record it in the design, and close I-0333

## Notes

Cost of delay inputs set by flai from I-0333. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T08:41:34Z, 0 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from I-0333's one instance and [server-performance.md](../../../design/system/server-performance.md) § 13 boards at once, measured again after S-0602 and S-0606 (S-0636). With five stories in review, `merge.queue` took 91 to 241 ms of each call. Most of it was git in `mergequeue.ComputeAt` (`flai/internal/mergequeue/eligible.go`):

- `verified` runs one `git rev-parse` per story in review, one after another. `mainDirt` runs two more (`MainBranch` and `mainHead`): 21 to 81 ms together.
- `mainDirt` runs `git status --porcelain` on the main checkout on every call: 52 to 148 ms.

The fix has three parts:

- Read every local branch's head and the checked-out branch in one `git for-each-ref` process, in place of the rev-parse calls.
- Run `git status` only when a story in review has a passing verify result, so its head can be held back by dirt. Run it at the same time as the heads read.
- `filling` takes main's head from the same read.

The queue's reasons and window stay as they are. A test that counts git processes reproduces the cause. Before the fix it counts one process per story in review, plus the status and the main-head reads. After it, the count is the same whatever the number of stories in review.

Caching git status on the watcher's word is left out. flai serve's file watcher covers the manifest and the `design`, `wip`, and `docs` folders, not the code tree, so it cannot say when the main checkout's dirt changed. Caching the heads read on the main-refs watcher's word, as S-0604 did for `release.pending`'s show-ref, is left for later. It would save one process of about 10 ms, and it would mean sharing `release`'s watch generation with `mergequeue`.

Touches, file by file, and where each came from:

| Touch | From | Why |
|-------|------|-----|
| `flai/internal/storygit/storygit.go` | layout | `MainBranch` and `BranchExists` live here; the one-process heads read joins them |
| `flai/internal/storygit/storygit_test.go` | layout | The heads read's real-git test |
| `flai/internal/mergequeue/eligible.go` | design | `ComputeAt`, `verified`, `mainDirt`, and `filling` (I-0333 names the file) |
| `flai/internal/mergequeue/eligible_test.go` | co-change (10%) | The git-process count that reproduces the cause |
| `design/system/server-performance.md` | design | The before and after, as S-0636 recorded them |
| `design/system/flai-cli.md` | co-change (29%), design | Its merge queue paragraph describes `Compute`, `ComputeAt`, and `filling` |
| `design/issues/I-0333-this-repository-s-board-get-takes-300-ms-or-more-in-half-its-calls-with-13.md` | design | Criterion 2: `flai issue close I-0333` |
| `design/issues/summary.md` | design | `flai issue close` rewrites the summary |

Co-change files left out: `design/system/flaiover-dashboard.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`, `flai/internal/mergequeue/land.go`, `flai/cmd/queue.go`, `flai/internal/mcpserver/server.go`, and the rest of what `flai touches suggest` lists. Neither the queue's output nor a command changes, and `mainHead` in `land.go` stays as it is for the landing. No folder touch.

No `after`: the story reads nothing another open story writes.

Forecast: 1h, raised from `flai forecast`'s 18m. The size-based figure leaves no room for a real-git test and a before-and-after `board.get` measure like S-0636's, which ran 130 calls a build on two builds. Delivery moves with it, from flai's 15:57 to 16:39.

Cost of delay: 12.50 USD a week, as `flai cod S-0698` works it out from the one input flai set from I-0333, `time_lost_per_cycle: 5m`, at 150 USD an hour. The value stands as computed. One occurrence is all the issue has recorded, and nothing gives a reason to weigh it more.
