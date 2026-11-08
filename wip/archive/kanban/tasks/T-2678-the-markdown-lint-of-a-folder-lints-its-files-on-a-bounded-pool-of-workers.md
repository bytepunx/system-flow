---
id: T-2678
type: task
nature: improvement
title: The markdown lint of a folder lints its files on a bounded pool of workers
status: done
parent: S-0600
owner: alex
created: 2026-11-07T23:50:25Z
updated: 2026-11-08T00:16:31Z
transitions:
  - to: ready
    at: 2026-11-07T23:55:20Z
    by: agent-S-0600
  - to: in-progress
    at: 2026-11-07T23:55:25Z
    by: agent-S-0600
  - to: done
    at: 2026-11-08T00:16:31Z
    by: agent-S-0600
stream: S-0600
tags: [cli]
touches: [flai/internal/mdlint/folder.go, flai/internal/mdlint/folder_test.go, flai/internal/check/check.go]
usage:
  source: log
  seconds: 1266
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 8508
      cache_read: 4115139
      cache_write: 126127
      cost: 1.7063
---
# T-2678 The markdown lint of a folder lints its files on a bounded pool of workers

## Work

The folder loop is in `(*checker).markdown()` (`flai/internal/check/check.go`, lines 1898 to 1928): it walks `wip`, reads each `.md` file, and calls `cfg.Lint` one file at a time. Move the walk into a new `flai/internal/mdlint/folder.go`, a function that takes the configuration and a folder and returns each file's findings in path order, reading and linting the files on a pool of at most `runtime.GOMAXPROCS(0)` workers, capped at 16. `(*Config).Lint` only reads the configuration, so the workers share it. Make `markdown()` call it and report the findings as it does today, in the same order.

First layer: it waits for no other task. The kept findings under `.flai-cache` build on it.

## Done when

- [ ] `flai check` reports the same `markdown.MDnnn` findings, in the same order, as before on a folder with findings in several files.
- [ ] A test in `flai/internal/mdlint/folder_test.go` lints a folder on more than one worker and gets each file's findings once, in path order; `go test -race ./internal/mdlint/ ./internal/check/` passes.

## Notes
