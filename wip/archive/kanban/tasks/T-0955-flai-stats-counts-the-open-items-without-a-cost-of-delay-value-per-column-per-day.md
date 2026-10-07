---
id: T-0955
type: task
nature: feature
title: flai stats counts the open items without a cost of delay value per column per day
status: done
parent: S-0213
owner: alex
created: 2026-10-05T05:46:25Z
updated: 2026-10-07T07:39:59Z
transitions:
  - to: ready
    at: 2026-10-07T07:37:56Z
    by: agent-S-0213
  - to: in-progress
    at: 2026-10-07T07:37:57Z
    by: agent-S-0213
  - to: done
    at: 2026-10-07T07:39:59Z
    by: agent-S-0213
stream: S-0213
tags: [flai]
touches: [flai/internal/metrics/costofdelay.go, flai/internal/metrics/costofdelay_test.go]
after: [T-0953]
usage:
  source: log
  seconds: 122
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 15
      output: 5804
      cache_read: 982985
      cache_write: 24491
      cost: 0.4597
---
# T-0955 flai stats counts the open items without a cost of delay value per column per day

## Work

Give `CostDay` in `flai/internal/metrics/costofdelay.go` the `without_value` map T-0953 defines. Fill it where `outstanding` is summed: per column, count the open items in it at the end of the day whose value is absent, every column present, with tasks left out as the definition says.

It waits for T-0953, whose definition it implements. T-0959 waits for it, because both change the cost of delay types in `costofdelay.go`.

## Done when

- fixture tests in `costofdelay_test.go` pin `without_value` per column across a day an item gains a value, and a day an item without one moves column
- every column is present on a day with none
- `go test ./internal/metrics/` passes

## Notes
