---
id: S-0127
type: story
nature: improvement
title: Threads should visually delineate between agent text and operator text
status: done
parent: E-0003
owner: alex
created: 2026-09-26T07:31:07Z
updated: 2026-09-26T07:46:34Z
transitions:
  - to: ready
    at: 2026-09-26T07:40:19Z
    by: alex
  - to: in-progress
    at: 2026-09-26T07:40:38Z
    by: agent-S-0127
  - to: review
    at: 2026-09-26T07:45:24Z
    by: agent-S-0127
  - to: done
    at: 2026-09-26T07:46:34Z
    by: alex
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
- [x] It’s visually clear which thread entries are operator vs agent.

## Tasks
- T-0460 Thread entries show the operator's on the right and agents' on the left, in different colours

## Notes

The operator is the manifest's `owner` (`designer` when none), as flai decides for its writes and inbox. Their entries sit on the right in the primary tint, agents' on the left in the neutral one. Verified by the component test (side and colour classes per entry) and by the classes appearing in a production build's CSS; not viewed in a running dashboard.
