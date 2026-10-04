---
id: T-0821
type: task
nature: improvement
title: planning replays a story's delivery with its own forecast duration
status: done
parent: S-0211
owner: alex
created: 2026-10-04T20:41:53Z
updated: 2026-10-04T20:49:31Z
transitions:
  - to: ready
    at: 2026-10-04T20:42:47Z
    by: agent-S-0211
  - to: in-progress
    at: 2026-10-04T20:42:49Z
    by: agent-S-0211
  - to: done
    at: 2026-10-04T20:49:31Z
    by: agent-S-0211
stream: S-0211
tags: [flai]
touches: [flai/internal/planning]
usage:
  source: log
  seconds: 402
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 26
      output: 877
      cache_read: 919681
      cache_write: 56893
      cost: 0.3936
---
# T-0821 planning replays a story's delivery with its own forecast duration

## Work

In `flai/internal/planning`, a replay beside `Forecast`: the same play-out of the board, but a story that has its own `forecast.duration` is played out with it rather than the duration history gives (ADR-0084: the deterministic replan keeps the planner's duration and moves only the delivery and the basis). Its basis says the duration is the story's own forecast. A story with no forecast duration is played out as `Forecast` does. Share the play-out with `Forecast` rather than copying it. It waits for nothing: the first layer, its own package.

## Done when

- [ ] A replay of a story with a forecast duration keeps that duration and gives the delivery the board plays out
- [ ] `Forecast` prints what it printed before
- [ ] Behavior tests in `flai/internal/planning` pass

## Notes
