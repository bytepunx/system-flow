---
id: S-0479
type: story
nature: remediation
title: A story's in-progress transition from flai story start was lost when a concurrent write rewrote the item file, leaving it ready while its agent worked it
status: in-progress
owner: alex
created: 2026-10-21T11:59:43Z
updated: 2026-10-21T13:12:46Z
transitions:
  - to: ready
    at: 2026-10-21T12:05:53Z
    by: alex
  - to: in-progress
    at: 2026-10-21T12:09:40Z
    by: agent-S-0479
tags: []
touches: [flai/internal/workitem/itemlock.go, flai/internal/workitem/itemlock_test.go, flai/internal/workitem/repo.go, flai/internal/itemedit/itemedit.go, flai/internal/itemedit/lock_test.go, design/system/work-hierarchy.md, design/adrs, design/issues/I-0175-a-story-s-in-progress-transition-from-flai-story-start-was-lost-when-a-concurrent-write-rewrote-the-item-file-leaving-it-ready-while-its-agent-worked-it.md, design/issues/summary.md, flai/cmd/accept_resume_test.go, flai/cmd/check_stats_test.go, flai/cmd/issue_test.go, flai/internal/workitem/item.go, design/issues/I-0176-a-usage-write-to-a-story-overwrites-its-move-to-in-progress-made-meanwhile-so-a-started-story-reads-as-ready-again.md, flai/internal/itemnew/itemnew_test.go, flai/internal/statsread/statsread_test.go, flai/internal/mcpserver/folder_test.go, design/issues/I-0186-testwaitforeventsacrossafolder-can-wake-on-the-new-story-s-creation-before-its-move-to-ready-and-finds-no-event.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2347
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-10-21
      ceremony: 5
      hand_edits: 6
      work: 49
  models:
    - model: claude-opus-5-5
      input: 264
      output: 1330
      cache_read: 16314079
      cache_write: 568976
      cost: 5.7203
  strategic:
    - kind: planner
      seconds: 349
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 110
          output: 347
          cache_read: 5095111
          cache_write: 225207
          cost: 0.3429
        - model: claude-opus-5-5
          input: 280
          output: 36889
          cache_read: 14604844
          cache_write: 513779
          cost: 6.1419
    - kind: orchestrator
      seconds: 32
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 10
          output: 121
          cache_read: 2106278
          cache_write: 9094
          cost: 0.4671
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-21T11:59:43Z
  value: 12.5
  by: planner-S-0479
  at: 2026-10-21T12:11:58Z
forecast:
  duration: 2h
  delivery: 2026-10-21T14:10:00Z
  basis: "flai forecast's 33m (89 s per unit over 33 large remediation stories, size 22) raised to 2h: S-0377, a lock fix of three writers on an existing pattern, took 37m, and this one adds a lock primitive, moves about a dozen writers onto it, restructures itemedit.Apply, and adds race tests and design."
  by: planner-S-0479
  at: 2026-10-21T12:11:58Z
finalized:
  by: alex
  at: 2026-10-21T12:04:48Z
---
# S-0479 A story's in-progress transition from flai story start was lost when a concurrent write rewrote the item file, leaving it ready while its agent worked it

## Goal

This story remediates [I-0175](../../../design/issues/I-0175-a-story-s-in-progress-transition-from-flai-story-start-was-lost-when-a-concurrent-write-rewrote-the-item-file-leaving-it-ready-while-its-agent-worked-it.md), "A story's in-progress transition from flai story start was lost when a concurrent write rewrote the item file, leaving it ready while its agent worked it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0175 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0175 is closed with `flai issue close I-0175 --reason` saying what fixed it
- [x] I-0176 is closed with `flai issue close I-0176 --reason` saying what fixed it

## Tasks
- T-1999 A work item write takes a project-wide items lock and reads the item again under it, with a test that reproduces the lost update
- T-2001 Moves, story start, blocks, acceptance, creation, and usage roll-ups write their items under the items lock from a fresh read
- T-2002 flai serve's usage writer, flai edit and item_edit, touches, and widening write their items under the items lock from a fresh read
- T-2005 The design records that every work item write takes the items lock, and I-0175 is closed with what fixed it
- T-2007 Repo.Save merges an item copy read before a concurrent write with the file on disk, key by key, under one lock, and refuses a key both changed
- T-2008 An item edit refused by the check puts back only its own change to an item file, keeping what was written meanwhile
- T-2009 The design and an ADR say how concurrent item writes merge, and I-0175 and I-0176 are closed

## Notes

Cost of delay inputs set by flai from I-0175. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-21T07:37:23Z, 0.2 days before this story; under one cycle counts as one).

### Planning

Proposed cause: no lock guards writes to work item files. `Repo.Save` in `flai/internal/workitem/repo.go` writes a whole file from a copy the caller read earlier. `flai serve` (the usage writer, and the replanner's forecast through `itemedit.Apply`) and the MCP server (`story_start`) run in different processes. A writer that read before the move and saved after it put `status: ready` back. `itemedit.Apply` keeps the widest window, because it runs `check.Run` between its read and its save. Proposed fix: one project-wide items lock under `.flai-cache`, taken with `filelock` as `mergequeue.Update` does, with every writer reading the item again under it.

Touches. `flai touches suggest` listed nothing, because the story declared no touches. Every touch comes from the code layout, read by the explorer:

| Touch | Source |
|-------|--------|
| `flai/internal/workitem/itemlock.go`, `itemlock_test.go` | layout: new lock primitive beside `repo.go` |
| `flai/internal/workitem/repo.go` | layout: `Save` |
| `flai/internal/workitem/transition.go`, `create.go`, `usage.go` | layout: unlocked saves in `TransitionAll`, create, `RollUp`, `ChargeStrategic` |
| `flai/internal/storystart/start.go`, `start_test.go` | layout: story start loads the story before the move (I-0175) |
| `flai/cmd/move.go`, `accept.go`, `touches.go` | layout: unlocked saves in block and unblock, accept, and touches |
| `flai/internal/serve/usage.go`, `usage_test.go` | layout: the usage writer (I-0176) |
| `flai/internal/itemedit/itemedit.go`, `widen.go`, `lock_test.go` | layout: `Apply` and `Widen` |
| `design/system/work-hierarchy.md` | design: who changes what after an item is made |
| `design/adrs` | design: a possible ADR on the lock |
| I-0175's file, `design/issues/summary.md` | goal: the issue is closed |

The one folder touch kept is `design/adrs`, because the ADR's number is not known until it is written. It is inside `claims.shared`, so it holds no ready story.

Figures:

- Forecast: 2h, raised from flai's 33m. S-0377 fixed a lock for three writers on an existing pattern in 37m. This story adds a lock primitive, moves about a dozen writers onto it, restructures `itemedit.Apply`, and adds race tests and design. Delivery is 2h after the 12:09:40Z start.
- Cost of delay: 12.50 USD a week, as `flai cod` gives from flai's input of 5m lost per 168h cycle. It stands. It counts I-0175 alone. With I-0176 (10m), the same cause would cost 37.50 USD a week.
