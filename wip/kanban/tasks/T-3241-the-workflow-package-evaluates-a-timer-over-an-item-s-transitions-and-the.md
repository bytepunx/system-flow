---
id: T-3241
type: task
nature: feature
title: The workflow package evaluates a timer over an item's transitions, and the baseline names lead, cycle, queue, and review time
status: backlog
parent: S-0750
owner: alex
created: 2026-11-12T16:08:58Z
updated: 2026-11-12T16:08:58Z
transitions: []
stream: S-0750
tags: [cli]
touches: [flai/internal/workflow/timers.go, flai/internal/workflow/timers_test.go, flai/internal/workflow/baseline/board.md]
after: [T-3236]
---
# T-3241 The workflow package evaluates a timer over an item's transitions, and the baseline names lead, cycle, queue, and review time

## Work

Waits for T-3236: it builds the timers the metrics ADR defines. Add `timers.go`: given an item's transitions and a timer, the seconds from its start point to its end point, or the sum of the intervals in the states it names, counting a reopened and a cancelled item as the ADR says. Give the baseline file lead, cycle, queue, and review time and time in state as timers. Test each on transitions taken from `flai/internal/metrics/testdata/good`.

## Done when

`timers_test.go` gives each baseline timer the value today's `Derive` gives on the fixture's items.

## Notes
