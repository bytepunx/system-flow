---
id: S-0156
type: story
nature: improvement
title: flai serve and flai mcp keep parsed work items between requests and read again only the files that changed
status: done
parent: E-0012
owner: alex
created: 2026-09-29T07:00:28Z
updated: 2026-09-29T19:41:05Z
transitions:
  - to: ready
    at: 2026-09-29T19:18:22Z
    by: alex
  - to: in-progress
    at: 2026-09-29T19:18:48Z
    by: agent-S-0156
  - to: review
    at: 2026-09-29T19:29:54Z
    by: agent-S-0156
  - to: done
    at: 2026-09-29T19:41:05Z
    by: alex
tags: [cli]
topics: [server-side, back-end]
touches: [flai/internal/workitem, flai/internal/mcpserver, design/system/server-performance.md, design/system/flai-cli.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 703
  models:
    - model: claude-opus-5-5
      input: 158
      output: 48849
      cache_read: 9703683
      cache_write: 163827
      cost: 4.229
---
# S-0156 flai serve and flai mcp keep parsed work items between requests and read again only the files that changed

## Goal

Every dashboard request and MCP tool call that lists work items reads and parses every item file from disk, the archive included: over 800 files and about 115 ms each time on this repository, and twice in one MCP `inbox`. Cause 1 of `design/system/server-performance.md`. Keep the parsed items in the long-running `flai serve` and `flai mcp` processes and read again only the files whose modification time or size changed, or that appeared or went, so that listing costs a directory walk.

## Acceptance criteria
- [x] `board.get`, `item.get`, and `items.list` answered by `flai serve` on this repository show `repo.list` under 15 ms in their `request answered` events once warm.
- [x] MCP `inbox`, `board`, `wait_for_work`, and `item_get` show the same, and `inbox` lists items once, not twice.
- [x] An item created, edited, moved, archived, or deleted by any process, flai or a hand edit, is seen by the next request: behaviour tests cover each.
- [x] Two requests at once share one read rather than each reading everything.

## Tasks
- T-0562 workitem keeps parsed items and reads again only the files that changed
- T-0563 MCP inbox lists work items once
- T-0564 Measure repo.list warm in flai serve and flai mcp, and record it

## Notes

Measured by S-0152: `repo.list` 114 ms of `board.get`'s 190 ms; `changes.read` 107 ms and `repo.list` 106 ms of MCP `inbox`'s 308 ms. `repo.List(false)` (open items only) takes 5 ms: the archive is the cost. Proposed: a cache in `workitem.Repo`, or beside it, keyed by path with modification time and size, filled on first use and checked by a walk on each list. The walk of all 1078 Markdown files takes 5.6 ms.

Delivered (S-0156): a process-wide store in `flai/internal/workitem/store.go` behind `Repo.List` and `Repo.Get`, and MCP `inbox` listing once. Warm `repo.list` is 3 to 5 ms for every request named above; figures in `design/system/server-performance.md` § After the stories.

How the first criterion was verified: the dashboard's methods were called five times in one process through `hostapi.Methods`, the table `flai serve` answers with, and timed by the same `perf` recorder that writes its `request answered` event. A second `flai serve` was not started beside the operator's, so the events were not read from a running `flai serve`. The MCP tools were measured through a real `flai mcp` over stdio, built from the story branch.

Behaviour tests: `flai/internal/workitem/store_test.go` (created, edited, moved, archived, and deleted by flai and by hand; a same-size write in the same tick; a file renamed over another; lists at once parsing each file once; a listed item is the caller's own) and `TestInboxListsItemsOnce` in `flai/internal/mcpserver/timing_test.go`.
