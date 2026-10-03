---
id: TH-0072
title: "Sync refuses uncommitted changes, so a task is committed before it is synced: may git.md say so?"
anchor:
  path: wip/kanban/stories/S-0197-agents-commit-each-task-after-flai-stream-sync-and-stream-sync-makes-conflicts-easy-to-resolve.md
  item: S-0197
status: resolved
participants: [agent-S-0197, alex]
created: 2026-10-03T00:39:27Z
updated: 2026-10-03T00:52:31Z
---

# TH-0072 Sync refuses uncommitted changes, so a task is committed before it is synced: may git.md say so?

On wip/kanban/stories/S-0197-agents-commit-each-task-after-flai-stream-sync-and-stream-sync-makes-conflicts-easy-to-resolve.md.

## Entries

### 2026-10-03T00:39:27Z agent-S-0197
S-0197's first criterion has `flai stream sync` refuse to run over uncommitted changes (today it uses `git rebase --autostash`). But the git convention, here and in the template, says to commit each task by running `flai stream sync` *first*, resolving, testing, and *then* committing the task's changes. Once sync refuses a dirty worktree, that order can't work: when a task finishes, its changes are uncommitted, so sync refuses.

What I recommend: the per-task cycle becomes **commit the task → `flai stream sync` → resolve any conflicts it reports, `git add` them, and `git rebase --continue` (which updates the task's commit) → run the tests for what the task changed → commit any fix they needed**. Before review: commit everything, sync, resolve, and run `scripts/close-out.sh`. That is what git does natively: a rebase replays commits, and resolving a conflict happens on committed work, so nothing is ever stashed or lost.

This means editing a baseline rule in `design/conventions/git.md` and `template/root/design/conventions/git.md` (the "Commit each task" bullet and the "Before moving a story to review" bullet), which the story doesn't list. ADR-0069's decision sentence says "runs flai stream sync, resolves…, runs the tests…, and commits", so I'd note in the narrative that the order changed, not the ADR (accepted ADRs are immutable).

May I (a) reword git.md's two bullets in both copies to commit → sync → resolve → test → commit fixes, and use the same order in the harness prompt, the work-management convention, close-out.sh, and the design/docs? Or (b) do you want a different order, for example sync at the *start* of each task on a clean worktree?

### 2026-10-03T00:52:11Z alex
a

### 2026-10-03T00:52:31Z agent-S-0197
Resolved: Designer chose (a): the per-task cycle is commit → flai stream sync → resolve (git add, git rebase --continue) → test → commit fixes; git.md's two bullets change in both copies; before review commit, sync, resolve, close-out.
