---
id: T-0503
type: task
nature: improvement
title: The board has a checkbox for each of epics, stories, and tasks in place of epics and tasks too
status: done
parent: S-0141
owner: alex
created: 2026-09-28T23:00:47Z
updated: 2026-09-28T23:03:46Z
transitions:
  - to: ready
    at: 2026-09-28T23:00:55Z
    by: agent-S-0141
  - to: in-progress
    at: 2026-09-28T23:02:09Z
    by: agent-S-0141
  - to: done
    at: 2026-09-28T23:03:46Z
    by: agent-S-0141
stream: S-0141
tags: [dashboard]
touches: [flaiover/src/lib/components, flaiover/src/routes/board]
---
# T-0503 The board has a checkbox for each of epics, stories, and tasks in place of epics and tasks too

## Work

Add `BoardTypes.svelte`: one checkbox each for epics, stories, and tasks, bound to the state from the first task. Put it on `/board` in place of the single "epics and tasks too" checkbox, and filter the columns' cards by it. WIP counts and reordering still count stories whether shown or not. Component test: three checkboxes, stories alone ticked by default, a click toggles the type and is stored.

## Done when

- The board shows the three checkboxes and hides or shows each type as it is ticked.
- `make flaiover-test` passes (prettier, eslint, svelte-check, vitest).

## Notes
