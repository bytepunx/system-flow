---
id: S-0847
type: story
nature: improvement
title: flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator
status: done
owner: alex
created: 2026-11-17T08:32:58Z
updated: 2026-11-27T23:38:24Z
transitions:
  - to: ready
    at: 2026-11-18T09:20:03Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-27T22:34:01Z
    by: agent-S-0847
  - to: review
    at: 2026-11-27T23:37:05Z
    by: agent-S-0847
  - to: done
    at: 2026-11-27T23:38:24Z
    by: orchestrator
blocked:
  - from: 2026-11-27T22:35:28Z
    until: 2026-11-27T22:57:15Z
    reason: awaiting the operator on TH-1323
tags: [flai, serve]
topics: [threads, ready, orchestration]
touches: [flai/internal/serve/starved.go, flai/internal/serve/starved_test.go, flai/internal/serve/readyempty_test.go, design/system/workflow.md, docs/users/flai.md, design/adrs/README.md, design/issues/I-0389-flai-serve-opens-a-new-board-starvation-thread-at-each-flip-of-the-board-48-in.md, design/issues/summary.md, design/adrs/0321-flai-serve-opens-writes-to-and-resolves-the-starved-board-s-thread-only-once.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2487
  priced_by: harness
  turns:
    - day: 2026-11-27
      ceremony: 1
      work: 50
  models:
    - model: claude-opus-5-5
      input: 172
      output: 55509
      cache_read: 14988309
      cache_write: 381724
      cost: 6.8999
  strategic:
    - kind: planner
      seconds: 5
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 6
          output: 18
          cache_read: 569949
          cache_write: 10963
          cost: 0.0397
        - model: claude-opus-5-5
          input: 9
          output: 835
          cache_read: 721386
          cache_write: 9307
          cost: 0.2541
    - kind: orchestrator
      seconds: 54
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 25
          output: 199
          cache_read: 744964
          cache_write: 130987
          cost: 0.2774
        - model: claude-sonnet-5-5
          input: 1
          output: 8
          cache_read: 7641
          cache_write: 5121
          cost: 0.0103
    - kind: analyzer
      seconds: 99
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 47
          output: 6584
          cache_read: 349360
          cache_write: 34306
          cost: 1.1031
        - model: claude-haiku-5-5
          input: 36
          output: 14843
          cache_read: 2398931
          cache_write: 47318
          cost: 0.1557
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-11-17T08:32:58Z
  value: 75
  by: planner-S-0847
  at: 2026-11-17T12:28:58Z
forecast:
  duration: 1h
  delivery: 2026-11-27T23:48:00Z
  basis: "Its own forecast of 1h; 1st in the pull order with an in-progress limit of 3, with nothing ahead of it."
  by: flai
  at: 2026-11-27T22:33:44Z
finalized:
  by: orchestrator
  at: 2026-11-17T12:35:34Z
---
# S-0847 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator

## Goal

This story remediates [I-0389](../../../design/issues/I-0389-flai-serve-opens-a-new-board-starvation-thread-at-each-flip-of-the-board-48-in.md), "flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator", which the analysis in [2026-11-16-all.md](../../../design/analysis/2026-11-16-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0389 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0389 is closed with `flai issue close I-0389 --reason` saying what fixed it

## Tasks
- T-3569 Propose a settle window for the starved board's thread and record it in an ADR refining ADR-0285
- T-3573 flai serve opens, resolves, and adds to the starved board's thread only once the board has held its state for the settle window
- T-3580 Close I-0389 with what fixed it

## Notes

The analyzer found I-0389: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0389. time_lost_per_cycle 30m carried over from I-0389's Impact section. Its Impact time_lost_per_cycle was taken rather than the 1m derived from its cost and count.

Strategic usage carried over from I-0389: analyzer 1m39s, 1.2588 USD.

### Planning

Touches, and where each came from:

| Touch | Source | Why |
| --- | --- | --- |
| `flai/internal/serve/starved.go` | layout | The starved board's thread is opened, written to, and resolved here (S-0691, ADR-0285). The story declared no touches. |
| `flai/internal/serve/starved_test.go` | co-change | `flai touches suggest` finds it changed with `starved.go` in S-0691's one commit. Its tests expect a thread at the first look and a new one at the next starving. |
| `flai/internal/serve/readyempty_test.go` | layout | `TestReadyEmptyAndTheStarvedBoardCountTheActiveWorkRoom` expects the starved thread at once. |
| `design/system/workflow.md` | design | § Cadence describes ADR-0285's thread and its resolve at the first look. |
| `docs/users/flai.md` | design | § flai serve describes the starved board's thread to the operator. |
| `design/adrs/README.md` | design | The ADR index flai generates for the ADR that refines ADR-0285 §4 and §5. `flai task done` widens T-3569 with the ADR's own file once `flai adr new` numbers it, so no folder touch is kept. |
| `design/issues/I-0389-…` and `design/issues/summary.md` | criteria | Criterion 2 closes I-0389 with `flai issue close`, which writes both. |

The story names no story in `after`: it reads nothing another open story writes. S-0845, from the same analysis's Finding 1, changes what holds the backlog, not the starved board's thread.

Figures:

- **Forecast.** `flai forecast` gave 21m: 124 s per unit of size over 42 done medium improvement stories on claude-opus-5-5, times size 10. It is raised to 1h. An ADR is proposed and written first, two existing tests change their expectations, and a close-out takes near 10m (S-0828's took 9m43s). The delivery moves from flai's 2026-11-17T20:48Z by the same 39m, to 21:27Z.
- **Cost of delay.** `flai cod` gave 75 USD a week from the input flai carried over from I-0389: 30m lost per 168h cycle at 150 USD an hour. It stands as given. The 48 threads in five days, a fifth of W46's threads to the operator, fit half an hour of reading a week.

The plan's assumptions are on the plan's thread: a fixed 15-minute settle window with no manifest key, and the `ready-empty` events out of scope.

### Accepted by the orchestrator

- Verified: 96a47c6e24f98bda908b0013058154d3fc0f14af
- At: 2026-11-27T23:38:24Z

Verdict: pass — both criteria are met by in-touches changes with a reproducing test, an ADR, and updated design and user docs, and the close-out passed at 96a47c6.
- 1: flai/internal/serve/starved.go, flai/internal/serve/starved_test.go, flai/internal/serve/readyempty_test.go, design/adrs/0321-flai-serve-opens-writes-to-and-resolves-the-starved-board-s-thread-only-once.md, design/system/workflow.md, docs/users/flai.md
- 2: design/issues/I-0389-flai-serve-opens-a-new-board-starvation-thread-at-each-flip-of-the-board-48-in.md, design/issues/summary.md
