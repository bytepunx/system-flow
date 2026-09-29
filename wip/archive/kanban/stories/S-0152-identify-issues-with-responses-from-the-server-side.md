---
id: S-0152
type: story
nature: research
title: Identify issues with responses from the server side
status: done
parent: E-0012
owner: alex
created: 2026-09-29T05:49:48Z
updated: 2026-09-29T07:03:59Z
transitions:
  - to: ready
    at: 2026-09-29T06:44:46Z
    by: alex
  - to: in-progress
    at: 2026-09-29T06:45:00Z
    by: agent-S-0152
  - to: review
    at: 2026-09-29T07:02:59Z
    by: agent-S-0152
  - to: done
    at: 2026-09-29T07:03:59Z
    by: alex
tags: [cli]
topics: [server-side]
touches: [flai/cmd, flai/internal/perf, flai/internal/channel, flai/internal/mcpserver, flai/internal/execx, flai/internal/hostapi, design/system/flai-cli.md, design/system/server-performance.md, design/system/README.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, design/issues]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1054
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 244
      output: 1506
      cache_read: 21641985
      cache_write: 254594
      cost: 8.6626
---
# S-0152 Identify issues with responses from the server side

## Goal

The server regularly hits long pauses when the board or another page is loading.

Let's make sure we instrument these calls in the modules handling them so that we can eliminate transports (MCP, WebSocket) from the telemetry collected.

Write stories to capture findings and recommended course of action.

## Acceptance criteria
- [x] Instrumentation is available to help track down hot spots or slow areas in the current implementation
- [x] Stories exist with proposed solutions for the issues identified

## Tasks
- T-0544 Time every dashboard and MCP request beneath its transport
- T-0545 Time the phases a request spends in the modules that answer it
- T-0546 Measure the dashboard's requests on this repository and record what is slow
- T-0547 Write the stories that fix what the measurement found

## Notes

Instrumentation (T-0544, T-0545): every request `flai serve` answers over the channel and every request `flai mcp` answers is timed inside flai and logged once as `request answered` (component `perf`) with method or tool, duration, size, and phases; `FLAI_SLOW_REQUEST` (default 500ms) sets when it is info. `flai hostapi --timing` times a method with no transport; `FLAI_PPROF_ADDR` offers Go's profiles from `flai serve` on loopback. Documented in `docs/users/flai.md#request-timing-s-0152` and `docs/operators/settings.md`.

Findings (T-0546): `design/system/server-performance.md`. The transport adds almost nothing: the dashboard's own durations are within a few milliseconds of flai's. Seven causes, each with a backlog story under E-0012 (T-0547): S-0156 (re-parsing every item, archive included, per request), S-0157 (25 git processes for pending releases), S-0158 (the whole check in the designer's inbox), S-0159 (polled reads that start flai), S-0160 (145 ms flai start from `atotto/clipboard` on a long `PATH`), S-0161 (the dashboard forgets every answer at any change), S-0162 (1.28 MB and 772 KB answers).
