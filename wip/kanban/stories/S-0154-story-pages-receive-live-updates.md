---
id: S-0154
type: story
nature: feature
title: Story pages receive live updates
status: ready
parent: E-0013
owner: alex
created: 2026-09-29T05:55:11Z
updated: 2026-09-29T06:04:14Z
transitions:
  - to: ready
    at: 2026-09-29T06:04:14Z
    by: alex
tags: [dashboard, cli]
topics: [client-side, server-side]
touches: [flaiover/src, flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0154 Story pages receive live updates

## Goal

Story pages currently go stale once they're opened. This means that the operator must refresh them to catch updates on agent status or updates to threads once they respond.

## Acceptance criteria
- [ ] The story gets updates as they occur and the page updates to reflect the changes
- [ ] The active thread pane needs to show that the agent is working on processing the operator response with a simple animation of some kind.

## Tasks

## Notes
