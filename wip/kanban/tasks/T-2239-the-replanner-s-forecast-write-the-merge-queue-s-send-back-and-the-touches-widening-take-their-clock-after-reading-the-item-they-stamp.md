---
id: T-2239
type: task
nature: remediation
title: The replanner's forecast write, the merge queue's send-back, and the touches widening take their clock after reading the item they stamp
status: done
parent: S-0522
owner: alex
created: 2026-10-28T14:21:56Z
updated: 2026-10-28T14:38:39Z
transitions:
  - to: ready
    at: 2026-10-28T14:25:08Z
    by: agent-S-0522
  - to: in-progress
    at: 2026-10-28T14:25:09Z
    by: agent-S-0522
  - to: done
    at: 2026-10-28T14:38:39Z
    by: agent-S-0522
stream: S-0522
tags: [cli]
touches: [flai/internal/serve/replan.go, flai/internal/serve/replan_test.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go, flai/internal/itemedit/widen.go, flai/internal/itemedit/widen_test.go]
usage:
  source: log
  seconds: 810
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 66
      output: 26671
      cache_read: 2844990
      cache_write: 151525
      cost: 2.0878
---
# T-2239 The replanner's forecast write, the merge queue's send-back, and the touches widening take their clock after reading the item they stamp

## Work

Three writers stamp `updated` from a clock read before they read the item, so a move made between the two leaves the stamp behind the transition. Change each to take its clock after the read, and add a test for each that moves the item between the writer's start and its write, with a clock the test controls, and checks that the file's `updated` is at or after the transition's `at`.

- `flai/internal/serve/replan.go`: `look` takes `now` once, before `repo.List` and before `replan` writes each moved forecast through `itemedit.Apply` with `Now: now`; the loop runs a `flai check` before and after each edit, so it spans minutes over many stories. Keep that `now` for `news`, the schedule, and the replay, and pass the clock read at each write, `r.o.Now().UTC()`, as the `Now` of each `Apply`. Test in `flai/internal/serve/replan_test.go`, beside the forecast-write tests: a story moved after the replanner's look began gets a forecast stamped no earlier than its move.
- `flai/internal/mergequeue/bisect.go`: the send-back of the story a red batch is bisected to takes `now` before `repo.Get(story)`. Read the item first, then the clock, for the thread it opens and the transition it records. Test in `flai/internal/mergequeue/bisect_test.go`.
- `flai/internal/itemedit/widen.go`: `WidenAll` evaluates `o.Now()` as an argument before `Widen` reads the item with `repo.Get`. Give `Widen` the clock to read after the `Get`, a `func() time.Time` in place of the `time.Time`, or read the item in the caller first; keep the notice's `At` the stamp written. Test in `flai/internal/itemedit/widen_test.go`.

Leave `cmd/edit.go` and `cmd/criteria.go` alone: `Apply` reads the item at once after they take their clock, and the task on `Repo.Save` keeps the order whatever the gap. Say in the narrative's `## Decisions` which of the three the story's agent finds wrote S-0482's stamp, when it can tell from the instance, and that the invariant in `Repo.Save` covers the rest.

This task waits for nothing. It shares no path with the task on `Repo.Save`, so the two run together.

## Done when

Each of the three writers reads its clock after it reads the item, the three tests exist, failed before the change and pass after it, `go test ./flai/internal/serve/... ./flai/internal/mergequeue/... ./flai/internal/itemedit/...` and `flai test` on the files changed pass, and `flai check --strict` is clean on them.

## Notes

The replanner's writes are the widest window: `inbox` on 2026-10-28 showed it stamping some thirty stories' forecasts over two minutes from one `now`. S-0482's own forecast kept its planner's stamp, so the replanner is not shown to be the writer of that instance; the main checkout's writes between S-0482's move at 12:14:19Z and the check at 12:18:49Z were not committed, so git does not say which it was.
