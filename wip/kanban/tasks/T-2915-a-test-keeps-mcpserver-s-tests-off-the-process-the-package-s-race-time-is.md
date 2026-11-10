---
id: T-2915
type: task
nature: improvement
title: A test keeps mcpserver's tests off the process, the package's -race time is measured against I-0234's, and I-0234 is closed
status: in-progress
parent: S-0684
owner: alex
created: 2026-11-10T08:47:57Z
updated: 2026-11-10T09:43:00Z
transitions:
  - to: ready
    at: 2026-11-10T09:42:59Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T09:43:00Z
    by: system-flow
stream: S-0684
tags: [flai, go, mcp, testing]
touches: [flai/internal/mcpserver/process_test.go, design/issues/I-0234-flai-internal-mcpserver-s-tests-still-take-136-to-162-s-under-race-since-the.md, design/issues/summary.md]
after: [T-2906, T-2910]
usage:
  source: log
  seconds: 20
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 6
      output: 19
      cache_read: 728610
      cache_write: 2965
      cost: 0.2693
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
