---
id: S-0723
type: story
nature: improvement
title: "flai check finds `item.blocked-closed` outside the story at close-out"
status: done
owner: alex
created: 2026-11-11T07:23:14Z
updated: 2026-11-13T07:23:40Z
transitions:
  - to: ready
    at: 2026-11-12T15:59:55Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-12T16:44:55Z
    by: agent-S-0723
  - to: review
    at: 2026-11-12T17:16:55Z
    by: agent-S-0723
  - to: done
    at: 2026-11-13T07:23:40Z
    by: alex
tags: []
touches: [flai/internal/workitem/rules.go, flai/internal/workitem/rules_test.go, flai/cmd/move.go, flai/cmd/move_test.go, flai/internal/check/check_items.go, flai/internal/check/items_test.go, design/system/workflow.md, docs/users/flai.md, design/issues/I-0347-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md, design/issues/summary.md, docs/users/flai-reference.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1932
  priced_by: harness
  turns:
    - day: 2026-11-12
      ceremony: 2
      work: 47
  models:
    - model: claude-opus-5-5
      input: 106
      output: 26099
      cache_read: 10619560
      cache_write: 250886
      cost: 4.6534
  strategic:
    - kind: planner
      seconds: 14
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 8
          output: 513
          cache_read: 899968
          cache_write: 7026
          cost: 0.2625
    - kind: orchestrator
      seconds: 73
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 21
          output: 193
          cache_read: 929158
          cache_write: 70743
          cost: 0.2851
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1m
    by: orchestrator
    at: 2026-11-12T15:57:45Z
  value: 2.5
  by: planner-S-0723
  at: 2026-11-12T15:58:13Z
forecast:
  duration: 22m
  delivery: 2026-11-12T16:49:00Z
  basis: "Its own forecast of 22m; 1st in the pull order with an in-progress limit of 3, behind S-0659, S-0683 and S-0709."
  by: flai
  at: 2026-11-12T16:04:53Z
finalized:
  by: alex
  at: 2026-11-12T07:06:11Z
---
# S-0723 flai check finds `item.blocked-closed` outside the story at close-out

## Goal

This story remediates [I-0347](../../../design/issues/I-0347-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md), "flai check finds `item.blocked-closed` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0347 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0347 is closed with `flai issue close I-0347 --reason` saying what fixed it

## Tasks
- T-3141 A move to done or cancelled closes the item's open blocked interval at the move
- T-3142 flai unblock closes an interval left open on a closed item at the item's close, archived items included, and item.blocked-closed names that fix
- T-3143 The design and the users guide say a closing move ends a block and flai unblock repairs a closed item, and I-0347 is closed

## Notes

### Planning

Cause, from I-0347's 32 instances: every one names the archived S-0678. It was blocked at 2026-11-11T05:26:16Z ("pause for urgent fix"), went to review at 06:13:48Z, and the orchestrator accepted it at 06:15:07Z with the interval still open. `Repo.Move` in `flai/internal/workitem/rules.go` neither refuses nor closes an open interval on a move to `done` or `cancelled`. So the item stays blocked in the archive, every close-out reports it outside its story, and `flai stats` counts its blocked time to now (`metrics.md` ends an open interval at now).

Proposed fix: a closing move closes the interval at the move and warns (T-3141). `flai unblock` on a closed item, archived ones included, closes one left open at the item's close, and `item.blocked-closed` names that command (T-3142). The design, the users guide, and the issue follow (T-3143). After acceptance, the operator runs `flai unblock S-0678` once in the main checkout.

Touches, all files, no folder:

| Touch | Source |
|-------|--------|
| `flai/internal/workitem/rules.go`, `rules_test.go` | layout: `Repo.Move` and `UnblockItem`, read in the code; co-change with the check and `move.go` |
| `flai/cmd/move.go`, `move_test.go` | layout: `flai unblock`; co-change with `rules.go` |
| `flai/internal/check/check_items.go`, `items_test.go` | layout: the `item.blocked-closed` rule |
| `design/system/workflow.md` | design: § Blocking; co-change 22% |
| `docs/users/flai.md` | design: § Blocking; co-change 33% |
| `design/issues/I-0347-...md`, `design/issues/summary.md` | goal and criterion 2: the issue closed, its summary regenerated |

Left out of the co-change list: `design/system/flai-cli.md` and `docs/users/flai-reference.md`, since no command or flag changes. `flai/cmd/accept.go` and `flai/internal/mcpserver/server.go` are left out too: they close items through `Repo.Move`, so they change nothing.

Forecast 22m stands as `flai forecast` gives it, at size 12 from 10 touches. Four of the touches are test files and two are the issue's own files, which take little. The three tasks run in series, as most of the stories behind the median do.

Cost of delay 2.50 USD a week stands as `flai cod` gives it. The input, `time_lost_per_cycle: 1m`, is the one recommended on TH-1096, set by the orchestrator. At the project's 150 USD an hour and one 168h cycle a week, it comes to 2.50 USD. I kept flai's cycle rather than counting each close-out as a cycle. The issue-made sibling stories, such as S-0709 and S-0730, are valued the same way, so the pull order compares like with like.

No `after`: the story reads no other story's output.
