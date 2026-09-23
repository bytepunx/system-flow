---
id: S-0104
type: story
nature: feature
title: Flai should spawn and manage agents in response to Kanban changes
status: done
parent: E-0008
owner: alex
created: 2026-09-23T16:50:50Z
updated: 2026-09-23T20:18:48Z
transitions:
  - to: ready
    at: 2026-09-23T16:52:04Z
    by: alex
  - to: in-progress
    at: 2026-09-23T17:55:10Z
    by: system-flow
  - to: review
    at: 2026-09-23T17:55:10Z
    by: system-flow
  - to: done
    at: 2026-09-23T20:18:48Z
    by: alex
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
---
# S-0104 Flai should spawn and manage agents in response to Kanban changes

## Goal

The flai server process should actively monitor stories for change to state ready so that it can spawn and manage agent workers that will do the work for every card in ready until a blocking scenario occurs (WIP limit, inbox queue, etc.)

## Acceptance criteria

- [x] flai should create agentic adapters so that it can correctly dispatch work to agents via their harness and specifying their model and other configuration values

- [x] a story moved to or created in ready is worked to completion. in the event input is needed from the user, the inbox should be used and then the MCP monitored for user response.

- [x] flaiover should indicate activity in both the story detail view and on the card to indicate when an agent is actively working on the story (green dot for working, yellow for waiting, red for failure)

## Tasks
- T-0367 Harness adapters: claude-code and the operator's command start a story's agent with its model and config
- T-0368 flai serve starts an agent for every story that enters ready, within the WIP limit, and knows each one's state: working, waiting, failed
- T-0369 agent.status reports each story's agent; the board card and the story page show a green, yellow, or red dot
- T-0370 Tried end to end with a real claude-code agent; ADR and docs

## Notes
