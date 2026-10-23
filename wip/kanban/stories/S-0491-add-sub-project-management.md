---
id: S-0491
type: story
nature: improvement
title: Add sub-project management
status: backlog
owner: alex
created: 2026-10-23T20:26:10Z
updated: 2026-10-23T20:26:10Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 3h
    by: alex
    at: 2026-10-23T20:26:10Z
---
# S-0491 Add sub-project management

## Goal

No flai command registers a sub-project, and no document in the project or the template gives the schema of projects. The operator is expected to add the entry to system-flow.yaml by hand. flai's own messages show that an entry has name, path, and kind, but there's little else to go off.

## Acceptance criteria
- [ ] The projects page supports sub-project management in addition to the projects themselves
- [ ] If no sub-folder exists, creating a sub-project adds that sub-folder to the project
- [ ] Add supporting API and MCP calls for this as well so that agents (orchestrator and story agents) can handle this directly without requiring the user to take direct steps

## Tasks

## Notes
