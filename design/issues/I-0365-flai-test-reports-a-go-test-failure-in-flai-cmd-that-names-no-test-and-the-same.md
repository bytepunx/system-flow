---
id: I-0365
title: flai test reports a go-test failure in flai/cmd that names no test, and the same paths pass when run again
class: defect
status: open
count: 1
cost: 6m
first_reported: 2026-11-13T20:17:49Z
last_reported: 2026-11-13T20:17:49Z
updated: 2026-11-13T20:17:49Z
---

# I-0365 flai test reports a go-test failure in flai/cmd that names no test, and the same paths pass when run again

## Description
flai test reports a go-test failure in flai/cmd that names no test, and the same paths pass when run again

## Instances

### 2026-11-13T20:17:49Z
Story: S-0730.
S-0730, T-3320, 2026-11-13 about 20:15Z. flai test on flai/internal/guard/items.go, items_test.go, and call.go failed its go-test tier after 5m4s with one finding, name github.com/bytepunx/system-flow/flai/cmd, path flai/cmd, and no test, line, or message, so the failure could not be read from the answer. flai test --json on flai/cmd/guard.go right after, with nothing changed, passed every tier. The host was running other stories' agents at the time. The finding should name the test or the package-level output (a panic, a timeout, or a TestMain exit) that failed.

## Remediation
