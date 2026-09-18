---
id: T-0154
type: task
nature: remediation
title: Opt-in setting for relative-path worktrees in flai stream open
status: done
parent: S-0050
owner: alex
created: 2026-09-18T19:47:24Z
updated: 2026-09-18T19:53:38Z
transitions:
  - to: ready
    at: 2026-09-18T19:51:46Z
    by: alex
  - to: in-progress
    at: 2026-09-18T19:51:46Z
    by: alex
  - to: done
    at: 2026-09-18T19:53:38Z
    by: alex
stream: S-0050
tags: []
touches: [flai/cmd, flai/internal/config]
---

# T-0154 Opt-in setting for relative-path worktrees in flai stream open

## Work
Add a per-user flai config key (proposed `worktrees.relative_paths`, boolean, default false) to `flai/internal/config` and `flai config`. In `openStoryBranch`, when the key is on: read the git version; at 2.48 or newer pass `--relative-paths` to `git worktree add`; below that, log a warning naming the version and the key and create an ordinary worktree. With the key off, pass nothing new. Parse the version from `git version` through the `execx.Runner`. Tests with a fake runner for off, on with new git, on with old git; a real-git test that the worktree is relative, which skips with a stated reason when the git on PATH is older than 2.48.

## Done when
- The fake-runner tests pass and pin the arguments
- The real-git test passes where git is 2.48 or newer and skips here with the reason printed
- `flai config` shows and sets the key

## Notes
