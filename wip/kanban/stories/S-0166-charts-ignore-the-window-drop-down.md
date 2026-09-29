---
id: S-0166
type: story
nature: remediation
title: Charts ignore the window drop-down.
status: backlog
parent: E-0013
owner: alex
created: 2026-09-29T22:45:43Z
updated: 2026-09-29T22:45:43Z
transitions: []
tags: [dashboard, cli]
topics: [server-side, client-side]
touches: [flaiover/src, flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0166 Charts ignore the window drop-down.

## Goal

Neither the time axis or data points are correctly re-rendered when the time window drop down is set. From what I can tell each chart seems to have a pre-determined time window it renders with and no changes in the UI have any affect.

## Acceptance criteria
- [ ] The axis matches the selected time window
- [ ] The data points are rendered correctly for the selected time window
- [ ] Changing the window changes the graph's time axis and plotted data to match

## Tasks

## Notes
