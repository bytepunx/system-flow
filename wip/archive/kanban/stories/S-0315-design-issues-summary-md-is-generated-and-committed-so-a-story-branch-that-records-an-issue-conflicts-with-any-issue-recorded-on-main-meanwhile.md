---
id: S-0315
type: story
nature: improvement
title: design/issues/summary.md is generated and committed, so a story branch that records an issue conflicts with any issue recorded on main meanwhile
status: done
owner: alex
created: 2026-10-07T18:59:45Z
updated: 2026-10-08T04:07:34Z
transitions:
  - to: ready
    at: 2026-10-07T23:53:23Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-07T23:53:29Z
    by: agent-S-0315
  - to: review
    at: 2026-10-08T00:10:16Z
    by: agent-S-0315
  - to: done
    at: 2026-10-08T04:07:34Z
    by: alex
tags: [cli]
topics: [cli, git, continuous-improvement]
touches: [flai/cmd/stream_sync_test.go, design/issues/I-0089-design-issues-summary-md-is-generated-and-committed-so-a-story-branch-that-records-an-issue-conflicts-with-any-issue-recorded-on-main-meanwhile.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1017
  turns:
    - day: 2026-10-07
      test_runs: 3
      hand_edits: 1
      work: 20
    - day: 2026-10-08
      work: 2
  models:
    - model: claude-opus-5-5
      input: 54
      output: 11089
      cache_read: 2743156
      cache_write: 149800
      cost: 1.969
  strategic:
    - kind: orchestrator
      seconds: 510
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 52
          output: 750
          cache_read: 13211756
          cache_write: 39807
          cost: 3.265
        - model: claude-sonnet-5-5
          input: 4
          output: 40
          cache_read: 11859
          cache_write: 27642
          cost: 0.0341
cost_of_delay:
  inputs:
    time_lost_per_cycle: 3m
    by: flai
    at: 2026-10-07T18:59:45Z
  value: 7.5
  by: planner-S-0315
  at: 2026-10-07T23:44:33Z
forecast:
  duration: 15m
  delivery: 2026-10-08T06:31:00Z
  basis: "Its own forecast of 15m; 21st in the pull order with an in-progress limit of 3, behind S-0232, S-0310, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0287, S-0288, S-0289, S-0290, S-0291, S-0297, S-0304, S-0305, S-0306, S-0309, S-0312 and S-0313."
  by: flai
  at: 2026-10-07T23:47:55Z
finalized:
  by: orchestrator
  at: 2026-10-07T23:53:18Z
---
# S-0315 design/issues/summary.md is generated and committed, so a story branch that records an issue conflicts with any issue recorded on main meanwhile

## Goal

This story remediates [I-0089](../../../design/issues/I-0089-design-issues-summary-md-is-generated-and-committed-so-a-story-branch-that-records-an-issue-conflicts-with-any-issue-recorded-on-main-meanwhile.md), "design/issues/summary.md is generated and committed, so a story branch that records an issue conflicts with any issue recorded on main meanwhile". The issue recommends this solution:

Directions to weigh: `flai stream sync` and `flai accept` regenerate `summary.md` themselves when it is the only conflict, since it is derived from the issue files; or the file is not committed and is generated where it is read.

S-0278 built the first direction ([ADR-0098](../adrs/0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md)) for I-0074, which has the same cause. The operator chooses at its acceptance whether that closes this issue too.

## Acceptance criteria
- [x] The cause I-0089 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0089 is closed with `flai issue close I-0089 --reason` saying what fixed it

## Tasks
- T-1293 A test reproduces I-0089: an issue recorded and one closed on main itself while a story branch records its own, and the sync and the acceptance regenerate design/issues/summary.md
- T-1294 Close I-0089 saying what fixed it

## Notes

Cost of delay inputs set by flai from I-0089. time_lost_per_cycle 3m: 3m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T10:32:27Z, 1.4 days before this story; under one cycle counts as one).

### Planning

The fix is already on main. S-0278 was accepted at 2026-10-06T20:07Z, after I-0089's only instance (S-0220, 06:07Z and 07:03Z that day). Its regeneration of `summary.md` on a rebase stop ([ADR-0098](../../../design/adrs/0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md)) covers I-0089's case too. Its tests reach main only through another story's acceptance, though. So this story adds a test of I-0089's own case (T-1293), then closes the issue (T-1294). It builds no new mechanism. If the test fails, the story's agent adds a fix task in `flai/internal/storygit/sync.go`.

Touches, file by file, none a folder:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/cmd/stream_sync_test.go` | layout | S-0278's tests of the regeneration are here, with the helpers T-1293 reuses. `flai touches suggest` from it and the summary lists issue files and docs, which this story does not change. |
| `design/issues/I-0089-…md` | design | Criterion 2 closes it. |
| `design/issues/summary.md` | design | `flai issue close` regenerates it. |

`design/issues` is in the manifest's `claims.shared`, so the two issue paths hold no story. `flai/cmd/stream_sync_test.go` is also in S-0297's touches (backlog). Whichever of the two is pulled second is held until the first is accepted.

Figures:

- **Forecast: 15m, up from flai's 4m.** flai sized it from 2 criteria and no touches, at 94 s per unit. S-0278's like tasks took longer: T-1012 (its sync test and fix) 526 s, and T-1015 (closing its issue) 224 s. One test reusing its helpers, one close, and the close-out come to about 15m. The delivery stays flai's, 2026-10-08T06:21:00Z. It comes from the story's 21st place in the pull order, which 11 more minutes do not move.
- **Cost of delay: 7.50 USD a week, as `flai cod` gives it.** That is 3m lost per 168h cycle at 150 USD an hour, from flai's input. It stands. Since S-0278, the conflict costs nothing new; what is left is an open issue with no test of its own case.
