---
id: S-0377
type: story
nature: remediation
title: "TestActivityWritersTakeTurns fails about one run in three: two writers that both find the activity lock stale can both take it, and an entry is lost"
status: done
owner: alex
created: 2026-10-12T16:53:00Z
updated: 2026-10-13T17:26:49Z
transitions:
  - to: ready
    at: 2026-10-13T16:49:07Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-13T16:49:14Z
    by: agent-S-0377
  - to: review
    at: 2026-10-13T17:26:00Z
    by: agent-S-0377
  - to: done
    at: 2026-10-13T17:26:49Z
    by: orchestrator
tags: [flai, concurrency, strategic-agents]
touches: [flai/internal/filelock/filelock.go, flai/internal/filelock/filelock_unix.go, flai/internal/filelock/filelock_windows.go, flai/internal/filelock/filelock_test.go, flai/go.mod, design/tech/go-libraries.md, flai/internal/workitem/activity.go, flai/internal/workitem/activity_test.go, flai/internal/guard/subagents.go, flai/internal/guard/subagents_test.go, flai/internal/serve/restart.go, flai/internal/serve/restart_test.go, design/system/strategic-agents.md, design/issues/I-0140-testactivitywriterstaketurns-fails-about-one-run-in-three-two-writers-that-both-find-the-activity-lock-stale-can-both-take-it-and-an-entry-is-lost.md, design/issues/summary.md, flai/cmd/guard_test.go, flai/cmd/promote_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2223
  turns:
    - day: 2026-10-13
      ceremony: 2
      test_runs: 1
      hand_edits: 3
      work: 33
  models:
    - model: claude-opus-5-5
      input: 188
      output: 61103
      cache_read: 8616369
      cache_write: 412314
      cost: 5.6265
  strategic:
    - kind: orchestrator
      seconds: 633
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 90
          output: 1535
          cache_read: 10337160
          cache_write: 37729
          cost: 2.6157
        - model: claude-sonnet-5-5
          input: 8
          output: 53
          cache_read: 85412
          cache_write: 50168
          cost: 0.1096
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-12T16:53:00Z
  value: 37.5
  by: planner-S-0377
  at: 2026-10-13T16:47:43Z
forecast:
  duration: 1h
  delivery: 2026-10-13T19:05:00Z
  basis: "flai forecast's 5m (131 s per unit over 8 small remediation stories, size 2 with no touches) raised to 1h for a new cross-platform lock package, three call sites moved to it, a crash test in a child process, and a close-out of about ten minutes; delivery is flai's 18:10Z start plus that hour"
  by: planner-S-0377
  at: 2026-10-13T16:47:43Z
finalized:
  by: orchestrator
  at: 2026-10-13T16:49:04Z
---
# S-0377 TestActivityWritersTakeTurns fails about one run in three: two writers that both find the activity lock stale can both take it, and an entry is lost

## Goal

This story remediates [I-0140](../../../design/issues/I-0140-testactivitywriterstaketurns-fails-about-one-run-in-three-two-writers-that-both-find-the-activity-lock-stale-can-both-take-it-and-an-entry-is-lost.md), "TestActivityWritersTakeTurns fails about one run in three: two writers that both find the activity lock stale can both take it, and an entry is lost". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0140 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0140 is closed with `flai issue close I-0140 --reason` saying what fixed it

## Tasks
- T-1510 flai/internal/filelock takes an OS file lock that its holder's death frees, on unix and Windows
- T-1511 lockActivity takes the activity document's lock through filelock, and TestActivityWritersTakeTurns passes every run
- T-1512 The guard's sub-agent record and flai serve's agents state lock through filelock, with no age break
- T-1513 Close I-0140, saying the OS file lock in flai/internal/filelock fixed it

## Notes

Cost of delay inputs set by flai from I-0140. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-12T16:22:16Z, 0 days before this story; under one cycle counts as one).

### Planning

