---
id: S-0150
type: story
nature: improvement
title: Agent status and updates on board need to be collapsible
status: done
parent: E-0013
owner: alex
created: 2026-09-29T05:31:24Z
updated: 2026-09-29T19:10:09Z
transitions:
  - to: ready
    at: 2026-09-29T06:44:24Z
    by: alex
  - to: in-progress
    at: 2026-09-29T19:04:53Z
    by: agent-S-0150
  - to: review
    at: 2026-09-29T19:09:19Z
    by: agent-S-0150
  - to: done
    at: 2026-09-29T19:10:09Z
    by: alex
tags: [dashboard]
topics: [front-end, board]
touches: [flaiover/src, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 290
  models:
    - model: claude-opus-5-5
      input: 76
      output: 19667
      cache_read: 3020278
      cache_write: 100107
      cost: 1.7986
---
# S-0150 Agent status and updates on board need to be collapsible

## Goal

The Agent notification/summary box on the board page contains valuable information but in some cases can get quite long. Allow the user to collapse the area using a toggle carrot so that only the bold title appears when collapsed.

## Acceptance criteria
- [x] The operator can collapse the box using the toggle
- [x] The operator can expand the box using the toggle
- [x] Changes to the status cause the box to auto-expand

## Tasks
- T-0557 The board's agent notice collapses to its bold title and expands again when its status changes
- T-0558 The dashboard design and user guide describe the collapsible agent notice

## Notes

- The box is `HostAgentNotice.svelte`. Every state now has a bold title (`Agent started`, `No agent could be started`, `Agent ended`, `A ready story is waiting`); the last two had none, so collapsing them would have left nothing to read.
- The collapse is kept in localStorage against what the notice says (command, running, last, waiting), so a reload or the board's 15 s re-ask that answers the same leaves it collapsed. Any change opens it and forgets the collapse. Per-story activity (`stories`) is not part of it: the notice does not show it, and it changes on every agent step.
- Verified by the component's behaviour tests (jsdom), not in a browser.
