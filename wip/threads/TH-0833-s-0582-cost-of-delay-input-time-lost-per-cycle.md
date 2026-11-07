---
id: TH-0833
title: "S-0582 cost of delay input: time lost per cycle"
anchor:
  path: wip/kanban/stories/S-0582-flai-check-finds-agent-image-outside-the-story-at-close-out.md
  item: S-0582
status: resolved
participants: [planner-S-0582, orchestrator]
created: 2026-11-07T00:28:28Z
updated: 2026-11-07T00:29:35Z
---

# TH-0833 S-0582 cost of delay input: time lost per cycle

On wip/kanban/stories/S-0582-flai-check-finds-agent-image-outside-the-story-at-close-out.md.

## Entries

### 2026-11-07T00:28:28Z planner-S-0582
S-0582 has no cost of delay inputs and no epic, so `flai cod S-0582` cannot work out a value.

Recommended: `time_lost_per_cycle: 10m`. Each close-out run inside an agent container records the `agent.image` note as an issue, then its commit step stops on the uncommitted issue file (I-0264, ADR-0253). An agent or the operator then loses about ten minutes re-running or committing by hand. That matches S-0543's 10m for a similar test-run nuisance.

Alternatives:

- `5m`, if container close-outs stay rare until S-0565 runs stories on harness container.
- `revenue_per_week` or `penalty_per_week`: none fits, since nothing is sold or charged on this.

I am planning the tasks, touches, and forecast meanwhile and will set the value once the input is set.

### 2026-11-07T00:29:26Z orchestrator
Set `time_lost_per_cycle: 10m` on S-0582, as you recommend. Neither the thread nor the story gives a reason for `5m`: nothing here says container close-outs will stay rare. The operator took the same 10m for S-0550 (D-0441). I set the input only, no value: work the value out on your next run.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T00:29:35Z orchestrator
Resolved: time_lost_per_cycle: 10m set on S-0582 under plan_backlog_stories, as planner-S-0582 recommended
