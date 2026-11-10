---
id: T-2915
type: task
nature: improvement
title: A test keeps mcpserver's tests off the process, the package's -race time is measured against I-0234's, and I-0234 is closed
status: done
parent: S-0684
owner: alex
created: 2026-11-10T08:47:57Z
updated: 2026-11-10T10:37:56Z
transitions:
  - to: ready
    at: 2026-11-10T09:42:59Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T09:43:00Z
    by: system-flow
  - to: done
    at: 2026-11-10T10:37:56Z
    by: agent-S-0684
stream: S-0684
tags: [flai, go, mcp, testing]
touches: [flai/internal/mcpserver/process_test.go, design/issues/I-0234-flai-internal-mcpserver-s-tests-still-take-136-to-162-s-under-race-since-the.md, design/issues/summary.md]
after: [T-2906, T-2910]
usage:
  source: log
  seconds: 3296
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 79
      output: 22634
      cache_read: 12768229
      cache_write: 343333
      cost: 4.9287
---
# T-2915 A test keeps mcpserver's tests off the process, the package's -race time is measured against I-0234's, and I-0234 is closed

## Work

This task waits for the third and fourth tasks, which leave the package's tests with nothing set on the process. It is the last layer.

- Add `flai/internal/mcpserver/process_test.go`, a parallel test that reproduces I-0234's cause. It reads the package's `*_test.go` files with `go/parser` and fails, naming the file and line, on any call to `t.Setenv`, `os.Setenv`, `os.Unsetenv`, or `gittest.Identity`, and on any assignment to a package-level variable. It passes over `TestMain` and `session_test.go`'s one test of `SessionFromEnv`. It also fails on a top-level `Test` function, other than `TestMain`, that does not call `t.Parallel()`. So a new serial test that would bring the serial phase back fails at once.
- Measure from `flai/`: run `go test -race -count=1 ./internal/mcpserver` three times, noting each wall time and the host's load average (`uptime`). Run it once more with `-short`. Compare with I-0234's 136 to 162 s, whose critical path was about 70 s of serial tests followed by the two prime-on-this-repository tests. Write the figures in this task's Notes. Never write a figure you did not observe.
- Close the issue with `flai issue close I-0234 --reason "<what fixed it and the measured times>"`. That writes the issue file and `design/issues/summary.md`.

## Done when

- [ ] `process_test.go` passes. Adding a `t.Setenv` to any other test file of the package makes it fail, which you check by hand and then revert.
- [ ] The three `-race` timings and the `-short` timing, each with its load, are in this task's Notes.
- [ ] I-0234 is closed with a reason naming the `Options` fields and the measured times.
- [ ] `flai test` passes on `flai/internal/mcpserver`.

## Notes

Measured by the story's agent on 2026-11-10 on a 24-core host. Each run is `go test -race -count=1 ./internal/mcpserver` from `flai/`. The time is what `go test` reports for the package; the load is the 1-minute average from `uptime` before and after the run.

The host's load was 53 to 156 throughout, against I-0234's 40 to 57. So main's package, at `d8b005efe`, was run alternately with the story's branch in a scratch worktree, as a baseline under the same load.

| Run | Story branch | Its load | main | Its load |
|-----|--------------|----------|------|----------|
| 1 | 317.7 s | 121.6 to 148.8 | 354.0 s | 146.3 to 64.7 |
| 2 | 200.0 s | 60.9 to 85.7 | 352.7 s | 81.7 to 94.7 |
| 3 | 155.5 s | 93.9 to 53.3 | 414.6 s | 62.4 to 90.0 |
| `-short` | 57.2 s | 53.3 to 62.3 | 177.6 s | 90.0 to 156.5 |

Every run passed, the real-git tests included, which skip under `-short`. Under the same load the branch took 155 to 318 s against main's 353 to 415 s, and 57 s against 178 s with `-short`. With no serial phase left, what remains is the two prime-on-this-repository tests. The figures are too noisy at this load to put a number on the gain.

`process_test.go` reproduces the cause. With a planted `t.Setenv`, `os.Setenv`, a package-variable assignment, and a test without `t.Parallel()` in `outdated_test.go`, it failed naming all four (`outdated_test.go:119` to `125`). It passed once they were reverted.
