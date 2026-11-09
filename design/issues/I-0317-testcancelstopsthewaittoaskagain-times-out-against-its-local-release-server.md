---
id: I-0317
title: TestCancelStopsTheWaitToAskAgain times out against its local release server under load and stops a close-out's integration tier
class: defect
status: open
count: 1
cost: 15m
first_reported: 2026-11-09T08:39:50Z
last_reported: 2026-11-09T08:39:50Z
updated: 2026-11-09T08:39:50Z
---

# I-0317 TestCancelStopsTheWaitToAskAgain times out against its local release server under load and stops a close-out's integration tier

## Description
TestCancelStopsTheWaitToAskAgain times out against its local release server under load and stops a close-out's integration tier

## Instances

### 2026-11-09T08:39:50Z
Story: S-0603.
S-0603's close-out (job S-0603-20261109T082735Z, at a6400dd) stopped at the integration tier after 9m21s. One test failed: flai/internal/selfupgrade/selfupgrade_test.go:536 TestCancelStopsTheWaitToAskAgain, "the deadline and the dropped answer: list the flai/v releases of o/r: Get <http://127.0.0.1:42234/repos/o/r/releases?per_page=100>: context deadline exceeded". S-0603 changes nothing in flai/internal/selfupgrade. The host's load average was 106 to 111 on 24 cores. The test's deadline is too tight for a loaded host, so a story that does not touch the package loses a close-out of about 12 minutes.

## Remediation
