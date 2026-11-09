---
id: T-2768
type: task
nature: remediation
title: Close I-0289 naming the guard's read forms of git tag, worktree, and branch
status: in-progress
parent: S-0623
owner: alex
created: 2026-11-09T21:56:08Z
updated: 2026-11-09T22:21:22Z
transitions:
  - to: ready
    at: 2026-11-09T22:21:21Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T22:21:22Z
    by: system-flow
stream: S-0623
tags: [issues]
touches: [design/issues/I-0289-flai-guard-refuses-the-orchestrator-git-tag-list-and-git-worktree-list-which.md, design/issues/summary.md]
after: [T-2766]
---
# T-2768 Close I-0289 naming the guard's read forms of git tag, worktree, and branch

## Work

Run `flai issue close I-0289 --reason "<reason>"`, the reason saying that the guard now passes `git tag` that lists, `git worktree list`, and `git branch` that lists for every role it judges git for, and naming the test T-2766 added. It rewrites the issue and `design/issues/summary.md`. Waits for T-2766, because the issue closes on the fix being in place.

## Done when

- I-0289 is closed with that reason, and `design/issues/summary.md` no longer lists it open.
- `flai check --strict` passes.

## Notes
