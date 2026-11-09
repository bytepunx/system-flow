---
id: S-0636
type: story
nature: improvement
title: 13 board.get calls at once answer under 300 ms on every served project, measured again after S-0602 and S-0606
status: backlog
parent: E-0033
owner: alex
created: 2026-11-09T07:36:47Z
updated: 2026-11-09T22:09:59Z
transitions: []
tags: [cli]
touches: [design/system/server-performance.md, design/issues/summary.md]
after: [S-0602, S-0606]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 30m
  delivery: 2026-11-10T00:12:00Z
  basis: "Its own forecast of 30m; 7th in the pull order with an in-progress limit of 5, behind S-0622, S-0623, S-0632, S-0634, S-0633, S-0625, S-0630, S-0628, S-0631 and S-0241."
  by: flai
  at: 2026-11-09T22:09:59Z
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

### Planning

Planned by planner-S-0636 on 2026-11-09. The run was asked to read a thread named "asked", which does not exist. S-0636 had no thread and no earlier plan, so this is a first plan.

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
- Delivery 2026-11-10T00:08:00Z: flai's 23:41 for 7th in the pull order with an in-progress limit of 5, moved by the 27m between the two durations.
- Cost of delay: no value yet. Neither S-0636 nor E-0033 has inputs, so `flai cod` refuses. The plan's thread on S-0636 asks the operator for an input.
