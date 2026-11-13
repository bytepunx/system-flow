---
id: I-0365
title: flai test reports a go-test failure in flai/cmd that names no test, and the same paths pass when run again
class: defect
status: closed
count: 1
cost: 6m
first_reported: 2026-11-13T20:17:49Z
last_reported: 2026-11-13T20:17:49Z
updated: 2026-11-13T23:08:28Z
---

# I-0365 flai test reports a go-test failure in flai/cmd that names no test, and the same paths pass when run again

## Description
flai test reports a go-test failure in flai/cmd that names no test, and the same paths pass when run again

## Instances

### 2026-11-13T20:17:49Z
Story: S-0730.
S-0730, T-3320, 2026-11-13 about 20:15Z. flai test on flai/internal/guard/items.go, items_test.go, and call.go failed its go-test tier after 5m4s with one finding, name github.com/bytepunx/system-flow/flai/cmd, path flai/cmd, and no test, line, or message, so the failure could not be read from the answer. flai test --json on flai/cmd/guard.go right after, with nothing changed, passed every tier. The host was running other stories' agents at the time. The finding should name the test or the package-level output (a panic, a timeout, or a TestMain exit) that failed.

## Remediation

Story S-0782 remediates this issue, created from it at 2026-11-13T22:46:52Z.
Closed 2026-11-13T23:08:28Z: S-0782 (T-3344): the go-test-json format now names each test left running when a package fails with no test failing and no timeout, as when its test binary dies on a goroutine panic, a fatal error, or os.Exit, with the crash's panic: or fatal error: line, placed by the test's frame in the crash's stacks; with no test running it gives the last lines of the package's whole output, its tests' included, and never an empty message. A fixture recorded from go1.26.0 (flai/internal/verify/testdata/go-test-crash.json) and two tests in parse_test.go reproduce the empty finding. The flake behind the instance kept no output, so a recurrence now names its test and is recorded as a new issue.
