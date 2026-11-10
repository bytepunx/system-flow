---
id: T-3005
type: task
nature: improvement
title: flai stream sync writes the generated files again at the branch's head after a rebase that wrote one at a stop, and commits what changed
status: backlog
parent: S-0702
owner: alex
created: 2026-11-10T10:31:20Z
updated: 2026-11-10T10:31:20Z
transitions: []
stream: S-0702
tags: [flai, git, issues]
touches: [flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/issues/generated.go, flai/internal/issues/generated_test.go]
after: [T-3002]
---
# T-3005 flai stream sync writes the generated files again at the branch's head after a rebase that wrote one at a stop, and commits what changed

## Work

Build the ADR's decision in `Sync` (`flai/internal/storygit/sync.go`). After the rebase is clean, and when `res.Regenerated` names a file, write each of `o.Files.Generated` again in the worktree. If one now differs from the head's, stage it and commit it on the branch with a subject such as `docs: [S-nnnn] write <paths> again after the rebase`, as `fold` commits its files. Record what was committed in the sync's result.

The summary's `updated` line changes on every write, so a write that changes nothing else is not a change. Give the generated file what it needs to say so, the same comparison `check_records.go` makes with `summaryUpdated`, in `flai/internal/issues/generated.go`. Restore the head's file when nothing else changed.

Reproduce I-0266 first in `sync_test.go`, in the style of `TestSyncRegeneratesGeneratedFilesItStopsOnAlone`:

- commit one changes the summary;
- commit two adds an issue file and leaves the summary alone;
- main changes the summary meanwhile;
- the sync must leave a head whose summary renders both issue files.

The test fails before the change. Add a test that a sync whose rebase wrote nothing again commits nothing.

It waits for T-3002, whose ADR it builds.

## Done when

- The reproducing test fails without the change and passes with it.
- A sync that wrote nothing again at a stop adds no commit.
- `flai test` passes on the changed files.

## Notes
