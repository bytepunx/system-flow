---
id: TH-0942
title: "Acceptance sent back: story/S-0616 stopped rebasing onto main on conflicts in design/workflows/formulas/README.md, template/root/design/workflows/formulas/README.md"
anchor:
  path: wip/kanban/stories/S-0616-testrenderprototypetemplate-refuses-any-double-curly-brace-in-a-rendered.md
  item: S-0616
status: resolved
participants: [flai, flai serve]
created: 2026-11-09T21:53:09Z
updated: 2026-11-09T22:09:58Z
---

# TH-0942 Acceptance sent back: story/S-0616 stopped rebasing onto main on conflicts in design/workflows/formulas/README.md, template/root/design/workflows/formulas/README.md

On wip/kanban/stories/S-0616-testrenderprototypetemplate-refuses-any-double-curly-brace-in-a-rendered.md.

## Entries

### 2026-11-09T21:53:09Z flai
The acceptance of S-0616 by alex rebased `story/S-0616` onto `main`, and the rebase stopped on conflicts that the merge of the issue files, the ADR index, the issue summary, and the shared markdown files did not settle, in:

- `design/workflows/formulas/README.md`
- `template/root/design/workflows/formulas/README.md`

flai aborted the rebase, so `story/S-0616` and `main` are as they were, and moved S-0616 back to in-progress. This thread is for S-0616's agent, which settles the conflicts; the operator does not. In its worktree, run `flai stream sync S-0616`, resolve each path it stops on, `git add` it, and run `git rebase --continue`; then close out and move S-0616 to review again. flai serve resolves this thread once S-0616 is back in review.

### 2026-11-09T22:09:58Z flai serve
Resolved: S-0616 is back in review since 2026-11-09T22:09:30Z, after the send back; its next acceptance takes it up again
