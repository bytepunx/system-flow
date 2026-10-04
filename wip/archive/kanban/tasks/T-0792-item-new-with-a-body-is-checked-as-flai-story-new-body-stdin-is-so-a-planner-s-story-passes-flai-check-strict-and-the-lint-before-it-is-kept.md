---
id: T-0792
type: task
nature: feature
title: item_new with a body is checked as flai story new --body-stdin is, so a planner's story passes flai check --strict and the lint before it is kept
status: done
parent: S-0209
owner: alex
created: 2026-10-04T04:02:44Z
updated: 2026-10-04T04:09:55Z
transitions:
  - to: ready
    at: 2026-10-04T04:03:31Z
    by: agent-S-0209
  - to: in-progress
    at: 2026-10-04T04:03:33Z
    by: agent-S-0209
  - to: done
    at: 2026-10-04T04:09:55Z
    by: agent-S-0209
stream: S-0209
tags: []
touches: [flai/internal/mcpserver]
usage:
  source: log
  seconds: 382
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 37
      output: 8297
      cache_read: 3169558
      cache_write: 36016
      cost: 1.0383
---
# T-0792 item_new with a body is checked as flai story new --body-stdin is, so a planner's story passes flai check --strict and the lint before it is kept

## Work

- MCP `item_new` (`flai/internal/mcpserver/items_write.go`) creates through `itemnew.Create`, with flai check run with the item in place and anything it introduces refusing it, when a body or `after` is given, as the CLI's `--body-stdin` and `--after` are; the tool's description in `folder.go` says so.
- Waits for nothing: the first layer.

## Done when

- [ ] A test in `flai/internal/mcpserver` shows `item_new` with a body that brings a check finding refused and nothing left behind, and a clean one kept.
- [ ] `go test ./internal/mcpserver` passes.

## Notes

Today only an `after` makes `item_new` checked; a body the planner writes is kept unchecked.
