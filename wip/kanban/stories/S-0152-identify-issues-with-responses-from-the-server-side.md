---
id: S-0152
type: story
nature: research
title: Identify issues with responses from the server side
status: in-progress
parent: E-0012
owner: alex
created: 2026-09-29T05:49:48Z
updated: 2026-09-29T06:45:00Z
transitions:
  - to: ready
    at: 2026-09-29T06:44:46Z
    by: alex
  - to: in-progress
    at: 2026-09-29T06:45:00Z
    by: agent-S-0152
tags: [cli]
topics: [server-side]
touches: [flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0152 Identify issues with responses from the server side

## Goal

The server regularly hits long pauses when the board or another page is loading.

Let's make sure we instrument these calls in the modules handling them so that we can eliminate transports (MCP, WebSocket) from the telemetry collected.

Write stories to capture findings and recommended course of action.

## Acceptance criteria
- [ ] Instrumentation is available to help track down hot spots or slow areas in the current implementation
- [ ] Stories exist with proposed solutions for the issues identified

## Tasks

## Notes
