---
id: T-3157
type: task
nature: improvement
title: The user docs and the git and work-management conventions, here and in the template, say to git add a resolved stop and run flai stream sync again
status: backlog
parent: S-0726
owner: alex
created: 2026-11-12T15:56:52Z
updated: 2026-11-12T15:56:52Z
transitions: []
stream: S-0726
tags: [flai, template]
touches: [docs/users/flai.md, design/conventions/git.md, template/root/design/conventions/git.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md]
after: [T-3154]
---
# T-3157 The user docs and the git and work-management conventions, here and in the template, say to git add a resolved stop and run flai stream sync again

## Work

Replace "resolve each path, `git add` it, and run `git rebase --continue`" with "resolve each path, `git add` it, and run `flai stream sync` again" (or call `task_done` again, or ask for the close-out again, as each place says), as T-3154's ADR decides:

- `docs/users/flai.md`: § Story branches, around lines 869–890, with its "To continue:" example, and § Closing a task, around line 970.
- `design/conventions/git.md`, the Rules bullets on `flai stream sync` and `flai task done`, and `template/root/design/conventions/git.md` the same. Changes to conventions land in both in this story, as `CLAUDE.md` says.
- `design/conventions/work-management.md`, the close-out bullet's "A sync that stops on a conflict stops the close-out" sentence, and `template/root/design/conventions/work-management.md` the same.
- Add a line under `template/CHANGELOG.md`'s top `## Unreleased` for the template's convention change.
- Keep `git rebase --abort` as the way back wherever it is named.

Waits for T-3154, whose ADR is what these documents describe. No path in common with T-3155, so the two run together in layer 2.

## Done when

- None of these files tells an agent to run `git rebase --continue` after a sync's stop.
- Each copy under `template/root/` matches its `design/conventions/` file in the changed bullets.
- The markdown lint passes on them with `flai test`.

## Notes
