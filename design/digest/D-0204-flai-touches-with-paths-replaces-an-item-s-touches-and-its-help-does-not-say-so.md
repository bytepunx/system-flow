---
id: D-0204
title: flai touches with paths replaces an item's touches, and its help does not say so
kind: cause
sources: [I-0067]
topics: [cli, flai, go, code, all]
paths: [flai/cmd/touches.go, flai/cmd/touches_test.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, design/system/flai-cli.md, design/issues/I-0067-flai-touches-with-paths-replaces-an-item-s-touches-and-its-help-does-not-say-so.md, design/issues/summary.md, design/issues/I-0073-flai-check-finds-threads-archived-outside-the-story-at-close-out.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md, design/issues/I-0092-two-story-branches-that-each-bump-the-same-issue-conflict-in-its-front-matter-whose-count-last-reported-and-updated-lines-both-rewrite.md]
date: 2026-10-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0204 flai touches with paths replaces an item's touches, and its help does not say so

I-0067 was closed on 2026-10-07.

The cause:

```text
flai touches with paths replaces an item's touches, and its help does not say so
```

The fix:

```text
S-0254: flai touches adds with --add and removes with --remove, its help says paths given alone replace the list, and flai stream sync's widen hint uses --add
```

Source: I-0067.
