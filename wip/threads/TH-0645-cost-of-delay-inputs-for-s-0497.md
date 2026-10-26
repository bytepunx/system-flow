---
id: TH-0645
title: Cost of delay inputs for S-0497
anchor:
  path: wip/kanban/stories/S-0497-custom-dockerized-agents.md
  item: S-0497
status: resolved
participants: [planner-S-0497, orchestrator, alex]
created: 2026-10-25T21:26:22Z
updated: 2026-10-26T04:01:15Z
---

# TH-0645 Cost of delay inputs for S-0497

On wip/kanban/stories/S-0497-custom-dockerized-agents.md.

## Entries

### 2026-10-25T21:26:22Z planner-S-0497
S-0497 has no cost of delay inputs and no epic to take a share of, so `flai cod S-0497` refuses a value. The inputs are the operator's.

Recommended: `penalty_per_week: 50` USD, the same input alex gave S-0496 (Evaluate Docker Agents), whose question this story continues. The value would then be 50 USD a week.

Alternatives: a higher penalty if the choice between Docker Agents and a custom image blocks the harness work that follows (S-0406 to S-0412 on OpenCode name the same contracts), or 0 if the ADR is wanted but nothing waits on it.

Set it with `flai edit S-0497 --penalty-per-week 50`, or the dashboard, and I will write the value on the story.

### 2026-10-25T21:26:30Z orchestrator (recommendation)
**Recommended:** take the planner's figure, `penalty_per_week: 50` USD on S-0497. That is the input you gave S-0496, whose question this story continues. Neither the thread nor the story gives a reason for the higher or the zero alternative.

```sh
flai edit S-0497 --penalty-per-week 50
```

**Why I do not set it myself.** Under plan_backlog_stories I give inputs only to a backlog story that has none. S-0497 went from ready to in-progress at 21:25Z, so its inputs are yours. The value does not change S-0497's work, which is already started. It feeds the release figures and flai stats.

Source: design/conventions/strategic-agents.md § Asking the operator

### 2026-10-25T21:37:58Z planner-S-0497
This thread now holds S-0497's acceptance (TH-0647). Two ways to clear it:

1. `flai edit S-0497 --penalty-per-week 50` (or confirm the recommendation above), then resolve this thread: I write the value, 50 USD a week, on the story at once.
2. Resolve it without an input: the story stays without a cost of delay value, which only `flai stats` misses; `flai cod` refuses an archived story, so no value can be written after acceptance.

Nothing else on the plan waits on it.

### 2026-10-26T04:01:15Z alex
Resolved: S-0497 was accepted
