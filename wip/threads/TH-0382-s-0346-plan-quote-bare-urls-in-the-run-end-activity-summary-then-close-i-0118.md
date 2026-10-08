---
id: TH-0382
title: "S-0346 plan: quote bare URLs in the run-end activity summary, then close I-0118"
anchor:
  path: wip/kanban/stories/S-0346-flai-check-finds-markdown-md034-outside-the-story-at-close-out.md
  item: S-0346
status: resolved
participants: [planner-S-0346, orchestrator, alex]
created: 2026-10-08T08:54:49Z
updated: 2026-10-08T09:12:34Z
---

# TH-0382 S-0346 plan: quote bare URLs in the run-end activity summary, then close I-0118

On wip/kanban/stories/S-0346-flai-check-finds-markdown-md034-outside-the-story-at-close-out.md.

## Entries

### 2026-10-08T08:54:49Z planner-S-0346
S-0346 is planned as three tasks in three layers.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1419 mdlint quotes the bare URLs md034 finds in a line | `flai/internal/mdlint/inline.go`, `quote.go`, `quote_test.go` | none |
| 2 | T-1420 A strategic run's end entry quotes the bare URLs in the summary flai takes from its final text | `flai/internal/serve/activity.go`, `activity_test.go`, `design/system/agent-narrative.md` | T-1419 |
| 3 | T-1421 Close I-0118 saying what fixed it | `design/issues/I-0118-…md`, `design/issues/summary.md` | T-1420 |

T-1420 calls T-1419's `QuoteBareURLs`. T-1421's close reason names T-1420's fix and test.

## Proposed cause

- All five instances are bare `www.` literals written between 04:23Z and 05:31Z today. Four are line 1529 of `wip/agents/orchestrator.md`; one is TH-0361's title.
- The installed flai that wrote them (1.38.1, then 1.39.3) predates S-0324's bare `www.` rule, published in 1.39.4. The host now runs 1.39.8, so each agent write to `wip/` is refused on one. `flai check --strict` on main finds no MD034 now.
- One path is left: the entry `flai serve` logs when a strategic run ends takes its summary from the run's final text. A bare URL there is refused, and the run's seconds and cost go unlogged with only a warning. T-1420 quotes it instead.

## Figures

- Forecast 40m, up from flai's 15m: the helper must match md034's three kinds of bare URL, and the close-out runs the integration tier. Delivery 2026-10-08T14:13:00Z.
- Cost of delay: none yet. TH-0381 asks for `time_lost_per_cycle`, recommending 5m.

## Assumptions

1. An agent's own `activity_log` summary is still refused on a finding, as `tooling.md` says. Only text flai takes from a run is quoted.
2. `threads.MirrorNarrative` writes without `LintGuard`, but the titles it copies are linted in the thread file. ADR-0123 leaves another story's narrative out of a close-out, so it is out of scope.
3. `design/issues` is in `claims.shared`, so the issue touches hold no other story. No folder touch is kept.

### 2026-10-08T08:55:03Z orchestrator
Plan approved: three tasks in three layers, all file touches, a 40m forecast. The cause matches what happened. Every instance was my own run-era bare `www.`, written while the host ran a flai older than S-0324's rule. The one remaining unchecked path, the run-end summary flai takes from a strategic run's final text, is the right thing to fix. Assumption 1 is sound: an agent's own `activity_log` summary stays refused on a finding. The cost of delay input is set on TH-0381 (5m).

Source: wip/kanban/stories/S-0346-flai-check-finds-markdown-md034-outside-the-story-at-close-out.md

### 2026-10-08T08:55:04Z orchestrator
Resolved: Plan approved; cost input set on TH-0381.

### 2026-10-08T08:55:48Z planner-S-0346
Cost of delay set: 12.50 USD a week, from `time_lost_per_cycle: 5m`, the input the orchestrator set on TH-0381. The plan above is otherwise unchanged.

### 2026-10-08T09:12:34Z alex
Resolved.
