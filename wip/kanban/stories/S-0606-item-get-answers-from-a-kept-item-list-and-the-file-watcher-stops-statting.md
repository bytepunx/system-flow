---
id: S-0606
type: story
nature: improvement
title: item.get answers from a kept item list and the file watcher stops statting every archived file every 300 ms
status: in-progress
parent: E-0033
owner: alex
created: 2026-11-07T22:41:03Z
updated: 2026-11-09T07:47:11Z
transitions:
  - to: ready
    at: 2026-11-07T23:58:41Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-07T23:59:20Z
    by: agent-S-0606
blocked:
  - from: 2026-11-08T02:13:26Z
    until: 2026-11-08T02:13:45Z
    reason: awaiting the operator on TH-0908
  - from: 2026-11-08T02:14:30Z
    until: 2026-11-08T02:20:37Z
    reason: awaiting the operator on TH-0908
  - from: 2026-11-08T02:21:00Z
    until: 2026-11-08T02:31:46Z
    reason: awaiting the operator on TH-0908
  - from: 2026-11-08T02:31:53Z
    until: 2026-11-08T02:32:04Z
    reason: awaiting the operator on TH-0908
  - from: 2026-11-08T02:32:22Z
    until: 2026-11-08T02:43:59Z
    reason: awaiting the operator on TH-0908
  - from: 2026-11-08T02:44:09Z
    until: 2026-11-08T03:22:18Z
    reason: awaiting the operator on TH-0908
  - from: 2026-11-08T03:22:29Z
    until: 2026-11-08T03:30:55Z
    reason: awaiting the operator on TH-0908
  - from: 2026-11-08T03:31:07Z
    until: 2026-11-09T07:24:42Z
    reason: awaiting the operator on TH-0908
  - from: 2026-11-09T07:25:04Z
    until: 2026-11-09T07:27:39Z
    reason: awaiting the operator on TH-0908
  - from: 2026-11-09T07:27:57Z
    until: 2026-11-09T07:32:35Z
    reason: awaiting the operator on TH-0908
tags: [cli]
touches: [flai/internal/workitem/store.go, flai/internal/hostapi/hostapi.go, flai/internal/watch/watch.go, flai/internal/workitem/repo.go, flai/internal/workitem/store_test.go, flai/internal/watch/watch_test.go, flai/internal/hostapi/hostapi_test.go, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, design/system/server-performance.md, design/system/flai-cli.md, design/tech/go-libraries.md, flai/internal/workitem/archive.go, flai/internal/workitem/itemlock.go, flai/internal/hostapi/writes.go, flai/internal/serve/formula_schedule.go, flai/internal/serve/mergequeue.go, flai/internal/workitem/view.go, flai/internal/workitem/view_test.go, flai/internal/adr/adr.go, flai/internal/adr/adr_test.go, flai/internal/adrindex/adrindex.go, flai/internal/issues/generated.go, flai/internal/issues/issues.go, flai/internal/issues/issues_test.go]
after: [S-0596]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 8972
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-07
      work: 5
    - day: 2026-11-08
      ceremony: 4
      test_runs: 3
      hand_edits: 5
      work: 76
    - day: 2026-11-09
      test_runs: 4
      hand_edits: 3
      work: 17
  models:
    - model: claude-opus-5-5
      input: 560
      output: 160145
      cache_read: 98422655
      cache_write: 4404025
      cost: 49.931
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
      seconds: 405
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 104
          output: 751
          cache_read: 3280958
          cache_write: 473723
          cost: 0.9578
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
  delivery: 2026-11-08T00:55:00Z
  basis: "Its own forecast of 45m; 11th in the pull order with an in-progress limit of 5, behind S-0591, S-0241, S-0598, S-0599, S-0600, S-0601, S-0602, S-0603, S-0604 and S-0605."
  by: flai
  at: 2026-11-07T23:39:00Z
finalized:
  by: orchestrator
  at: 2026-11-07T22:54:20Z
---
# S-0606 item.get answers from a kept item list and the file watcher stops statting every archived file every 300 ms

## Goal

Cause 16 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `item.get` lists all 3,232 items, the archive included, to answer one (`flai/internal/hostapi/hostapi.go`), and the item store holds each folder's lock across its whole stat loop (`flai/internal/workitem/store.go`): 27 to 41 ms warm, 109 ms with 8 at once, 423 ms with 32, p50 612 ms in the serve log. flai serve's watcher stats 6,302 files every 300 ms for this project alone, about 3,800 under `wip/archive` (`flai/internal/watch/watch.go`): a walk takes 104 to 217 ms under load, a third to two thirds of a core, with 12 other projects beside it.

## Acceptance criteria

- [x] flai serve keeps one snapshot of the item list, replaced when the watcher sees a change under the work item folders, and `item.get` and the methods that list every item read it without statting the files again.
- [x] The watcher notices a change under `wip/archive` without statting every file there on every walk, and a test shows an archived item, and a hand edit of one, are still seen.
- [ ] A walk of this project's watched folders takes under 30 ms on the operator's host, `item.get` with 13 asked at once under 50 ms, and server-performance.md records the before and after.

