---
id: D-0203
title: flai stream sync's trial merge blames a story for conflicts between main and another story's stale branch
kind: cause
sources: [I-0064]
topics: [cli, flai, go, code, all]
paths: [flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/flai-cli.md, design/issues/I-0064-flai-stream-sync-s-trial-merge-blames-a-story-for-conflicts-between-main-and-another-story-s-stale-branch.md, design/issues/summary.md, design/issues/I-0073-flai-check-finds-threads-archived-outside-the-story-at-close-out.md, design/issues/I-0076-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md]
date: 2026-10-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0203 flai stream sync's trial merge blames a story for conflicts between main and another story's stale branch

I-0064 was closed on 2026-10-07.

The cause:

```text
flai stream sync's trial merge blames a story for conflicts between main and another story's stale branch
```

The fix:

```text
S-0251: flai stream sync's trial merge reports a pair's conflict only on a path both stories changed since they left main (changedByBoth in flai/cmd/stream_sync.go), so what main brought against a stale branch is no longer blamed on the story that synced
```

Source: I-0064.
