---
id: T-0563
type: task
nature: improvement
title: MCP inbox lists work items once
status: done
parent: S-0156
owner: alex
created: 2026-09-29T19:20:03Z
updated: 2026-09-29T19:25:34Z
transitions:
  - to: ready
    at: 2026-09-29T19:20:09Z
    by: agent-S-0156
  - to: in-progress
    at: 2026-09-29T19:23:47Z
    by: agent-S-0156
  - to: done
    at: 2026-09-29T19:25:34Z
    by: agent-S-0156
stream: S-0156
tags: []
touches: [flai/internal/mcpserver]
usage:
  source: log
  seconds: 107
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 27
      output: 8271
      cache_read: 1643003
      cache_write: 27739
      cost: 0.716
---
# T-0563 MCP inbox lists work items once

## Work

`inbox` lists every item for its board view and again for the changes since the agent last looked. List once and give the same items to both.

## Done when

- `inbox` shows one `repo.list` phase and no listing inside `changes.read`; its tests pass.

## Notes
