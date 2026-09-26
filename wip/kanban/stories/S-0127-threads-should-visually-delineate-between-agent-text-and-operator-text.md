---
id: S-0127
type: story
nature: improvement
title: Threads should visually delineate between agent text and operator text
status: backlog
parent: E-0003
owner: alex
created: 2026-09-26T07:31:07Z
updated: 2026-09-26T07:34:30Z
transitions: []
tags: [dashboard]
touches: [flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0127 Threads should visually delineate between agent text and operator text

## Goal

When an operator is viewing a thread, it should be visually obvious which entries are the agent and which belong to the operator. Typical visualization would include a right alignment for operator and left alignment for agent with an additional change in colors. Since this is stored in markdown, the options may be limited but we should find a way to make the difference easy to spot.

## Acceptance criteria
- [ ] It’s visually clear which thread entries are operator vs agent. 

## Tasks

## Notes
