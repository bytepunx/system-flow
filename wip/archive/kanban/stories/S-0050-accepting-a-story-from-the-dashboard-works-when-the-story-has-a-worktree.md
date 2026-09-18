---
id: S-0050
type: story
nature: remediation
title: Accepting a story from the dashboard works when the story has a worktree
status: done
parent: E-0006
owner: alex
created: 2026-09-18T18:36:10Z
updated: 2026-09-18T20:25:58Z
transitions:
  - to: ready
    at: 2026-09-18T19:45:54Z
    by: alex
  - to: in-progress
    at: 2026-09-18T19:46:37Z
    by: alex
  - to: review
    at: 2026-09-18T20:19:29Z
    by: alex
  - to: done
    at: 2026-09-18T20:25:58Z
    by: alex
tags: [cli, dashboard]
touches: [flai/cmd, flai/internal/check, flaiover, design/adrs]
---

# S-0050 Accepting a story from the dashboard works when the story has a worktree

## Goal
Dropping a story on done, or pressing accept on its page, works from the dashboard when the story was worked on a branch in a worktree, which is every story since ADR-0019. Today it fails before the confirmation is shown. The fix is unconditional: `flai dashboard` mounts the repository at its host path, so the absolute paths git records for a worktree resolve in the container with any git version. Relative-path worktrees are added as an explicit opt-in for setups the mount cannot cover.

## Acceptance criteria
- [x] `flai dashboard` always mounts the repository in the container at the same absolute path it has on the host and points `PROJECT_DIR` at it, whatever git version is installed; nothing about the mount depends on a version check
- [x] With flaiover running through `flai dashboard`, the acceptance preview for a story in review that has a `story/S-nnnn` worktree shows the branch to be merged and the release plan, and confirming it merges, moves to done, archives, commits, and releases, the same as `flai accept` from a host shell
- [x] Every git command flai runs for a story worktree from inside the container succeeds: the release plan, `flai stream sync`, the rebase and merge, and removing the worktree; git in the container no longer lists host worktrees as prunable; worktrees created before this change work without being recreated
- [x] Relative-path worktrees are an explicit opt-in setting, off by default. With it off, flai never passes `--relative-paths` and never causes `extensions.relativeWorktrees` to be set. flai never turns it on by detecting a git version
- [x] With the setting on and git 2.48 or newer, `flai stream open` creates the worktree with `--relative-paths`, and the worktree works on the host, in the container, and after the repository folder is moved. With the setting on and an older git, `flai stream open` warns, naming the git version and the setting, and creates an ordinary worktree, which the mount covers
- [x] `flai check` warns when the clone has `extensions.relativeWorktrees` set and the git on `PATH` is older than 2.48, read from the config file because that git cannot open the repository, and the message says how to recover
- [x] When the host path cannot be used as a container path, `flai dashboard` says so at start, names the opt-in setting as the way to make acceptance work from the board, and still serves the dashboard. When git cannot be run for a worktree anyway, the card shows what to do (`flai accept S-nnnn` from a shell) in place of the raw git error
- [x] Tests that fail without the fix: the `docker run` arguments carry the host path for the volume and `PROJECT_DIR`; a story with a worktree is accepted through the path the dashboard uses; the opt-in cases run where git is 2.48 or newer and skip with a stated reason where it is not, and CI runs them
- [x] An ADR refining ADR-0019 records the mount contract and the opt-in; `design/system/flaiover-dashboard.md`, `design/system/flai-cli.md`, `design/tech`, `docs/users`, and the operator settings index document the mount, the setting, what turning it on does to the clone, and how to turn it back off

## Tasks
- T-0151 ADR-0022 and living design: repository mounted at its host path, relative-path worktrees as an opt-in
- T-0152 flai dashboard mounts the repository at its host path and sets PROJECT_DIR to it
- T-0153 flaiover image and server hold no assumption that the project is at /project
- T-0154 Opt-in setting for relative-path worktrees in flai stream open
- T-0155 flai check warns when the clone needs a newer git than the one on PATH
- T-0156 Acceptance that cannot run git for a worktree says what to do, not the raw git error
- T-0157 User, operator, and tech documentation for the mount and the setting
- T-0158 Verify: all tiers, and acceptance of a story with a worktree through a real container