## Tasks
- T-2615 The item store keeps one snapshot of a root's item list, replaced on a change it is told of or on the process's own save, and lists and gets from it without statting
- T-2616 The watcher stats a slow path's folders on every walk and its files on a slower interval, and a test sees an archived item arrive and a hand edit of one
- T-2617 item.get finds its item, children, and plan in one list of the kept snapshot, with no glob and no second list
- T-2618 flai serve turns on each served project's item snapshot, replaces it from the watcher's changes under the item folders, and watches wip/archive on the slow tier
- T-2619 Measure a watcher walk and item.get with 13 at once before and after on the operator's host, and record them with the snapshot and the slow tier in the design
- T-2697 item.get copies only the item and its children it answers, and reads the rest of the list in place
- T-2698 Measure item.get with 13 at once again after T-2697, and record it under cause 16
- T-2699 flai writes issues and ADRs by write and rename, so a change moves its folder's modification time
- T-2700 flai serve watches threads, conversations, digest, issues, and ADRs on the slow tier, and a flai issue and ADR write are seen at the next fast look
- T-2701 Measure the watcher walk again with five more folders on the slow tier, and record them and the reason under cause 16

## Notes

- Proposed fix for the watcher: stat the archive's folders on every walk, which shows an archived item arriving, and its files on a slower interval, such as every 10 s, which shows a hand edit in place. A folder's modification time does not move when a file in it is edited in place, so it cannot stand alone. Or use inotify through fsnotify, which adds a dependency that `design/tech` would record.

### Planning

Planned by planner-S-0606 on 2026-11-07, first at 22:50 (TH-0886, approved and resolved by the orchestrator), then again at 23:46, when the operator asked for every E-0033 story to be planned. The second run found no operator entry on S-0606's threads and nothing in the story's words changed, so it revisited the plan against what changed since: S-0596 was accepted at 23:38, and its section of `design/system/server-performance.md`, cause 16 included, is now on main. All five tasks stand unchanged; none was added, split, merged, or dropped. The second run's thread is TH-0889.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/workitem/store.go` | declared; the store from S-0156, whose `list` holds each folder's lock across its stat loop |
| `flai/internal/hostapi/hostapi.go` | declared; `item.get` calls `repo.Get`, which globs, then `repo.List(true)` again |
| `flai/internal/watch/watch.go` | declared; `snapshot` walks and stats every file under `wip`, the archive included |
| `flai/internal/workitem/repo.go` | layout: `Repo.List`, `Repo.Get`, and `Repo.Save` are where a snapshot is read and where the process's own save replaces it |
| `flai/internal/workitem/store_test.go` | layout: the store's tests, which the snapshot extends |
| `flai/internal/watch/watch_test.go` | layout: criterion 2's test of an archived item and a hand edit of one |
| `flai/internal/hostapi/hostapi_test.go` | co-change: `item.get`'s tests |
| `flai/internal/serve/serve.go` | layout: the watcher's callback and `watchedPaths`, where a project's snapshot is replaced and the archive is watched on the slower tier |
| `flai/internal/serve/serve_test.go` | layout: the wiring's test |
| `design/system/server-performance.md` | design: criterion 3 records the before and after under S-0596's section |
| `design/system/flai-cli.md` | co-change; design: its Internal structure describes the store and the watcher's 300 ms walk |
| `design/tech/go-libraries.md` | design: it says the watcher stats every file every 300 ms and why fsnotify was not taken |

Not taken from `flai touches suggest`, which now ranks documentation first: `docs/users/flai.md` (63%), `docs/users/flai-reference.md` (22%), and `docs/operators/index.md` (26%), since no command, setting, or operator step changes; `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md`, since the dashboard asks the same methods; `flai/internal/hostapi/writes.go` (7%), since a write saves through `Repo.Save`, which replaces the snapshot itself; `flai/cmd/serve_actions.go` and `flaiover/src/lib/server/agent.ts`, which neither list items through the store nor watch. No folder touch.

`after`: S-0596, now done. It wrote the section of `design/system/server-performance.md` this story's goal links, cause 16 and its figures; criterion 3 records the before and after there.

S-0600 (cause 10) also changes `flai/internal/workitem/store.go`, and S-0604 (cause 14) `flai/internal/hostapi/hostapi.go`. Neither writes output this story reads, so neither is in `after`; the trial merge judges the overlap.

Figures:

- Forecast 45m, kept against flai's 24m now (94 s per unit of size over 124 done improvement stories on claude-opus-5-5, times size 15). Comparable stories that kept state between requests and measured it took 703 s (S-0156), 953 s (S-0500), 1170 s (S-0157), and 2191 s (S-0531). This one builds two mechanisms, wires both into flai serve, and measures on the operator's host, about 2700 s.
- Delivery 2026-11-08T00:55:00Z: flai's own delivery for its queue position, 10th behind S-0591 to S-0605 with an in-progress limit of 5, with the 45m duration.
- Cost of delay 30 USD a week, as `flai cod` works it out from the one input, `penalty_per_week: 30`, which the orchestrator set at 22:51 on TH-0886. It stands unadjusted: no revenue rides on the story, and the penalty is the cost of a slow dashboard. It sits on TH-0888's scale for E-0033, between S-0605's and S-0607's 20 and S-0604's 40.
