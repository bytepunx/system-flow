---
id: S-0690
type: story
nature: remediation
title: flai serve takes a project's watcher's first snapshot after its first looks, so a change made in between waits up to a minute for the next look
status: ready
owner: alex
created: 2026-11-10T08:49:01Z
updated: 2026-11-11T01:46:11Z
transitions:
  - to: ready
    at: 2026-11-11T01:44:25Z
    by: alex
tags: [flai, serve]
touches: [flai/internal/watch/watch.go, flai/internal/watch/watch_test.go, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, design/system/flai-cli.md, design/issues/I-0273-flai-serve-takes-a-project-s-watcher-s-first-snapshot-after-its-first-looks-so.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 9
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 4
          cache_read: 120237
          cache_write: 676
          cost: 0.0332
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-10T08:49:01Z
  value: 12.5
  by: planner-S-0690
  at: 2026-11-10T10:30:31Z
forecast:
  duration: 40m
  delivery: 2026-11-11T03:25:00Z
  basis: "Its own forecast of 40m; 3rd in the pull order with an in-progress limit of 3, behind S-0674, S-0676, S-0688, S-0695 and S-0692."
  by: flai
  at: 2026-11-11T01:46:11Z
finalized:
  by: orchestrator
  at: 2026-11-10T11:04:46Z
---
# S-0690 flai serve takes a project's watcher's first snapshot after its first looks, so a change made in between waits up to a minute for the next look

## Goal

This story remediates [I-0273](../../../design/issues/I-0273-flai-serve-takes-a-project-s-watcher-s-first-snapshot-after-its-first-looks-so.md), "flai serve takes a project's watcher's first snapshot after its first looks, so a change made in between waits up to a minute for the next look". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0273 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0273 is closed with `flai issue close I-0273 --reason` saying what fixed it

## Tasks
- T-2981 watch.Watcher takes its first snapshot on a call of its own before Run, and Run compares its first tick with it
- T-2984 flai serve takes a served project's watcher's first snapshot before its item snapshot and its first looks
- T-2988 flai-cli.md says the watcher's first snapshot is taken before a served project's first looks
- T-2989 Close I-0273 with the fix that takes the watcher's first snapshot before the first looks

## Notes

Cost of delay inputs set by flai from I-0273. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T05:23:04Z, 3.1 days before this story; under one cycle counts as one).

### Planning

The proposed fix, from I-0273's instance: `watch.Watcher` gains a call that takes its first snapshot at once, and `reconcile` in `flai/internal/serve/serve.go` makes it right after the watcher is made, before `workitem.SnapshotItems` and the first looks. `Run` then compares its first tick with that snapshot, so a change written while the first looks run is reported within a tick, not at the next minute.

Touches, all files, no folder:

| Touch | From |
|-------|------|
| `flai/internal/watch/watch.go` | layout: `Watcher.Run` takes the first snapshot there |
| `flai/internal/watch/watch_test.go` | layout: the watch package's tests |
| `flai/internal/serve/serve.go` | layout: `reconcile` starts the watcher after the first looks |
| `flai/internal/serve/serve_test.go` | co-change (13 of 61 commits with the two sources) and the instance's test, `TestServeStartsADefinedRoleOnAnEventAndOnASchedule`, with its probe wait |
| `design/system/flai-cli.md` | co-change (5 of 61) and design: its item snapshot paragraph names a change before the first look as one the watcher misses |
| `design/issues/I-0273-...md` | criterion 2 |
| `design/issues/summary.md` | criterion 2: `flai issue close` regenerates it |

The story declared no touches before this plan. It has no `after`: it reads nothing another open story writes.

Forecast: `flai forecast` gave 21m, 134 s per unit of size over 54 done remediation stories at size 9. Raised to 40m: the serve test must reproduce a race that showed in 5 of 25 runs under load, and is run 20 times. The close-out also runs the integration and smoke tiers that a `flai/` change selects. The delivery moves with it, to 14:40Z.

Cost of delay: `flai cod` gives 12.50 USD a week, 5m lost per 168h cycle at 150 USD an hour. The value stands as flai worked it out.

Tasks, in three layers:

1. T-2981, the watcher's first snapshot on a call of its own.
2. T-2984, serve takes it before the item snapshot and the first looks, with the reproducing test (after T-2981).
3. T-2988, `design/system/flai-cli.md` (after T-2984); T-2989, close I-0273 (after T-2984). They share no path and run together.
