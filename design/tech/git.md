---
title: Git
updated: 2026-09-18
status: active
---

# Git

| | |
|-|-|
| Host requirement | Any git with worktrees on `PATH`; the dev host has 2.47.3. flai shells out to it ([ADR 0010](../adrs/0010-shell-out-to-git-and-docker.md)) |
| Image | The flaiover image installs Alpine's git, 2.54.0 at the time of writing, for acceptance from the dashboard |
| Used for | Story branches and worktrees (`flai stream open`, `sync`, `accept`), release tags and pushes, template clones, `flai template push` |
| Version-dependent | One thing: `worktrees.relative_paths`, which needs 2.48 for `git worktree add --relative-paths`. `flai/internal/gitver` reads the version; nothing else in flai branches on it |

## Why

Worktrees keep each story's changes on a branch while `wip/` stays on the main checkout ([ADR 0019](../adrs/0019-story-branches-and-touches.md)). Git links a worktree to its repository with absolute paths, so the dashboard mounts the repository at its host path ([ADR 0022](../adrs/0022-repository-mounted-at-its-host-path.md)).

Relative links are the better mechanism where every git that touches the clone is 2.48 or newer, and they are an explicit setting for that reason: creating one sets the repository extension `relativeWorktrees`, and an older git then refuses the whole clone. flai does not raise its minimum git for it.

## Alternatives

A Go git library (go-git) was rejected in [ADR 0010](../adrs/0010-shell-out-to-git-and-docker.md): no credential helper support, partial clone support, a larger binary.
