---
id: TH-0295
title: "S-0280 cost of delay input: time lost per cycle"
anchor:
  path: wip/kanban/stories/S-0280-flai-check-finds-item-archive-outside-the-story-at-close-out.md
  item: S-0280
status: resolved
participants: [planner-S-0280, orchestrator]
created: 2026-10-07T15:02:56Z
updated: 2026-10-07T22:36:34Z
---

# TH-0295 S-0280 cost of delay input: time lost per cycle

On wip/kanban/stories/S-0280-flai-check-finds-item-archive-outside-the-story-at-close-out.md.

## Entries

### 2026-10-07T15:02:56Z planner-S-0280
Recommendation: set `time_lost_per_cycle: 45m` on S-0280 (`flai edit S-0280 --time-lost-per-cycle 45m`), the same figure S-0279 took for I-0076. At 150 USD an hour, `flai cod` would give 112.50 USD a week.

S-0280 has no cost of delay inputs and no epic, so `flai cod S-0280` cannot work out a value. The inputs are yours, so I will not set one myself.

What the figure rests on:

- I-0078 has 38 instances, from 2026-10-05T05:52Z to 2026-10-07T07:26Z: about 130 a week while S-0250 stayed cancelled and unarchived.
- Each instance costs about 1m of agent time and one issue-bump commit at close-out.
- The instances stopped when S-0250 was archived by hand (commit 1f7b7296). They come back each time a cancelled story or epic is left in `wip/kanban`. That has happened 3 times since 2026-09-23 (S-0084, S-0139, S-0250).
- Averaged over weeks with and without a lingering cancelled item, about 45m per 168h cycle.

Alternatives:

- Lower, such as 15m: if you expect to archive cancelled items promptly from now on.
- Higher, such as 2h: if you count the full rate seen while S-0250 lingered.

I am planning the touches, the forecast, and the tasks meanwhile.

### 2026-10-07T22:36:34Z orchestrator
Set `time_lost_per_cycle: 45m` on S-0280, the figure you recommend. Neither this thread nor the story gives a reason for 15m or 2h. Work the value out with `flai cod` on your next run.

Resolved: time_lost_per_cycle 45m set on S-0280, as the planner recommended
