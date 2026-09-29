---
id: S-0141
type: story
nature: improvement
title: Create checkboxes for displaying each work item type
status: done
parent: E-0003
owner: alex
created: 2026-09-27T04:06:55Z
updated: 2026-09-29T00:33:03Z
transitions:
  - to: ready
    at: 2026-09-28T22:56:39Z
    by: alex
  - to: in-progress
    at: 2026-09-28T23:00:14Z
    by: agent-S-0141
  - to: review
    at: 2026-09-28T23:07:45Z
    by: agent-S-0141
  - to: done
    at: 2026-09-29T00:33:03Z
    by: alex
tags: [dashboard]
touches: [flaiover/src, design/issues/I-0044-agents-flai-serve-starts-inherit-the-host-s-address-and-token-so-a-flai-serve-an-agent-runs-by-hand-takes-over-the-operator-s-mcp-servers.md, design/issues/summary.md, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0141 Create checkboxes for displaying each work item type

## Goal

Include a checkbox for each work item type in the board. Use local browser storage to remember the user's last configuration choices so it remains the same between page refreshes/navigation.

## Acceptance criteria
- [x] a set of checkboxes for tasks, stories, and epics are visible
- [x] toggling the checkbox for a particular work type toggles its visibility
- [x] stories are on by default until a user changes the toggle

## Tasks
- T-0502 The board remembers which work item types it shows, per browser, with stories on by default
- T-0503 The board has a checkbox for each of epics, stories, and tasks in place of epics and tasks too
- T-0504 Describe the type checkboxes, look at them in a browser, and pass every tier

## Notes

- The choice is kept per browser in `localStorage` (`flaiover-board-types`), not per project, as the theme is (`$lib/boardtypes.svelte.ts`). A stored value that is not a valid choice falls back to stories alone; a storage that refuses, as in a private window, keeps the choice for the page.
- WIP counts and reordering count stories whichever boxes are ticked, as before.
- Tried on 2026-09-28 (T-0504) in headless Chrome against a flaiover dev server built from `story/S-0141` on 127.0.0.1:5199 with `FLAIOVER_AUTH=off`, with `/api/board` answered in the browser by a stubbed board of one epic, two stories, and one task. No flai host or serve was started and the operator's dashboard was not touched; the dev server was stopped by PID.
  - A fresh browser (nothing stored) showed three checkboxes, epics, stories, and tasks, with stories alone ticked, and the two story cards only.
  - Ticking epics added E-0001, ticking tasks added T-0003, unticking stories removed S-0002 and S-0004; each click was stored at once. In progress read 1/2 throughout.
  - After a reload, and after going to /inbox and back by the Board link, the boxes and cards were as left.
