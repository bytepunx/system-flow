---
id: S-0277
type: story
nature: improvement
title: "flai check finds `threads.archived` outside the story at close-out"
status: ready
owner: alex
created: 2026-10-05T04:40:46Z
updated: 2026-10-07T01:43:45Z
transitions:
  - to: ready
    at: 2026-10-07T01:12:26Z
    by: alex
tags: [flai]
topics: [cli, threads]
touches: [flai/internal/threads/archived.go, flai/internal/threads/archived_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/cmd/accept.go, flai/internal/preview/accept.go, flai/cmd/accept_threads_test.go, flai/cmd/archive.go, flai/cmd/archive_test.go, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, docs/users/flai-reference.md, design/adrs, design/issues/I-0073-flai-check-finds-threads-archived-outside-the-story-at-close-out.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h35m
    by: planner-S-0277
    at: 2026-10-07T01:21:52Z
  value: 237.5
  by: planner-S-0277
  at: 2026-10-07T01:21:57Z
forecast:
  duration: 45m
  delivery: 2026-10-07T04:09:00Z
  basis: "Its own forecast of 45m; 6th in the pull order with an in-progress limit of 3, behind S-0269, S-0294, S-0307, S-0270, S-0271, S-0274, S-0275 and S-0245."
  by: flai
  at: 2026-10-07T01:43:45Z
finalized:
  by: alex
  at: 2026-10-07T01:12:18Z
---
# S-0277 flai check finds `threads.archived` outside the story at close-out

## Goal

This story remediates [I-0073](../../../design/issues/I-0073-flai-check-finds-threads-archived-outside-the-story-at-close-out.md), "flai check finds `threads.archived` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0073 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0073 is closed with `flai issue close I-0073 --reason` saying what fixed it

## Tasks
- T-1142 The threads package resolves the threads left open on items being archived, and the threads.archived finding names flai thread resolve
- T-1144 flai accept resolves the threads still open on what it archives, in the acceptance commit, and its dry run lists them
- T-1146 flai archive resolves the threads still open on the items it archives, and its dry run lists them
- T-1147 The design, the user guide, the reference, and an ADR say that archiving an item resolves the threads still open on it
- T-1148 The threads left open on archived items are resolved, flai check finds no threads.archived, and I-0073 is closed

## Notes

### Planning

Planned by planner-S-0277 on 2026-10-07. The plan is summarised on TH-0236.

Proposed solution, from the issue's 27 instances. Each instance is a story accepted, and so archived, while a thread on it was still `open` or `answered`: TH-0112 on S-0249, TH-0127 on S-0276, TH-0194 on S-0227, TH-0203 on S-0229, and TH-0232 on S-0272. Neither `flai accept` (`flai/cmd/accept.go`) nor `flai archive` (`flai/cmd/archive.go`) looks at threads. Only the orchestrator's acceptance refuses while one is open (ADR-0093, `flai/internal/preview/orchestrator.go`). So flai `check` warns `threads.archived` at every close-out until someone runs `flai thread resolve` by hand. The fix: archiving an item, by acceptance or by `flai archive`, resolves the threads still open on it, with an entry naming why, in the same commit.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/threads/archived.go`, `archived_test.go` | layout: new helpers next to `Resolve` in `flai/internal/threads/threads.go` |
| `flai/internal/check/check.go`, `check_test.go` | design: the `threads.archived` rule and its test |
| `flai/cmd/accept.go`, `flai/internal/preview/accept.go`, `flai/cmd/accept_threads_test.go` | layout: the acceptance flow and its result type |
| `flai/cmd/archive.go`, `flai/cmd/archive_test.go` | layout: the archive command, which has no test file yet |
| `design/system/flai-cli.md`, `design/system/workflow.md`, `docs/users/flai.md`, `docs/users/flai-reference.md` | design: where accept, archive, and threads are described |
| `design/adrs` | design: the decision gets an ADR whose number is not known yet |
| `design/issues/I-0073-...md`, `design/issues/summary.md` | goal: the criterion that closes I-0073 |

`flai touches suggest` listed only documents that change often with any flai change, such as `design/system/flaiover-dashboard.md` and `docs/operators/settings.md`. None of them bears on this story, so none was added.

Folder touch kept: `design/adrs`. T-1147 adds an ADR whose file name is not known until `flai adr new` numbers it. The folder is in the manifest's `claims.shared`, so it holds no other story.

Holds and overlaps:

- `flai/cmd/accept.go` and `flai/internal/preview/accept.go` are also touched by S-0286 (in progress). This story is held until S-0286 moves to review.
- S-0307 (ready) touches the same two files.
- S-0279's planned tasks touch `flai/internal/check/check.go`.

Forecast: 45m, delivery 2026-10-07T04:30:00Z. `flai forecast` gave 24m, at 78 s per unit over 32 done improvement stories, times size 18 (2 criteria, 16 touches). I raised it because the five tasks include two command changes with new tests, an ADR, and a regenerated reference. Similar close-out remediations took 12m (S-0276), 29m (S-0249), and 43m (S-0266). The delivery is flai's 04:05Z, moved back by the extra 21m and rounded. The story is ninth in the pull order.

Cost of delay: 237.50 USD a week, as `flai cod S-0277` computes it and unadjusted. The input is `time_lost_per_cycle: 1h35m`, which the operator chose on TH-0233. It comes from I-0073's 27 instances between 2026-10-05T01:38Z and 2026-10-07T01:04Z, about 95 a week, at about 1m each, with `hour_rate` 150 and a 168h cycle.
