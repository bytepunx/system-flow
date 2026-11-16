---
id: S-0739
type: story
nature: feature
title: workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name
status: review
parent: E-0034
owner: alex
created: 2026-11-12T15:59:07Z
updated: 2026-11-15T07:42:39Z
transitions:
  - to: ready
    at: 2026-11-14T17:47:28Z
    by: alex
  - to: in-progress
    at: 2026-11-14T23:49:15Z
    by: agent-S-0739
  - to: review
    at: 2026-11-15T01:15:44Z
    by: agent-S-0739
  - to: in-progress
    at: 2026-11-15T07:30:26Z
    by: flai
  - to: review
    at: 2026-11-15T07:42:39Z
    by: agent-S-0739
tags: [cli]
topics: [cli, planning]
touches: [flai/internal/workitem/boardview.go, flai/internal/workitem/boardview_test.go, flai/internal/workitem/board.go, flai/internal/workitem/workflow_fixture_test.go, flai/internal/workitem/follow.go, flai/internal/workitem/follow_test.go, flai/internal/workitem/taskplan.go, flai/internal/workitem/taskplan_test.go, flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/workitem/order.go, flai/internal/workitem/plancandidates.go, flai/internal/workitem/promote.go, flai/internal/workitem/promote_test.go, flai/internal/workitem/promotable.go, flai/internal/workitem/finalizable.go, flai/internal/workitem/changes.go, flai/internal/workitem/share.go, flai/internal/workitem/policy.go, flai/internal/workitem/narrative.go, flai/internal/workitem/flow.go, design/system/flai-cli.md, flai/internal/workitem/order_test.go, flai/internal/workitem/plancandidates_test.go, flai/internal/workitem/changes_test.go, flai/internal/workitem/share_test.go, flai/internal/workitem/policy_test.go, flai/internal/workitem/narrative_test.go, flai/internal/workitem/finalizable_test.go, design/issues/I-0376-flai-guard-counts-a-sub-agent-that-ended-on-an-api-error-as-still-running-and.md, design/issues/summary.md]
after: [S-0737]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 5968
  priced_by: harness
  turns:
    - day: 2026-11-14
      work: 25
    - day: 2026-11-15
      ceremony: 2
      test_runs: 1
      work: 61
  models:
    - model: claude-opus-5-5
      input: 494
      output: 192594
      cache_read: 24776161
      cache_write: 1456811
      cost: 17.7225
  strategic:
    - kind: orchestrator
      seconds: 359
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 87
          output: 628
          cache_read: 3050926
          cache_write: 420022
          cost: 1.0398
        - model: claude-sonnet-5-5
          input: 16
          output: 122
          cache_read: 111468
          cache_write: 97188
          cost: 0.164
cost_of_delay:
  value: 80.17
  by: system-flow
  at: 2026-11-14T17:43:24Z
forecast:
  duration: 45m
  delivery: 2026-11-15T01:09:00Z
  basis: "Its own forecast of 45m; 4th in the pull order with an in-progress limit of 3, behind S-0738, S-0741 and S-0746."
  by: flai
  at: 2026-11-14T23:48:54Z
finalized:
  by: orchestrator
  at: 2026-11-14T17:44:31Z
---
# S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name

## Goal

The `workitem` package stops naming states. An epic follows its stories (ADR-0076), the board view lays out its columns, a story's task plan finds its layers, a ready story is held, the pull order is kept, and a story is promoted or finalized by the kind of the state each item is in, and by the WIP limit and order the workflow gives that state, so they work on a workflow they have never seen.

## Acceptance criteria

- [x] `follow.go`, `boardview.go`, `board.go`, `taskplan.go`, `hold.go`, `order.go`, `plancandidates.go`, `promote.go`, `promotable.go`, `finalizable.go`, `changes.go`, `share.go`, `policy.go`, and `narrative.go` name no state: each reads a state's kind, WIP limit, or order from the workflow
- [x] The board view gives its columns in the workflow's order, with each column's kind and WIP limit, so a reader of the view needs no list of its own
- [x] An epic follows its stories by kind: to the first queue state with the first story, to active work with the first started, to the wait state with the last open one, and to accepted when that one is accepted
- [x] Tests run each of these against a workflow with a state renamed and a state added between `in-progress` and `review`, and give the same result as the baseline wherever the kinds agree
- [x] With the baseline, every existing `workitem` test passes unchanged

## Tasks

- T-3203 The board view gives the workflow's columns in order with each one's kind and WIP limit, with a shared test workflow that renames a state and adds one
- T-3206 An epic follows its stories by the kinds of their states
- T-3208 The task plan, the holds, the pull order, and the plan candidates read a state's kind
- T-3210 Promotion, finalizing, the change feed, shares, the policy order, and the narrative's checks read a state's kind

## Notes

- These files hold about 80 of the package's references to a state's name (`follow.go` 18, `taskplan.go` 11, `boardview.go` 10, `plancandidates.go` 8, `hold.go` 7, `order.go` 7). `item.go` keeps the constants for the baseline's names.

### Planning

- Touches, all files, from layout: a count of each file's references to a state's name in `flai/internal/workitem`, plus the test files beside them and `workflow_fixture_test.go`, a new test file for the workflow the package's tests share. `item.go` was dropped from the first draft's touches: the constants stay.
- Left out from `flai touches suggest`: `mcpserver/server.go`, `check/check.go`, `cmd/board.go`, `cmd/edit.go`, and `itemedit` (12% to 19%). Each is another story's under E-0034.
- `after`: S-0737. This story reads the kinds, WIP limits, and order from the workflow `Repo` carries from S-0737 on, and its tests open a repository with another workflow through S-0737's test hook.
- Tasks: T-3203 (the board view and the shared fixture) first; then T-3206, T-3208, and T-3210 together, since each uses the fixture and none shares a path with another.
- Forecast 45m, as `flai forecast`'s rate gives it for size 25 once the touches were settled (it gave 44m for 24). The rewrite is mechanical, but `follow.go` and `taskplan.go` carry logic that must give the same epics and layers.
- Cost of delay: 80.17 USD a week, as `flai cod` gives it: this story's share of E-0034's 1500.00 USD a week, the operator's `time_lost_per_cycle: 10h` (TH-1115) at 150 USD an hour, by its forecast of 45m over the 14h2m of the epic's 19 open stories. Kept: the epic's outcome needs every story, so a share by size is the split that fits.
- 2026-11-15T07:30:26Z: moved to in-progress: alex asked its agent to rebase story/S-0739 onto main and resolve what holds its acceptance, 1 blocker: S-0739's verify does not cover the main branch it would land on: its last flai verify, at 9f905b8, ran over main at 449ee65, and the 58 paths main brought in si...; its agent syncs it with flai stream sync S-0739, takes each blocker's recommended action, closes out again, and moves it to review
