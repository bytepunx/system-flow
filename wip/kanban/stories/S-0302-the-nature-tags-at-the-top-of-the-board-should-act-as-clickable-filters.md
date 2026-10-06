---
id: S-0302
type: story
nature: improvement
title: The nature tags at the top of the board should act as clickable filters
status: backlog
owner: alex
created: 2026-10-06T23:17:17Z
updated: 2026-10-06T23:17:17Z
transitions: []
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: alex
    at: 2026-10-06T23:17:17Z
---
# S-0302 The nature tags at the top of the board should act as clickable filters

## Goal

The nature tags along the top of the board should work like filters. Their default coloring should represent 'filtered out' while a brighter color set should represent 'shown' with clicking each toggling them on or off. All natures should be toggled on (shown) by default.

## Acceptance criteria
- [ ] When viewing a board, all natures are toggled on so that all natures of work items are displayed
- [ ] When toggling a nature off, it reverts to its default coloring and cards with that nature are hidden
- [ ] When toggling a nature on, it shows a highlighted version of its default color theme, and cards with that nature are shown
- [ ] A user's choices are preserved in browser local storage so that navigating away from the board does not reset the nature filtering

## Tasks

## Notes
