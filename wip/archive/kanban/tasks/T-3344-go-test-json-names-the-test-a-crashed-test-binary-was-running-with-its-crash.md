---
id: T-3344
type: task
nature: remediation
title: go-test-json names the test a crashed test binary was running, with its crash and exit, when a package fails with no test failing
status: done
parent: S-0782
owner: alex
created: 2026-11-13T22:53:13Z
updated: 2026-11-13T23:08:16Z
transitions:
  - to: ready
    at: 2026-11-13T22:56:00Z
    by: agent-S-0782
  - to: in-progress
    at: 2026-11-13T22:56:01Z
    by: agent-S-0782
  - to: done
    at: 2026-11-13T23:08:16Z
    by: agent-S-0782
stream: S-0782
tags: [flai, testing]
touches: [flai/internal/verify/gotest.go, flai/internal/verify/parse_test.go, flai/internal/verify/testdata/go-test-crash.json, design/system/project-manifest.md]
usage:
  source: log
  seconds: 735
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 27
      output: 7577
      cache_read: 1863402
      cache_write: 91280
      cost: 1.1476
---
# T-3344 go-test-json names the test a crashed test binary was running, with its crash and exit, when a package fails with no test failing

## Work

Criterion 1. In `flai/internal/verify/gotest.go`, the `fail` case for a package with no failed test and no timeout panic gives one finding named for the package, whose message is `packageMessage` of the package's own lines, the output events with no `Test`. When the test binary dies mid-run, `go test -json` writes the crash as output of the test that was running and sends that test no `fail`, so those lines hold only `FAIL<tab><package>` and the finding names no test, line, or message (I-0365).

1. Record the shape first. In a scratch module outside the checkout, write a package whose test panics in a goroutine and one whose test calls `os.Exit(1)`, run each with `go test -json -count=1`, and keep the events as `flai/internal/verify/testdata/go-test-crash.json`, its paths rewritten to the fixture root as `go-test-timeout.json`'s are. Note in the task's `## Notes` what go names for each: the test left without an end, the `panic:` or `fatal error:` line, and the `exit status n` line.
2. Track each test's `run` and its end (`pass`, `fail`, `skip`) per package. For a package that fails with no failed test, after `timedOut` finds none: give one finding per test that started and never ended, a subtest standing for its parent as `timedOut` does, named for the test, placed by the crash's first stack frame under the root as `testFrame` and `frame` place one, else at the package's folder, with the crash line and the exit line as its message. With no such test, keep the package's finding, its message the last lines of the package's whole output (`g.all`), its tests' included, never empty.
3. Add tests to `flai/internal/verify/parse_test.go`: the fixture's two packages, a binary killed with `signal: killed` and no panic, and the existing `TestGoTestReportsAPackageThatFailsWithNoTestFailing`, which must still pass.
4. Add a sentence to the `format` paragraph of `design/system/project-manifest.md`, after the S-0540 one, saying what `go-test-json` names for a package that fails with no test failing (S-0782, I-0365).

It waits for no task: it is the story's first layer.

## Done when

- A package whose test binary panics in a goroutine, or whose test calls `os.Exit`, gives a finding naming that test, its crash line, and the exit, from the recorded fixture.
- A package that fails with no test started and unended gives a finding whose message is not empty.
- `flai test flai/internal/verify` passes in the worktree.
- `design/system/project-manifest.md` says what `go-test-json` names for such a package.

## Notes
