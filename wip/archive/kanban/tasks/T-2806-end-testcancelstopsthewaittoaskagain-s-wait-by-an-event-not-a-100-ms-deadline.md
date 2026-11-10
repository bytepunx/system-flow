---
id: T-2806
type: task
nature: remediation
title: End TestCancelStopsTheWaitToAskAgain's wait by an event, not a 100 ms deadline, and reproduce a slow first answer
status: done
parent: S-0638
owner: alex
created: 2026-11-09T22:08:00Z
updated: 2026-11-09T23:00:14Z
transitions:
  - to: ready
    at: 2026-11-09T22:51:01Z
    by: agent-S-0638
  - to: in-progress
    at: 2026-11-09T22:51:01Z
    by: agent-S-0638
  - to: done
    at: 2026-11-09T23:00:14Z
    by: agent-S-0638
stream: S-0638
tags: [flai, testing, flaky-test]
touches: [flai/internal/selfupgrade/selfupgrade.go, flai/internal/selfupgrade/selfupgrade_test.go]
usage:
  source: log
  seconds: 553
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 951
      cache_read: 1487504
      cache_write: 11247
      cost: 0.4065
---
# T-2806 End TestCancelStopsTheWaitToAskAgain's wait by an event, not a 100 ms deadline, and reproduce a slow first answer

## Work

`TestCancelStopsTheWaitToAskAgain` in `flai/internal/selfupgrade/selfupgrade_test.go` cancels `List` with a 100 ms wall-clock deadline. Under load the first attempt takes longer than that, so the deadline cuts it before the dropped answer and the test fails (I-0317). Make the cancel follow the first attempt's failure instead:

- In `get` in `flai/internal/selfupgrade/selfupgrade.go`, add a package-level test seam, nil in production, called once the pause between attempts begins, before the `select`. Say in its comment that tests set it.
- Rewrite the test to cancel its context from that seam, with `retryPause` still an hour, and to assert the dropped answer (`unexpected EOF`), the context's error (`context.Canceled`, or `DeadlineExceeded` if a cause is used), and exactly one request. Reset the seam with `t.Cleanup`.
- Add a test that reproduces the load: a front handler that holds the first answer for longer than the old 100 ms deadline, say 300 ms, before it drops it, and still gets the same three results.

Waits for no task: it is the first layer.

## Done when

- Both tests pass with `flai test flai/internal/selfupgrade/selfupgrade_test.go`.
- They pass with `go test -run 'TestCancelStopsTheWaitToAskAgain|<the new test>' -count 50 ./internal/selfupgrade/` under `-race`.
- No wall-clock deadline is left in the package's retry tests.

## Notes