## Notes
Found by the operator on 2026-09-18 accepting S-0049 from the board. The card showed `git log --format=%H --fixed-strings --grep=[S-0049]: exit status 128 fatal: not a git repository: (null)`. Recorded as I-0017. S-0049 was accepted from a host shell.

Cause: the accept dry run plans the release from the story's worktree (`flai/cmd/accept.go`, `planRoot`). The worktree's `.git` file holds the absolute host path of the main repository's `.git/worktrees/S-0049`, and `flai dashboard` mounts the repository at `/project` (`flai/cmd/dashboard.go`, `--volume <root>:/project`, `PROJECT_DIR=/project`), so that path does not exist in the container. Reproduced with `docker exec` running `git log -1` in `/project/.flai-cache/worktrees/S-0049`. The real merge would fail the same way after the dry run. `git worktree list` in the container marked the worktree prunable; flai never runs `git worktree prune`.

Decided by the operator on 2026-09-18: the mount at the host path is unconditional, and relative-path worktrees are part of this story as an explicit opt-in. The operator first proposed choosing between the two by detecting the git version; that was set aside for the reasons below.

Observed on 2026-09-18. In the image's git 2.54.0, `git worktree add --relative-paths` wrote `gitdir: ../r/.git/worktrees/wt` in the worktree and a relative back-pointer, the worktree survived moving the parent folder, and creating it changed the clone from `core.repositoryformatversion=0` with no extension to `1` with `extensions.relativeWorktrees=true`. The host's git 2.47.3 refuses a repository with that extension: `fatal: unknown repository extension found: relativeworktrees`. The refusal covers the whole clone, not the one worktree.

Why not switch on the detected version. The version detected is the git that runs `flai stream open`; the extension binds every git that touches the clone afterwards, including tools that bundle or reimplement git (support in libgit2, JGit, and go-git was not checked). A host upgrade from 2.47 to 2.48 would change the repository format silently at the next story, and a later downgrade or an older tool would lose the repository. Worktrees made before an upgrade keep absolute paths, so the mount is needed at every version anyway. Hence the mount always, and the format change only when someone chooses it.

What the opt-in is for: a host whose path cannot be mirrored into a Linux container, a native Windows path being the case in view, and a clone that gets moved. The operator's host is WSL2 with Linux paths and does not need it.

To settle when refining. Where the setting lives: the repository format is local to a clone, so a per-user flai config key fits better than the committed manifest, which would impose it on every contributor whatever their git; name to be chosen. Whether `/project` stays as a second mount for anything that still assumes it; S-0043 (hub) should know the mount path is no longer a constant. Whether turning the setting on offers `git worktree repair --relative-paths` for existing worktrees; not required, since the mount covers them. The warn-and-fall-back behaviour on an older git is the agent's choice, on the workflow principle that a finished story is worth more than a hard stop; the operator may prefer a refusal. The exact recovery steps after opting in by mistake need to be tried before they are documented.

S-0046 shipped acceptance from the board and S-0041 extends it; neither was exercised against a story with a worktree inside the container.

Not in the board `order`; the operator has not placed it against S-0040 to S-0043. Until it lands, accept stories from a host shell.

Verification, 2026-09-18. A local image was built from the story branch and a second dashboard started with the branch's `flai dashboard` against a scratch git project, on its own name and port; the operator's dashboard was not touched. The repository was mounted at its host path with `PROJECT_DIR` equal to it; `git worktree list` in the container showed the story worktree, made by the host's git 2.47.3 with absolute links, and nothing prunable. `GET /api/items/S-0001/acceptance` returned the branch and a release plan built from the worktree's commit, with no blockers, which is the call that failed for S-0049. `POST /api/items/S-0001/move` to done merged the branch, moved the story to done, archived it, committed, removed the worktree, and deleted the branch. The scratch project has no components and no remote, so the release step ran and released nothing and nothing was pushed; tagging and pushing were not exercised in the container and are unchanged by this story. The same image started by hand with the old `/project` mount returned the new blocker in the preview and refused the move with the same advice, leaving the story in review. The opt-in's real-git test skips on the host and passed inside the image's git 2.54.0; CI's runner reported git 2.55.0 in its latest run, so it runs there. The board itself was not driven in a browser; the API calls are the ones its card and confirmation make.
