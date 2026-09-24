---
id: S-0116
type: story
nature: remediation
title: when a story moves to ready, if there is available capacity, it should be assigned to an agent
status: in-progress
parent: E-0003
owner: alex
created: 2026-09-24T08:37:20Z
updated: 2026-09-24T08:46:28Z
transitions:
  - to: ready
    at: 2026-09-24T08:37:26Z
    by: alex
  - to: in-progress
    at: 2026-09-24T08:44:21Z
    by: agent-S-0116
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
- T-0411 A ready story whose agent is changed after its agent ended gets another
- T-0412 The story page says changing the agent of a ready story starts another
- T-0413 A story moved from backlog to ready with room in progress is started, per the designer's answer on TH-0010
- T-0414 Record the added case in the living design and the operator and user docs

## Notes
