---
id: T-2528
type: task
nature: remediation
title: TestServeStartsADefinedRoleOnAnEventAndOnASchedule waits for the watcher's run for the move within a bound a host at load 76 does not reach
status: backlog
parent: S-0569
owner: alex
created: 2026-11-07T03:47:46Z
updated: 2026-11-07T04:16:38Z
transitions: []
stream: S-0569
tags: [flai, tests]
touches: [flai/internal/serve/serve_test.go, flai/internal/serve/agents_test.go]
after: [T-2525]
---
# T-2528 TestServeStartsADefinedRoleOnAnEventAndOnASchedule waits for the watcher's run for the move within a bound a host at load 76 does not reach

## Work

I-0247's first, second, fourth, and fifth instances: `TestServeStartsADefinedRoleOnAnEventAndOnASchedule` (`flai/internal/serve/serve_test.go`) failed at line 1009, "never happened: the watcher's run for the move ended", at load averages of 64 to 76. The test waits for the watcher to be armed (line 999) before the move (line 1006). So it does not race itself: `waitFor` (`flai/internal/serve/agents_test.go`, 30 seconds, raised by S-0548 for I-0226) ran out while the stub role's run started and ended on a loaded host.

First find what the 30 seconds are spent on: time each step from the move to the run's end under `GOMAXPROCS=2 go test -race -count=10 -run TestServeStartsADefinedRoleOnAnEventAndOnASchedule ./internal/serve` beside a CPU load, and record the times in the narrative. If one step waits on a fixed interval, such as a watcher poll, a debounce, or a run's restart delay, make the test shorten that interval through its options rather than wait longer. Otherwise give this wait a bound a host at load 76 does not reach. If you change `waitFor` itself, check every caller in the package still passes. Waits for T-2525: both edit test files of `flai/internal/serve`, one package, so a half-done edit of one would break the other's build and its stress runs.

## Done when

- [ ] The narrative records where the time from the move to the run's end goes on a loaded host
- [ ] `GOMAXPROCS=2 go test -race -count=10 -run TestServeStartsADefinedRoleOnAnEventAndOnASchedule ./internal/serve` passes beside a CPU load, recorded in the narrative
- [ ] `go test -race ./internal/serve` passes
- [ ] The test asserts no less than before

## Notes
