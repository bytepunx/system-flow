---
id: T-0742
type: task
nature: feature
title: MCP item_new, item_edit, item_move and hostapi item.new, item.edit, item.move carry the planning fields
status: done
parent: S-0199
owner: alex
created: 2026-10-03T05:41:25Z
updated: 2026-10-03T06:12:14Z
transitions:
  - to: ready
    at: 2026-10-03T06:01:09Z
    by: agent-S-0199
  - to: in-progress
    at: 2026-10-03T06:01:09Z
    by: agent-S-0199
  - to: done
    at: 2026-10-03T06:12:14Z
    by: agent-S-0199
stream: S-0199
tags: []
touches: [flai/internal/mcpserver, flai/internal/hostapi]
after: [T-0741]
usage:
  source: log
  seconds: 665
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 91
      output: 33531
      cache_read: 5421675
      cache_write: 136775
      cost: 2.5621
---
# T-0742 MCP item_new, item_edit, item_move and hostapi item.new, item.edit, item.move carry the planning fields

## Work

MCP `item_new` takes `draft` for a story; `item_edit` takes `draft`, `cost_of_delay` (inputs and value), `forecast`, and their clears, through `itemedit` with the calling agent as `by`; `item_move` refuses a draft story to ready with "finalize it first", since no agent has the orchestrator's permission yet (S-0218 defines it). hostapi `item.new` takes `draft`, `item.edit` takes the same fields as flags of `flai edit`, and `item.move` refuses a draft to ready unless `finalize` is given, which passes `--yes`. Waits for T-0741, whose flags and `Change` fields it calls.

## Done when

Each tool and method sets and clears each field and the draft move is refused through both, with tests.

## Notes
