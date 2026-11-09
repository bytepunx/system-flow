---
id: D-0536
title: publish.preview and the publish that tags keep reading git through PendingBatch, not the main-refs watcher's kept history, since the operator reads the preview right before a publish
kind: decision
sources: [S-0604]
topics: [server-side, back-end, flai, go]
paths: [flai/internal/release/history.go, flai/internal/release/release.go, flai/internal/serve/mergequeue.go, design/system/server-performance.md]
date: 2026-11-09
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0536 publish.preview and the publish that tags keep reading git through PendingBatch, not the main-refs watcher's kept history, since the operator reads the preview right before a publish

S-0604 (cause 14 of E-0033) made `board.get`, the inbox, and the orchestrator's figures reuse the kept release history while flai serve's main-refs watcher covers the root and no ref moved since it was read (`PendingBatchWatched`). It left `publish.preview` on `PendingBatch`, which reads git at every call. From S-0604's `## Decisions` (`wip/archive/agents/S-0604.md`):

```text
`publish.preview` keeps reading git (`PendingBatch`), not the watcher's word: the operator reads it right before a publish. It gains from the per-root lock already.
```

`design/system/server-performance.md` records that `publish.preview`, `release.EvaluateRepo`, and the publish that tags and pushes still call `PendingBatch`, but not why. A later performance story should not move the preview onto the watched history: what the operator reads before a publish must be what git holds at that moment, and the per-root lock already keeps one project's git from holding another's.

Sources: S-0604.
