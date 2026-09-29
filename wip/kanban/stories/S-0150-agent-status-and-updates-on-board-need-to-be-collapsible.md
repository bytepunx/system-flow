---
id: S-0150
type: story
nature: improvement
title: Agent status and updates on board need to be collapsible
status: ready
parent: E-0013
owner: alex
created: 2026-09-29T05:31:24Z
updated: 2026-09-29T06:44:24Z
transitions:
  - to: ready
    at: 2026-09-29T06:44:24Z
    by: alex
tags: [dashboard]
topics: [front-end, board]
touches: [flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0150 Agent status and updates on board need to be collapsible

## Goal

The Agent notification/summary box on the board page contains valuable information but in some cases can get quite long. Allow the user to collapse the area using a toggle carrot so that only the bold title appears when collapsed.

## Acceptance criteria
- [ ] The operator can collapse the box using the toggle
- [ ] The operator can expand the box using the toggle
- [ ] Changes to the status cause the box to auto-expand

## Tasks

## Notes
