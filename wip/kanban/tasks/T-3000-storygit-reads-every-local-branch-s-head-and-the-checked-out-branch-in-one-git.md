---
id: T-3000
type: task
nature: improvement
title: storygit reads every local branch's head and the checked-out branch in one git for-each-ref
status: done
parent: S-0698
owner: alex
created: 2026-11-10T10:31:03Z
updated: 2026-11-11T08:31:25Z
transitions:
  - to: ready
    at: 2026-11-11T08:31:17Z
    by: agent-S-0698
  - to: in-progress
    at: 2026-11-11T08:31:20Z
    by: agent-S-0698
  - to: done
    at: 2026-11-11T08:31:25Z
    by: agent-S-0698
stream: S-0698
tags: [cli]
touches: [flai/internal/storygit/storygit.go, flai/internal/storygit/storygit_test.go, flai/internal/storygit/trial.go]
usage:
  source: log
  seconds: 5
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 48
      cache_read: 707536
      cache_write: 535
      cost: 0.2618
---
# T-3000 storygit reads every local branch's head and the checked-out branch in one git for-each-ref

## Work

Add a heads read to `flai/internal/storygit/storygit.go`, beside `MainBranch` and `BranchExists`. It runs one `git for-each-ref --format='%(HEAD) %(refname) %(objectname)' refs/heads` in the main checkout. It answers the checked-out branch, its head commit, and each local branch's head commit, by branch name.

- A detached HEAD answers the error `MainBranch` gives for it, so callers keep their wording.
- A branch whose loose ref resolves to no commit (I-0137) is left out of the answer, as `git rev-parse --verify` names no commit for it. Its caller reports it as unread, as `verified` does now.

This task waits for no other. It is the first layer, and `mergequeue` reads through it in the next task.

## Done when

- The heads read runs exactly one git process, whatever the number of branches.
- A real-git test in `flai/internal/storygit/storygit_test.go` checks it. The test makes a repository with a main branch and two `story/S-nnnn` branches at different commits, and a detached HEAD. It skips under `-short`, as the package's other real-git tests do.
- `flai test flai/internal/storygit` passes.

## Notes
