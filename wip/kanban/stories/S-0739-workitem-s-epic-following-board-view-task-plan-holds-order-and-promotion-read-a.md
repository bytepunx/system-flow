---
id: S-0739
type: story
nature: feature
title: workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T15:59:07Z
updated: 2026-11-12T15:59:07Z
transitions: []
tags: [cli]
topics: [cli, planning]
touches: [flai/internal/workitem/follow.go, flai/internal/workitem/follow_test.go, flai/internal/workitem/boardview.go, flai/internal/workitem/boardview_test.go, flai/internal/workitem/board.go, flai/internal/workitem/taskplan.go, flai/internal/workitem/taskplan_test.go, flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/workitem/order.go, flai/internal/workitem/plancandidates.go, flai/internal/workitem/promote.go, flai/internal/workitem/promotable.go, flai/internal/workitem/finalizable.go, flai/internal/workitem/changes.go, flai/internal/workitem/share.go, flai/internal/workitem/policy.go, flai/internal/workitem/narrative.go, flai/internal/workitem/item.go]
after: [S-0737]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
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

## Notes

- These files hold about 80 of the package's references to a state's name (`follow.go` 18, `taskplan.go` 11, `boardview.go` 10, `plancandidates.go` 8, `hold.go` 7, `order.go` 7). `item.go` keeps the constants for the baseline's names.
