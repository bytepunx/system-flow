---
id: T-2555
type: task
nature: remediation
title: wait_for_work and the folder's waits take their snapshot before their first check, so a change made between the two wakes the wait
status: done
parent: S-0560
owner: alex
created: 2026-11-07T08:19:45Z
updated: 2026-11-07T08:45:33Z
transitions:
  - to: ready
    at: 2026-11-07T08:35:23Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T08:35:24Z
    by: system-flow
  - to: done
    at: 2026-11-07T08:45:33Z
    by: agent-S-0560
stream: S-0560
tags: [flai, tests]
touches: [flai/internal/mcpserver/work.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/timing_test.go]
usage:
  source: log
  seconds: 608
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 51
      output: 430
      cache_read: 6289533
      cache_write: 64341
      cost: 2.2751
---
# T-2555 wait_for_work and the folder's waits take their snapshot before their first check, so a change made between the two wakes the wait

## Work

I-0246's instances of S-0553 and S-0555: `TestWaitForWorkPassesOverAHeldStory` and `TestAHeldWaitFailsOnATruncatedStory` answered `timed_out` after the change was made. Reading the code gives a cause in the product, not only in the tests. `waitForWork` in `flai/internal/mcpserver/work.go` runs its first `s.work` and only then takes `before := s.snapshot()`. A change that lands between the two is in `before`, so `diff` never shows it and the wait holds until its deadline. The folder's `waitForWork` and `events` in `flai/internal/mcpserver/folder.go` have the same gap: they run `decide` or `catchUp`, and `hold` takes its snapshot afterwards. A loaded host widens the gap, because `s.work` runs trial merges. `waitForEvents` in `flai/internal/mcpserver/server.go` already has the right order, with `before` taken ahead of `catchUp`.

Confirm this first: write the test below and see it fail against the current order. Then take the snapshot before the first check in all three waits, which means `hold` receives its `before` from the caller. Also put `keepAlive` after the snapshot in `waitForWork`, as `hold` and `waitForEvents` have it, so a first heartbeat means the wait has its fingerprint. The tests the story's next mcpserver task changes rely on that.

Add the test to `flai/internal/mcpserver/timing_test.go`: a change made after the first check and before the snapshot must wake the wait within its hold. Use a seam the server already has if one reaches that point, such as `Options.Now` or the trial judge. Otherwise add the smallest test hook and say so in the commit. This task waits for no other.

## Done when

- [ ] A test in `flai/internal/mcpserver/timing_test.go` fails on the old order and passes on the new one, for `wait_for_work` and for the folder's `wait_for_work` and `wait_for_events`
- [ ] The first heartbeat of each wait comes after its snapshot
- [ ] `flai test flai/internal/mcpserver` passes

## Notes
