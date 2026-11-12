---
id: TH-1088
title: "S-0729's close-out stopped at sync: the branch is 8 commits behind main"
anchor:
  path: wip/kanban/stories/S-0729-publish-what-is-accepted-and-not-yet-released.md
  item: S-0729
status: resolved
participants: [orchestrator, alex]
created: 2026-11-12T03:44:57Z
updated: 2026-11-12T03:45:20Z
---

# TH-1088 S-0729's close-out stopped at sync: the branch is 8 commits behind main

On wip/kanban/stories/S-0729-publish-what-is-accepted-and-not-yet-released.md.

## Entries

### 2026-11-12T03:44:57Z orchestrator
**What's missing:** a passing close-out at the head of `story/S-0729`. Until there is one, I cannot accept the story.

- Its first close-out passed at 9d16de6. The head is now 6df18e9, the commit that close-out made.
- I asked for the close-out again at 6df18e9 (job S-0729-20261112T034354Z). It stopped at `sync`. Meanwhile the operator had accepted S-0717, S-0701, and S-0681, so the branch does not contain main. It is 8 commits behind, and they change paths outside `wip/`, among them `design/issues/I-0347-…md` and `design/issues/summary.md`, which this story also changes.
- `flai accept --dry-run` names no other blocker. The three steps T-3132, T-3133, and T-3134 are done, and v1.86.21 is tagged and pushed.

**What clears it:** run `flai stream sync S-0729`, settle the issue file and `summary.md` against main, and ask for the close-out again with `verify`. Once it passes at the head, I accept the story on my next run.

### 2026-11-12T03:45:20Z alex
Resolved: S-0729 was accepted
