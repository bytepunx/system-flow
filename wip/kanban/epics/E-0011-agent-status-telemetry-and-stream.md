---
id: E-0011
type: epic
nature: feature
title: Agent Status, Telemetry, and Stream
status: ready
owner: alex
created: 2026-09-28T22:51:06Z
updated: 2026-09-29T00:33:59Z
transitions:
  - to: ready
    at: 2026-09-29T00:33:59Z
    by: alex
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
---
# E-0011 Agent Status, Telemetry, and Stream

## Outcome

The operator can view a feed window for each active agent, each task and story now get token and pricing (where available) front-matter. New charts help operators understand how tokens are user per unit of time in relation to story/task completion and the relative cost (or estimated cost).

## Stories
- S-0142 Add a stream window to the activity page agent panes
- S-0143 Stories and Tasks track tokens and cost

## Notes
