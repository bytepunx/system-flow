---
id: S-0833
type: story
nature: remediation
title: flai release --pending fails when another process holds the main checkout's git index lock, and the publish formula's step then stops on exit 1
status: done
owner: alex
created: 2026-11-17T08:17:27Z
updated: 2026-11-17T09:46:17Z
transitions:
  - to: ready
    at: 2026-11-17T08:34:27Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-17T09:15:12Z
    by: agent-S-0833
  - to: review
    at: 2026-11-17T09:45:01Z
    by: agent-S-0833
  - to: done
    at: 2026-11-17T09:46:17Z
    by: orchestrator
tags: [flai, release, git]
topics: [release]
touches: [flai/cmd/release.go, flai/cmd/release_lock_test.go, docs/users/flai.md, design/issues/I-0359-flai-release-pending-fails-when-another-process-holds-the-main-checkout-s-git.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1804
  priced_by: harness
  turns:
    - day: 2026-11-17
      ceremony: 2
      work: 36
  models:
    - model: claude-opus-5-5
      input: 134
      output: 23781
      cache_read: 9618927
      cache_write: 316643
      cost: 4.7658
  strategic:
    - kind: planner
      seconds: 30
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 10
          output: 25
          cache_read: 902548
          cache_write: 20450
          cost: 0.063
        - model: claude-opus-5-5
          input: 61
          output: 5710
          cache_read: 4498452
          cache_write: 69861
          cost: 1.5422
    - kind: orchestrator
      seconds: 19
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 48
          cache_read: 354304
          cache_write: 7664
          cost: 0.1123
cost_of_delay:
  inputs:
    time_lost_per_cycle: 16m
    by: flai
    at: 2026-11-17T08:17:27Z
  value: 40
  by: planner-S-0833
  at: 2026-11-17T08:30:26Z
forecast:
  duration: 18m
  delivery: 2026-11-17T09:30:00Z
  basis: "Its own forecast of 18m; 1st in the pull order with an in-progress limit of 3, behind S-0752, S-0806 and S-0831."
  by: flai
  at: 2026-11-17T09:04:31Z
finalized:
  by: orchestrator
  at: 2026-11-17T08:33:51Z
---
# S-0833 flai release --pending fails when another process holds the main checkout's git index lock, and the publish formula's step then stops on exit 1

## Goal

This story remediates [I-0359](../../../design/issues/I-0359-flai-release-pending-fails-when-another-process-holds-the-main-checkout-s-git.md), "flai release --pending fails when another process holds the main checkout's git index lock, and the publish formula's step then stops on exit 1". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0359 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0359 is closed with `flai issue close I-0359 --reason` saying what fixed it

## Tasks
- T-3535 The publish commit of flai release --pending waits out a held index lock, with a test that holds it
- T-3538 docs/users/flai.md says flai release --pending waits out a held index lock as acceptance does
- T-3539 Close I-0359 with the reason that the publish commit waits out a held index lock

## Notes

Cost of delay inputs set by flai from I-0359. time_lost_per_cycle 16m: 8m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-11-13T08:58:25Z, 4 days before this story; under one cycle counts as one).

### Planning

Proposed solution, from I-0359's two instances: `computeApplyAndTagPending` in `flai/cmd/release.go` runs the publish's `git add -A` and `git commit` straight through the runner. I-0100 fixed the same race in `flai accept` with `storygit.RunPastIndexLock` (`flai/internal/storygit/indexlock.go`), which waits about nine seconds and never removes the lock. Both instances saw the lock gone within seconds, so routing those two calls through the helper fixes the cause. `flai release --pending` is the one publish path: `release_publish` and the dashboard's publish run it, so the fix holds once the publish formula is retired (ADR-0316, proposed).

Touches, all files, no folder:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/cmd/release.go` | layout | The two git calls of the publish commit, at `computeApplyAndTagPending`. `flai touches suggest S-0833 flai/cmd/release.go` was the start; the story declared no touches. |
| `flai/cmd/release_lock_test.go` | layout | A new test, after the `lockingRunner` pattern of `flai/cmd/accept_lock_test.go`, kept apart from `release_pending_test.go`. |
| `docs/users/flai.md` | co-change (29%) | § Accept and release describes acceptance's wait for the lock and not the publish's. |
| `design/issues/I-0359-…md`, `design/issues/summary.md` | design | `flai issue close`, criterion 2. |

Co-change also lists `flai/internal/release/release.go` (46%) and `docs/users/flai-reference.md` (42%). Neither is predicted: the tagging in `release.go` takes no index lock, and the change adds no flag, so the generated reference does not move.

Forecast 18m stands as `flai forecast` gives it: the change is two calls and a test, and the close-out's integration and smoke tiers for a change under `flai/` take most of that time.

Cost of delay value 40 USD a week stands as `flai cod` works it from the operator's input `time_lost_per_cycle` 16m at 150 USD an hour.

No `after`: the story reads nothing another open story writes. S-0811 (ready) also changes the release's publish. A trial merge judges any overlap.

### Accepted by the orchestrator

- Verified: bb41b9fbeabf54194f23c4a8b9fd059d7ad8925b
- At: 2026-11-17T09:46:17Z

Verdict: S-0833 meets both acceptance criteria and breaks no convention; close-out passed at bb41b9fbeabf54194f23c4a8b9fd059d7ad8925b.
- 1: flai/cmd/release.go, flai/cmd/release_lock_test.go, docs/users/flai.md
- 2: design/issues/I-0359-flai-release-pending-fails-when-another-process-holds-the-main-checkout-s-git.md, design/issues/summary.md
