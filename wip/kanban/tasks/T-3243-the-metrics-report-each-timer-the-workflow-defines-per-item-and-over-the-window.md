---
id: T-3243
type: task
nature: feature
title: The metrics report each timer the workflow defines, per item and over the window, past items included
status: backlog
parent: S-0752
owner: alex
created: 2026-11-12T16:09:04Z
updated: 2026-11-12T16:09:04Z
transitions: []
stream: S-0752
tags: [cli]
touches: [flai/internal/metrics/timers.go, flai/internal/metrics/timers_test.go, flai/internal/metrics/metrics.go, flai/internal/metrics/window.go]
---
# T-3243 The metrics report each timer the workflow defines, per item and over the window, past items included

## Work

Add to `flai/internal/metrics` a section, named as the timers' ADR says, that gives every timer the workflow defines beyond the baseline's: per item its seconds, and over the window its count, mean, p50, and p85, as the baseline's durations are given. Compute it from the transitions each item carries, so items completed before the timer was added are reported. Test a timer from the first move to `review` to the acceptance, and a sum of the time in two states, on fixture items with known transitions.

## Done when

`timers_test.go` gives both timers' values, and the baseline's fixtures and golden file are unchanged.

## Notes
