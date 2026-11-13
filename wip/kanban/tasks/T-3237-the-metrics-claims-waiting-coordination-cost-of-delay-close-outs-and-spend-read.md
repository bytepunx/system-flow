---
id: T-3237
type: task
nature: feature
title: The metrics' claims, waiting, coordination, cost of delay, close-outs, and spend read a state's kind instead of its name
status: backlog
parent: S-0750
owner: alex
created: 2026-11-12T16:08:37Z
updated: 2026-11-12T16:08:37Z
transitions: []
stream: S-0750
tags: [cli]
touches: [flai/internal/metrics/claims.go, flai/internal/metrics/waiting.go, flai/internal/metrics/coordination.go, flai/internal/metrics/costofdelay.go, flai/internal/metrics/costorder.go, flai/internal/metrics/closeouts.go, flai/internal/metrics/forecast.go, flai/internal/metrics/digest.go, flai/internal/metrics/usage.go, flai/internal/metrics/strategic.go, flai/internal/metrics/spend.go]
---
# T-3237 The metrics' claims, waiting, coordination, cost of delay, close-outs, and spend read a state's kind instead of its name

## Work

Replace the references to a state's name in these files, 11 in `claims.go` and 8 in `waiting.go` among them, with reads of kinds: claims are the active-work stories', waits are time in wait states, the cost of delay counts per column in the workflow's order, and completion is reaching an accepted state. Change no figure: the existing fixtures and `windows.golden.json` stay as they are.

## Done when

These files name no state, and every `metrics` test and the golden file pass unchanged.

## Notes
