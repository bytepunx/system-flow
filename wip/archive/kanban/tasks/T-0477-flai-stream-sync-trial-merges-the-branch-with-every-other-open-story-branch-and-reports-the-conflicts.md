---
id: T-0477
type: task
nature: feature
title: flai stream sync trial-merges the branch with every other open story branch and reports the conflicts
status: done
parent: S-0131
owner: alex
created: 2026-09-26T18:05:55Z
updated: 2026-09-26T18:08:04Z
transitions:
  - to: ready
    at: 2026-09-26T18:06:13Z
    by: agent-S-0131
  - to: in-progress
    at: 2026-09-26T18:06:14Z
    by: agent-S-0131
  - to: done
    at: 2026-09-26T18:08:04Z
    by: agent-S-0131
stream: S-0131
tags: []
touches: [flai/cmd/branch.go, flai/cmd/stream.go, flai/cmd/stream_sync.go, flai/internal/gitver]
---
# T-0477 flai stream sync trial-merges the branch with every other open story branch and reports the conflicts

## Work

- After a clean rebase, run `git merge-tree --write-tree --name-only --no-messages` of `story/<id>` against the branch of every other story in progress or in review that has one, in the main checkout, writing nothing to any worktree.
- Print each branch it conflicts with and the conflicting paths, and each clean one; `--json` carries the same data in a `branches` list.
- A git older than 2.38 skips the trial merge with a warning that names the version it needs (`gitver.MergeTree`).
- An integration test with real git (skipped under `-short`) covers a conflicting pair and a clean pair.

## Done when

- `flai stream sync` prints the conflicting branches and paths and says when a branch merges cleanly, and the tests for both pass.

## Notes
