---
id: TH-0891
title: "S-0599 plan: five tasks in four layers, and its cost of delay input"
anchor:
  path: wip/kanban/stories/S-0599-flai-serve-s-replanner-writes-every-forecast-it-moves-in-one-pass-with-one.md
  item: S-0599
status: answered
participants: [planner-S-0599, orchestrator, agent-S-0599]
created: 2026-11-07T23:49:56Z
updated: 2026-11-08T00:25:08Z
---

# TH-0891 S-0599 plan: five tasks in four layers, and its cost of delay input

On wip/kanban/stories/S-0599-flai-serve-s-replanner-writes-every-forecast-it-moves-in-one-pass-with-one.md.

## Entries

### 2026-11-07T23:49:56Z planner-S-0599
**Question, recommended answer first:** set `penalty_per_week: 60` on S-0599, as on TH-0888. If you count the 80 s wait after each acceptance as yours, 100 fits better, as the orchestrator says there. An answer on either thread settles it. Until then S-0599 has no cost of delay value: neither it nor E-0033 has inputs, so `flai cod` refuses.

**Plan.** Five tasks in four layers:

| Layer | Task | Touches | After |
|------:|------|---------|-------|
| 1 | T-2628 Save many work items under one hold of the items lock | `flai/internal/workitem/itemlock.go`, `itemlock_test.go` | none |
| 1 | T-2632 Check the work-item rules and the markdown lint over a given set of items | `flai/internal/check/check.go`, `items_test.go` | none |
| 2 | T-2641 Apply many forecast changes in one pass with one check and one index write | `flai/internal/itemedit/forecasts.go`, `forecasts_test.go` | T-2628, T-2632 |
| 3 | T-2647 The replanner writes its moved forecasts through the batch entry | `flai/internal/serve/replan.go`, `replan_test.go` | T-2641 |
| 4 | T-2655 Measure a 20-forecast replan in a scratch copy and record cause 9's before and after | `design/system/server-performance.md`, `design/system/strategic-agents.md` | T-2647 |

T-2628 and T-2632 share no path, so they run together.

**Figures.**

- Touches: 10 files, no folder. Where each came from is under `### Planning` in the story's Notes.
- Forecast: 2h30m, not `flai forecast`'s 6m, which counts only 3 criteria and 1 touch. Delivery 2026-11-08T02:30Z.
- `after`: none. S-0598 (cause 8) scopes the checked edit, which this plan does not call.

**Assumptions.**

1. Criterion 1's "one check over the written items" is the per-item rules (`oneItem`) and the markdown lint of the written files, not a whole `check.Run`. The story's Notes propose this: forecasts are values flai owns. A whole `Run` before and after would cost about 5 to 6 s alone and miss criterion 2's 5 s.
2. Under ADR-0184, a conflict on one item refuses that item alone. The other moved forecasts are still written, and the refused one is warned of, as a failed `Apply` is today.
3. The batch still records one edit notice per item, so `flai serve`'s news and `replan-agent` see the same edits as before.
4. `design/system/strategic-agents.md` § The replan policy changes too. It describes the per-item `Apply` path and its cost.

Nothing to split, merge, or drop.

### 2026-11-07T23:50:06Z planner-S-0599
S-0599 has been in progress since 23:49:39. agent-S-0599 started it once you finalized it, and the five tasks are there for it to review. The cost of delay value is still unset. The question above still stands, here or on TH-0888: once you set `penalty_per_week`, `flai cod S-0599` can work out the value. I am ending this run, because a story in progress is no longer mine to plan.

### 2026-11-07T23:53:44Z orchestrator (recommendation)
**Recommended:** set `penalty_per_week: 60` on S-0599. That is the planner's figure and TH-0888's table.

**Why:**

- 60 puts S-0599 beside S-0600 (60), which you also confirm on TH-0888. It is below S-0603 (100), which I set at that thread's figure.
- 100 fits only if you count the 80 s wait after each acceptance as a cost to you, as I noted on TH-0888. Nothing on the story or this thread decides that, so I take the planner's figure.
- The `throughput` policy orders ready by forecast, so neither figure changes an order today.

**Why I recommend rather than set it:** S-0599 is in progress, not in the backlog. A backlog story's first input is the only one `plan_backlog_stories` lets me set. Once you confirm this, or answer TH-0888, I set the input under `set_cost_of_delay`. I set no value.

**The plan:** its five tasks in four layers fit S-0599's criteria. agent-S-0599 is reviewing them now.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-08T00:05:48Z agent-S-0599
Plan reviewed and kept: five tasks in four layers. One change: `check.Items` goes in a new `flai/internal/check/items.go`, not `check.go`, so this story stays off S-0598's changes there. T-2628 and T-2632 are done.

The cost of delay input is the operator's call. I am leaving it to your recommendation and TH-0888, and the work does not wait on it.

### 2026-11-08T00:25:08Z orchestrator (recommendation)
**Recommended:** set `penalty_per_week: 60` on S-0599. That is the planner's figure here and in TH-0888's table.

The plan is settled. agent-S-0599 reviewed and kept it, and moved `check.Items` into `flai/internal/check/items.go`. Three of its five tasks are done. Only the cost of delay input is still open.

**Why 60:** it puts S-0599 beside S-0600 (60) and below S-0603 (100). 100 fits only if you count the 80 s replan after each acceptance as a cost to you, and nothing on the story or this thread says so. Under the `throughput` policy neither figure changes any order.

**Why I do not set it:** S-0599 is in progress, not in the backlog. A backlog story's first input is the only one `plan_backlog_stories` lets me set. Once you confirm this here or answer TH-0888, I set the input under `set_cost_of_delay`. I set no value.

Source: design/conventions/strategic-agents.md § As the orchestrator
