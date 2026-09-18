---
id: ADR-0022
title: The dashboard mounts the repository at its host path; relative-path worktrees are an opt-in
status: accepted
date: 2026-09-18
supersedes: []
superseded_by: []
---

# ADR-0022 The dashboard mounts the repository at its host path; relative-path worktrees are an opt-in

## Context

[ADR-0019](0019-story-branches-and-touches.md) works each story on a branch in a git worktree under `.flai-cache/worktrees/`. Git links a worktree to its repository with two absolute paths: the worktree's `.git` file names `<repo>/.git/worktrees/<id>`, and that folder's `gitdir` names the worktree. `flai dashboard` mounted the repository in the flaiover container at `/project`, so inside the container neither path existed. Every git command run in a worktree failed with `fatal: not a git repository`, and accepting a story from the board failed at its preview, which plans the release from the worktree (I-0017). Git in the container also listed the host's worktrees as prunable.

Git 2.48 added `git worktree add --relative-paths`, which writes both links relative. Observed with git 2.54.0: the worktree then works wherever the repository is mounted and after the repository folder is moved, and creating it changes the clone from `core.repositoryformatversion=0` with no extensions to version `1` with `extensions.relativeWorktrees=true`. Observed with git 2.47.3: a repository carrying that extension is refused outright, `fatal: unknown repository extension found: relativeworktrees`. The refusal covers the whole clone, not the one worktree.

The operator proposed choosing between the two by the git version found. The version found is the git that runs `flai stream open`; the extension binds every git that touches the clone afterwards, including tools that bundle or reimplement git. A host upgrade from 2.47 to 2.48 would change the repository format at the next story with nobody having chosen it, and a later downgrade or an older tool would lose the repository. Worktrees made before such an upgrade keep absolute paths, so the mount is needed at every version anyway.

## Decision

`flai dashboard` mounts the repository in the container at the same absolute path it has on the host and sets `PROJECT_DIR` to that path, at every git version and with no version check. Absolute paths git recorded on the host therefore resolve in the container, for worktrees that already exist as well as new ones. When the host path cannot be a container path, a Windows drive path being the case in view, `flai dashboard` keeps the `/project` mount, says at start that stories with a branch cannot be accepted from the board, and names the two ways out: the opt-in below, or `flai accept` from a shell.

Relative-path worktrees are an explicit setting, the per-user flai config key `worktrees.relative_paths`, off by default. flai never turns it on by detecting a git version. With it on and git 2.48 or newer, `flai stream open` creates the worktree with `--relative-paths`. With it on and an older git, `flai stream open` warns, naming the version and the key, and creates an ordinary worktree, which the mount covers. `flai check` warns when the clone carries `extensions.relativeWorktrees` and the git on `PATH` is older than 2.48, reading the config file because that git cannot open the repository.

## Consequences

- Acceptance from the board works for stories with a branch, which is every story since ADR-0019, with whatever git the host has.
- The container's project path is no longer the constant `/project`. flaiover already reads `PROJECT_DIR` and derives flai's config and cache paths from it; anything that assumed `/project` is corrected in S-0050, and the hub work in S-0043 must treat the path as a per-project value. The image keeps `/project` as its default for anyone running it by hand.
- Absolute paths inside `.flai-cache/config.json`, such as `cache_dir`, now resolve in the container too.
- The repository format changes only when someone sets the key. The setting is per user because the format is local to a clone; a committed manifest field would impose it on every contributor whatever their git.
- A user who opts in takes on a constraint: every git that touches that clone must be 2.48 or newer. The documentation says so and says how to go back.
- Two worktree layouts exist, so the opt-in's real-git tests run only where git is 2.48 or newer and skip with a stated reason elsewhere; CI runs them.

## Alternatives considered

- Choose the strategy by detecting the git version: a silent repository format change on a package upgrade, decided by one git out of the several that touch the clone, and the mount is needed anyway.
- Relative paths only: locks out any host whose git is older than 2.48, including the one this was found on.
- Mount at the host path only, no opt-in: leaves hosts whose path cannot be mirrored into the container with no way to accept from the board.
- Rewrite the worktree's links inside the container, or run `git worktree repair` there: the links are shared with the host through the mount, so repairing for one side breaks the other.
- Stop planning the release from the worktree: fixes the preview and not the rebase, merge, sync, and removal that also run in it.
