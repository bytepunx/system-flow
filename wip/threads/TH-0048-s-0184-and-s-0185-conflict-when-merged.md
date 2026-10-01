---
id: TH-0048
title: S-0184 and S-0185 conflict when merged
anchor:
  path: wip/kanban/stories/S-0185-flai-new-refuses-an-empty-required-variable-before-writing-and-flai-upgrade-renders-a-fork-s-own-template-variables.md
  item: S-0185
status: answered
participants: [flai, agent-S-0185]
created: 2026-10-01T09:10:18Z
updated: 2026-10-01T09:10:22Z
---

# TH-0048 S-0184 and S-0185 conflict when merged

On wip/kanban/stories/S-0185-flai-new-refuses-an-empty-required-variable-before-writing-and-flai-upgrade-renders-a-fork-s-own-template-variables.md.

## Entries

### 2026-10-01T09:10:18Z flai
A trial merge of story/S-0184 with story/S-0185 at flai stream sync conflicts in:

- `design/issues/summary.md`

Whichever of S-0184 and S-0185 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-01T09:10:22Z agent-S-0185
This is the same kind of conflict as TH-0045, which S-0185 resolved this way after S-0179 landed. The only conflict is `design/issues/summary.md`, which flai regenerates from the issue files. S-0184 closes I-0025 and I-0029, and S-0185 closes I-0040 and I-0041, so the issue files themselves don't conflict.

Whichever story rebases second resolves it by running `flai issue summary` and `git add design/issues/summary.md`, then `git rebase --continue`, with no change to either story's scope. S-0185 is about to move to review. If S-0184 is accepted first, S-0185 will do this at its next sync.
