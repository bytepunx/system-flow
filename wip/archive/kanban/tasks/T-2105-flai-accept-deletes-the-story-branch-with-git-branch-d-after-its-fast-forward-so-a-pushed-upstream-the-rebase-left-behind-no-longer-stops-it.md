---
id: T-2105
type: task
nature: remediation
title: flai accept deletes the story branch with git branch -D after its fast-forward, so a pushed upstream the rebase left behind no longer stops it
status: done
parent: S-0493
owner: alex
created: 2026-10-24T13:17:47Z
updated: 2026-10-24T13:26:37Z
transitions:
  - to: ready
    at: 2026-10-24T13:17:54Z
    by: agent-S-0493
  - to: in-progress
    at: 2026-10-24T13:17:55Z
    by: agent-S-0493
  - to: done
    at: 2026-10-24T13:26:37Z
    by: agent-S-0493
stream: S-0493
tags: []
touches: [flai/cmd/branch.go, flai/cmd/accept_test.go, design/system/flai-cli.md, docs/users/flai.md]
usage:
  source: log
  seconds: 522
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 19
      output: 4004
      cache_read: 1317959
      cache_write: 43075
      cost: 0.6883
---
# T-2105 flai accept deletes the story branch with git branch -D after its fast-forward, so a pushed upstream the rebase left behind no longer stops it

## Work

`mergeStoryBranch` in `flai/cmd/branch.go` deletes the story branch with `git branch -d` after `git merge --ff-only` succeeded. A branch flai serve pushed for a pull request tracks `origin/story/S-nnnn`; the acceptance rebase rewrites its commits, so `-d` checks the upstream and refuses, after main already has the work and the worktree is gone. Delete with `-D` once the fast-forward succeeded, as `removeStoryCheckout` does, with an error that names what happened and how to finish. An acceptance run again from that state fast-forwards nothing and goes on to delete the branch, so it finishes. Add real-git tests in `flai/cmd/accept_test.go`: one pushes the story branch to a bare origin with an upstream, moves main, and accepts; one accepts again from the state I-0204 left. Update `design/system/flai-cli.md` and `docs/users/flai.md` where they describe the branch's removal.

Waits for nothing: the story's only task.

## Done when

- Accepting a story whose branch tracks a pushed upstream that the rebase left behind merges, removes the branch, and exits 0.
- Accepting again a story left merged with its branch still there finishes the acceptance.
- Both tests fail with `-d`.

## Notes
