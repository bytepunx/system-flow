---
id: T-0502
type: task
nature: improvement
title: The board remembers which work item types it shows, per browser, with stories on by default
status: done
parent: S-0141
owner: alex
created: 2026-09-28T23:00:47Z
updated: 2026-09-28T23:02:09Z
transitions:
  - to: ready
    at: 2026-09-28T23:00:55Z
    by: agent-S-0141
  - to: in-progress
    at: 2026-09-28T23:00:55Z
    by: agent-S-0141
  - to: done
    at: 2026-09-28T23:02:09Z
    by: agent-S-0141
stream: S-0141
tags: [dashboard]
touches: [flaiover/src/lib]
---
# T-0502 The board remembers which work item types it shows, per browser, with stories on by default

## Work

Add `flaiover/src/lib/boardtypes.svelte.ts`: the three types a board card can be (epic, story, task), whether each is shown, read from and written to `localStorage` under one key the way `theme.svelte.ts` keeps the theme. With nothing stored, or something stored that is not a valid choice, stories are shown and epics and tasks are not. A storage that throws (a private window) leaves the choice for the page. Behaviour tests with a fake storage.

## Done when

- A fresh browser shows stories only; a stored choice is read back as stored; a toggle is written at once.
- `pnpm test:unit` passes with the new tests.

## Notes
