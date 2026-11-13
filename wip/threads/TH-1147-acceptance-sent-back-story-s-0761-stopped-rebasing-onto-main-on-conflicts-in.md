---
id: TH-1147
title: "Acceptance sent back: story/S-0761 stopped rebasing onto main on conflicts in design/system/flai-cli.md"
anchor:
  path: wip/kanban/stories/S-0761-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md
  item: S-0761
status: open
participants: [flai]
created: 2026-11-13T09:19:42Z
updated: 2026-11-13T09:19:42Z
---

# TH-1147 Acceptance sent back: story/S-0761 stopped rebasing onto main on conflicts in design/system/flai-cli.md

On wip/kanban/stories/S-0761-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md.

## Entries

### 2026-11-13T09:19:42Z flai
The acceptance of S-0761 by orchestrator rebased `story/S-0761` onto `main`, and the rebase stopped on conflicts that the merge of the issue files, the ADR index, the issue summary, and the shared markdown files did not settle, in:

- `design/system/flai-cli.md`

flai aborted the rebase, so `story/S-0761` and `main` are as they were, and moved S-0761 back to in-progress. This thread is for S-0761's agent, which settles the conflicts; the operator does not. In its worktree, run `flai stream sync S-0761`, resolve each path it stops on, `git add` it, and run `git rebase --continue`; then close out and move S-0761 to review again. flai serve resolves this thread once S-0761 is back in review.
