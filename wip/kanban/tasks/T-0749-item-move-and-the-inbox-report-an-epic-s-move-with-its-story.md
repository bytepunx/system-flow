---
id: T-0749
type: task
nature: improvement
title: item_move and the inbox report an epic's move with its story
status: done
parent: S-0200
owner: alex
created: 2026-10-03T07:11:21Z
updated: 2026-10-03T07:32:13Z
transitions:
  - to: ready
    at: 2026-10-03T07:11:52Z
    by: agent-S-0200
  - to: in-progress
    at: 2026-10-03T07:19:14Z
    by: agent-S-0200
  - to: done
    at: 2026-10-03T07:32:01Z
    by: agent-S-0200
stream: S-0200
tags: []
touches: [flai/internal/mcpserver]
after: [T-0745]
usage:
  source: log
  seconds: 767
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 47
      output: 16703
      cache_read: 2890830
      cache_write: 64777
      cost: 1.3083
---
# T-0749 item_move and the inbox report an epic's move with its story

## Work

The `item_move` result carries the epic's move, and `inbox` and `wait_for_events` describe an epic's change as following its story. Check that the dashboard's `item.move` and `item.accept` pass the epic's move through.

It waits for T-0745, whose `MoveResult.Followed` and `Change.Follows` it reports.

## Done when

- [x] `ItemMoveOut` has `followed`, and the change summary names the story an epic followed
- [x] Tests in `flai/internal/mcpserver` cover both

## Notes
