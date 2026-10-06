---
id: S-0300
type: story
nature: improvement
title: The planner agent should be available on epic work items
status: backlog
owner: alex
created: 2026-10-06T21:42:35Z
updated: 2026-10-06T21:42:35Z
transitions: []
tags: [dashboard, cli]
topics: [planner]
touches: [flaiover/src, flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 4h
    by: alex
    at: 2026-10-06T21:42:35Z
---
# S-0300 The planner agent should be available on epic work items

## Goal

The epic page and the epic context menu should provide interfaces for invoking the planner to write its stories and their tasks.

## Acceptance criteria
- [ ] The epic page has a plan button that invokes the planner for the epic
- [ ] The epic item's context menu has a plan option
- [ ] The planner agent is capable of writing draft stories and draft tasks for the epic

## Tasks

## Notes
