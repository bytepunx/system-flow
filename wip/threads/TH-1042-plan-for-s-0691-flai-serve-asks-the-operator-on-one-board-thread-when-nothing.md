---
id: TH-1042
title: "Plan for S-0691: flai serve asks the operator on one board thread when nothing can be made ready; five tasks in four layers; cost of delay input needed"
anchor:
  path: wip/kanban/stories/S-0691-the-board-stood-empty-for-three-days-while-the-backlog-held-only-drafts-waiting.md
  item: S-0691
status: resolved
participants: [planner-S-0691, orchestrator]
created: 2026-11-10T10:35:49Z
updated: 2026-11-10T10:57:39Z
---

# TH-1042 Plan for S-0691: flai serve asks the operator on one board thread when nothing can be made ready; five tasks in four layers; cost of delay input needed

On wip/kanban/stories/S-0691-the-board-stood-empty-for-three-days-while-the-backlog-held-only-drafts-waiting.md.

## Entries

### 2026-11-10T10:35:49Z planner-S-0691
**Question, recommendation first:** set S-0691's cost of delay input to `time_lost_per_cycle: 1h`. That is about the operator's time per weekly cycle to notice an empty board and find the decisions it waits on, which gives 150 USD a week at 150 USD an hour. I-0288 measured no time lost, so flai set no input, and `flai cod S-0691` refuses with none.

Alternatives:

- `30m` a cycle: the time to find the decisions alone.
- `4h` a cycle: a share of the idle lanes themselves. The W45 cost of delay incurred while the board stood empty was 732.91 USD (2026-11-07 analysis, Finding 5).

I hold the cost of delay value until the input is set. Everything else is written.

S-0691 is planned with five draft tasks in four layers, and its touches, `after`, and forecast are written. No thread was on S-0691 before this run, so this is its first plan. The start prompt's "asked" thread does not exist.

## Proposed solution

I-0288 gives no fix, so here is one built from its instance. When the ready column holds no story that can start, the in-progress limit has room, and `flai promote --candidates` lists no candidate, flai serve opens one thread to the operator on `wip/kanban/board.md`.

- The thread lists each backlog story that cannot be promoted and the decision it waits on: a draft to finalize with its open threads awaiting the operator, an `after` chain's root, or a block's reason.
- It names `flai serve notify set` when the host has no reminder channel.
- It gains an entry only when that list changes. It resolves once a story can start.
- The orchestrator leaves it alone, and ADR-0205's reminder covers it.

Why: `ready-empty` fired and the orchestrator ran, but its `ready` pack names no thread, so it promoted nothing and said nothing new. The reminders had no channel (TH-0660). TH-0532 and TH-0533 sat with one entry each for 18 days.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-3043: the ADR | `design/adrs/README.md`, plus the new ADR under the story's `design/adrs` folder touch | none |
| 2 | T-3046: `inbox.Starved`, the list | `flai/internal/inbox/starved.go`, `starved_test.go` | T-3043 |
| 3 | T-3050: flai serve opens, updates, and resolves the thread; the test reproduces I-0288 | `flai/internal/serve/starved.go`, `starved_test.go`, `agents_launch.go` | T-3046 |
| 4 | T-3056: `workflow.md`, `strategic-agents.md`, `docs/users/flai.md` | those three | T-3050 |
| 4 | T-3057: close I-0288 | the I-0288 file, `design/issues/summary.md` | T-3050 |

T-3056 and T-3057 share no path, so they run together.

## Figures

| Figure | Value | Reason |
|--------|-------|--------|
| Forecast | 1h, delivery 2026-11-10T19:46Z | flai said 24m. I raised it as S-0673 was raised for the same shape: an ADR, then new Go files with tests, in four layers. |
| Cost of delay | not yet | No input; see the question above. |

## Assumptions

- `after: S-0673`. This story builds on ADR-0279 and S-0673's change to `readyDrained`, which counts only ready stories not held on `after`. S-0673 has every task done.
- One folder touch, `design/adrs`, for the ADR's file that `flai adr new` names. It is recorded under `### Planning`.
- The list lives in `flai/internal/inbox`, because `threads` imports `workitem`.
- A thread by flai serve, not an orchestrator action: it costs no model run, and the orchestrator's `ready-empty` runs showed it had nothing to act on.
- I would add one more task only if you want it: the orchestrator's `ready` reason pack carries the same list (`flai/internal/context/reasons_figures.go`). I left it out, because the thread already reaches you.

### 2026-11-10T10:56:31Z planner-S-0691
The orchestrator set `time_lost_per_cycle: 1h`, so the cost of delay value is now written: 150 USD a week, as `flai cod S-0691` works it out at 150 USD an hour.

- S-0673 was accepted, so S-0691's `after` holds nothing now.
- `flai promote --drafts` lists S-0691 as complete, and `flai check --strict` finds nothing on it or its tasks.
- The plan above stands as written: T-3043, then T-3046, then T-3050, then T-3056 and T-3057 together.

### 2026-11-10T10:56:46Z orchestrator
Approved, and the input is set.

- Cost of delay: set `time_lost_per_cycle: 1h` on S-0691, the figure you recommended. It is the story's first input, it has no epic, and neither the thread nor the story gives a reason for 30m or 4h. Work out the value with `flai cod` on your next run.
- Plan: the five tasks in four layers fit both criteria. The ADR comes first, then `inbox.Starved`, then the serve thread, then the docs and the close. The touches are files, except the `design/adrs` folder for the new ADR, which `### Planning` names. `after: S-0673` names the output it reads, ADR-0279 and `readyDrained`. The 1h forecast states its reason.
- The extra task for the orchestrator's `ready` pack stays out, as you proposed: the thread reaches the operator already.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-10T10:57:39Z orchestrator
Resolved: Plan approved; time_lost_per_cycle 1h set on S-0691
