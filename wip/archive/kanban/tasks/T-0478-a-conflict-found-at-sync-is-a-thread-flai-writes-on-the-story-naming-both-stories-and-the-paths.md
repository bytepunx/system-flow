---
id: T-0478
type: task
nature: feature
title: A conflict found at sync is a thread flai writes on the story, naming both stories and the paths
status: done
parent: S-0131
owner: alex
created: 2026-09-26T18:05:55Z
updated: 2026-09-26T18:09:52Z
transitions:
  - to: ready
    at: 2026-09-26T18:06:13Z
    by: agent-S-0131
  - to: in-progress
    at: 2026-09-26T18:08:04Z
    by: agent-S-0131
  - to: done
    at: 2026-09-26T18:09:52Z
    by: agent-S-0131
stream: S-0131
tags: []
touches: [flai/cmd/stream_sync.go, flai/internal/threads]
---
# T-0478 A conflict found at sync is a thread flai writes on the story, naming both stories and the paths

## Work

- A conflict opens one thread per pair of stories, written by `flai`, on the story that synced, titled with both IDs and naming the conflicting paths. A thread written by `flai` awaits every agent and the designer, so both stories' agents see it in the MCP `inbox` and the designer sees it in the dashboard's inbox.
- A later sync that finds the same pair and the same paths writes nothing; a different set of paths adds an entry to the open thread; a pair that merges cleanly again resolves it.
- The integration test checks the thread appears in the MCP inbox for the other story's agent and in the designer inbox, and that a clean sync resolves it.

## Done when

- Tests show the thread in both inboxes, deduplicated across syncs, resolved when the pair is clean.

## Notes
