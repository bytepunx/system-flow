---
id: S-0739
type: story
nature: feature
title: workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T15:59:07Z
updated: 2026-11-13T12:53:45Z
transitions: []
tags: [cli]
topics: [cli, planning]
touches: [flai/internal/workitem/boardview.go, flai/internal/workitem/boardview_test.go, flai/internal/workitem/board.go, flai/internal/workitem/workflow_fixture_test.go, flai/internal/workitem/follow.go, flai/internal/workitem/follow_test.go, flai/internal/workitem/taskplan.go, flai/internal/workitem/taskplan_test.go, flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/workitem/order.go, flai/internal/workitem/plancandidates.go, flai/internal/workitem/promote.go, flai/internal/workitem/promote_test.go, flai/internal/workitem/promotable.go, flai/internal/workitem/finalizable.go, flai/internal/workitem/changes.go, flai/internal/workitem/share.go, flai/internal/workitem/policy.go, flai/internal/workitem/narrative.go]
after: [S-0737]
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
      seconds: 12
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 14
          cache_read: 108314
          cache_write: 18020
          cost: 0.0361
draft: true
forecast:
  duration: 45m
  delivery: 2026-11-13T18:04:00Z
  basis: "Its own forecast of 45m; 13th in the pull order with an in-progress limit of 3, behind S-0661, S-0740, S-0769, S-0760, S-0745, S-0759, S-0751, S-0757, S-0663, S-0662, S-0241, S-0730, S-0736, S-0737 and S-0738."
  by: flai
  at: 2026-11-13T12:53:45Z
---
# S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name

## Goal

The `workitem` package stops naming states. An epic follows its stories (ADR-0076), the board view lays out its columns, a story's task plan finds its layers, a ready story is held, the pull order is kept, and a story is promoted or finalized by the kind of the state each item is in, and by the WIP limit and order the workflow gives that state, so they work on a workflow they have never seen.

## Acceptance criteria

- [ ] `follow.go`, `boardview.go`, `board.go`, `taskplan.go`, `hold.go`, `order.go`, `plancandidates.go`, `promote.go`, `promotable.go`, `finalizable.go`, `changes.go`, `share.go`, `policy.go`, and `narrative.go` name no state: each reads a state's kind, WIP limit, or order from the workflow
- [ ] The board view gives its columns in the workflow's order, with each column's kind and WIP limit, so a reader of the view needs no list of its own
- [ ] An epic follows its stories by kind: to the first queue state with the first story, to active work with the first started, to the wait state with the last open one, and to accepted when that one is accepted
- [ ] Tests run each of these against a workflow with a state renamed and a state added between `in-progress` and `review`, and give the same result as the baseline wherever the kinds agree
- [ ] With the baseline, every existing `workitem` test passes unchanged

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
- Cost of delay: not set yet, pending TH-1097, as on S-0736.
