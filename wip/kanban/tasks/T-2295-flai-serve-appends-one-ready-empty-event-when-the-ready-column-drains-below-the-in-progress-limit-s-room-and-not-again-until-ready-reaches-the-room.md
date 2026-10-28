---
id: T-2295
type: task
nature: feature
title: flai serve appends one ready-empty event when the ready column drains below the in-progress limit's room, and not again until ready reaches the room
status: in-progress
parent: S-0535
owner: alex
created: 2026-10-28T21:38:16Z
updated: 2026-10-28T22:45:29Z
transitions:
  - to: ready
    at: 2026-10-28T22:45:27Z
    by: system-flow
  - to: in-progress
    at: 2026-10-28T22:45:29Z
    by: system-flow
stream: S-0535
tags: [flai, events, serve, board]
touches: [flai/internal/serve/readyempty.go, flai/internal/serve/readyempty_test.go, flai/internal/serve/agents.go]
after: [T-2292]
---
# T-2295 flai serve appends one ready-empty event when the ready column drains below the in-progress limit's room, and not again until ready reaches the room

## Work

The launcher's `look()` in `flai/internal/serve/agents.go` reads the ready stories and the in-progress room (`readyStories`, `free`) each minute and on every change under the kanban folder. That is where the column draining is seen.

Criterion 2 and ADR-0211 §1 define the event: the number of ready stories, held ones counted, falls below the in-progress limit's room (the limit less the stories in progress), once per draining, and not again until the ready count reaches the room again. With no in-progress limit there is no room to fill and no event.

- Add `flai/internal/serve/readyempty.go`: per project, keep whether the last look found the ready count below the room. Append a `ready-empty` event through T-2292's log on the look that finds it below after one that did not. Give the ready count, the room, and the limit.
- Keep the flag across a restart of flai serve, in the project's `.flai-cache`, so that a restart with ready still below emits nothing.
- Call it from `look()` in `agents.go`, before any early return on an empty ready column.

It waits for T-2292, whose log it appends to. It shares no path with T-2294, T-2296, or T-2297, so the four run together.

## Done when

- Tests show: one event when ready goes from at least the room to below it; none while it stays below over several looks; none with the in-progress limit full (room zero) or with no limit; one again after ready reaches the room and drains again.
- A test restarts the launcher with ready still below and reads no second event.
- `flai test` on the changed paths passes.

## Notes

ADR-0209 §1 starts the orchestrator on `ready-empty` to promote stories into the room. The planner's draft fired only on an empty column; S-0535's agent took criterion 2's wording, "below the in-progress limit's room", and recorded it in ADR-0211.
