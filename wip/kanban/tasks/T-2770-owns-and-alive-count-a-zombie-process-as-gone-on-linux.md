---
id: T-2770
type: task
nature: remediation
title: Owns and Alive count a zombie process as gone on Linux
status: done
parent: S-0632
owner: alex
created: 2026-11-09T21:57:34Z
updated: 2026-11-09T22:12:11Z
transitions:
  - to: ready
    at: 2026-11-09T22:02:48Z
    by: agent-S-0632
  - to: in-progress
    at: 2026-11-09T22:02:51Z
    by: agent-S-0632
  - to: done
    at: 2026-11-09T22:12:11Z
    by: agent-S-0632
stream: S-0632
tags: [flai, serve]
touches: [flai/internal/serve/proc_linux.go, flai/internal/serve/proc_unix.go, flai/internal/serve/proc_other.go, flai/internal/serve/proc_linux_test.go]
usage:
  source: log
  seconds: 560
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 14
      output: 59
      cache_read: 771354
      cache_write: 7549
      cost: 0.2848
---
# T-2770 Owns and Alive count a zombie process as gone on Linux

## Work

`process` (`flai/internal/serve/proc_linux.go`) reads the session and the start time from `/proc/<pid>/stat`, but not the state, field 3. `Alive` and `Owns` (`flai/internal/serve/proc_unix.go`) send signal 0, which succeeds on a zombie. So `running()` counts a zombie as running and `settleOrphans` never settles its run (I-0299).

Split the parsing of a stat line out of `process` into a function that takes the line's text. Have it return the state too. `Alive` and `Owns` then report false when the state is `Z` (zombie) or `X` (dead). Keep `process`'s callers in step with the change. In `proc_other.go`, where no state is known without a subprocess, return an empty state, so its behaviour does not change.

Add `flai/internal/serve/proc_linux_test.go`, a new file with a `linux` build tag. Test the parser on fixture lines: a running process, a zombie, and a command name holding spaces and a `)`. Also test that `Owns` and `Alive` are false for a real child left unreaped after it exits (start it with `os/exec`, wait for its state to read `Z`, check, then reap it).

Waits for no task: it shares no file with the reaping task, and the two run together.

## Done when

- `Owns` and `Alive` return false for a process whose `/proc/<pid>/stat` state is `Z` or `X`, and are unchanged otherwise.
- The new tests pass on Linux, and the package still builds for darwin and windows (`GOOS=darwin go vet ./internal/serve/`, `GOOS=windows go vet ./internal/serve/`).
- `flai test` on the four touched files passes.

## Notes