Proposed remedy, from I-0140's instance: replace each lock file that a writer creates with `O_EXCL` and breaks by age (`os.Stat` then `os.Remove`) with an operating-system lock on a file that stays in place. On unix that is `flock(2)` through `syscall.Flock`. On Windows it is `LockFileEx` through `golang.org/x/sys/windows`, since flai ships for Windows too. The kernel frees such a lock when its holder dies, so no lock is ever stale, no age is needed, and two writers can never both hold it. A rename-to-a-unique-name break narrows the race but still loses one when a third writer creates a lock between the rename and the check, so it is not proposed.

The same Stat-then-Remove break is in `lockRecord` (`flai/internal/guard/subagents.go`, the guard's record of a session's sub-agents) and `lockFile` (`flai/internal/serve/restart.go`, the agents state). The plan moves all three to one package, `flai/internal/filelock`, so the cause is gone wherever it occurs.

Touches, where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/workitem/activity.go` | I-0140's instance names `lockActivity` there |
| `flai/internal/workitem/activity_test.go` | layout: `TestActivityWritersTakeTurns`; co-change with `activity.go` |
| `flai/internal/filelock/filelock.go`, `filelock_unix.go`, `filelock_windows.go`, `filelock_test.go` | layout: the new package, split by build tag as `flai/internal/serve/proc_unix.go` and `proc_windows.go` are |
| `flai/go.mod` | layout: `golang.org/x/sys` moves from indirect to direct for `LockFileEx` |
| `design/tech/go-libraries.md` | design: the documentation convention keeps `design/tech` current when a dependency is added |
| `flai/internal/guard/subagents.go`, `subagents_test.go` | layout: `lockRecord` has the same break |
| `flai/internal/serve/restart.go`, `restart_test.go` | layout: `lockFile` has the same break |
| `design/system/strategic-agents.md` | design: § The orchestrator says "a lock older than ten seconds is taken as left behind", which this changes |
| `design/issues/I-0140-…md`, `design/issues/summary.md` | criterion 2: `flai issue close` writes both; `design/issues` is a shared path, so neither holds another story |

`flai touches suggest` with `activity.go` and `activity_test.go` listed `design/system/agent-narrative.md` and `flai/internal/workitem/narrative.go` (2 of 5 commits); they were changed with the activity documents' first feature, not with the lock, so they are left out. With `subagents.go` and `restart.go` it listed serve and harness files that change with agent features, not with the locks, so those are left out too. No folder touch is kept: each task names its files.

Forecast: flai's 5m rests on size 2, two criteria and no touches, since the story declared none. It is raised to 1h. The work is a new package for two platforms, three call sites, a test that kills a child process holding the lock, and a close-out whose tiers take about ten minutes. Delivery is flai's 18:10Z, the start it played out from the pull order, plus that hour.

Cost of delay: 37.50 USD a week stands. It is flai cod's figure from flai's 15m time-lost input: 15m per 168h cycle at 150 USD an hour. The input was left as flai set it.

### Accepted by the orchestrator

- Verified: d3a204fc86c68b33ffafdc7014b425c50573f8a6
- At: 2026-10-13T17:26:49Z

Verdict: meets all criteria (verifier at d3a204fc86c68b33ffafdc7014b425c50573f8a6; flai verify passed every step at that commit)

- 1: flai/internal/filelock/filelock.go, flai/internal/filelock/filelock_unix.go, flai/internal/filelock/filelock_windows.go, flai/internal/filelock/filelock_test.go, flai/internal/workitem/activity.go, flai/internal/workitem/activity_test.go, flai/internal/guard/subagents.go, flai/internal/guard/subagents_test.go, flai/internal/serve/restart.go, flai/internal/serve/restart_test.go, flai/go.mod, design/tech/go-libraries.md, design/system/strategic-agents.md
- 2: design/issues/I-0140-testactivitywriterstaketurns-fails-about-one-run-in-three-two-writers-that-both-find-the-activity-lock-stale-can-both-take-it-and-an-entry-is-lost.md, design/issues/summary.md
