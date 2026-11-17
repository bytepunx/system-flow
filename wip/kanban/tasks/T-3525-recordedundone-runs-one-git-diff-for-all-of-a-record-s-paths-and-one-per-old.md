---
id: T-3525
type: task
nature: improvement
title: RecordedUndone runs one git diff for all of a record's paths and one per old base, not one or two per path
status: in-progress
parent: S-0828
owner: alex
created: 2026-11-17T08:28:54Z
updated: 2026-11-17T11:50:13Z
transitions:
  - to: ready
    at: 2026-11-17T11:50:12Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T11:50:13Z
    by: system-flow
stream: S-0828
tags: [flai, git, performance]
touches: [flai/internal/storygit/undone.go, flai/internal/storygit/undone_test.go]
---
# T-3525 RecordedUndone runs one git diff for all of a record's paths and one per old base, not one or two per path

## Work

In `flai/internal/storygit/undone.go`, change `RecordedUndone` so that it no longer runs `readHunks` once or twice per recorded path:

- Run one `git diff main story/S-nnnn -- <every recorded path>` and split its output by file, so each path gets the lines the branch adds and removes as `readHunks` gives them today.
- For the paths whose branch still removes lines, group them by `old_base` and run one `git diff <old_base> main -- <those paths>` per group, split the same way.
- Read the path of each file from its own header, with `core.quotePath=false`, `--no-renames`, and explicit `--src-prefix`/`--dst-prefix` so a host's `diff.noprefix` cannot change it; a deleted or added file has `/dev/null` on one side; mind git's trailing tab after a path with a space. A path git prints no diff for has no lines.
- Keep everything else as it is: the record rewritten without the paths that drop out, the error texts naming the record to delete, `undoneLines`, `atPositions`, and `notKept`.

Add a test in `flai/internal/storygit/undone_test.go` that reproduces I-0349: a record of several paths with one old base, read through a runner that counts the git processes, runs at most two `git diff` and the one `git rev-parse` of the main branch, and gives the same undone paths and lines as before. Cover a record with two old bases, a path the branch no longer changes, and a path with a space.

Waits for no task: it is the fix itself.

## Done when

- `RecordedUndone` runs at most one `git diff` for the branch and one per distinct `old_base`, whatever the record's number of paths, and the new test shows it.
- The existing tests of the undone record pass: `flai/internal/storygit`, `flai/internal/mergequeue`, `flai/internal/hostapi`, and `flai/cmd`'s `accept_undone_test.go` and `stream_restore_test.go`, with `flai test` on the changed paths.

## Notes
