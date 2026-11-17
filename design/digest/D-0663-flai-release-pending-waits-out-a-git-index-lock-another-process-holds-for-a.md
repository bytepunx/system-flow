---
id: D-0663
title: flai release --pending waits out a git index lock another process holds for a moment, as flai accept does, and fails with git's message only when the lock stays held about nine seconds
kind: cause
sources: [I-0359, S-0833]
topics: [release, flai, cli, go, code, git]
paths: [flai/cmd/release.go, flai/cmd/release_lock_test.go, flai/internal/storygit/indexlock.go, docs/users/flai.md, design/issues/I-0359-flai-release-pending-fails-when-another-process-holds-the-main-checkout-s-git.md]
date: 2026-11-17
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0663 flai release --pending waits out a git index lock another process holds for a moment, as flai accept does, and fails with git's message only when the lock stays held about nine seconds

The cause of [I-0359](../issues/I-0359-flai-release-pending-fails-when-another-process-holds-the-main-checkout-s-git.md): `flai release --pending` staged and committed the main checkout with a plain `git add -A` and `git commit`, so when the orchestrator, the curator, or another git process held `.git/index.lock` for a moment, the publish exited 1 with git's "File exists" and the publish formula's step stopped. It happened in S-0764 and S-0793 ([D-0603](D-0603-a-publish-step-that-exits-1-because-another-git-process-briefly-held-the-main.md)).

S-0833 fixed it the way S-0307 fixed the same race in `flai accept` for I-0100 ([D-0208](D-0208-the-acceptance-s-commit-step-fails-when-another-process-holds-git-s-index-lock.md)): the two git calls in `computeApplyAndTagPending` in `flai/cmd/release.go` run through `storygit.RunPastIndexLock`, which runs them again for about nine seconds while git reports the lock held, and never removes the lock. The issue closed with: "flai release --pending runs its publish's git add -A and git commit through storygit.RunPastIndexLock, as I-0100's fix does for flai accept, so a lock another git process holds for a moment is waited out (about nine seconds, never removed)."

`flai/cmd/release_lock_test.go` holds both cases: `TestReleasePendingWaitsOutAnIndexLockHeldForAMoment` failed with "index.lock': File exists" before the change and commits, tags, and pushes after; `TestReleasePendingFailsWhenTheIndexLockStaysHeld` shows that a lock held throughout still fails, with git's message and the helper's hint, and tags and pushes nothing. That second test waits the helper's full nine seconds, since its wait count is not settable from package `cmd`, and runs in parallel rather than adding a test hook to `storygit`.

So a later index-lock race in a flai command that commits the main checkout is fixed by routing its git calls through `storygit.RunPastIndexLock`, not by a retry of its own.
