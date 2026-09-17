---
id: T-0083
type: task
nature: feature
title: "Repo reader: manifest, items, docs tree and file, with a mtime cache and a watcher emitting SSE"
status: done
parent: S-0011
owner: alex
created: 2026-09-17T04:02:08Z
updated: 2026-09-17T04:09:26Z
transitions:
  - to: ready
    at: 2026-09-17T04:09:25Z
    by: agent
  - to: in-progress
    at: 2026-09-17T04:09:26Z
    by: agent
  - to: done
    at: 2026-09-17T04:09:26Z
    by: agent
stream: S-0011
tags: [dashboard]
---

# T-0083 Repo reader: manifest, items, docs tree and file, with a mtime cache and a watcher emitting SSE

## Work
src/lib/server/repo: resolve PROJECT_DIR, load system-flow.yaml, list work items from kanban and archive with front matter parsed by the yaml package, walk the documentation trees (conventions, adrs, system, tech, issues, docs, wip) with front matter, read one file raw, cache by path and mtime, watch with chokidar and expose an event emitter for SSE.

## Done when
Unit tests on the metrics fixture and this repository: item counts, status consistency, tree shape.

## Notes
