---
id: T-1478
type: task
nature: remediation
title: flai stream sync and flai task done warn on each path where the rebased branch undoes main's change, and say how to restore it
status: backlog
parent: S-0371
owner: alex
created: 2026-10-12T17:15:49Z
updated: 2026-10-12T17:15:49Z
transitions: []
stream: S-0371
tags: [cli, go]
touches: [flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go]
after: [T-1477]
---
# T-1478 flai stream sync and flai task done warn on each path where the rebased branch undoes main's change, and say how to restore it

## Work

It waits for T-1477, which adds `SyncResult.UndoesMain`. It shares no path with T-1479, so the two run together.

- In `flai/cmd/stream_sync.go`, `printSyncChecks` prints, for each entry of `res.UndoesMain`, a warning such as `story/S-nnnn removes <n> lines main gained in <path> since <old base>; restore main's file in the worktree with git checkout <main> -- <path> and commit, unless the story means to remove them`, before the outside-touches note. It prints it whether or not the path is in the story's claim. `flai task done` prints through the same function, so it needs no change of its own.
- In `flai/cmd/stream.go`, the help of `flai stream sync` says what the warning means and that `--json` gives it as `undoes_main`.
- In `flai/cmd/stream_sync_test.go`, a test with the existing helpers (`syncProject`, `openSyncStory`, `commitIn`) that runs I-0124's sequence through `flai stream sync` and asserts the warning's wording and the `undoes_main` JSON, for a path inside the claim.

## Done when

- `flai stream sync` on I-0124's sequence prints the warning and exits 0, and `--json` lists the path under `undoes_main`.
- `flai test` on the three files passes.

## Notes

`docs/users/flai-reference.md`, generated from the help, is regenerated in T-1480, after this task and T-1479, so the two do not both write it.
