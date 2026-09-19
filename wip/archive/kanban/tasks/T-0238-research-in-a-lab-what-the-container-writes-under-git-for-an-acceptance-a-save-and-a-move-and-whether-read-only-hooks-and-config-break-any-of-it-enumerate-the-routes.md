---
id: T-0238
type: task
nature: feature
title: "Research in a lab: what the container writes under .git for an acceptance, a save, and a move, and whether read-only hooks and config break any of it; enumerate the routes"
status: done
parent: S-0064
owner: alex
created: 2026-09-19T10:18:10Z
updated: 2026-09-19T10:22:07Z
transitions:
  - to: ready
    at: 2026-09-19T10:18:11Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T10:18:11Z
    by: system-flow
  - to: done
    at: 2026-09-19T10:22:07Z
    by: system-flow
stream: S-0064
tags: []
---

# T-0238 Research in a lab: what the container writes under .git for an acceptance, a save, and a move, and whether read-only hooks and config break any of it; enumerate the routes

## Work
In the scratchpad, never this repository: a scratch project with a story branch in a worktree, run the published image as `flai dashboard` runs it, and record which paths under `.git` change for an acceptance with a merge and a worktree removal, a document save, a board move, and a push with a key. Then mount `.git/hooks` and `.git/config` read-only and repeat, noting what fails. Try each route by which the container could leave something that git later executes on the host: hooks, `core.hooksPath`, `core.fsmonitor`, `core.sshCommand`, editor and pager, credential helpers, aliases, includes, `url.*.insteadOf` and the remote URLs, filter, diff, and merge drivers through config and `.git/info/attributes`, per-worktree config, and anything else found. Separate the routes git never shows (outside the work tree, or ignored) from changes to tracked files, which `git status` shows and acceptance refuses to include unasked.

## Done when
- The routes are written down with how each was tried and whether read-only mounts close it
- What acceptance needs to write under `.git` is known

## Notes
