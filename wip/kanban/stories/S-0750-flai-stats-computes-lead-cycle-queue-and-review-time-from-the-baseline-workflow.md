---
id: S-0750
type: story
nature: feature
title: flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:01:11Z
updated: 2026-11-12T16:01:11Z
transitions: []
tags: [cli]
topics: [cli, analysis]
touches: [design/adrs/README.md, design/system/metrics.md, flai/internal/workflow/timers.go, flai/internal/workflow/timers_test.go, flai/internal/workflow/baseline/board.md, flai/internal/metrics/metrics.go, flai/internal/metrics/metrics_test.go, flai/internal/metrics/timers.go, flai/internal/metrics/timers_test.go, flai/internal/metrics/claims.go, flai/internal/metrics/waiting.go, flai/internal/metrics/coordination.go, flai/internal/metrics/costofdelay.go, flai/internal/metrics/costorder.go, flai/internal/metrics/closeouts.go, flai/internal/metrics/forecast.go, flai/internal/metrics/digest.go, flai/internal/metrics/usage.go, flai/internal/metrics/strategic.go, flai/internal/metrics/spend.go]
after: [S-0737]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
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

## Notes

- Per item values are in `flai/internal/metrics/metrics.go`'s `Derive`; time in state, queue, lead, and cycle time are worked out there. `metrics.go` holds 14 of the package's references, `claims.go` 11, `waiting.go` 8.
- The ADR is this story's, as the epic's Notes ask; S-0736's ADR fixes only the timer's format in the workflow file.
