---
id: TH-0767
title: Merge queue batch 20261030-182125 fails with S-0529 at integration
anchor:
  path: wip/kanban/stories/S-0529-a-story-the-operator-accepts-outside-the-merge-queue-lands-without-the-queue-tiers-so-a-stale-generated-file-on-main-fails-every-later-batch-at-integration.md
  item: S-0529
status: resolved
participants: [flai serve]
created: 2026-10-30T18:30:48Z
updated: 2026-10-30T19:04:06Z
---

# TH-0767 Merge queue batch 20261030-182125 fails with S-0529 at integration

On wip/kanban/stories/S-0529-a-story-the-operator-accepts-outside-the-merge-queue-lands-without-the-queue-tiers-so-a-stale-generated-file-on-main-fails-every-later-batch-at-integration.md.

## Entries

### 2026-10-30T18:30:48Z flai serve
Merge queue batch `20261030-182125` fails with S-0529 at `integration`. The queue verified S-0529 alone on main at `a59c5fb99d48`, at `a1fc087d13fe`, the copy of its verified head `45e640d96b73` rebased in the batch, and `integration` failed:

```text
flai/internal/mcpserver/verify_test.go:275 TestVerifyLeavesTheQueueTiersToTheMergeQueueWhileItCanBatch: all: passed true, steps rebase=passed sync=passed narrative=passed check=passed ok=reused late=passed, want rebase=passed sync=passed narrative=passed check=passed ok=passed late=passed
```

The whole output of each tier it ran is in `.flai-cache/merge-queue/batches/20261030-182125`.

flai moved S-0529 back to in-progress and took it out of the merge queue. Fix it on `story/S-0529`, run flai verify S-0529, move it to review, and enqueue it again with flai queue add S-0529.

### 2026-10-30T19:04:06Z flai serve
Resolved: S-0529 is back in review at 3ef6d4b45261, a new commit; it is batched again once it is enqueued
