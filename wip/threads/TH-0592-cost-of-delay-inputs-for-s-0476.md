---
id: TH-0592
title: Cost of delay inputs for S-0476
anchor:
  path: wip/kanban/stories/S-0476-the-components-page-cannot-show-a-component-that-asks-the-app-s-server-in-its-answered-state-since-the-harness-answers-every-such-request-with-501.md
  item: S-0476
status: resolved
participants: [planner-S-0476, orchestrator]
created: 2026-10-22T20:57:12Z
updated: 2026-10-22T22:24:39Z
---

# TH-0592 Cost of delay inputs for S-0476

On wip/kanban/stories/S-0476-the-components-page-cannot-show-a-component-that-asks-the-app-s-server-in-its-answered-state-since-the-harness-answers-every-such-request-with-501.md.

## Entries

### 2026-10-22T20:57:12Z planner-S-0476
S-0476 has no cost of delay inputs and no epic to take a share of, so `flai cod S-0476` cannot work out a value. I-0164's cost is "0m", so none was derived from it. The inputs are yours.

Recommended: `time_lost_per_cycle: 20m`. Each time someone previews one of the fourteen components I-0164 names on the Components page (a quarter of flaiover's catalogue: AcceptConfirm, ActivityView, AnalyzerPanel, CancelConfirm, ImportPrompt, InboxView, Messages, NewAdrForm, NewItemForm, OrchestratorPanel, PlannerPanel, PublishBanner, Threads, WorktreeUncommitted), its answered state cannot be seen there, and the agent or operator falls back to the running dashboard against a live project to see it, which costs about that long per occurrence. The other issue stories of 2026-10-21 use the same input (S-0470: 5m, S-0480: 10m, S-0481: 15m).

Alternatives:

- `time_lost_per_cycle: 5m`, if previewing those components in their answered state is rare.
- `penalty_per_week` in the project's currency, if you would rather price the gap in the component library directly.

Set them with `flai edit S-0476 --time-lost-per-cycle 20m`, or `item_edit` with `cost_of_delay` inputs. I plan the touches, forecast, and tasks meanwhile and write the value when the inputs are there.

### 2026-10-22T21:07:26Z orchestrator (recommendation)
Recommended for the operator: take the planner's figure. Run `flai edit S-0476 --time-lost-per-cycle 20m`.

- It prices about 20 minutes a cycle lost to checking those fourteen components against a live dashboard. That is in line with the other issue stories of 2026-10-21: S-0470 at 5m, S-0480 at 10m, S-0481 at 15m.
- Nothing on the thread or the story argues for 5m, or for a penalty instead.

I would set this myself for a backlog story with no input, under plan_backlog_stories. S-0476 is already in ready, so the input is yours.

Source: design/system/strategic-agents.md § Cost of delay

### 2026-10-22T22:24:39Z orchestrator
Resolved: S-0476 is finished and in review, so its cost of delay no longer orders any work. Resolved under plan_backlog_stories so the merge queue can take the story. The input stays yours: `flai edit S-0476 --time-lost-per-cycle 20m` still records it for the stats if you want it.
