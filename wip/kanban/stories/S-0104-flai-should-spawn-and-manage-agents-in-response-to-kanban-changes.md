---
id: S-0104
type: story
nature: feature
title: Flai should spawn and manage agents in response to Kanban changes
status: ready
parent: E-0008
owner: alex
created: 2026-09-23T16:50:50Z
updated: 2026-09-23T16:52:04Z
transitions:
  - to: ready
    at: 2026-09-23T16:52:04Z
    by: alex
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
---
# S-0104 Flai should spawn and manage agents in response to Kanban changes

## Goal

The flai server process should actively monitor stories for change to state ready so that it can spawn and manage agent workers that will do the work for every card in ready until a blocking scenario occurs (WIP limit, inbox queue, etc.)

## Acceptance criteria

- [ ] flai should create agentic adapters so that it can correctly dispatch work to agents via their harness and specifying their model and other configuration values

- [ ] a story moved to or created in ready is worked to completion. in the event input is needed from the user, the inbox should be used and then the MCP monitored for user response.

- [ ] flaiover should indicate activity in both the story detail view and on the card to indicate when an agent is actively working on the story (green dot for working, yellow for waiting, red for failure)

## Tasks

## Notes
