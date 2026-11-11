---
id: T-2984
type: task
nature: remediation
title: flai serve takes a served project's watcher's first snapshot before its item snapshot and its first looks
status: done
parent: S-0690
owner: alex
created: 2026-11-10T10:29:34Z
updated: 2026-11-11T02:32:14Z
transitions:
  - to: ready
    at: 2026-11-11T02:17:19Z
    by: agent-S-0690
  - to: in-progress
    at: 2026-11-11T02:17:20Z
    by: agent-S-0690
  - to: done
    at: 2026-11-11T02:32:14Z
    by: agent-S-0690
stream: S-0690
tags: [flai, serve]
touches: [flai/internal/serve/serve.go, flai/internal/serve/serve_test.go]
after: [T-2981]
usage:
  source: log
  seconds: 894
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 51
      output: 15252
      cache_read: 4570021
      cache_write: 130557
      cost: 2.2637
---
# T-2984 flai serve takes a served project's watcher's first snapshot before its item snapshot and its first looks

## Work

It waits for T-2981, whose method it calls.

In `flai/internal/serve/serve.go`, where `reconcile` serves a new project, call the method T-2981 adds on the project's `watcher` right after the watcher is made. Call it before `workitem.SnapshotItems(e.Root)` and before `starter.look(cctx, false)` and the other first looks. Then a change written while those looks run, or between the item snapshot and the watcher's first tick, is reported within a tick. The item snapshot is marked stale by it and the looks run again on it, not at the next minute.

Keep the once-a-minute `workitem.ItemsChanged` and `starter.look` as they are: they still cover an agent that outlived an earlier flai serve. Rewrite the loop's comment that says a change made before the watcher's first look is never reported.

In `flai/internal/serve/serve_test.go`, add a test that reproduces I-0273. Serve a project, wait for its first look, and at once write a change the looks act on, such as a story moved to review. The change must reach the dashboard as a `change` notification, and the run it starts must be queued, well inside a minute: give it a deadline of a few seconds. The window is a race, which the instance saw in 5 of 25 runs under a pinned CPU load. Widen it in the test where the lab allows, such as by writing the change from a hook the first look calls, so that the test fails on the code before this task in most runs. Remove the `watcherReports` probe wait from `TestServeStartsADefinedRoleOnAnEventAndOnASchedule`, and the helper with it if nothing else calls it, since the fix makes it needless. Run that test and the new one under `-count` of at least 20 to show neither waits for the minute.

## Done when

- A served project's watcher takes its first snapshot before the item snapshot and the first looks.
- The new serve test fails on the code before this task in at least one of 20 runs, and the task's Notes record how many; it passes in 20 of 20 after it.
- `TestServeStartsADefinedRoleOnAnEventAndOnASchedule` passes without the probe wait, 20 runs in a row.
- `flai test flai/internal/serve` passes.

## Notes

`TestAChangeMadeDuringServesFirstLooksIsActedOnWithinSeconds` moves the story to review from the `Agent` hook of the first look after the role starter's, the role scheduler's, so the window is not a race: against `main`'s `serve.go` it failed 20 of 20 runs (`go test -short -count=20`, "no watcher run within 10 s of the move made in the first looks"). With the fix, it and `TestServeStartsADefinedRoleOnAnEventAndOnASchedule`, without the `watcherReports` probe wait (the helper is removed), passed 25 of 25 runs each under `-race -count=25`. A first version ended at the run's start and failed 4 of 20 on its own `TempDir` cleanup, with the role's stub still writing (D-0427); it now waits for the run's end and its activity entry. The dashboard's `change` notification is not asserted: `channel.Client.Notify` drops one sent before the dashboard connects, and the watcher's first tick can come first. The run, which only the watcher's report or the minute's look can start, shows the report.
