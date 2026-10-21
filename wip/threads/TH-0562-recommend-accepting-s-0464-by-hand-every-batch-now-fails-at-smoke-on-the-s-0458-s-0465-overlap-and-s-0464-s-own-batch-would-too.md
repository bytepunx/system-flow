---
id: TH-0562
title: "Recommend accepting S-0464 by hand: every batch now fails at smoke on the S-0458/S-0465 overlap, and S-0464's own batch would too"
anchor:
  path: wip/kanban/stories/S-0464-the-merge-queue-s-batch-runs-smoke-with-no-close-out-story-so-scripts-check-sh-runs-an-unscoped-flai-check-strict-and-fails-every-batch-on-findings-outside-its-stories.md
  item: S-0464
status: resolved
participants: [orchestrator, alex]
created: 2026-10-21T08:19:29Z
updated: 2026-10-21T08:22:46Z
---

# TH-0562 Recommend accepting S-0464 by hand: every batch now fails at smoke on the S-0458/S-0465 overlap, and S-0464's own batch would too

On wip/kanban/stories/S-0464-the-merge-queue-s-batch-runs-smoke-with-no-close-out-story-so-scripts-check-sh-runs-an-unscoped-flai-check-strict-and-fails-every-batch-on-findings-outside-its-stories.md.

## Entries

### 2026-10-21T08:19:29Z orchestrator
Recommended: accept S-0464 yourself at 3b3eb4c6. I would then publish it and you would restart flai serve on that release. Until then, keep stories out of the queue.

Why:
- The queue sent back S-0456 (TH-0560, 08:12Z) and S-0463 (TH-0561, 08:19Z) at smoke. Both verifies passed. Each batch failed on the same warning, which is in neither story: `wip.overlap: S-0458 touches design/system/flaiover-dashboard.md and docs/users/flaiover.md, which S-0465 (in progress) also touches`. That overlap is S-0458 started beside S-0465 (MS-0162). It lasts until one of them leaves in-progress.
- S-0464 is the fix (I-0174, ADR-0183). But its scoping only works when flai serve's batch sets MERGE_QUEUE_STORIES. The serve running now is the installed 1.67.6, which does not set it. So S-0464's own batch would run the same unscoped check and fail the same way, and be sent back.
- S-0464 is ready otherwise:
  - flai verify passed every step at its head 3b3eb4c6.
  - `flai accept --dry-run` names no blocker.
  - The verifier matched criterion 1 to the batch and bisection code and tests (TestTheBatchsTiersAreGivenItsStoriesAndNoCloseOut, TestCheckStoryTakesSeveralStories), plus scripts/check.sh, lint-md.sh and ADR-0183, which is accepted.
  - It matched criterion 2 to I-0174, which is closed.
  - The verifier noted one non-blocking gap: template/root/scripts/lint-md.sh is unchanged.

Alternatives:
1. Enqueue S-0464 anyway. Its batch likely fails and sends it back.
2. Wait until S-0465 leaves in-progress. The overlap then ends and batches can pass again, S-0464's included.

I have not enqueued S-0464. The merge queue refuses my direct accept.

### 2026-10-21T08:22:18Z alex
Resolved: S-0464 was accepted

### 2026-10-21T08:22:43Z orchestrator
Thank you. Released as flai 1.67.7 (tag flai/v1.67.7) and template 1.0.107, both pushed. The batches scope to their own stories only once flai serve runs 1.67.7. That means `flai self-upgrade` and a restart of serve, which is yours to do. Until then, a batch can still fail on the S-0458/S-0465 overlap, which lasts while both are in progress. S-0456 and S-0463 were sent back on it and will need enqueuing again.

### 2026-10-21T08:22:46Z orchestrator
Resolved: The operator accepted S-0464 at 08:22Z, and it was released as flai 1.67.7 and template 1.0.107.
