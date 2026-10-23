---
id: TH-0625
title: Merge queue batch 20261023-194240 fails with S-0489 at integration
anchor:
  path: wip/kanban/stories/S-0489-host-page-under-configuration.md
  item: S-0489
status: resolved
participants: [flai serve]
created: 2026-10-23T19:50:23Z
updated: 2026-10-23T20:04:43Z
---

# TH-0625 Merge queue batch 20261023-194240 fails with S-0489 at integration

On wip/kanban/stories/S-0489-host-page-under-configuration.md.

## Entries

### 2026-10-23T19:50:23Z flai serve
Merge queue batch `20261023-194240` fails with S-0489 at `integration`. The queue verified S-0489 alone on main at `0a1648e65dc1`, at `1521ddbc2050`, the copy of its verified head `d75b8189e04b` rebased in the batch, and `integration` failed:

```text
flai/cmd/reference_test.go:127 TestReferenceIsCurrent: docs/users/flai-reference.md is stale: run make flai-reference
```

The whole output of each tier it ran is in `.flai-cache/merge-queue/batches/20261023-194240`.

flai moved S-0489 back to in-progress and took it out of the merge queue. Fix it on `story/S-0489`, run flai verify S-0489, move it to review, and enqueue it again with flai queue add S-0489.

### 2026-10-23T20:04:43Z flai serve
Resolved: S-0489 is back in review at c5c265c34d4f, a new commit; it is batched again once it is enqueued
