---
id: S-0792
type: story
nature: remediation
title: Every story has to go back to be rebased and retested
status: done
owner: alex
created: 2026-11-15T07:37:37Z
updated: 2026-11-15T18:51:49Z
transitions:
  - to: ready
    at: 2026-11-15T07:37:38Z
    by: alex
  - to: in-progress
    at: 2026-11-15T16:42:04Z
    by: agent-S-0792
  - to: review
    at: 2026-11-15T17:20:36Z
    by: agent-S-0792
  - to: in-progress
    at: 2026-11-15T18:40:30Z
    by: flai
  - to: review
    at: 2026-11-15T18:50:19Z
    by: agent-S-0792
  - to: done
    at: 2026-11-15T18:51:49Z
    by: orchestrator
tags: []
topics: [acceptance, verify]
touches: [design/adrs, design/adrs/README.md, flai/internal/verify/overlap.go, flai/internal/verify/overlap_test.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/cmd/accept_stale_verify_test.go, flai/cmd/accept.go, docs/users/flai-reference.md, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0378-an-acceptance-sends-a-story-back-for-a-verify-that-missed-one-main-commit.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2905
  priced_by: harness
  turns:
    - day: 2026-11-15
      ceremony: 3
      work: 58
  models:
    - model: claude-opus-5-5
      input: 230
      output: 71982
      cache_read: 17387384
      cache_write: 785294
      cost: 10.5378
  strategic:
    - kind: planner
      seconds: 380
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 42
          output: 23055
          cache_read: 1368547
          cache_write: 113999
          cost: 0.0966
        - model: claude-opus-5-5
          input: 42
          output: 16509
          cache_read: 2586213
          cache_write: 165540
          cost: 2.1719
    - kind: orchestrator
      seconds: 207
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 49
          output: 360
          cache_read: 1698572
          cache_write: 309340
          cost: 0.6039
        - model: claude-sonnet-5-5
          input: 6
          output: 49
          cache_read: 48472
          cache_write: 46569
          cost: 0.0748
cost_of_delay:
  inputs:
    time_lost_per_cycle: 12h
    by: alex
    at: 2026-11-15T07:37:37Z
  value: 1800
  by: planner-S-0792
  at: 2026-11-15T07:42:54Z
forecast:
  duration: 1h30m
  delivery: 2026-11-15T10:23:00Z
  basis: "Its own forecast of 1h30m; 5th in the pull order with an in-progress limit of 3, behind S-0741, S-0749, S-0748 and S-0750."
  by: flai
  at: 2026-11-15T07:47:13Z
---
# S-0792 Every story has to go back to be rebased and retested

## Goal

Recently every single story in review is getting the following error requiring rebase and re-verify. This is adding fat to much time.

Error: “verify does not cover the main branch it would land on: its last flai verify …”

We must correct this immediately.

## Acceptance criteria
- [x] Every story in review should not have to rebase and reverify for every story that is approved before it.
- [x] The only time we need to reverify is when changes between a story accepted overlap with the same change set any given story in review such that the behavior of the PR is undoubtedly impacted by the changes in main that will be rebased on.

## Tasks
- T-3376 Record in an ADR that verify_behind_main sends back only the tiers main's paths overlapping the story's change set select
- T-3377 Add the overlap of main's new paths with a story's change set to flai/internal/verify
- T-3378 Make verify_behind_main send back only the tiers the overlapping main paths select
- T-3379 Describe the narrowed verify_behind_main gate in the design, the user guide, and flai accept's help

## Notes

### Planning

The blocker is `verify_behind_main` ([ADR-0291](../../../design/adrs/0291-every-acceptance-but-the-merge-queue-s-landing-waits-for-a-flai-verify-over.md)), computed in `verifyBehindMain` in `flai/internal/preview/accept.go`. Today any path main brought in since the story's last passing verify sends back every tier it selects, and this repository's `integration` and `smoke` tiers select all of `flai/**`, so every acceptance that changes flai sends back every other story in review (TH-1213, TH-1215, TH-1216; [I-0378](../../../design/issues/I-0378-an-acceptance-sends-a-story-back-for-a-verify-that-missed-one-main-commit.md)).

The plan narrows the gate to main's paths that overlap the story's change set, the paths its last verify recorded (`Report.Paths`): a path the story changes too, or a Go package that a package the story changes imports, directly or through others (the existing `importers` closure in `flai/internal/verify/importers.go`). Only the tiers those overlapping paths select are sent back. The I-0229 case that ADR-0291 was written for (main changed `flai/internal/roles`, the story's test in `flai/internal/guard` imports it) stays blocked. This is an assumption, stated on the plan's thread: if the operator wants file overlap alone, only T2's rule and T1's ADR change.

Touches:

| Touch | Source |
|-------|--------|
| `design/adrs` | Folder kept: the new ADR that refines ADR-0291 gets its file name from `flai adr new`, which no task can name yet. |
| `design/adrs/README.md` | Layout: the ADR index flai regenerates. |
| `flai/internal/verify/overlap.go`, `overlap_test.go` | Layout: new files for the overlap of main's paths with the story's change set, beside `importers.go`. |
| `flai/internal/preview/accept.go`, `accept_test.go` | Co-change and the explorer's map: the gate and its test, `TestAnAcceptanceWaitsForAVerifyOverMainsNewPaths`. |
| `flai/cmd/accept_stale_verify_test.go` | Explorer: the end-to-end refusal test, whose same-directory case needs checking against the new rule. |
| `flai/cmd/accept.go`, `docs/users/flai-reference.md` | Co-change (24%): the accept help text and the reference generated from it. |
| `design/system/workflow.md`, `design/system/flai-cli.md`, `docs/users/flai.md` | Design and co-change: the documents that describe the gate. |

Figures:

- Forecast: 1h30m, adjusted from `flai forecast`'s 6m, which counted 2 criteria and no touches. S-0683 built this gate, with an ADR, the same tests, and the same docs, in 55m of agent time. Delivery 2026-11-15T10:00Z plays the longer duration from the 6th place in the pull order that `flai forecast` gave.
- Cost of delay: 1800 USD a week, as `flai cod` gives it from the operator's 12h lost per cycle at 150 USD an hour, one cycle a week. It stands: each send back costs a sync and a ten-minute close-out per story in review.
- No `after`: the story reads no output of an open story.
- 2026-11-15T18:40:30Z: moved to in-progress: alex asked its agent to rebase story/S-0792 onto main and resolve what holds its acceptance, 1 blocker: S-0792's verify does not cover the main branch it would land on: its last flai verify, at 3231fa9, ran over main at 0d2f772, and the 2 paths main brought in sin...; its agent syncs it with flai stream sync S-0792, takes each blocker's recommended action, closes out again, and moves it to review

### Accepted by the orchestrator

- Verified: ea2ab57ac980d2f7eda49c2cf3606bc34d2ee149
- At: 2026-11-15T18:51:49Z

Verdict: S-0792 meets both acceptance criteria and breaks no convention, verified at ea2ab57ac980d2f7eda49c2cf3606bc34d2ee149.
- 1: flai/internal/preview/accept.go, flai/internal/verify/overlap.go, flai/internal/verify/overlap_test.go, flai/internal/preview/accept_test.go, design/adrs/0303-an-acceptance-s-verify-behind-main-sends-back-only-the-tiers-selected-by-the.md
- 2: flai/internal/verify/overlap.go, flai/internal/verify/overlap_test.go, flai/cmd/accept_stale_verify_test.go, flai/internal/preview/accept.go, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, flai/cmd/accept.go, docs/users/flai-reference.md
