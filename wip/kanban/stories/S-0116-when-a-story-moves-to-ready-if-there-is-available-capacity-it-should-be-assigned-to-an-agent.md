---
id: S-0116
type: story
nature: remediation
title: when a story moves to ready, if there is available capacity, it should be assigned to an agent
status: backlog
parent: E-0003
owner: alex
created: 2026-09-24T08:37:20Z
updated: 2026-09-24T08:37:20Z
transitions: []
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0116 when a story moves to ready, if there is available capacity, it should be assigned to an agent

## Goal

When stories are moved from the backlog into ready, if there is available capacity in the in process WIP, the story should be assigned to an agent.

## Acceptance criteria
- [ ] a story moved from backlog to ready should be assigned to an agent if in process WIP isn't full
- [ ] a story that had agent information changed in ready state should be assigned to an agent if in process WIP isn't full

## Tasks

## Notes
