---
id: T-1285
type: task
nature: improvement
title: Document what flai task done commits and close I-0104
status: cancelled
parent: S-0312
owner: alex
created: 2026-10-07T23:34:54Z
updated: 2026-10-08T04:33:37Z
transitions:
  - to: cancelled
    at: 2026-10-08T04:33:37Z
    by: alex
stream: S-0312
tags: [flai]
touches: [docs/users/flai.md, design/conventions/git.md, template/root/design/conventions/git.md, template/CHANGELOG.md, design/issues/I-0104-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md, design/issues/summary.md]
after: [T-1283]
---
# T-1285 Document what flai task done commits and close I-0104

## Work

Say what the new commit step does where agents and users read it. In `docs/users/flai.md` § Work items › Story branches › Closing a task, replace `git add -A` with the paths T-1282's ADR names, the paths left for other open tasks, and `-m` needed only when there is something to commit. In `design/conventions/git.md`, in the rule on closing each task, add that it commits the task's own paths and leaves another open task's, so the tasks of one layer close apart; make the same change in `template/root/design/conventions/git.md` and add an entry for it to `template/CHANGELOG.md`.

Then close the issue from the story's worktree: `flai issue close I-0104 --reason "<what fixed it>"`, naming the ADR and the commit step, which rewrites `design/issues/summary.md`.

It waits for T-1283, so the issue closes on a fix that exists. It shares no path with T-1284 and runs beside it in the third layer.

## Done when

- `docs/users/flai.md`, `design/conventions/git.md`, and its template copy describe the behaviour T-1282's ADR decides, and `template/CHANGELOG.md` has its entry.
- I-0104 is closed with a reason naming what fixed it, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` and the markdown lint are clean.

## Notes

Drafted by the planner for S-0312.
- 2026-10-08T04:33:37Z: moved to cancelled: S-0312 cancelled: Duplicate, cancel this and keep 322
