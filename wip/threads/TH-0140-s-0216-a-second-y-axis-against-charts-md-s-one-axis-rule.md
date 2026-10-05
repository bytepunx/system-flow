---
id: TH-0140
title: "S-0216: a second y-axis against charts.md's one-axis rule?"
anchor:
  path: wip/kanban/stories/S-0216-strategic-cost-and-strategic-use-charts-show-what-the-planner-orchestrator-and-analyzer-add-against-delivery.md
  item: S-0216
status: open
participants: [planner-S-0216]
created: 2026-10-05T05:45:10Z
updated: 2026-10-05T05:45:10Z
---

# TH-0140 S-0216: a second y-axis against charts.md's one-axis rule?

On wip/kanban/stories/S-0216-strategic-cost-and-strategic-use-charts-show-what-the-planner-orchestrator-and-analyzer-add-against-delivery.md.

## Entries

### 2026-10-05T05:45:10Z planner-S-0216
The criteria ask for the mean cost per story "as a line on a second axis". `design/tech/charts.md` makes "one y-axis per chart" a rule in every builder, and `flaiover/src/lib/viz/charts.test.ts` (line 294, "must not use two y-axes") enforces it.

**Recommended:** keep one y-axis. Both charts plot one unit against itself: Strategic Cost has dollars in its bars and its line, and Strategic Use has hours in its bars and its lines. So one shared axis reads true, and the rule and the test stand. The criterion would then read "as a line on the same axis".

Alternative: allow a second y-axis for the two Strategic kinds only. That needs an ADR, an exception in `charts.md`, and an exemption in the test. It adds a task and about 15m.

Assumptions I am planning on unless you say otherwise:

1. Buckets are the days of the window, as `strategic_days[]` gives them; weeks would need flai to lay them out.
2. The mean waiting per story completed each day is computed in the dashboard from `items[].wait_threads_seconds` and `wait_review_seconds` of the stories completed that day, so no metrics contract change or ADR is needed.
3. "The ratio stated" is the strategic agents' cost over the window, per story completed, as a share of the mean agent cost per story.
