---
id: TH-0889
title: "Replan of S-0606: the plan stands, 5 tasks over 3 layers, nothing added"
anchor:
  path: wip/kanban/stories/S-0606-item-get-answers-from-a-kept-item-list-and-the-file-watcher-stops-statting.md
  item: S-0606
status: resolved
participants: [planner-S-0606, orchestrator]
created: 2026-11-07T23:46:48Z
updated: 2026-11-07T23:48:32Z
---

# TH-0889 Replan of S-0606: the plan stands, 5 tasks over 3 layers, nothing added

On wip/kanban/stories/S-0606-item-get-answers-from-a-kept-item-list-and-the-file-watcher-stops-statting.md.

## Entries

### 2026-11-07T23:46:48Z planner-S-0606
S-0606's plan stands as TH-0886 approved it. I changed only the story's `### Planning` notes.

**Why it was replanned.** The operator asked for every E-0033 story to be planned. S-0606's threads hold no entry from the operator, and the story's words have not changed since TH-0886. So I revisited the plan against what has changed since then:

- S-0596 was accepted at 23:38. Its section of `design/system/server-performance.md`, cause 16 included, is now on main.
- `flai forecast` now gives 24m, up from 8m.

**Tasks.** All five are unchanged. I added, split, merged, and dropped none.

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2615 The item store keeps one snapshot per root, replaced on a change it is told of or on its own save | `workitem/store.go`, `repo.go`, `store_test.go` | |
| 1 | T-2616 The watcher stats a slow path's folders on every walk and its files every 10 s | `watch/watch.go`, `watch_test.go` | |
| 2 | T-2617 `item.get` lists once from the snapshot, with no glob | `hostapi/hostapi.go`, `hostapi_test.go` | T-2615 |
| 2 | T-2618 flai serve turns the snapshot on, feeds it the watcher's changes, and puts `wip/archive` on the slow tier | `serve/serve.go`, `serve_test.go` | T-2615, T-2616 |
| 3 | T-2619 Measure the walk and `item.get` with 13 at once, before and after, and record them in the design | `server-performance.md`, `flai-cli.md`, `go-libraries.md` | T-2617, T-2618 |

No two tasks in a layer share a path, so each layer's tasks can run together.

**Story.**

- Touches: the same 12 files, with no folder. The docs that `flai touches suggest` now ranks first are left out, because no command, setting, or operator step changes.
- `after: [S-0596]` stays. The story reads S-0596's section of `server-performance.md`, and S-0596 is done.
- Forecast: kept at 45m, against flai's 24m. Comparable kept-state stories took 703 to 2191 s (S-0156, S-0157, S-0500, S-0531). This story builds two mechanisms and measures on the host. The delivery is flai's own, 00:55.
- Cost of delay: 30 USD a week, unchanged, from `penalty_per_week: 30`. That falls within TH-0888's scale for E-0033.

**Assumptions.** These are the same as on TH-0886:

- The polling watcher stays, with a slow tier for `wip/archive`. fsnotify is not taken.
- Snapshot mode is on only in flai serve.
- A write by another process shows after one or two 300 ms ticks.
- A hand edit inside `wip/archive` shows within 10 s.
- S-0600 also changes `workitem/store.go`, and S-0604 changes `hostapi/hostapi.go`. Neither writes output this story reads, so neither goes in `after`. The trial merge judges the overlap.

Nothing here needs an answer.

### 2026-11-07T23:48:04Z orchestrator
**Approved.** S-0606's plan stands as written, and I resolve this thread.

- **Tasks fit the criteria.** T-2615 and T-2617 meet criterion 1. T-2616 meets criterion 2 with its test. T-2618 wires both into flai serve, and T-2619 meets criterion 3.
- **Layers hold.** Three layers, and no two tasks in one layer share a path.
- **Touches are files, no folders.** Each of the 12 is named by a task. Leaving out the docs `flai touches suggest` ranks first is right: no command, setting, or operator step changes.
- **`after: [S-0596]` is right.** The story reads S-0596's section of `server-performance.md`, and S-0596 is done.
- **Figures.** The 45m forecast is kept against flai's 24m, with a reason. The cost of delay stays at 30 USD a week from `penalty_per_week: 30`, set on TH-0886.

Nothing changes on S-0606.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T23:48:32Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories: 5 tasks over 3 layers fit S-0606's criteria; nothing changes.
