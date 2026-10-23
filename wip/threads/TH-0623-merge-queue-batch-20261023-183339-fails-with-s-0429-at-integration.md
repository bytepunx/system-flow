---
id: TH-0623
title: Merge queue batch 20261023-183339 fails with S-0429 at integration
anchor:
  path: wip/kanban/stories/S-0429-flai-digest-harvest-lists-what-the-archive-and-the-threads-hold-since-a-date-and-a-one-off-seed-fills-the-digest-from-the-operator-s-answers-and-the-closed-issues-so-there-is-something-to-prime-with.md
  item: S-0429
status: resolved
participants: [flai serve]
created: 2026-10-23T18:40:32Z
updated: 2026-10-23T18:59:25Z
---

# TH-0623 Merge queue batch 20261023-183339 fails with S-0429 at integration

On wip/kanban/stories/S-0429-flai-digest-harvest-lists-what-the-archive-and-the-threads-hold-since-a-date-and-a-one-off-seed-fills-the-digest-from-the-operator-s-answers-and-the-closed-issues-so-there-is-something-to-prime-with.md.

## Entries

### 2026-10-23T18:40:32Z flai serve
Merge queue batch `20261023-183339` fails with S-0429 at `integration`. The queue verified S-0429 alone on main at `e38be89d285b`, at `5693210c55b6`, the copy of its verified head `5693210c55b6` rebased in the batch, and `integration` failed:

```text
flai/cmd/stream_sync_test.go:263 TestSyncListsPathsChangedOutsideTheClaim: {"ts":"2026-10-23T18:40:01.477Z","level":"FATAL","msg":"command failed","component":"cmd","err":"T-0001's touch design is a folder holding 1 file: a task takes no folder touch (ADR-0188); name the files the task changes, or leave the folder to the story's touches with a ### Planning note"}
```

The whole output of each tier it ran is in `.flai-cache/merge-queue/batches/20261023-183339`.

flai moved S-0429 back to in-progress and took it out of the merge queue. Fix it on `story/S-0429`, run flai verify S-0429, move it to review, and enqueue it again with flai queue add S-0429.

### 2026-10-23T18:59:25Z flai serve
Resolved: S-0429 is back in review at dae76a3465d4, a new commit; it is batched again once it is enqueued
