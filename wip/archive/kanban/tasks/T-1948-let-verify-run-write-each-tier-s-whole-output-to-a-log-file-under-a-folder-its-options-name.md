---
id: T-1948
type: task
nature: improvement
title: Let verify.Run write each tier's whole output to a log file under a folder its options name
status: done
parent: S-0457
owner: alex
created: 2026-10-21T06:09:38Z
updated: 2026-10-21T06:34:37Z
transitions:
  - to: ready
    at: 2026-10-21T06:13:37Z
    by: agent-S-0457
  - to: in-progress
    at: 2026-10-21T06:13:37Z
    by: agent-S-0457
  - to: done
    at: 2026-10-21T06:34:37Z
    by: agent-S-0457
stream: S-0457
tags: [flai]
touches: [flai/internal/verify/run.go, flai/internal/verify/run_test.go, flai/internal/verify/story.go, flai/internal/verify/verify.go]
usage:
  source: log
  seconds: 1260
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 9422
      cache_read: 1227679
      cache_write: 61896
      cost: 0.7985
---
# T-1948 Let verify.Run write each tier's whole output to a log file under a folder its options name

## Work

- Give `RunOptions` in `flai/internal/verify/run.go` a folder for logs, such as `LogDir string`. When it is set, `runTier` writes the tier's combined output, as it runs, to `<LogDir>/<tier>.log`, creating the folder, and `TierResult` records the file's path, such as `Log`. When it is empty nothing is written, so `flai verify`, `flai test`, and the close-out behave as now.
- Carry the path into the step `Tiers` reports in `flai/internal/verify/story.go`, so a caller of `verify.Tiers` sees which file holds each tier's output.
- A failure to write the log is not a tier's failure: the tier's result stands and the log is left out.
- Test it in `flai/internal/verify/run_test.go`: a failing tier's whole output is in its log, a tier not reached has none, and no `LogDir` writes nothing.

It waits for nothing: it changes no path the other first-layer task changes, so the two run together.

## Done when

- With `LogDir` set, every tier that ran has its whole output in its own file and its step names the file; without it, nothing is written.
- `flai test` on the changed paths passes.

## Notes
