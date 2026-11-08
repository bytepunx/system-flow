---
id: TH-0901
title: "S-0600's plan: six tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0600-a-flai-process-parses-work-items-on-a-pool-of-workers-and-keeps-parses-and-lint.md
  item: S-0600
status: resolved
participants: [planner-S-0600, agent-S-0600, alex]
created: 2026-11-07T23:53:02Z
updated: 2026-11-07T23:57:21Z
---

# TH-0901 S-0600's plan: six tasks in three layers

On wip/kanban/stories/S-0600-a-flai-process-parses-work-items-on-a-pool-of-workers-and-keeps-parses-and-lint.md.

## Entries

### 2026-11-07T23:53:02Z planner-S-0600
S-0600 has six tasks in three layers, 15 touches, and a forecast of 1h30m. Its cost of delay value waits on TH-0888.

| Layer | Task | Touches | After |
|------:|------|---------|-------|
| 1 | T-2674 `store.list` parses on a bounded pool | `flai/internal/workitem/store.go`, `store_test.go` | none |
| 1 | T-2678 The folder lint runs on a bounded pool | `flai/internal/mdlint/folder.go`, `folder_test.go`, `flai/internal/check/check.go` | none |
| 2 | T-2682 Parsed items kept under `.flai-cache` | `flai/internal/workitem/kept.go`, `kept_test.go`, `store.go`, `transition.go` | T-2674 |
| 2 | T-2685 Lint findings kept under `.flai-cache` | `flai/internal/mdlint/config.go`, `kept.go`, `kept_test.go`, `folder.go`, `flai/internal/check/check.go` | T-2678 |
| 3 | T-2686 Design and backup runbook name both cache files | `design/system/flai-cli.md`, `docs/operators/runbooks/backup.md` | T-2682, T-2685 |
| 3 | T-2688 Measure cold `move` and `check` before and after | `design/system/server-performance.md` | T-2682, T-2685 |

The items track and the lint track share no path, so each layer's two tasks can run together.

## Assumptions

- The pool is bounded by `GOMAXPROCS`, capped at 16, with no setting. The story's figures show little gain past 16.
- The lint configuration's hash covers the configuration file's bytes and flai's rule set with a format version. A flai upgrade that changes a rule then misses, not just an edited configuration.
- The kept items restore `Path` and `Archived` from where the file is. They recompute `readAs`, which `itemlock.go` merges on (I-0175), and keep `Unknown` fields.
- In `transition.go`, the second `List(false)` for `WriteIndex` is taken from the first list, but only if the index stays byte for byte the same. The criteria do not ask for this; the declared touch invites it.
- `docs/users/flai.md` is not touched: no command or setting changes.

## Figures

- Forecast: 1h30m, against flai's 29m. S-0531 kept one cache file in 36m31s; this story keeps two, adds two pools, and measures two binaries.
- Cost of delay: `flai cod` refuses for want of inputs. TH-0888 recommends `penalty_per_week: 60` on S-0600. Once it is set, the value is `flai cod S-0600`'s answer.

S-0600 was finalized and moved to ready while I planned. I changed only its touches, its forecast, and its `### Planning` notes, and left its goal and criteria as they were.

### 2026-11-07T23:54:30Z agent-S-0600
Plan taken as drafted: six tasks, three layers, and the assumptions as written. Recorded in the narrative's Decisions. Starting with layer 1, T-2674 and T-2678, which run together.

### 2026-11-07T23:57:21Z alex
Resolved.
