---
id: T-1585
type: task
nature: improvement
title: flai stream sync, flai task done, and flai accept merge the ADR index themselves, and the trial merge leaves it out of a pair's conflicts
status: done
parent: S-0306
owner: alex
created: 2026-10-14T23:04:12Z
updated: 2026-10-14T23:27:31Z
transitions:
  - to: ready
    at: 2026-10-14T23:07:55Z
    by: agent-S-0306
  - to: in-progress
    at: 2026-10-14T23:07:56Z
    by: agent-S-0306
  - to: done
    at: 2026-10-14T23:27:31Z
    by: agent-S-0306
stream: S-0306
tags: []
touches: [flai/internal/issues/generated.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go, docs/users/flai-reference.md]
after: [T-1583]
usage:
  source: log
  seconds: 1175
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 112
      output: 695
      cache_read: 7181089
      cache_write: 220173
      cost: 3.416
---
# T-1585 flai stream sync, flai task done, and flai accept merge the ADR index themselves, and the trial merge leaves it out of a pair's conflicts

## Work

In `issues.SyncFiles`, widen `Merged` so that `Match` also takes `adrindex.Path(r)` and `Merge` writes it in the worktree with `adrindex.Merge` over `storygit.Stages`, as `mergeStopped` does for an issue file. Every sync caller (`flai stream sync`, `flai task done`, `flai accept`, the conflict log) then merges it, and the trial merge leaves it out of a pair's conflicts. Update the doc comments there and the `flai stream sync` help text in `flai/cmd/stream.go` that lists what a stop does not wait for, and regenerate `docs/users/flai-reference.md` with `make flai-reference`.

Add a test in `flai/cmd/stream_sync_test.go` beside the summary.md ones that reproduces I-0099: two story branches each add an ADR row to the index; one is merged to main; the other's sync finishes clean, keeping both rows in number order; and a trial merge of the two open branches reports no conflict on the index. It fails without the change.

Waits for T-1583, whose `adrindex` package it calls.

## Done when

- The test passes, and fails against `main`'s `generated.go`.
- `flai test` on the changed paths passes.

## Notes
