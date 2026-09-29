---
id: S-0152
type: story
nature: research
title: Identify issues with responses from the server side
status: in-progress
parent: E-0012
owner: alex
created: 2026-09-29T05:49:48Z
updated: 2026-09-29T06:56:28Z
transitions:
  - to: ready
    at: 2026-09-29T06:44:46Z
    by: alex
  - to: in-progress
    at: 2026-09-29T06:45:00Z
    by: agent-S-0152
tags: [cli]
topics: [server-side]
touches: [flai/cmd, flai/internal/perf, flai/internal/channel, flai/internal/mcpserver, flai/internal/workitem, flai/internal/execx, flai/internal/hostapi, design/system, docs, design/issues]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 879
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 224
      output: 1360
      cache_read: 19048908
      cache_write: 240173
      cost: 7.6365
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
- T-0544 Time every dashboard and MCP request beneath its transport
- T-0545 Time the phases a request spends in the modules that answer it
- T-0546 Measure the dashboard's requests on this repository and record what is slow
- T-0547 Write the stories that fix what the measurement found

## Notes
