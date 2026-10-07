---
id: S-0307
type: story
nature: remediation
title: The acceptance's commit step fails when another process holds git's index lock, and the acceptance cannot then be finished by flai
status: in-progress
owner: alex
created: 2026-10-07T01:07:13Z
updated: 2026-10-07T01:42:04Z
transitions:
  - to: ready
    at: 2026-10-07T01:10:44Z
    by: alex
  - to: in-progress
    at: 2026-10-07T01:42:04Z
    by: system-flow
tags: [flai]
topics: [cli, git]
touches: [flai/internal/storygit/indexlock.go, flai/internal/storygit/indexlock_test.go, flai/internal/preview/accept.go, flai/cmd/accept.go, flai/cmd/accept_lock_test.go, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0100-the-acceptance-s-commit-step-fails-when-another-process-holds-git-s-index-lock-and-the-acceptance-cannot-then-be-finished-by-flai.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2m
    by: flai
    at: 2026-10-07T01:07:13Z
  value: 5
  by: planner-S-0307
  at: 2026-10-07T01:13:34Z
forecast:
  duration: 40m
  delivery: 2026-10-07T04:05:00Z
  basis: "flai forecast's 22m (107 s per unit of size over 15 medium-band remediation stories, times size 12), raised 18m for a real-git test of a held index lock and the resumed commit's overlap notices, which size does not count; delivery shifted by the same."
  by: planner-S-0307
  at: 2026-10-07T01:13:34Z
finalized:
  by: alex
  at: 2026-10-07T01:10:39Z
---
# S-0307 The acceptance's commit step fails when another process holds git's index lock, and the acceptance cannot then be finished by flai

## Goal

This story remediates [I-0100](../../../design/issues/I-0100-the-acceptance-s-commit-step-fails-when-another-process-holds-git-s-index-lock-and-the-acceptance-cannot-then-be-finished-by-flai.md), "The acceptance's commit step fails when another process holds git's index lock, and the acceptance cannot then be finished by flai". The issue recommends this solution:

Directions to weigh: retry the commit a few times when git reports the index lock held, since the other writer is brief; or make an acceptance that is done and archived but whose commit is missing resumable, as one that is done but not archived already is ("completed from step 0 without a second transition"); and record the commit's failure in the journal.

## Acceptance criteria
- [ ] The cause I-0100 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0100 is closed with `flai issue close I-0100 --reason` saying what fixed it

## Tasks
- T-1137 A storygit helper runs a git command again while another process holds the index lock, and gives up with git's whole error
- T-1138 flai accept commits through the index-lock retry, and an acceptance whose commit failed is finished by running it again
- T-1139 The design and the users' guide say that acceptance waits out a held index lock and that a failed acceptance commit is finished with flai accept
- T-1140 I-0100 is closed with flai issue close, saying that acceptance retries a held index lock and finishes a failed commit on a second run

## Notes

Cost of delay inputs set by flai from I-0100. time_lost_per_cycle 2m: 2m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T23:31:35Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0307 on 2026-10-07. The plan thread on S-0307 lists the tasks, their layers, and the assumptions.

The plan takes all three of the issue's directions:

- Retry the commit while the index lock is held.
- Make a done, archived acceptance whose commit is missing resumable.
- Fail with git's whole error and the command that finishes the acceptance. The run exits non-zero, so the host's journal records a failure.

Layers:

1. T-1137, the index-lock retry helper.
2. T-1138, the acceptance's commit through it, the resumed commit, and the test that reproduces I-0100.
3. T-1139, the design and the guides; T-1140, closing I-0100.

Touches, file by file, with no folder touch:

- **Declared:** none. The story declared no touches.
- **Layout:**
  - `flai/internal/storygit/indexlock.go` and its test: a new file beside the package's other git helpers. It is not `commit.go`, which S-0275 plans.
  - `flai/cmd/accept_lock_test.go`: accept tests live in `flai/cmd` as `accept_*_test.go`.
- **Design and layout:**
  - `flai/cmd/accept.go`: holds step 3, the `git add -A` and `git commit`.
  - `flai/internal/preview/accept.go`: holds the refusal of a closed item that the resumed commit narrows.
- **Design:**
  - `design/system/flai-cli.md`: the `flai accept` row.
  - `docs/users/flai.md`: Accept and release.
- **Co-change and design:** `docs/users/flai-reference.md`, generated from the changed help.
- **Criterion 2:** the I-0100 file and `design/issues/summary.md`, which `flai issue close` writes.
- **Not taken from `flai touches suggest`:** what it listed is what any flai story's docs co-change with, none of it named by the goal. That is the dashboard design, the operators' guides, the ADR index, and `flai/internal/hostapi/writes.go`. The journal entry already takes its outcome from the run's error, so T-1138 confirms it rather than changing `writes.go`. That keeps the story clear of S-0269 to S-0275, which all touch `writes.go`.

`flai/cmd/accept.go` and `flai/internal/preview/accept.go` overlap S-0286, which is in progress, through its T-1120. So S-0307 is held until S-0286 moves to review. No other path avoids it. The docs and issue paths are shared paths, which never hold.

Forecast 40m, delivery 2026-10-07T04:05Z.

- `flai forecast` gave 22m: 107 s per unit of size over 15 done medium-band remediation stories on claude-opus-5-5, times size 12 (2 criteria, 10 touches). It delivered at 03:47Z, 7th in the pull order.
- Raised by 18m. T-1138's real-git test holds the lock and releases it on a timer. The resumed commit takes its overlap notices from the story's commits. Size counts neither.
- Delivery is shifted by the same 18m. It does not count the hold behind S-0286.

Cost of delay 5 USD a week, as `flai cod` gives from flai's input from I-0100: 2m lost per 168h cycle at 150 USD an hour. It stands as computed. One occurrence so far. A failed commit leaves main with a staged acceptance that someone must notice and commit, which the figure does not price.
