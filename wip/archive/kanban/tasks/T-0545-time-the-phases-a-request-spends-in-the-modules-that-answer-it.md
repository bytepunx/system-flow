---
id: T-0545
type: task
nature: research
title: Time the phases a request spends in the modules that answer it
status: done
parent: S-0152
owner: alex
created: 2026-09-29T06:46:36Z
updated: 2026-09-29T06:56:24Z
transitions:
  - to: ready
    at: 2026-09-29T06:50:42Z
    by: agent-S-0152
  - to: in-progress
    at: 2026-09-29T06:50:42Z
    by: agent-S-0152
  - to: done
    at: 2026-09-29T06:56:24Z
    by: agent-S-0152
stream: S-0152
tags: []
touches: [flai/internal/workitem, flai/internal/execx, flai/internal/hostapi, flai/cmd]
usage:
  source: log
  seconds: 342
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 90
      output: 651
      cache_read: 8800495
      cache_write: 58072
      cost: 3.5046
---
# T-0545 Time the phases a request spends in the modules that answer it

## Work

Record phases where the time goes, without the transport: opening the repository, listing items (kanban and archive), loading the board, each subprocess `execx.System` runs (git, flai) by program and subcommand, threads, the documentation tree, and search. `flai hostapi --timing` prints the same breakdown to stderr, so a method is measured in-process with no connection at all. `flai serve` gains an opt-in profiler address for CPU and heap profiles.

## Done when

- A slow `board.get` or `item.get` names, in its event, which phases took the time and how many subprocesses it ran.
- `flai hostapi --timing` and the profiler are documented in `docs/users/flai.md` and `docs/operators/settings.md`.
- Tests cover the phases recorded by a method and by a subprocess.

## Notes
