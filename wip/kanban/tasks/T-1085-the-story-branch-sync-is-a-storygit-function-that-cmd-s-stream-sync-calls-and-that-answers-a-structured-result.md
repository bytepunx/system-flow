---
id: T-1085
type: task
nature: improvement
title: The story-branch sync is a storygit function that cmd's stream sync calls and that answers a structured result
status: backlog
parent: S-0269
owner: alex
created: 2026-10-06T22:52:50Z
updated: 2026-10-06T22:52:50Z
transitions: []
stream: S-0269
tags: [cli, git]
touches: [flai/cmd/stream.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go]
---
# T-1085 The story-branch sync is a storygit function that cmd's stream sync calls and that answers a structured result

## Work

At present `flai stream sync` lives in methods of cmd's `app`: its RunE is `newStreamSyncCmd` in `flai/cmd/stream.go`, and its helpers are `checkSync`, `outsideClaim`, `trialMerge`, and `reportConflicts` in `flai/cmd/stream_sync.go`. The MCP server cannot import cmd, so it cannot reach them. Move the sync into `flai/internal/storygit`.

- Add a function in `flai/internal/storygit/sync.go` that syncs a story's branch and answers a result struct. The struct gives:
  - synced, or refused: uncommitted changes, or a rebase unfinished;
  - the conflicting paths, and how to continue;
  - what lies outside the story's claim;
  - the trial-merge conflicts with other stories.
- It takes the runner, the repository, the story, and the clock as arguments. It does not read them from cmd's `app`.
- Leave `flai stream sync`'s text output, exit codes, and conflict threads as they are, built from that result. Its existing tests in `flai/cmd/stream_sync_test.go` keep passing unchanged, save where they reach the moved helpers.
- Test the function in `flai/internal/storygit/sync_test.go` with the `gittest` fixtures: a clean sync, one refused for uncommitted changes, one refused for a rebase in progress, and a conflict.

This task waits for nothing and touches no file the other layer 1 tasks touch, so it runs with them.

## Done when

- [ ] A storygit function syncs a story branch and answers a structured result, with tests for clean, refused, and conflicting syncs.
- [ ] `flai stream sync` calls it, and its behaviour and tests are unchanged.
- [ ] `scripts/flai-test.sh` passes for `flai/cmd` and `flai/internal/storygit`.

## Notes

Drafted by the planner from the explorer's reading of `flai/cmd/stream.go` and `flai/cmd/stream_sync.go`.
