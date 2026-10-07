---
id: T-0965
type: task
nature: feature
title: flai stats prints the pull order's projected saving and the items without a cost of delay value
status: done
parent: S-0213
owner: alex
created: 2026-10-05T05:47:06Z
updated: 2026-10-07T07:59:43Z
transitions:
  - to: ready
    at: 2026-10-07T07:51:09Z
    by: agent-S-0213
  - to: in-progress
    at: 2026-10-07T07:51:10Z
    by: agent-S-0213
  - to: done
    at: 2026-10-07T07:59:43Z
    by: agent-S-0213
stream: S-0213
tags: [flai]
touches: [flai/cmd/stats.go, flai/cmd/check_stats_test.go, docs/users/flai.md]
after: [T-0959]
usage:
  source: log
  seconds: 513
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 22
      output: 8754
      cache_read: 1482523
      cache_write: 36937
      cost: 0.6933
---
# T-0965 flai stats prints the pull order's projected saving and the items without a cost of delay value

## Work

`flai stats` prints the cost of delay outstanding now and incurred over the window (`flai/cmd/stats.go`). Under them, add two lines, each only when there is something to show:

- the items without a value now, per column
- the ready column's projected cost under the pull order, by cost of delay, and by WSJF, with the saving and the order that gives it, and how many ready stories were left out

Describe both lines, and the `without_value` and `order` keys under `cost_of_delay` in `--json`, in the `flai stats` section of `docs/users/flai.md`, linking `metrics.md`.

It waits for T-0959, which computes what it prints. It runs with the chart builders' task, whose paths it does not share.

## Done when

- a test in `check_stats_test.go` pins both lines on a fixture, and their absence when there is nothing to show
- `docs/users/flai.md` describes them
- `go test ./cmd/` and the markdown lint pass

## Notes
