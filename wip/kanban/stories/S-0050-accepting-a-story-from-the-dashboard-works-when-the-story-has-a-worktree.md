---
id: S-0050
type: story
nature: remediation
title: Accepting a story from the dashboard works when the story has a worktree
status: backlog
parent: E-0006
owner: alex
created: 2026-09-18T18:36:10Z
updated: 2026-09-18T18:36:10Z
transitions: []
tags: [cli, dashboard]
touches: [flai/cmd, flaiover]
---

# S-0050 Accepting a story from the dashboard works when the story has a worktree

## Goal
Dropping a story on done, or pressing accept on its page, works from the dashboard when the story was worked on a branch in a worktree, which is every story since ADR-0019. Today it fails before the confirmation is shown.

## Acceptance criteria
- [ ] With flaiover running in its container through `flai dashboard`, the acceptance preview for a story in review that has a `story/S-nnnn` worktree shows the branch to be merged and the release plan, and confirming it merges, moves to done, archives, commits, and releases, the same as `flai accept` from a host shell
- [ ] Every git command flai runs for a story worktree from inside the container succeeds: the release plan, `flai stream sync`, the rebase and merge, and removing the worktree; and git in the container no longer lists host worktrees as prunable
- [ ] The host is unaffected: worktrees created before and after the change work from a host shell, with the host's git version, and `flai stream open` needs no newer git than it does today unless `design/tech` and the docs say so
- [ ] When git cannot be run for the worktree anyway, the dashboard shows what to do (accept from a shell with `flai accept S-nnnn`) in place of the raw git error
- [ ] A test that fails without the fix: a story with a worktree, accepted through the same path the dashboard uses, with the repository at a different path than the one the worktree was created under
- [ ] `design/system/flaiover-dashboard.md`, `design/system/flai-cli.md`, and `docs/users` say how the repository is mounted; ADR-0019 gets a refinement ADR if the worktree layout or the mount contract changes

## Tasks

## Notes
Found by the operator on 2026-09-18 accepting S-0049 from the board. The card showed `git log --format=%H --fixed-strings --grep=[S-0049]: exit status 128 fatal: not a git repository: (null)`. Recorded as I-0017.

Cause: the accept dry run plans the release from the story's worktree (`flai/cmd/accept.go`, `planRoot`). The worktree's `.git` file holds the absolute host path of the main repository's `.git/worktrees/S-0049`, and `flai dashboard` mounts the repository at `/project` (`flai/cmd/dashboard.go`, `--volume <root>:/project`, `PROJECT_DIR=/project`), so that path does not exist in the container. Reproduced with `docker exec flaiover-system-flow sh -c 'cd /project/.flai-cache/worktrees/S-0049 && git log -1'`. The real merge would fail the same way after the dry run. `git worktree list` in the container marks the worktree prunable; flai never runs `git worktree prune`, and nothing else in the image should.

Two candidate fixes, to choose when refining. Mount the repository at its host path, in addition to or in place of `/project`, and set `PROJECT_DIR` to it: one line in `dashboard.go`, works with any git, and the mount path stops being a constant, which the hub work in S-0043 should know about. Or create worktrees with relative paths (`git worktree add --relative-paths`): needs git 2.48 or newer on every machine that touches the repository, because it sets a repository extension older git refuses; the host here has 2.47.3 and the image 2.54.0.

S-0046 shipped acceptance from the board and S-0041 extends it; neither was exercised against a story with a worktree inside the container.

Not in the board `order`; the operator has not placed it against S-0040 to S-0043. Until it lands, accept stories from a host shell.
