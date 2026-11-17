---
id: T-3535
type: task
nature: remediation
title: The publish commit of flai release --pending waits out a held index lock, with a test that holds it
status: in-progress
parent: S-0833
owner: alex
created: 2026-11-17T08:29:55Z
updated: 2026-11-17T09:16:08Z
transitions:
  - to: ready
    at: 2026-11-17T09:16:08Z
    by: agent-S-0833
  - to: in-progress
    at: 2026-11-17T09:16:08Z
    by: agent-S-0833
stream: S-0833
tags: [flai, release, git]
touches: [flai/cmd/release.go, flai/cmd/release_lock_test.go]
---
# T-3535 The publish commit of flai release --pending waits out a held index lock, with a test that holds it

## Work

`computeApplyAndTagPending` in `flai/cmd/release.go` runs `git add -A` and `git commit -q -m "chore: publish …"` in the main checkout directly through `a.runner.Run`. When another git process holds `.git/index.lock`, as the orchestrator, the curator, or `flai serve`'s replan may, either call fails with `index.lock': File exists`, and `flai release --pending` exits 1 having tagged and pushed nothing (I-0359). I-0100 fixed the same race in `flai accept` with `storygit.RunPastIndexLock` (`flai/internal/storygit/indexlock.go`), which retries for about nine seconds and never removes the lock.

- Run those two calls through `storygit.RunPastIndexLock`, as `flai/cmd/accept.go` does. Leave the other git calls of the path as they are: `git status`, `git tag`, and `git push` take no index lock.
- Add `flai/cmd/release_lock_test.go`. With a runner like `lockingRunner` in `flai/cmd/accept_lock_test.go`, which creates the lock file on the first `git add -A` and the first `git commit -q` and removes it on the second, run a pending release end to end, as `TestReleasePendingEndToEnd` in `flai/cmd/release_pending_test.go` does, and show that it commits, tags, and pushes. Show too that a lock held past the waits fails with git's message and the helper's hint, and leaves no release tag.
- Run the test before the change, to see that it reproduces the failure, and after.

This task waits for none.

## Done when

- `flai/cmd/release.go` runs the publish's `git add -A` and `git commit` through `storygit.RunPastIndexLock`.
- `flai/cmd/release_lock_test.go` fails without the change and passes with it.
- `flai test flai/cmd/release.go flai/cmd/release_lock_test.go` passes.

## Notes
