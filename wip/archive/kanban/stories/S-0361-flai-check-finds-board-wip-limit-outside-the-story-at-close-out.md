---
id: S-0361
type: story
nature: improvement
title: "flai check finds `board.wip-limit` outside the story at close-out"
status: done
owner: alex
created: 2026-10-12T14:12:27Z
updated: 2026-10-12T14:31:39Z
transitions:
  - to: ready
    at: 2026-10-12T14:20:12Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-12T14:20:18Z
    by: agent-S-0361
  - to: review
    at: 2026-10-12T14:26:07Z
    by: agent-S-0361
  - to: done
    at: 2026-10-12T14:31:39Z
    by: orchestrator
tags: []
touches: [flai/internal/check/scope_test.go, design/issues/I-0122-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 360
  turns:
    - day: 2026-10-12
      ceremony: 1
      hand_edits: 1
      work: 17
  models:
    - model: claude-opus-5-5
      input: 40
      output: 5737
      cache_read: 2119766
      cache_write: 127377
      cost: 1.5579
  strategic:
    - kind: orchestrator
      seconds: 1099
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 96
          output: 1295
          cache_read: 7642250
          cache_write: 39182
          cost: 1.9366
        - model: claude-sonnet-5-5
          input: 6
          output: 46
          cache_read: 40248
          cache_write: 30904
          cost: 0.0575
cost_of_delay:
  inputs:
    penalty_per_week: 0
    by: orchestrator
    at: 2026-10-12T14:14:59Z
  value: 0
  by: planner-S-0361
  at: 2026-10-12T14:15:33Z
forecast:
  duration: 9m
  delivery: 2026-10-13T01:51:00Z
  basis: "flai forecast: median 100 s per unit of size over 9 done small improvement stories on claude-opus-5-5, times size 5 (2 criteria, 3 touches); 25th in the pull order with an in-progress limit of 5. Kept: the fix is already on main (S-0348), so the work is a confirmation, a comment, and an issue close."
  by: planner-S-0361
  at: 2026-10-12T14:13:58Z
finalized:
  by: orchestrator
  at: 2026-10-12T14:20:05Z
---
# S-0361 flai check finds `board.wip-limit` outside the story at close-out

## Goal

This story remediates [I-0122](../../../design/issues/I-0122-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md), "flai check finds `board.wip-limit` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0122 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0122 is closed with `flai issue close I-0122 --reason` saying what fixed it

## Tasks
- T-1435 Confirm S-0348's scoped-check test covers I-0122's instance and name I-0122 in it
- T-1436 Close I-0122 saying S-0348 fixed it

## Notes

### Planning

The cause is already fixed. I-0122's one instance, S-0291's close-out at 2026-10-08T08:30:23Z, has the same title and cause as I-0123. S-0348 fixed that at fd5c2970 (2026-10-08T09:35:05Z, flai 1.39.15) under [ADR-0133](../../../design/adrs/0133-a-close-out-records-no-board-wip-limit-and-a-check-scoped-to-a-story-leaves-out.md), with `TestScopedCheckLeavesOutABoardWIPLimitOutsideTheStory`. TH-0388 recommends closing I-0122 as a duplicate. So the plan confirms the test covers the instance and closes the issue; it builds no fix.

Touches, all files, no folder:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/check/scope_test.go` | design | ADR-0133's test; T-1435 names I-0122 in its comment |
| I-0122's issue file under `design/issues/` | layout | `flai issue close` rewrites its status |
| `design/issues/summary.md` | layout | regenerated when the issue closes |

`flai touches suggest` listed co-changed issue files, `design/system/flai-cli.md`, and the docs. None is kept: the story changes no behaviour, so no design or docs change.

Forecast 9m, delivery 2026-10-13T01:51:00Z, as `flai forecast` gave it: 100 s per unit of size times size 5. Kept: the work is a test comment and an issue close.

Cost of delay 0.00 USD a week, as `flai cod` gave it from the orchestrator's input on TH-0399: a penalty of 0 a week. Kept: the cause stopped occurring with flai 1.39.15, so waiting costs only an open row in `design/issues/summary.md`.

### Accepted by the orchestrator

- Verified: bf9682f34307fdb95b8b2df0f87b628db3f08bfe
- At: 2026-10-12T14:31:39Z

Verdict: pass at bf9682f34307fdb95b8b2df0f87b628db3f08bfe; flai verify passed every step at that commit, and the verifier matched both criteria to the diff.

- 1: flai/internal/check/scope_test.go
- 2: design/issues/I-0122-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md, design/issues/summary.md
