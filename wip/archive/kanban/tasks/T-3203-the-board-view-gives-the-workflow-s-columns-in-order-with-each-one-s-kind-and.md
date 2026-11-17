---
id: T-3203
type: task
nature: feature
title: The board view gives the workflow's columns in order with each one's kind and WIP limit, with a shared test workflow that renames a state and adds one
status: done
parent: S-0739
owner: alex
created: 2026-11-12T16:06:23Z
updated: 2026-11-15T00:25:43Z
transitions:
  - to: ready
    at: 2026-11-14T23:56:42Z
    by: agent-S-0739
  - to: in-progress
    at: 2026-11-14T23:56:43Z
    by: agent-S-0739
  - to: done
    at: 2026-11-15T00:25:43Z
    by: agent-S-0739
stream: S-0739
tags: [cli]
touches: [flai/internal/workitem/boardview.go, flai/internal/workitem/boardview_test.go, flai/internal/workitem/board.go, flai/internal/workitem/workflow_fixture_test.go, flai/internal/workitem/flow.go, design/system/flai-cli.md]
usage:
  source: log
  seconds: 1740
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 82
      output: 32807
      cache_read: 4164605
      cache_write: 221722
      cost: 2.7929
---
# T-3203 The board view gives the workflow's columns in order with each one's kind and WIP limit, with a shared test workflow that renames a state and adds one

## Work

Add `workflow_fixture_test.go` with the test workflow the package's tests share: the baseline with `review` renamed and a state of its own between `in-progress` and `review`. Make `boardview.go` and `board.go` lay out the columns in the workflow's order and give each column its kind, WIP limit, and whether it is ordered, naming no state. Test the view on the baseline and on the fixture.

## Done when

The board view names no state, gives each column's kind and limit, and its tests pass on both workflows.

## Notes
