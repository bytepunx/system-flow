---
id: S-0156
type: story
nature: improvement
title: flai serve and flai mcp keep parsed work items between requests and read again only the files that changed
status: backlog
parent: E-0012
owner: alex
created: 2026-09-29T07:00:28Z
updated: 2026-09-29T07:00:28Z
transitions: []
tags: [cli]
topics: [server-side, back-end]
touches: [flai/internal/workitem, flai/internal/hostapi, flai/internal/mcpserver]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0156 flai serve and flai mcp keep parsed work items between requests and read again only the files that changed

## Goal

Every dashboard request and MCP tool call that lists work items reads and parses every item file from disk, the archive included: over 800 files and about 115 ms each time on this repository, and twice in one MCP `inbox`. Cause 1 of `design/system/server-performance.md`. Keep the parsed items in the long-running `flai serve` and `flai mcp` processes and read again only the files whose modification time or size changed, or that appeared or went, so that listing costs a directory walk.

## Acceptance criteria
- [ ] `board.get`, `item.get`, and `items.list` answered by `flai serve` on this repository show `repo.list` under 15 ms in their `request answered` events once warm.
- [ ] MCP `inbox`, `board`, `wait_for_work`, and `item_get` show the same, and `inbox` lists items once, not twice.
- [ ] An item created, edited, moved, archived, or deleted by any process, flai or a hand edit, is seen by the next request: behaviour tests cover each.
- [ ] Two requests at once share one read rather than each reading everything.

## Tasks

## Notes

Measured by S-0152: `repo.list` 114 ms of `board.get`'s 190 ms; `changes.read` 107 ms and `repo.list` 106 ms of MCP `inbox`'s 308 ms. `repo.List(false)` (open items only) takes 5 ms: the archive is the cost. Proposed: a cache in `workitem.Repo`, or beside it, keyed by path with modification time and size, filled on first use and checked by a walk on each list. The walk of all 1078 Markdown files takes 5.6 ms.
