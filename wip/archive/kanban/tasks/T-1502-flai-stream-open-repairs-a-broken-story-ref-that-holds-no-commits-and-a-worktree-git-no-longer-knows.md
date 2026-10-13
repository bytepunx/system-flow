---
id: T-1502
type: task
nature: remediation
title: flai stream open repairs a broken story ref that holds no commits and a worktree git no longer knows
status: done
parent: S-0375
owner: alex
created: 2026-10-13T16:36:46Z
updated: 2026-10-13T16:58:21Z
transitions:
  - to: ready
    at: 2026-10-13T16:39:14Z
    by: agent-S-0375
  - to: in-progress
    at: 2026-10-13T16:39:14Z
    by: agent-S-0375
  - to: done
    at: 2026-10-13T16:58:21Z
    by: agent-S-0375
stream: S-0375
tags: [flai, git]
touches: [flai/internal/storygit/storygit.go, flai/internal/storygit/open.go, flai/internal/storygit/open_test.go, flai/cmd/stream.go, flai/cmd/stream_reopen_test.go]
usage:
  source: log
  seconds: 1147
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 89
      output: 31127
      cache_read: 4095987
      cache_write: 168725
      cost: 2.5021
---
# T-1502 flai stream open repairs a broken story ref that holds no commits and a worktree git no longer knows

## Work

After I-0137's crash, `.git/refs/heads/story/S-0337` and its reflog were empty files. `BranchExists` (`flai/internal/storygit/storygit.go`) reads a broken ref as no branch, so `Open` (`flai/internal/storygit/open.go`) ran `git worktree add -b`, which git refuses ("cannot lock ref ... reference broken"). The agent may not write under `.git`, so the story sat in progress for four days until the operator deleted the ref by hand (TH-0390). The worktree was zero bytes too, and `Open` takes any existing folder at the worktree path as a worktree.

- Add a helper in `storygit.go` that tells a broken story ref apart from a missing one: the loose ref file exists under `git rev-parse --git-common-dir`, but `git rev-parse --verify --quiet refs/heads/<branch>` resolves nothing. A ref that resolves to nothing holds no commits, so removing it loses no work.
- In `Open`, before choosing between `worktree add <branch>` and `worktree add -b`, remove a broken story ref and its reflog file under `logs/refs/heads/`, then log at warn what was removed. Only for the story's own `story/S-nnnn` ref, never another.
- In `Open`, take an existing folder as the worktree only when git knows it: its `.git` file names a gitdir that `git worktree list --porcelain` lists. A folder git does not know, or whose `.git` file is empty, is a crashed start: remove it only when the story's branch is missing or was broken, so it holds no commits, then run `git worktree prune` and add the worktree again. Otherwise refuse, naming the folder and the command that removes it, rather than delete what may be work.
- Pass `-c core.fsync=reference` to the `git worktree add` that creates the branch, so a crash after it returns leaves the ref whole. Git 2.36 added `core.fsync`; the host runs 2.47.
- In `flai/cmd/stream.go`, make `canReopen` treat a worktree folder git does not know as absent, using the same helper, so `flai stream open` takes up a narrative with content over a crashed worktree.
- Test with real git, beside `TestOpenReusesAnExistingWorktree` in `open_test.go`: a zero-byte story ref with a zero-byte reflog, and a worktree folder whose files and `.git` are zero bytes. `Open` creates the branch from main and a working worktree. A further case: a worktree folder git does not know while the branch holds a commit is refused, and nothing is removed.

Waits for nothing: no other task touches these files.

## Done when

- `Open` replaces a broken story ref and a crashed worktree that hold no commits, and refuses one whose branch holds commits.
- The new tests reproduce I-0137's empty ref and worktree, and pass.
- `flai test flai/internal/storygit flai/cmd` passes.

## Notes

Drafted by the planner from I-0137, TH-0390, and `open.go` lines 50-105. The safety line is that flai removes only what resolves to no commit.
