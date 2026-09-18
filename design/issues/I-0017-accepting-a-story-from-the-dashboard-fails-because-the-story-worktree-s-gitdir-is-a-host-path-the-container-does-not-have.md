---
id: I-0017
title: Accepting a story from the dashboard fails because the story worktree's gitdir is a host path the container does not have
class: defect
status: closed
count: 1
cost: 10m
first_reported: 2026-09-18T18:35:50Z
last_reported: 2026-09-18T18:35:50Z
updated: 2026-09-18T21:09:10Z
---

# I-0017 Accepting a story from the dashboard fails because the story worktree's gitdir is a host path the container does not have

## Description
Accepting a story from the dashboard fails because the story worktree's gitdir is a host path the container does not have

## Instances

### 2026-09-18T18:35:50Z
S-0049: the operator dropped the card on done. The accept dry run plans the release from .flai-cache/worktrees/S-0049, whose .git file points at /home/alex/git/bytepunx/system-flow/.git/worktrees/S-0049. The container mounts the repository at /project, so git in the worktree fails with 'fatal: not a git repository: (null)' and the card shows the raw git error. Reproduced with docker exec. Nothing was changed; the story stayed in review. git in the container lists the worktree as prunable. flai accept from a host shell plans correctly.

## Remediation
Queued as S-0050. A git worktree records absolute paths in both directions: the worktree's `.git` file names the main repository's `.git/worktrees/<id>`, and that folder's `gitdir` names the worktree. `flai dashboard` mounts the repository at `/project` (`flai/cmd/dashboard.go`), so neither path exists in the container. Candidate fixes: mount the repository at its host path as well as, or instead of, `/project` and point `PROJECT_DIR` at it; or create worktrees with relative paths (`git worktree add --relative-paths`, git 2.48 or newer; the host here has 2.47.3, the image 2.54.0, and the setting makes the repository unreadable to older git). Until fixed, accept stories that have a branch from a host shell with `flai accept`. Decided by the operator on 2026-09-18: S-0050 mounts the repository at its host path unconditionally and adds relative-path worktrees as an explicit opt-in, never chosen by detecting the git version, because creating one sets `extensions.relativeWorktrees` on the whole clone and git older than 2.48 then refuses it (observed with 2.47.3).
Closed 2026-09-18T21:09:10Z: fixed by S-0050 (flai 1.2.3): flai dashboard mounts the repository at its host path. Confirmed on 2026-09-18 when the operator accepted S-0048 from the board after the dashboard was restarted by the new flai
