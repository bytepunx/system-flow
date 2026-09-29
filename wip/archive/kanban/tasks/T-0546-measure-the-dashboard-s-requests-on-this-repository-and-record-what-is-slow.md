---
id: T-0546
type: task
nature: research
title: Measure the dashboard's requests on this repository and record what is slow
status: done
parent: S-0152
owner: alex
created: 2026-09-29T06:46:36Z
updated: 2026-09-29T06:59:29Z
transitions:
  - to: ready
    at: 2026-09-29T06:56:24Z
    by: agent-S-0152
  - to: in-progress
    at: 2026-09-29T06:56:25Z
    by: agent-S-0152
  - to: done
    at: 2026-09-29T06:59:29Z
    by: agent-S-0152
stream: S-0152
tags: []
touches: [design/system]
usage:
  source: log
  seconds: 184
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 51
      output: 17931
      cache_read: 4678966
      cache_write: 51917
      cost: 1.7099
---
# T-0546 Measure the dashboard's requests on this repository and record what is slow

## Work

Use the instrumentation to measure the methods a board, story page, and documents page ask for, and the MCP tools agents call most, on this repository. Record the numbers, the hot spots, and their causes in a design document for E-0012.

## Done when

- `design/system/server-performance.md` lists each measured method with its time, size, and phases, and names the causes of the slow ones.

## Notes
