---
id: S-0606
type: story
nature: improvement
title: item.get answers from a kept item list and the file watcher stops statting every archived file every 300 ms
status: backlog
parent: E-0033
owner: alex
created: 2026-11-07T22:41:03Z
updated: 2026-11-07T22:54:20Z
transitions: []
tags: [cli]
touches: [flai/internal/workitem/store.go, flai/internal/hostapi/hostapi.go, flai/internal/watch/watch.go, flai/internal/workitem/repo.go, flai/internal/workitem/store_test.go, flai/internal/watch/watch_test.go, flai/internal/hostapi/hostapi_test.go, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, design/system/server-performance.md, design/system/flai-cli.md, design/tech/go-libraries.md]
after: [S-0596]
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
    - kind: planner
      seconds: 581
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 66
          output: 37240
          cache_read: 2919837
          cache_write: 141523
          cost: 0.2272
        - model: claude-opus-5-5
          input: 84
          output: 29871
          cache_read: 9235195
          cache_write: 266116
          cost: 4.5737
    - kind: orchestrator
      seconds: 160
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 50
          output: 364
          cache_read: 1564195
          cache_write: 252630
          cost: 0.4618
cost_of_delay:
  inputs:
    penalty_per_week: 30
    by: orchestrator
    at: 2026-11-07T22:51:02Z
  value: 30
  by: planner-S-0606
  at: 2026-11-07T22:52:00Z
forecast:
  duration: 45m
  delivery: 2026-11-08T00:15:00Z
  basis: "Comparable kept-state stories took 703 s (S-0156), 953 s (S-0500), 1170 s (S-0157), and 2191 s (S-0531); this one builds two mechanisms, a snapshot and a slower archive tier, wires both into flai serve, and measures on the operator's host, so about 2700 s."
  by: planner-S-0606
  at: 2026-11-07T22:48:46Z
finalized:
  by: orchestrator
  at: 2026-11-07T22:54:20Z
---
# S-0606 item.get answers from a kept item list and the file watcher stops statting every archived file every 300 ms

## Goal

Cause 16 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `item.get` lists all 3,232 items, the archive included, to answer one (`flai/internal/hostapi/hostapi.go`), and the item store holds each folder's lock across its whole stat loop (`flai/internal/workitem/store.go`): 27 to 41 ms warm, 109 ms with 8 at once, 423 ms with 32, p50 612 ms in the serve log. flai serve's watcher stats 6,302 files every 300 ms for this project alone, about 3,800 under `wip/archive` (`flai/internal/watch/watch.go`): a walk takes 104 to 217 ms under load, a third to two thirds of a core, with 12 other projects beside it.

## Acceptance criteria

- [ ] flai serve keeps one snapshot of the item list, replaced when the watcher sees a change under the work item folders, and `item.get` and the methods that list every item read it without statting the files again.
- [ ] The watcher notices a change under `wip/archive` without statting every file there on every walk, and a test shows an archived item, and a hand edit of one, are still seen.
- [ ] A walk of this project's watched folders takes under 30 ms on the operator's host, `item.get` with 13 asked at once under 50 ms, and server-performance.md records the before and after.

## Tasks
- T-2615 The item store keeps one snapshot of a root's item list, replaced on a change it is told of or on the process's own save, and lists and gets from it without statting
- T-2616 The watcher stats a slow path's folders on every walk and its files on a slower interval, and a test sees an archived item arrive and a hand edit of one
- T-2617 item.get finds its item, children, and plan in one list of the kept snapshot, with no glob and no second list
- T-2618 flai serve turns on each served project's item snapshot, replaces it from the watcher's changes under the item folders, and watches wip/archive on the slow tier
- T-2619 Measure a watcher walk and item.get with 13 at once before and after on the operator's host, and record them with the snapshot and the slow tier in the design

## Notes

- Proposed fix for the watcher: stat the archive's folders on every walk, which shows an archived item arriving, and its files on a slower interval, such as every 10 s, which shows a hand edit in place. A folder's modification time does not move when a file in it is edited in place, so it cannot stand alone. Or use inotify through fsnotify, which adds a dependency that `design/tech` would record.

### Planning

Planned by planner-S-0606. No earlier plan and no thread on the story existed. The plan's thread is TH-0886, which the orchestrator approved and resolved.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/workitem/store.go` | declared; the store from S-0156, whose `list` holds each folder's lock across its stat loop |
| `flai/internal/hostapi/hostapi.go` | declared; `item.get` calls `repo.Get`, which globs, then `repo.List(true)` again |
| `flai/internal/watch/watch.go` | declared; `snapshot` walks and stats every file under `wip`, the archive included |
| `flai/internal/workitem/repo.go` | layout: `Repo.List`, `Repo.Get`, and `Repo.Save` are where a snapshot is read and where the process's own save replaces it |
| `flai/internal/workitem/store_test.go` | layout: the store's tests, which the snapshot extends |
| `flai/internal/watch/watch_test.go` | layout: criterion 2's test of an archived item and a hand edit of one |
| `flai/internal/hostapi/hostapi_test.go` | co-change, 19%: `item.get`'s tests |
| `flai/internal/serve/serve.go` | layout: the watcher's callback and `watchedPaths`, where a project's snapshot is replaced and the archive is watched on the slower tier |
| `flai/internal/serve/serve_test.go` | layout: the wiring's test |
| `design/system/server-performance.md` | design: criterion 3 records the before and after under S-0596's section |
| `design/system/flai-cli.md` | co-change, 14%; design: its Internal structure describes the store and the watcher's 300 ms walk |
| `design/tech/go-libraries.md` | design: it says the watcher stats every file every 300 ms and why fsnotify was not taken |

Not taken from `flai touches suggest`: `flai/internal/hostapi/writes.go` (38%), since a write saves through `Repo.Save`, which replaces the snapshot itself; `flai/cmd/serve_actions.go`, `flaiover/src/lib/server/agent.ts`, and `flai/internal/mcpserver/server.go`, which neither list items through the store nor watch; and `docs/users/flai.md` and `docs/users/flai-reference.md`, since no command or setting changes. No folder touch.

`after`: S-0596. It writes the section of `design/system/server-performance.md` this story's goal links, cause 16 and its figures, on its branch; criterion 3 records the before and after there.

Figures:

- Forecast 45m, adjusted from flai's 8m (80 s per unit of size over 19 done improvement stories, times size 6). Comparable stories that kept state between requests and measured it took 703 s (S-0156), 953 s (S-0500), 1170 s (S-0157), and 2191 s (S-0531). This one builds two mechanisms, wires both into flai serve, and measures on the operator's host, about 2700 s.
- Delivery 2026-11-08T00:15:00Z: flai's delivery for its queue position, 23:36, moved by the 37 minutes added to the duration.
- Cost of delay 30 USD a week, as `flai cod` works it out from the one input, `penalty_per_week: 30`, which the orchestrator set at 22:51 on TH-0886. It stands unadjusted: no revenue rides on the story, and the penalty is the cost of a slow dashboard.
