---
id: ADR-0128
title: "flai task done commits the paths the closing task covers and those no other open task covers, leaves another open task's, and needs a message only when it commits"
status: accepted
date: 2026-10-08
supersedes: []
superseded_by: []
refines: [ADR-0107]
topics: [cli, conventions]
---

# ADR-0128 flai task done commits the paths the closing task covers and those no other open task covers, leaves another open task's, and needs a message only when it commits

## Context

[ADR-0107](0107-flai-task-done-closes-a-task-in-one-call-commit-sync-move-log-widen-touches.md) made `flai task done` close a task in one call. Its first step runs `git add -A` and `git commit -m` in the story's worktree, so it commits every uncommitted file there. A story's agent often runs the tasks of one layer together in one worktree and closes each when it is finished. Closing the first commits every sibling's files under its message, and the touches step widens the closing task's touches with the siblings' paths. The siblings then close with nothing to commit, and `-m` is still required for each.

Two issues count it:

- [I-0104](../issues/I-0104-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md): S-0212, S-0215, S-0298, S-0328, and S-0314. Each agent amended a message, committed tasks' files by hand before closing them with nothing to commit, or reset touches with `flai touches`.
- [I-0108](../issues/I-0108-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md): S-0274's layers 3 and 4, and S-0309's layer 2, where three tasks' edits went into one commit.

The operator confirmed the remedy below on TH-0337, when it was planned for S-0312, since cancelled as a duplicate of S-0322.

## Decision

`flai task done` commits the changed paths the closing task covers and those no other open task of its story covers, leaves uncommitted the paths only another open task covers and reports them, and needs a message only when there is something to commit.

1. **Commit.** Each path changed in the story's worktree is sorted by the `touches` of the story's tasks. An open task is one whose status is neither `done` nor `cancelled`. A touch covers a path as touches cover a path elsewhere in flai: it names the file or a folder above it.
   - A path the closing task covers is committed, even when another open task covers it too: it goes with the task that declared it and closed first.
   - A path no other open task covers is committed, so a file the task added without declaring it goes with it.
   - A path that only another open task covers is left uncommitted.

   The step stages and commits the paths as `storygit.CommitPaths` does, with `git add -A -- <paths>` and `git commit -- <paths>`, and leaves everything else in the worktree as it was.
2. **Paths left.** The answer lists the paths left uncommitted, so the agent sees what waits for another task's close. A path left is not a failure: the run goes on. The sync refuses a worktree with uncommitted changes; when those changes are only the paths left, its refusal does not stop the run, and the branch is synced by the close that leaves none.
3. **Message.** `-m` is needed only when there is something to commit. With nothing to commit the close goes on without it. The narrative log entry is then `--log` when given, and otherwise says the task was closed: `Closed T-nnnn: <title>`. The MCP tool `task_done` and the host method `task.done` take the message on the same terms, and `task.done` passes it to `flai task done` only when one is given.
4. **Touches.** The touches step widens the task's and the story's touches with the committed paths only. A path left never reaches the closing task's touches.

Every other step of ADR-0107, and its order, stands, the sync's refusal for the paths left aside.

## Consequences

- The tasks of one layer close apart in one worktree, each with a commit of its own paths under its own message, and with no hand commit or touches reset.
- A file no task declares goes with the first task to close after it is written, and widens that task's touches.
- A path two open tasks both declare goes with the first of them to close; the second may close with nothing of it left to commit.
- The agent reads the paths left in each answer. The story's last open task has no other open task beside it, so its close commits every changed path.
- A task whose touches are wrong leaves its own files for a sibling or takes a sibling's: the touches the planner writes now decide each commit.

## Alternatives considered

- **Commit only the closing task's touches.** A file the task added without declaring it would be left behind with no task to take it.
- **A `--paths` option listing what to commit.** The agent would write out what the task's touches already say, on every close.
- **Keep committing everything and ask agents to close each task before the next is written.** It forbids running the tasks of a layer together, which sub-agents do.
