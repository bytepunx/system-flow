---
id: S-0636
type: story
nature: improvement
title: 13 board.get calls at once answer under 300 ms on every served project, measured again after S-0602 and S-0606
status: backlog
parent: E-0033
owner: alex
created: 2026-11-09T07:36:47Z
updated: 2026-11-09T07:36:47Z
transitions: []
tags: [cli]
touches: [design/system/server-performance.md]
after: [S-0602, S-0606]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0636 13 board.get calls at once answer under 300 ms on every served project, measured again after S-0602 and S-0606

## Goal

S-0604's criterion 4 asked that 13 `board.get` calls at once on the operator's host, one per served project, each answer under 300 ms. S-0604 removed cause 14's git processes and shared lock, and the other 12 projects now answer under 300 ms. This repository's board still took 342 to 761 ms (`design/system/server-performance.md` § Cause 14 fixed). Its remaining time was its thread read with stories in review (cause 13, S-0602) and its item list and holds (cause 16, S-0606). The operator moved the criterion's original words here on TH-0912, to be measured again once those two land.

## Acceptance criteria

- [ ] 13 `board.get` calls at once on the operator's host, one per served project, each answer under 300 ms, and server-performance.md records the before and after.

## Tasks

## Notes

- `after: S-0602, S-0606`: this story measures what S-0602 (threads kept parsed, cause 13) and S-0606 (item snapshot, cause 16) deliver for this repository's board; it cannot be measured before both land.
- Measure as S-0604's T-2681 did: an in-process program calling `board.get` through `hostapi.MethodsFor` for the projects in `~/.flai/serve/projects.json`, with serve's main-refs watcher running, `main` against the change.
- Comes from S-0604 criterion 4, moved on TH-0912 (operator, 2026-11-09), as D-0465 did for S-0424.
