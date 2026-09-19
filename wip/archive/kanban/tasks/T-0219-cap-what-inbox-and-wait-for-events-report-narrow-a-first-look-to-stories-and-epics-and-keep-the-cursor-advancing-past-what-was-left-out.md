---
id: T-0219
type: task
nature: feature
title: Cap what inbox and wait_for_events report, narrow a first look to stories and epics, and keep the cursor advancing past what was left out
status: done
parent: S-0061
owner: alex
created: 2026-09-19T08:15:26Z
updated: 2026-09-19T08:19:17Z
transitions:
  - to: ready
    at: 2026-09-19T08:15:27Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T08:15:27Z
    by: system-flow
  - to: done
    at: 2026-09-19T08:19:17Z
    by: system-flow
stream: S-0061
tags: []
---

# T-0219 Cap what inbox and wait_for_events report, narrow a first look to stories and epics, and keep the cursor advancing past what was left out

## Work
In `flai/internal/mcpserver/cursor.go`, `catchUp` keeps the newest changes up to a constant cap and returns how many older ones it left out; the constant carries its reason (the 209-change first look was 68 KB and did not fit a tool result). With no cursor it reports story and epic changes only, not task transitions. The cursor advances past everything, reported, filtered, or omitted, including the keys of changes stamped with the current second, so nothing left out comes back. `inbox` gains `changes_omitted` and `wait_for_events` gains `events_omitted`. The pull-order event is not counted against the cap. Ready work and threads are untouched. Tests: more changes than the cap with and without a cursor, a first look with task transitions, and a second look that is empty.

## Done when
- The tests pass and fail without the change
- `make flai-test` passes

## Notes
