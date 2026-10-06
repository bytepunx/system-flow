---
id: S-0303
type: story
nature: improvement
title: The type specifiers act as filters for the board
status: ready
owner: alex
created: 2026-10-06T23:23:11Z
updated: 2026-10-06T23:23:12Z
transitions:
  - to: ready
    at: 2026-10-06T23:23:12Z
    by: alex
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
    at: 2026-10-06T23:23:11Z
---
# S-0303 The type specifiers act as filters for the board

## Goal

The type legend at the top of the board acts as a filter for card types on the board. When toggled on, their coloring changes from the default to a highlighted version. When toggled off, their coloring reverts to its original. Toggled on types are shown on the board while toggled off types are hidden. This replaces the checkbox set that currently provides this filtering.

## Acceptance criteria
- [ ] The type legend provides a border around each type so that it's clearer that it's a clickable entity (in the same way nature's have a border)
- [ ] When a type is clicked, this toggles whether that type is shown or hidden
- [ ] When a type is "on" and that type of card should be shown, change its coloring to a highlighted version of the default
- [ ] Store a user's toggle selections in the browser's local storage so that navigating away from the page does not reset their selections
- [ ] by default, all types are selected

## Tasks

## Notes
