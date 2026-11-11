---
id: S-0716
type: story
nature: remediation
title: One malformed work item file in the main checkout's wip makes every flai read of the project fail, a story's close-out and wait_for_events among them
status: ready
owner: alex
created: 2026-11-11T01:40:46Z
updated: 2026-11-11T08:58:19Z
transitions:
  - to: ready
    at: 2026-11-11T08:51:27Z
    by: alex
tags: [flai]
touches: [flai/internal/workitem/store.go, flai/internal/workitem/repo.go, flai/internal/workitem/store_test.go, flai/internal/workitem/narrative.go, flai/internal/workitem/narrative_test.go, flai/internal/check/check.go, flai/internal/check/check_items.go, flai/internal/check/scope.go, flai/internal/check/items_test.go, flai/internal/check/scope_test.go, flai/internal/mcpserver/inbox_test.go, flai/internal/mcpserver/stream_test.go, flai/internal/verify/story_test.go, design/system/work-hierarchy.md, docs/users/flai.md, design/issues/I-0344-one-malformed-work-item-file-in-the-main-checkout-s-wip-makes-every-flai-read.md, design/issues/summary.md]
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
    - kind: planner
      seconds: 69
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 18
          output: 4212
          cache_read: 1182011
          cache_write: 33986
          cost: 0.5926
cost_of_delay:
  inputs:
    time_lost_per_cycle: 3m
    by: flai
    at: 2026-11-11T01:40:46Z
  value: 7.5
  by: planner-S-0716
  at: 2026-11-11T08:57:45Z
forecast:
  duration: 1h30m
  delivery: 2026-11-11T12:37:00Z
  basis: "Its own forecast of 1h30m; 4th in the pull order with an in-progress limit of 3, behind S-0694, S-0698, S-0702, S-0706, S-0719 and S-0724."
  by: flai
  at: 2026-11-11T08:58:19Z
finalized:
  by: alex
  at: 2026-11-11T08:51:23Z
---
# S-0716 One malformed work item file in the main checkout's wip makes every flai read of the project fail, a story's close-out and wait_for_events among them

## Goal

This story remediates [I-0344](../../../design/issues/I-0344-one-malformed-work-item-file-in-the-main-checkout-s-wip-makes-every-flai-read.md), "One malformed work item file in the main checkout's wip makes every flai read of the project fail, a story's close-out and wait_for_events among them". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0344 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0344 is closed with `flai issue close I-0344 --reason` saying what fixed it

## Tasks
- T-3107 Repo.List skips a work item file that does not parse and keeps its path and error for the check
- T-3108 flai check reports a work item file that does not parse as one finding on its path and checks the rest
- T-3109 stream_state, wait_for_events, and inbox go on past a work item file that does not parse, with tests
- T-3110 Document how flai reads past a work item file that does not parse, and close I-0344

## Notes

Cost of delay inputs set by flai from I-0344. time_lost_per_cycle 3m: 3m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T11:05:08Z, 0.6 days before this story; under one cycle counts as one).

### Planning

The proposed fix, from the instance: `Repo.List` stops at the first file that does not parse. In `flai/internal/workitem/store.go`, `lookAll` returns `l.err` for the whole folder, and `Repo.List` in `repo.go` returns that error to all of its 90 or so callers: `flai check` (`check.go` `Run`), the close-out (`verify/story.go`), `task_done`, `wait_for_events` (`mcpserver/cursor.go` `catchUp`), `inbox`, and `stream_state`'s index write (`narrative.go` `WriteStreamState`, then `WriteIndex`). Fix it at the loader. `List` skips a file that does not parse and keeps its path and error, which callers can read (for example `Repo.Unreadable()`). `flai check` reports each such file as one finding on its path, and every other read goes on without it. The roles, conventions, and narrative loaders work this way already (`roles.Load`, `conventions.Load`, `ActiveNarratives`, which `WriteIndex` lists as unreadable). `NextID` globs file names, so an unreadable file's ID is still never given out again.

Touches. None were declared, and `flai touches suggest` had no path to start from. They come from:

- Layout, from the explorer's read of the code: `flai/internal/workitem/store.go`, `repo.go`, `store_test.go` (`TestAListOfADamagedItemFailsWithItsParseError` pins today's hard failure and changes with the fix), `narrative.go`, `narrative_test.go`, `flai/internal/check/check.go`, `check_items.go`, `scope.go`, `items_test.go`, `scope_test.go`, `flai/internal/mcpserver/inbox_test.go`, `stream_test.go`, and `flai/internal/verify/story_test.go`.
- Design: `design/system/work-hierarchy.md` § Front matter schema, which covers lenient reads but not a file that does not parse, and `docs/users/flai.md` § Check the repository, which lists the check's rules.
- The criteria: `design/issues/I-0344-…md` and `design/issues/summary.md`, which `flai issue close` writes.

No folder touch is kept. If the story's agent records the decision as an ADR, `flai task done` widens the touches to its file.

Forecast. With the 17 touches in place, flai gives 35m: a median 109 s per unit over 66 done remediation stories in the large band, times a size of 19. I raised it to 1h30m. That size counts touches, not the plan's three serial layers (T-3107, then T-3108 and T-3109, then T-3110). It also leaves out the close-out, which runs the integration and smoke tiers because the story changes `flai/`. Delivery is flai's 10:36Z plus the 55 minutes added.

Cost of delay. The value of 7.50 USD a week stands, as flai worked it out from the operator's input of 3m lost per 168h cycle at 150 USD an hour. Nothing in the issue argues for more. The instance cost 3m once.

No `after`: the story reads no other open story's output.
