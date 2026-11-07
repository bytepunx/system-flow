---
id: T-2522
type: task
nature: remediation
title: "TestWaitForEventsAcrossAFolder makes its move only once the wait_for_events call is holding, and passes under -race -count=10 with GOMAXPROCS=2"
status: done
parent: S-0569
owner: alex
created: 2026-11-07T03:47:28Z
updated: 2026-11-07T04:40:24Z
transitions:
  - to: ready
    at: 2026-11-07T04:17:05Z
    by: agent-S-0569
  - to: in-progress
    at: 2026-11-07T04:17:08Z
    by: agent-S-0569
  - to: done
    at: 2026-11-07T04:40:24Z
    by: agent-S-0569
stream: S-0569
tags: [flai, tests]
touches: [flai/internal/mcpserver/folder_test.go]
usage:
  source: log
  seconds: 1396
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 51
      output: 13468
      cache_read: 6361537
      cache_write: 213958
      cost: 2.7344
---
# T-2522 TestWaitForEventsAcrossAFolder makes its move only once the wait_for_events call is holding, and passes under -race -count=10 with GOMAXPROCS=2

## Work

I-0247's third and fifth instances: `TestWaitForEventsAcrossAFolder` (`flai/internal/mcpserver/folder_test.go`) failed at line 371, "changed paths are relative to the folder: []", and at line 362, `timed_out:true`. The test starts the call in a goroutine (`callLater`), sleeps 150 ms, then moves the story. When the call has not taken its snapshot and catch-up (`holdFor`, `flai/internal/mcpserver/server.go`) before the move lands, catch-up answers the event with `Changed: []`. When the hold is slow to see the move, its 3-second `timeout_seconds` and `MaxWait` run out. This is the same race as I-0226 and I-0186.

First reproduce it: run the test with `GOMAXPROCS=2 go test -race -count=20 -run TestWaitForEventsAcrossAFolder ./internal/mcpserver` beside a CPU load, and record the failures in the narrative. Then replace the sleep with a wait for the call to be holding. Use a signal the server already gives if there is one. Otherwise add a test hook in `server.go` and say so in the task's commit. Raise the hold and the test's own 5-second select so that a loaded host does not reach them: a longer bound costs nothing on a fast host. If the other `wait_for_events` tests in the package that take a fixed `timeout_seconds: 3` share the same sleep, fix them the same way and name them in the narrative. Waits for no other task.

## Done when

- [ ] The race is reproduced before the change, or the narrative says how many runs found none
- [ ] `GOMAXPROCS=2 go test -race -count=10 -run TestWaitForEventsAcrossAFolder ./internal/mcpserver` passes beside a CPU load, recorded in the narrative
- [ ] The test asserts no less than before

## Notes
