---
id: T-2326
type: task
nature: remediation
title: The go-test-json format names each test still running when go test's timeout panics, with its running time
status: backlog
parent: S-0540
owner: alex
created: 2026-10-29T07:36:52Z
updated: 2026-10-29T07:36:52Z
transitions: []
stream: S-0540
tags: [flai, testing]
touches: [flai/internal/verify/gotest.go, flai/internal/verify/parse_test.go, flai/internal/verify/testdata/go-test-timeout.json, design/system/project-manifest.md]
---
# T-2326 The go-test-json format names each test still running when go test's timeout panics, with its running time

## Work

When `go test -json` times out, Go 1.26 writes `panic: test timed out after <d>` and a `running tests:` list as `output` events of the test that is running. It sends no `fail` event for that test, only a `fail` for the package. `parseGoTest` in `flai/internal/verify/gotest.go` then takes the `len(g.failed[ev.Package]) == 0` branch. The package's own output is empty, so the finding names only the package. This is I-0227's second cause.

- On a package `fail` with no failed test, look through the output kept for that package's tests for the timeout panic. Give one finding for each test its `running tests:` list names. Each finding's name is the test, its path is the package's directory, and its message is the panic line with the test's running time, as in `test timed out after 10m0s; TestX had run 1s`. Without a timeout panic, keep today's package finding.
- Add `flai/internal/verify/testdata/go-test-timeout.json`, the events a real timed-out run gives. Record them from a scratch module with two tests, one passing and one sleeping, run with `go test -json -timeout 2s`, and trim the stack frames to a few. Test it in `flai/internal/verify/parse_test.go`, beside `TestGoTestReportsAPackageThatFailsWithNoTestFailing`.
- In `design/system/project-manifest.md` § the `tests` list, under `format`, say that `go-test-json` names the tests running at a timeout.

It waits for no task: no other task of this story touches these files.

## Done when

- A failing `go-test-json` tier whose output is a timeout panic gives one finding per running test, naming it and its running time, and the new test passes.
- A package failure with no failed test and no timeout panic still gives the package finding, and that test still passes.
- `flai test flai/internal/verify` passes.

## Notes
