---
id: T-0953
type: task
nature: feature
title: metrics.md defines the pull order's projected cost of delay and the count of items without a value, and an ADR records them
status: done
parent: S-0213
owner: alex
created: 2026-10-05T05:46:16Z
updated: 2026-10-07T07:37:41Z
transitions:
  - to: ready
    at: 2026-10-07T07:32:46Z
    by: agent-S-0213
  - to: in-progress
    at: 2026-10-07T07:32:47Z
    by: agent-S-0213
  - to: done
    at: 2026-10-07T07:37:41Z
    by: agent-S-0213
stream: S-0213
tags: [flai]
touches: [design/system/metrics.md, design/adrs]
usage:
  source: log
  seconds: 294
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 42
      output: 16700
      cache_read: 2828273
      cache_write: 70466
      cost: 1.3226
---
# T-0953 metrics.md defines the pull order's projected cost of delay and the count of items without a value, and an ADR records them

## Work

`metrics.md` is the contract between `flai stats` and the dashboard, so the two values the charts need that S-0205 left out are defined there first, with an ADR. Add them under `Planning, waiting, and claims (S-0205) › Cost of delay`:

- `cost_of_delay.days[].without_value`: per column, `backlog`, `ready`, `in-progress`, and `review`, the count of open stories and epics in it at the end of the day that have no `cost_of_delay.value`, every column present. Name whether tasks count; they carry no value, so they should not.
- `cost_of_delay.order`: for the ready column now, the cost of delay its stories will incur until each is pulled, projected under three orders: `current` (the board's pull order), `cod` (value, highest first), and `wsjf` (value over forecast duration in hours, highest first), as S-0217's `flai order --by` orders them. Define:
  - the lanes: the board's in-progress limit, each freed when a story in progress reaches its forecast duration, as `flai forecast` lays the pull order out (`planning.playOut`)
  - the series: for each order, points of cumulative incurred cost from now to the horizon, the time the last ready story is pulled under the slowest order, at each pull
  - each story's incurred cost, its value times the seconds from now to its projected pull over 604800, consistent with `cost_of_delay_incurred`
  - `total` per order, and `saving`: `current` minus the lower of `cod` and `wsjf`, with which one
  - `left_out`: the ready stories with no value or no forecast duration, by ID, which no order places
- The rounding of each amount, as the section's other values are rounded.

Write the ADR from the template. It records the additions, the choice to project waiting until the pull rather than until delivery, and that the orders are the policy orders of S-0217. Link it from the section.

This task waits for nothing: it is the contract the Go tasks implement.

## Done when

- `metrics.md` defines `without_value` and `order` with their keys, units, and rounding, and links the ADR
- the ADR is proposed in `design/adrs` from the template
- `flai check --strict` and the markdown lint pass

## Notes
