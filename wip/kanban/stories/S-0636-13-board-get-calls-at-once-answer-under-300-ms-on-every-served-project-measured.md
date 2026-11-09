---
id: S-0636
type: story
nature: improvement
title: 13 board.get calls at once answer under 300 ms on every served project, measured again after S-0602 and S-0606
status: in-progress
parent: E-0033
owner: alex
created: 2026-11-09T07:36:47Z
updated: 2026-11-09T23:11:22Z
transitions:
  - to: ready
    at: 2026-11-09T22:33:00Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-09T22:59:24Z
    by: agent-S-0636
blocked:
  - from: 2026-11-09T23:11:22Z
    reason: awaiting the operator on TH-0961
tags: [cli]
touches: [design/system/server-performance.md, design/issues/summary.md, design/issues/I-0326-this-repository-s-board-get-takes-300-ms-or-more-in-half-its-calls-with-13.md]
after: [S-0602, S-0606]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 718
  priced_by: harness
  turns:
    - day: 2026-11-09
      ceremony: 2
      hand_edits: 3
      work: 26
  models:
    - model: claude-opus-5-5
      input: 126
      output: 34682
      cache_read: 16263090
      cache_write: 346373
      cost: 6.5877
  strategic:
    - kind: planner
      seconds: 138
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 20
          output: 760
          cache_read: 543470
          cache_write: 57646
          cost: 0.2078
    - kind: orchestrator
      seconds: 358
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 31
          output: 363
          cache_read: 1642278
          cache_write: 57451
          cost: 0.453
        - model: claude-sonnet-5-5
          input: 1
          output: 11
          cache_read: 14233
          cache_write: 8148
          cost: 0.0166
cost_of_delay:
  inputs:
    penalty_per_week: 20
    by: orchestrator
    at: 2026-11-09T22:14:16Z
  value: 20
  by: planner-S-0636
  at: 2026-11-09T22:21:32Z
forecast:
  duration: 30m
  delivery: 2026-11-09T23:38:00Z
  basis: "Its own forecast of 30m; 4th in the pull order with an in-progress limit of 5, behind S-0622, S-0623, S-0632, S-0633, S-0634, S-0645, S-0644 and S-0638."
  by: flai
  at: 2026-11-09T22:33:17Z
finalized:
  by: orchestrator
  at: 2026-11-09T22:30:03Z
---
# S-0636 13 board.get calls at once answer under 300 ms on every served project, measured again after S-0602 and S-0606

## Goal

S-0604's criterion 4 asked that 13 `board.get` calls at once on the operator's host, one per served project, each answer under 300 ms. S-0604 removed cause 14's git processes and shared lock, and the other 12 projects now answer under 300 ms. This repository's board still took 342 to 761 ms (`design/system/server-performance.md` § Cause 14 fixed). Its remaining time was its thread read with stories in review (cause 13, S-0602) and its item list and holds (cause 16, S-0606). The operator moved the criterion's original words here on TH-0912, to be measured again once those two land.

## Acceptance criteria

- [ ] 13 `board.get` calls at once on the operator's host, one per served project, each answer under 300 ms, and server-performance.md records the before and after.

## Tasks
- T-2803 Measure 13 board.get calls at once before and after S-0602 and S-0606 on the operator's host, and record them under a new section of server-performance.md
- T-2804 Where a board call still takes 300 ms or more, file the shortfall as an issue and ask the operator on a thread how the criterion stands

## Notes

- `after: S-0602, S-0606`: this story measures what S-0602 (threads kept parsed, cause 13) and S-0606 (item snapshot, cause 16) deliver for this repository's board; it cannot be measured before both land.
- Measure as S-0604's T-2681 did: an in-process program calling `board.get` through `hostapi.MethodsFor` for the projects in `~/.flai/serve/projects.json`, with serve's main-refs watcher running, `main` against the change.
- Comes from S-0604 criterion 4, moved on TH-0912 (operator, 2026-11-09), as D-0465 did for S-0424.
- Measured on 2026-11-09 (T-2803, `design/system/server-performance.md` § 13 boards at once, measured again after S-0602 and S-0606 (S-0636)): the other 12 projects answered under 300 ms in all 120 calls; this repository took 300 ms or more in 5 of 10, from the merge queue's git while stories are in review. The shortfall is recorded as [I-0326](../../../design/issues/I-0326-this-repository-s-board-get-takes-300-ms-or-more-in-half-its-calls-with-13.md), and the operator is asked how criterion 1 stands on TH-0961 (D-0532, D-0535).

### Planning

Planned by planner-S-0636 on 2026-11-09. The plan's thread is TH-0954, where the orchestrator approved it and set the cost of delay input. A second run on 2026-11-09 wrote the cost of delay value and changed nothing else.

Tasks, in two layers:

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-2803 measure 13 `board.get` calls at once before and after, and record them | `design/system/server-performance.md` |
| 2 | T-2804 where a call still takes 300 ms or more, file the shortfall as an issue and ask on a thread (after T-2803) | `design/issues/summary.md` |

Touches, by where each came from:

| Touch | From |
|-------|------|
| `design/system/server-performance.md` | declared; the criterion records the before and after there |
| `design/issues/summary.md` | design: D-0532 and D-0535 send an absolute target a story cannot meet to an issue, and `flai issue new` updates the summary. The issue's own file is named when T-2804 writes it |

No folder touch. `flai touches suggest` also ranks `design/system/flai-cli.md` (36%), `design/system/flaiover-dashboard.md`, `docs/users/flai.md`, `design/tech/go-libraries.md`, and four more docs. They are left out because the story measures and changes no code, command, setting, or library.

`after`: S-0602 and S-0606, both done. This story reads what they wrote: S-0602's kept thread store and S-0606's item snapshot, on `main`, and their sections of `server-performance.md`. The measure is of their effect on this repository's board.

Figures:

- Forecast 30m, against flai's 3m. flai's figure is 74 s a unit of size times size 2, from one criterion and one touch. That size misses the work: two builds in the same hour, a timing program, four runs of 65 calls, and a section of figures. S-0604's T-2681 did the same as part of a larger story.
- Delivery: flai's replan from the 30m duration, 2026-11-10T00:28:00Z at 10th in the pull order with an in-progress limit of 5. It stands: it already counts the 30m.
- Cost of delay 20 USD a week, as `flai cod S-0636` works it out from the one input, `penalty_per_week: 20`, which the orchestrator set on TH-0954. It stands: it is the bottom of E-0033's scale (S-0605, S-0607), since the story changes no code and S-0602 and S-0606 have landed.
