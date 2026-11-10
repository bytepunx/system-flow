---
id: T-2829
type: task
nature: remediation
title: Remove the stray open copy of T-2760 that the archive race left in wip/kanban/tasks
status: backlog
parent: S-0648
owner: alex
created: 2026-11-10T06:37:35Z
updated: 2026-11-10T06:37:35Z
transitions: []
stream: S-0648
tags: [wip]
touches: [wip/kanban/tasks/T-2760-fetch-the-remote-s-branch-and-tags.md]
---
# T-2829 Remove the stray open copy of T-2760 that the archive race left in wip/kanban/tasks

## Work

First layer: it waits for no task and shares no path with T-2827, so the two can run together.

- Delete `wip/kanban/tasks/T-2760-fetch-the-remote-s-branch-and-tags.md` with `git rm` in the story's worktree. Its archived copy, `wip/archive/kanban/tasks/T-2760-fetch-the-remote-s-branch-and-tags.md`, is the one S-0642's acceptance wrote. Leave it as it is: `wip/archive` is not edited.
- The open copy carries the harness-priced usage, cost 0.3414, where the archived copy carries the estimate, cost 0.2198. Write both figures in the task's Notes, so the 0.1216 USD the archive does not count is on record.
- Before deleting, check that the main checkout's copy is still the same as at commit 436e8f612. If anything has written it since, say what in the Notes.

## Done when

- `git ls-files wip/kanban/tasks | grep T-2760` prints nothing on the story branch.
- `flai check --strict` run in the story's worktree reports no `item.duplicate-id`.
- The two usage figures are in the task's Notes.

## Notes
