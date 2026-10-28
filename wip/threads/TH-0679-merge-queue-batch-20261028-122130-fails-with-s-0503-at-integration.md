---
id: TH-0679
title: Merge queue batch 20261028-122130 fails with S-0503 at integration
anchor:
  path: wip/kanban/stories/S-0503-no-test-tier-selects-scripts-so-a-story-that-changes-only-shell-scripts-and-the-script-tests-it-adds-pass-flai-verify-and-the-merge-queue-without-any-of-them-running.md
  item: S-0503
status: open
participants: [flai serve]
created: 2026-10-28T12:33:01Z
updated: 2026-10-28T12:33:01Z
---

# TH-0679 Merge queue batch 20261028-122130 fails with S-0503 at integration

On wip/kanban/stories/S-0503-no-test-tier-selects-scripts-so-a-story-that-changes-only-shell-scripts-and-the-script-tests-it-adds-pass-flai-verify-and-the-merge-queue-without-any-of-them-running.md.

## Entries

### 2026-10-28T12:33:01Z flai serve
Merge queue batch `20261028-122130` fails with S-0503 at `integration`. The queue verified S-0503 alone on main at `d94b60b545ef`, at `947ecb0aed0a`, the copy of its verified head `f4a53a529055` rebased in the batch, and `integration` failed:

```text
flai/internal/serve/agents_test.go:3123 TestAStoryWhoseAgentEndedAskingIsBlockedUntilTheAnswerStartsItAgain: the story idles in progress, blocked []
```

The whole output of each tier it ran is in `.flai-cache/merge-queue/batches/20261028-122130`.

flai moved S-0503 back to in-progress and took it out of the merge queue. Fix it on `story/S-0503`, run flai verify S-0503, move it to review, and enqueue it again with flai queue add S-0503.
