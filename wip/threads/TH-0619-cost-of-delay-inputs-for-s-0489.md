---
id: TH-0619
title: Cost of delay inputs for S-0489
anchor:
  path: wip/kanban/stories/S-0489-host-page-under-configuration.md
  item: S-0489
status: resolved
participants: [planner-S-0489, orchestrator]
created: 2026-10-23T18:05:56Z
updated: 2026-10-23T19:30:06Z
---

# TH-0619 Cost of delay inputs for S-0489

On wip/kanban/stories/S-0489-host-page-under-configuration.md.

## Entries

### 2026-10-23T18:05:56Z planner-S-0489
S-0489 has no cost of delay inputs and no epic to take a share of, so `flai cod S-0489` refuses a value. The inputs are yours.

Recommended: `time_lost_per_cycle: 1h`, which at the project's `planning.hour_rate` of 150 USD and a 168h cycle gives 150 USD a week: about an hour a week spent reaching the dashboard from the host alone, or working around the tunnel's refusal, until a custom domain is accepted.

Alternatives: `penalty_per_week: 100` if you would rather state the cost directly; `time_lost_per_cycle: 30m` (75 USD a week) if the tunnel is rarely needed; or an epic for the story, whose value it would take a share of.

Set it with `flai edit S-0489 --time-lost-per-cycle 1h`, or `item_edit` `cost_of_delay` on the story, inputs only; I work out the value from them.

### 2026-10-23T18:06:07Z orchestrator (recommendation)
Recommended for the operator: take the planner's figure, `time_lost_per_cycle: 1h`. That is 150 USD a week at `planning.hour_rate`. Set it with `flai edit S-0489 --time-lost-per-cycle 1h`.

Why I have not set it: S-0489 is in ready, not the backlog. Under `plan_backlog_stories` I give inputs only to a backlog story that has none, so a ready story's inputs are yours.

Nothing on S-0489 points to one of the alternatives (30m, or a 100 USD penalty), so the recommended figure stands. Once you set it, the planner works out the value.

Source: design/conventions/strategic-agents.md § Asking the operator

### 2026-10-23T19:30:06Z orchestrator
Resolved: Under plan_backlog_stories: S-0489 is built and in review, so its cost of delay input no longer orders or holds it. The input stays the operator's to set with flai edit S-0489 --time-lost-per-cycle 1h if wanted for the release figures.
