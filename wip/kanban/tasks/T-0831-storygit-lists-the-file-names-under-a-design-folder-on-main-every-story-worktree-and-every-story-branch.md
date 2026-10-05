---
id: T-0831
type: task
nature: remediation
title: storygit lists the file names under a design folder on main, every story worktree, and every story branch
status: done
parent: S-0252
owner: alex
created: 2026-10-04T23:27:29Z
updated: 2026-10-04T23:59:02Z
transitions:
  - to: ready
    at: 2026-10-04T23:57:04Z
    by: agent-S-0252
  - to: in-progress
    at: 2026-10-04T23:57:04Z
    by: agent-S-0252
  - to: done
    at: 2026-10-04T23:59:02Z
    by: agent-S-0252
stream: S-0252
tags: [flai, git]
touches: [flai/internal/storygit]
usage:
  source: log
  seconds: 116
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 32
      output: 145
      cache_read: 748409
      cache_write: 54556
      cost: 0.3257
---
# T-0831 storygit lists the file names under a design folder on main, every story worktree, and every story branch

## Work

Add a helper to `flai/internal/storygit` that takes an `execx.Runner`, the repository, and a repository-relative folder such as `design/issues`. It returns the base names of the files under that folder, gathered from four places:

- the main checkout's working tree;
- the working tree of every linked worktree, so that an issue written but not yet committed in another story's worktree counts;
- every local `story/*` branch, read with `git ls-tree --name-only`, as `preview.resultsBlocker` already does for one branch;
- the main branch.

A branch or worktree that cannot be read is skipped, not fatal, so numbering never fails for want of one. `git for-each-ref refs/heads/story/` lists the branches, and `git worktree list --porcelain` lists the worktrees. S-0245, the same defect for ADR numbers (I-0063), can reuse this helper for `design/adrs`, so it names a folder rather than issues.

This task waits for nothing: it is the first layer.

## Done when

- A test in `flai/internal/storygit` builds a repository with a main branch, two story branches each holding a different file under the folder, and a worktree holding an uncommitted one. The helper returns all of these names.
- A missing folder, or a repository with no story branches, gives the main checkout's names and no error.
- `scripts/flai-test.sh` passes.

## Notes
