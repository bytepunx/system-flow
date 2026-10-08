---
id: T-1329
type: task
nature: remediation
title: Document what flai task done commits, and close I-0104 and I-0108
status: backlog
parent: S-0322
owner: alex
created: 2026-10-08T04:32:45Z
updated: 2026-10-08T04:38:31Z
transitions: []
stream: S-0322
tags: [flai]
touches: [docs/users/flai.md, design/conventions/git.md, template/root/design/conventions/git.md, template/CHANGELOG.md, design/issues/I-0104-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md, design/issues/I-0108-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md, design/issues/summary.md]
after: [T-1328]
---
# T-1329 Document what flai task done commits, and close I-0104 and I-0108

## Work

Say what the new commit step does where agents and users read it. In `docs/users/flai.md` § Work items › Story branches › Closing a task, replace `git add -A` with the paths T-1330's ADR names, the paths left for other open tasks, and `-m` needed only when there is something to commit; say there too that `task_done` and the host method `task.done` take a message only then. In `design/conventions/git.md`, in the rule on closing each task, add that it commits the task's own paths and leaves another open task's, so the tasks of one layer close apart. Make the same change in `template/root/design/conventions/git.md`, and add an entry for it to `template/CHANGELOG.md`.

Then close both issues from the story's worktree, each with a reason naming T-1330's ADR, the commit step in `flai/internal/taskdone/taskdone.go`, and the T-1328 test that reproduces it:

```bash
flai issue close I-0104 --reason "<what fixed it>"
flai issue close I-0108 --reason "<what fixed it>"
```

Each rewrites `design/issues/summary.md`, so both closes are in this one task.

It waits for T-1328, so the issues close on a fix a test shows. It shares no path with T-1331 and runs beside it in the third layer.

## Done when

- `docs/users/flai.md`, `design/conventions/git.md`, and its template copy describe the behaviour T-1330's ADR records, and `template/CHANGELOG.md` has its entry.
- I-0104 and I-0108 are closed with reasons naming what fixed them, and `design/issues/summary.md` lists neither.
- `flai check --strict` and the markdown lint are clean.

## Notes

Drafted by the planner for S-0322, from S-0312's cancelled T-1285 and this story's first T-1329, which closed I-0108 alone.
