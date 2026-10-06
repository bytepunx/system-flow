---
id: T-1013
type: task
nature: improvement
title: flai stream sync's trial merge does not report design/issues/summary.md as a conflict between two open story branches
status: done
parent: S-0278
owner: alex
created: 2026-10-06T11:32:24Z
updated: 2026-10-06T19:59:54Z
transitions:
  - to: ready
    at: 2026-10-06T19:56:21Z
    by: agent-S-0278
  - to: in-progress
    at: 2026-10-06T19:56:22Z
    by: agent-S-0278
  - to: done
    at: 2026-10-06T19:59:54Z
    by: agent-S-0278
stream: S-0278
tags: [flai]
touches: [flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go]
after: [T-1012]
usage:
  source: log
  seconds: 212
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 42
      output: 13538
      cache_read: 1589339
      cache_write: 58884
      cost: 0.9467
---
# T-1013 flai stream sync's trial merge does not report design/issues/summary.md as a conflict between two open story branches

## Work

In the first and third instances of I-0074, `checkSync`'s trial merge opened a thread (TH-0126, TH-0130) because two open branches both rewrote `summary.md`. Once T-1012 makes sync and acceptance regenerate the file, that conflict costs nothing, so the trial merge should not report it.

In `flai/cmd/stream_sync.go`, drop the generated paths of T-1012's list from what `trialMerge` returns. A pair whose only conflicts are generated files then counts as clean. `reportConflicts` opens no thread for that pair, and resolves the pair's open thread as it does for any clean pair. `--json` shows the pair `clean`.

Add tests to `flai/cmd/stream_sync_test.go`:

- Two open branches that conflict only in `summary.md` report clean and open no thread.
- Two open branches that also conflict in another file report that file alone.

This task waits for T-1012: it reuses that task's list of generated paths, and both add tests to `stream_sync_test.go`.

## Done when

- The trial merge reports no conflict, and opens no thread, for two branches that conflict only in `design/issues/summary.md`.
- Any other conflicting path is still reported.
- The new tests and `scripts/flai-test.sh` pass.

## Notes
