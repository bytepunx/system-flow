---
id: S-0142
type: story
nature: feature
title: Add a stream window to the activity page agent panes
status: ready
parent: E-0011
owner: alex
created: 2026-09-28T23:26:54Z
updated: 2026-09-29T03:17:59Z
transitions:
  - to: ready
    at: 2026-09-29T03:17:59Z
    by: alex
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0142 Add a stream window to the activity page agent panes

## Goal

flai serve should provide active agent streams via command for the dashboard. the dashboard should provide live feeds from the agent under the activity page so that agent activity is more observable.

## Acceptance criteria
- [ ] Active agent streams are available
- [ ] Active agent streams are visible on the agent activity cards

## Tasks

## Notes
