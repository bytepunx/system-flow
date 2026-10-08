---
id: T-1321
type: task
nature: remediation
title: A sync, a task's close, and an acceptance fold the branch's duplicate issue after the rebase, with a test that reproduces I-0112
status: backlog
parent: S-0326
owner: alex
created: 2026-10-08T00:33:54Z
updated: 2026-10-08T00:33:54Z
transitions: []
stream: S-0326
tags: [cli, go, git]
touches: [flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/stream.go, flai/cmd/branch.go, flai/internal/taskdone/taskdone.go, flai/cmd/stream_sync_test.go]
after: [T-1320]
---
# T-1321 A sync, a task's close, and an acceptance fold the branch's duplicate issue after the rebase, with a test that reproduces I-0112

## Work

Run T-1320's fold once a rebase of a story branch onto the main branch completes, in each place that rebases one: `flai stream sync` (`flai/cmd/stream.go`), the sync `flai task done` runs (`flai/internal/taskdone/taskdone.go`), and `flai accept`'s rebase through `syncStoryBranch` (`flai/cmd/branch.go`). Each passes `issues.Generated` to `storygit.SyncOptions` today; give the options a hook the issues package fills the same way, so that `storygit` stays free of the issues package, and have `storygit` run it in the worktree after the rebase, regenerate the summary, and commit the fold on the branch.

Reproduce I-0112 in `flai/cmd/stream_sync_test.go` with the helpers its summary tests use: two story branches each record a `narrative.state` finding outside the story under one title, the first is accepted, and the second syncs. Without the fold the second branch holds two open issues of that title; with it, one, carrying both instances.

It waits for T-1320 because it calls the fold.

## Done when

- The reproduction test fails without the fold and passes with it
- `flai stream sync --json`, and the report without it, name each issue folded and the issue it went into
- `flai test` on the paths this task changed passes

## Notes
