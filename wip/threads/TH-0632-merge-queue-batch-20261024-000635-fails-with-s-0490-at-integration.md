---
id: TH-0632
title: Merge queue batch 20261024-000635 fails with S-0490 at integration
anchor:
  path: wip/kanban/stories/S-0490-flai-test-runs-go-test-with-short-so-a-real-git-test-a-task-adds-is-skipped-until-the-integration-tier-runs-in-the-merge-queue-s-batch.md
  item: S-0490
status: answered
participants: [flai serve, alex, agent-S-0490]
created: 2026-10-24T00:15:07Z
updated: 2026-10-24T00:24:57Z
---

# TH-0632 Merge queue batch 20261024-000635 fails with S-0490 at integration

On wip/kanban/stories/S-0490-flai-test-runs-go-test-with-short-so-a-real-git-test-a-task-adds-is-skipped-until-the-integration-tier-runs-in-the-merge-queue-s-batch.md.

## Entries

### 2026-10-24T00:15:07Z flai serve
Merge queue batch `20261024-000635` fails with S-0490 at `integration`. The queue verified S-0490 alone on main at `be3af0d161ce`, at `c3e46ddc8bef`, the copy of its verified head `c3e46ddc8bef` rebased in the batch, and `integration` failed:

```text
flai/cmd/reference_test.go:127 TestReferenceIsCurrent: docs/users/flai-reference.md is stale: run make flai-reference
```

The whole output of each tier it ran is in `.flai-cache/merge-queue/batches/20261024-000635`.

flai moved S-0490 back to in-progress and took it out of the merge queue. Fix it on `story/S-0490`, run flai verify S-0490, move it to review, and enqueue it again with flai queue add S-0490.

### 2026-10-24T00:19:44Z alex
Resolved.

### 2026-10-24T00:24:57Z agent-S-0490
Fixed on `story/S-0490`. The failure came from main, not from S-0490: `docs/users/flai-reference.md` has been stale on main since S-0491 added `orchestration.permissions.edit_subprojects` to `flai manifest set`'s help. S-0491 was accepted outside the queue, so `integration` never ran on it, and every batch fails until main has the regenerated file.

- T-2104 regenerated the reference (`f4a069b2c`).
- `flai test` on `flai/cmd` passes, `TestReferenceIsCurrent` included (`go-test` 7m28s).
- The close-out passed again. It left `integration` and `smoke` to the queue.
- I-0203 records the cause.

S-0490 moves to review now; it needs `flai queue add S-0490` to be batched again.
