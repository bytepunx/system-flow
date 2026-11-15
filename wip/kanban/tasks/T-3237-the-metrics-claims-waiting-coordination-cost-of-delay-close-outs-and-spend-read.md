---
id: T-3237
type: task
nature: feature
title: The metrics' claims, waiting, coordination, cost of delay, close-outs, and spend read a state's kind instead of its name
status: done
parent: S-0750
owner: alex
created: 2026-11-12T16:08:37Z
updated: 2026-11-15T20:56:18Z
transitions:
  - to: ready
    at: 2026-11-15T20:32:05Z
    by: agent-S-0750
  - to: in-progress
    at: 2026-11-15T20:32:05Z
    by: agent-S-0750
  - to: done
    at: 2026-11-15T20:56:18Z
    by: agent-S-0750
stream: S-0750
tags: [cli]
touches: [flai/internal/metrics/claims.go, flai/internal/metrics/waiting.go, flai/internal/metrics/coordination.go, flai/internal/metrics/costofdelay.go, flai/internal/metrics/costorder.go, flai/internal/metrics/forecast.go, flai/internal/metrics/digest.go, flai/internal/metrics/usage.go, flai/internal/metrics/strategic.go, flai/internal/metrics/spend.go, flai/internal/metrics/metrics.go]
usage:
  source: log
  seconds: 1443
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 94
      output: 525
      cache_read: 5524627
      cache_write: 317433
      cost: 2.2083
---
# T-3237 The metrics' claims, waiting, coordination, cost of delay, close-outs, and spend read a state's kind instead of its name

## Work

Replace the references to a state's name in these files, 11 in `claims.go` and 8 in `waiting.go` among them, with reads of kinds: claims are the active-work stories', waits are time in wait states, the cost of delay counts per column in the workflow's order, and completion is reaching an accepted state. Change no figure: the existing fixtures and `windows.golden.json` stay as they are.

## Done when

These files name no state, and every `metrics` test and the golden file pass unchanged.

## Notes
