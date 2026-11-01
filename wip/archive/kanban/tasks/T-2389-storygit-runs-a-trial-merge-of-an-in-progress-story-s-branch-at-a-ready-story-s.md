---
id: T-2389
type: task
nature: improvement
title: storygit runs a trial merge of an in-progress story's branch at a ready story's base and returns the conflicting paths
status: done
parent: S-0551
owner: alex
created: 2026-10-31T07:42:57Z
updated: 2026-10-31T08:56:59Z
transitions:
  - to: ready
    at: 2026-10-31T08:42:38Z
    by: agent-S-0551
  - to: in-progress
    at: 2026-10-31T08:42:39Z
    by: agent-S-0551
  - to: done
    at: 2026-10-31T08:56:59Z
    by: agent-S-0551
stream: S-0551
tags: [flai, go]
touches: [flai/internal/storygit/trial.go, flai/internal/storygit/trial_test.go]
usage:
  source: log
  seconds: 860
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 29
      output: 9203
      cache_read: 1809685
      cache_write: 66372
      cost: 0.916
---
# T-2389 storygit runs a trial merge of an in-progress story's branch at a ready story's base and returns the conflicting paths

## Work

Add a function in a new `flai/internal/storygit/trial.go` that answers, for a ready story and an in-progress story, whether the in-progress story's branch, `story/<S-nnnn>`, merges cleanly at the ready story's base, the main branch's head, and which paths conflict when it does not. Build it on `TrialMerge` in `sync.go`, which runs `git merge-tree --write-tree`, and keep its git 2.38 gate: below it, or with no branch for the in-progress story, answer that the merge could not run, so that the caller holds as today. It waits for nothing: it reads git only, and the hold verdict that uses it takes it as a function.

Keep `sync.go` unchanged so that this task shares no file with the sync code.

## Done when

- A test in `trial_test.go`, on a temporary repository, shows a clean trial merge answering clean, a conflicting one answering the conflicting paths, and a missing branch answering that it could not run.
- `flai test flai/internal/storygit` passes.

## Notes

The ready story has no branch until it starts, so its base is the main branch's head.
