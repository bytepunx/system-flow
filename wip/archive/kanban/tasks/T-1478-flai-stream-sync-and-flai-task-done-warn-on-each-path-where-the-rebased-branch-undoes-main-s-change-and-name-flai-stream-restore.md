---
id: T-1478
type: task
nature: remediation
title: flai stream sync and flai task done warn on each path where the rebased branch undoes main's change, and name flai stream restore
status: done
parent: S-0371
owner: alex
created: 2026-10-12T17:15:49Z
updated: 2026-10-14T22:30:27Z
transitions:
  - to: ready
    at: 2026-10-14T22:14:48Z
    by: agent-S-0371
  - to: in-progress
    at: 2026-10-14T22:14:49Z
    by: agent-S-0371
  - to: done
    at: 2026-10-14T22:30:27Z
    by: agent-S-0371
stream: S-0371
tags: [cli, go]
touches: [flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go, flai/internal/storygit/undone.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go]
after: [T-1574]
usage:
  source: log
  seconds: 938
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 42
      output: 15800
      cache_read: 2492767
      cache_write: 116840
      cost: 1.4562
---
# T-1478 flai stream sync and flai task done warn on each path where the rebased branch undoes main's change, and name flai stream restore

## Work

It waits for T-1574: both change `flai/cmd/stream.go`, and the warning names the command T-1574 builds. T-1574 waits for T-1477, which adds `SyncResult.UndoesMain`.

- In `flai/cmd/stream_sync.go`, `printSyncChecks` prints, for each entry of `res.UndoesMain`, a warning such as `story/S-nnnn removes <n> lines main gained in <path> since <old base>; restore them with flai stream restore S-nnnn --path <path>, or keep the removal at acceptance with --keep-removal and a reason`, before the outside-touches note. It prints it whether or not the path is in the story's claim. `flai task done` prints through the same function, so it needs no change of its own.
- In `flai/cmd/stream.go`, the help of `flai stream sync` says what the warning means and that `--json` gives it as `undoes_main`.
- In `flai/cmd/stream_sync_test.go`, a test with the existing helpers (`syncProject`, `openSyncStory`, `commitIn`) that runs I-0124's sequence through `flai stream sync` and asserts the warning's wording and the `undoes_main` JSON, for a path inside the claim.

## Done when

- `flai stream sync` on I-0124's sequence prints the warning naming `flai stream restore` and exits 0, and `--json` lists the path under `undoes_main`.
- `flai test` on the three files passes.

## Notes

`docs/users/flai-reference.md`, generated from the help, is regenerated in the docs task, once, after every help changed.
