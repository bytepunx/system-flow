---
id: S-0750
type: story
nature: feature
title: flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today
status: ready
parent: E-0034
owner: alex
created: 2026-11-12T16:01:11Z
updated: 2026-11-15T07:18:01Z
transitions:
  - to: ready
    at: 2026-11-15T02:00:20Z
    by: orchestrator
tags: [cli]
topics: [cli, analysis]
touches: [design/adrs/README.md, design/system/metrics.md, flai/internal/workflow/timers.go, flai/internal/workflow/timers_test.go, flai/internal/workflow/baseline/board.md, flai/internal/metrics/metrics.go, flai/internal/metrics/metrics_test.go, flai/internal/metrics/claims.go, flai/internal/metrics/waiting.go, flai/internal/metrics/coordination.go, flai/internal/metrics/costofdelay.go, flai/internal/metrics/costorder.go, flai/internal/metrics/closeouts.go, flai/internal/metrics/forecast.go, flai/internal/metrics/digest.go, flai/internal/metrics/usage.go, flai/internal/metrics/strategic.go, flai/internal/metrics/spend.go]
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
      seconds: 20
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 31
          cache_read: 219655
          cache_write: 32346
          cost: 0.0733
cost_of_delay:
  value: 106.89
  by: system-flow
  at: 2026-11-14T17:44:11Z
forecast:
  duration: 1h
  delivery: 2026-11-15T08:56:00Z
  basis: "Its own forecast of 1h; 4th in the pull order with an in-progress limit of 3, behind S-0741, S-0749 and S-0748."
  by: flai
  at: 2026-11-15T07:18:01Z
finalized:
  by: orchestrator
  at: 2026-11-14T17:45:56Z
---
# S-0750 flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today

## Goal

Time tracking reads the workflow. The metrics definitions are the contract between `flai stats` and the dashboard and change only with an ADR, so this story writes that ADR: lead, cycle, queue, and review time become timers the baseline workflow names, each a reading of an item's recorded transitions (ADR-0004), and `flai stats` computes them from the file. Every figure the baseline reports today stays byte for byte, tested against the existing fixtures.

## Acceptance criteria

- [ ] An accepted ADR changes `metrics.md`'s definitions of lead, cycle, queue, and review time, and time in state, to readings of the workflow's timers, and says how a timer reads an item's transitions: its start and end points, or the sum of the time in the states it names, and how a reopened or cancelled item counts
- [ ] `flai/internal/workflow` evaluates a timer over an item's transitions, and the baseline workflow names lead, cycle, queue, and review time as the ADR gives them
- [ ] `flai/internal/metrics` computes each per item value and aggregate from the baseline's timers, and its other figures, claims, waiting, coordination, cost of delay, and close-outs included, read a state's kind instead of its name
- [ ] `flai stats --json` gives the same bytes as before on `flai/internal/metrics/testdata/good` and on `testdata/windows.golden.json`, which this story does not change, and flaiover's tests on those fixtures pass unchanged
- [ ] Tests compute the timers against a workflow with a state renamed and one added between `in-progress` and `review`, and show cycle time spanning the added state
- [ ] `design/system/metrics.md` defines each figure from its timer and links the ADR

## Tasks

- T-3236 Write and accept the metrics ADR that defines lead, cycle, queue, and review time and time in state as readings of the workflow's timers
- T-3237 The metrics' claims, waiting, coordination, cost of delay, close-outs, and spend read a state's kind instead of its name
- T-3241 The workflow package evaluates a timer over an item's transitions, and the baseline names lead, cycle, queue, and review time
- T-3242 metrics.md defines lead, cycle, queue, and review time and time in state from their timers and links the ADR
- T-3246 The metrics derive each item's durations and the aggregates from the baseline's timers, with flai stats byte for byte as before

## Notes

- Per item values are in `flai/internal/metrics/metrics.go`'s `Derive`; time in state, queue, lead, and cycle time are worked out there. `metrics.go` holds 14 of the package's references, `claims.go` 11, `waiting.go` 8.
- The ADR is this story's, as the epic's Notes ask; S-0736's ADR fixes only the timer's format in the workflow file.

### Planning

- Touches, all files:
  - `design/adrs/README.md`: co-change; `flai adr new` writes it. No folder touch on `design/adrs`.
  - `design/system/metrics.md`: design (criteria 1 and 6).
  - `flai/internal/workflow/timers.go`, `timers_test.go`, `baseline/board.md`: layout, beside S-0737's parser and baseline.
  - `flai/internal/metrics/*.go` named here: layout, each file of the package that names a state, and `metrics_test.go` for the byte-for-byte test. The fixtures under `testdata/` are not touched on purpose: criterion 4 holds them still.
  - `metrics/timers.go` was dropped from the first draft: a workflow's own timers are S-0752's file.
  - Left out from `flai touches suggest`: `flai-cli.md`, the guides, and `flaiover-dashboard.md` (7% to 22%). No figure a reader sees changes.
- `after`: S-0737. The timers are evaluated over the parser's `Timer` type and written into the baseline S-0737 embeds, and the metrics read kinds from the workflow it gives.
- Tasks: T-3236 (the ADR) and T-3237 (kinds in the other metrics) together, no path in common; T-3241 (timer evaluation) and T-3242 (`metrics.md`) after T-3236; T-3246 (`Derive` from timers) after T-3241.
- Forecast 1h, up from the 44m `flai forecast`'s rate gives for size 24. The story writes an ADR the operator answers, as S-0736 does, and the byte-for-byte rule on the fixtures leaves no room for a near miss in `Derive`.
- Cost of delay: 106.89 USD a week, as `flai cod` gives it: this story's share of E-0034's 1500.00 USD a week, the operator's `time_lost_per_cycle: 10h` (TH-1115) at 150 USD an hour, by its forecast of 1h over the 14h2m of the epic's 19 open stories. Kept: the epic's outcome needs every story, so a share by size is the split that fits.
