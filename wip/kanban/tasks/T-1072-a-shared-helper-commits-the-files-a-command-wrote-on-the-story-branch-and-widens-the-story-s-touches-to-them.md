---
id: T-1072
type: task
nature: improvement
title: A shared helper commits the files a command wrote on the story branch and widens the story's touches to them
status: done
parent: S-0275
owner: alex
created: 2026-10-06T22:52:18Z
updated: 2026-10-07T08:20:19Z
transitions:
  - to: ready
    at: 2026-10-07T08:15:11Z
    by: agent-S-0275
  - to: in-progress
    at: 2026-10-07T08:15:11Z
    by: agent-S-0275
  - to: done
    at: 2026-10-07T08:20:19Z
    by: agent-S-0275
stream: S-0275
tags: [cli]
touches: [flai/internal/storygit/commit.go, flai/internal/storygit/commit_test.go, flai/internal/itemedit/widen.go, flai/internal/itemedit/widen_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go]
usage:
  source: log
  seconds: 308
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 52
      output: 292
      cache_read: 2184356
      cache_write: 121602
      cost: 1.0361
---
# T-1072 A shared helper commits the files a command wrote on the story branch and widens the story's touches to them

## Work

Write the one helper that `--commit` on `flai issue bump`, `new`, and `close` and on `flai adr new` calls. It is layer 1 and waits for no task: the CLI, MCP, and host channel tasks build on it.

- In `flai/internal/storygit/commit.go`, add a function that takes the runner, the worktree root, the story ID, the paths a command wrote, and a subject. It stages and commits only those paths on the story branch, with the story's commit prefix as `git.md` gives it. When nothing changed it commits nothing and says so. It refuses when the checkout is not the story's worktree (`InWorkTree`, `story/S-nnnn` branch). It returns the commit hash and the paths it committed.
- In `flai/internal/itemedit/widen.go`, add a function that widens a story's touches to paths it does not already cover, as `flai touches` does when it adds paths. It keeps every existing touch, adds nothing a folder touch already covers, and returns what it added and any overlaps `ClaimWatch.Grown` reports.
- Look at S-0269 (`flai task done` widens touches too) before writing the widening: if it has landed, reuse its helper rather than writing a second one.

## Done when

- Tests in `commit_test.go` cover a commit on a story branch with the prefix, a commit with nothing to commit, and a refusal outside a story worktree. They use a temporary repository.
- Tests in `widen_test.go` cover adding a new path, skipping a path a folder touch covers, and keeping existing touches.
- `scripts/flai-test.sh` passes.

## Notes

Written by the planner. The file names are predicted; if the package already has a better home for either function, use it and widen the task's touches.
