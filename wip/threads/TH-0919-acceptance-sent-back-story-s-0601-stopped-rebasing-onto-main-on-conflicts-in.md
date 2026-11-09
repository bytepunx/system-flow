---
id: TH-0919
title: "Acceptance sent back: story/S-0601 stopped rebasing onto main on conflicts in design/system/server-performance.md"
anchor:
  path: wip/kanban/stories/S-0601-flai-check-runs-its-rules-at-once-and-keeps-the-component-props-node-read-so-a.md
  item: S-0601
status: resolved
participants: [flai, flai serve]
created: 2026-11-09T08:07:20Z
updated: 2026-11-09T08:21:47Z
---

# TH-0919 Acceptance sent back: story/S-0601 stopped rebasing onto main on conflicts in design/system/server-performance.md

On wip/kanban/stories/S-0601-flai-check-runs-its-rules-at-once-and-keeps-the-component-props-node-read-so-a.md.

## Entries

### 2026-11-09T08:07:20Z flai
The acceptance of S-0601 by alex rebased `story/S-0601` onto `main`, and the rebase stopped on conflicts that the merge of the issue files, the ADR index, the issue summary, and the shared markdown files did not settle, in:

- `design/system/server-performance.md`

flai aborted the rebase, so `story/S-0601` and `main` are as they were, and moved S-0601 back to in-progress. This thread is for S-0601's agent, which settles the conflicts; the operator does not. In its worktree, run `flai stream sync S-0601`, resolve each path it stops on, `git add` it, and run `git rebase --continue`; then close out and move S-0601 to review again. flai serve resolves this thread once S-0601 is back in review.

### 2026-11-09T08:21:47Z flai serve
Resolved: S-0601 is back in review since 2026-11-09T08:21:37Z, after the send back; its next acceptance takes it up again
