---
id: S-0125
type: story
nature: research
title: Determine how to prime context with relevant documentation only
status: done
parent: E-0010
owner: alex
created: 2026-09-26T07:23:21Z
updated: 2026-09-26T17:48:17Z
transitions:
  - to: ready
    at: 2026-09-26T07:40:35Z
    by: alex
  - to: in-progress
    at: 2026-09-26T08:03:30Z
    by: agent-S-0125
  - to: review
    at: 2026-09-26T17:47:07Z
    by: agent-S-0125
  - to: done
    at: 2026-09-26T17:48:17Z
    by: alex
tags: [cli]
touches: [flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0125 Determine how to prime context with relevant documentation only

## Goal

Research methods for determining the relevant documentation (adrs, conventions, etc) and excepts from them in order to build up the right context for an agent that is receiving instructions from completing coding tasks.

Use threads to communicate with the operator on decisions and direction. Author an ADR with the outcomes and then craft stories for epic 10 to implement the decisions made.

## Acceptance criteria
- [x] New ADR capturing the solution to focused agent context building
- [x] Stories for epic 10 capturing the work necessary to implement it

## Tasks
- T-0461 Survey how an agent's context can be primed with only the documentation its story needs
- T-0462 Decide the design with the designer and record it in an ADR
- T-0463 Write the E-0010 stories that implement the decision

## Notes

Findings: `design/system/agent-context.md`. Decision: ADR-0047 (TH-0020, TH-0021). Stories: S-0134 to S-0138 under E-0010, in backlog for the designer to refine and order; each has a goal and checkbox criteria.
