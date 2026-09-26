---
id: S-0125
type: story
nature: research
title: Determine how to prime context with relevant documentation only
status: ready
parent: E-0010
owner: alex
created: 2026-09-26T07:23:21Z
updated: 2026-09-26T07:40:35Z
transitions:
  - to: ready
    at: 2026-09-26T07:40:35Z
    by: alex
tags: [cli]
touches: [flai/cmd]
agent:
  harness: claude-code
  model: claude-opus5
  config:
    effort: high
---
# S-0125 Determine how to prime context with relevant documentation only

## Goal

Research methods for determining the relevant documentation (adrs, conventions, etc) and excepts from them in order to build up the right context for an agent that is receiving instructions from completing coding tasks.

Use threads to communicate with the operator on decisions and direction. Author an ADR with the outcomes and then craft stories for epic 10 to implement the decisions made.

## Acceptance criteria
- [ ] New ADR capturing the solution to focused agent context building
- [ ] Stories for epic 10 capturing the work necessary to implement it

## Tasks

## Notes
