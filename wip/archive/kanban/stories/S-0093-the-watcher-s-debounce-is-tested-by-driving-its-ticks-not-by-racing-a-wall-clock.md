---
id: S-0093
type: story
nature: remediation
title: The watcher's debounce is tested by driving its ticks, not by racing a wall clock
status: done
parent: E-0003
owner: alex
created: 2026-09-21T22:48:55Z
updated: 2026-09-21T23:01:34Z
transitions:
  - to: ready
    at: 2026-09-21T22:49:01Z
    by: system-flow
  - to: in-progress
    at: 2026-09-21T22:55:12Z
    by: system-flow
  - to: review
    at: 2026-09-21T22:55:12Z
    by: system-flow
  - to: done
    at: 2026-09-21T23:01:34Z
    by: alex
tags: [cli]
touches: [flai/internal]
---
# S-0093 The watcher's debounce is tested by driving its ticks, not by racing a wall clock

## Goal
`TestAFileStillBeingWrittenIsReportedOnce` (`flai/internal/watch`, from S-0073) writes every 15 ms against a 40 ms tick and asserts one report. It depends on the test's writer goroutine never stalling for two ticks, which a shared CI runner running `-race` does not promise. It has failed the flai workflow on main twice: 2026-09-20 after S-0078 and 2026-09-21 after S-0087 (I-0030, count 2), each time unrelated to the change just pushed, each time green on rerun, and so a red main that says nothing true about the code.

The watcher itself is right: a runner that pauses for two ticks between writes makes it, correctly, report the file as settled midway. The defect is in how the test observes it. `Run`'s package comment even says the package "can be tested with a clock", but the per-tick logic lives inside `Run`'s loop where only a real ticker can reach it.

## Acceptance criteria
- [x] The debounce (what a tick does with the previous snapshot, the pending files, and the new snapshot) is a unit that takes snapshots and returns what is ready to report, with no ticker, timer, or filesystem in it; `Run` only supplies snapshots and ticks
- [x] The still-being-written, added, changed, removed, and removed-while-pending cases are tested by feeding that unit snapshots directly, so they pass or fail the same on every machine regardless of how slow it is
- [x] `TestAFileStillBeingWrittenIsReportedOnce` no longer sleeps to pace writes; nothing in the package's tests asserts a report count that a stalled goroutine could change
- [x] The tests that still run the real ticker (the end-to-end path through `snapshot` and `Run`) keep working and cannot see a half-written file: the test helper writes a file in one step
- [x] `go test -race -short -count=200 ./internal/watch` passes on this host, including while every core is busy
- [x] I-0030 is closed with this story as its remediation

## Tasks
- T-0319 Move the debounce out of Run into a unit that takes snapshots, and test it by feeding it snapshots
- T-0320 Make the test helper write a file in one step, prove it under load, and close I-0030

## Notes
Recorded as I-0030. The deferred fix from the first failure ("drive the clock or the ticks itself instead of sleeping") is this story. Not reproducible on this host even with 48 busy processes on 24 cores, so the test is fixed by construction rather than by observing the failure go away; the 200-run criterion is a regression guard, not the proof.
