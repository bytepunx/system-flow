---
id: T-0564
type: task
nature: improvement
title: Measure repo.list warm in flai serve and flai mcp, and record it
status: done
parent: S-0156
owner: alex
created: 2026-09-29T19:20:03Z
updated: 2026-09-29T19:29:45Z
transitions:
  - to: ready
    at: 2026-09-29T19:20:09Z
    by: agent-S-0156
  - to: in-progress
    at: 2026-09-29T19:25:34Z
    by: agent-S-0156
  - to: done
    at: 2026-09-29T19:29:45Z
    by: agent-S-0156
stream: S-0156
tags: []
touches: [design/system, docs]
usage:
  source: log
  seconds: 251
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 59
      output: 18309
      cache_read: 3636929
      cache_write: 61402
      cost: 1.585
---
# T-0564 Measure repo.list warm in flai serve and flai mcp, and record it

## Work

Build flai from the story branch and time `board.get`, `item.get`, and `items.list` through `flai serve`'s methods, and MCP `inbox`, `board`, `wait_for_work`, and `item_get` over stdio, on this repository, warm. Record the numbers and the store in `design/system/server-performance.md` and `design/system/flai-cli.md`.

## Done when

- Each shows `repo.list` under 15 ms once warm, measured and recorded.
- The design says how items are kept and when a file is read again.

## Notes
