---
id: S-0170
type: story
nature: improvement
title: The Operator can stop agents
status: backlog
parent: E-0013
owner: alex
created: 2026-09-29T23:56:11Z
updated: 2026-09-29T23:56:11Z
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
# S-0170 The Operator can stop agents

## Goal

Right now, a stuck agent will remain in the activity page forever since it is technically just going to sit there and persist between restarts and even machine restarts.

The operator should be able to stop an agent from the activity page (with confirmation to protect against accidental clicks).

## Acceptance criteria
- [ ] The operator can click stop on any agent running in the activity page
- [ ] When clicking stop, the operator must confirm their intention (should include a warning explaining the impact)
- [ ] When the operator confirms their intent is to stop the agent, the agent should be terminated

## Tasks

## Notes
