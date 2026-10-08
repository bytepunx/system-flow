---
id: TH-0343
title: S-0318 cost of delay input
anchor:
  path: wip/kanban/stories/S-0318-flai-check-finds-markdown-md038-outside-the-story-at-close-out.md
  item: S-0318
status: resolved
participants: [planner-S-0318, orchestrator]
created: 2026-10-08T00:06:16Z
updated: 2026-10-08T00:06:25Z
---

# TH-0343 S-0318 cost of delay input

On wip/kanban/stories/S-0318-flai-check-finds-markdown-md038-outside-the-story-at-close-out.md.

## Entries

### 2026-10-08T00:06:16Z planner-S-0318
S-0318 has no cost of delay inputs and no epic, so `flai cod S-0318` cannot work out a value. I need one input from you.

Recommended: `time_lost_per_cycle: 15m`, which makes the value 37.5 at the rate S-0279 and S-0280 imply (45m gave 112.5).

Why 15m: I-0096 has one instance (count 1). It cost one close-out note and an issue recorded against S-0227, and S-0229's agent spent about 15 minutes fixing its own narrative (commit c9544681). The siblings S-0279 (wip.overlap, 18 instances) and S-0280 (item.archive, 38 instances) were each given 45m. I-0096 recurs far less, so a third of that fits.

Alternatives:

- `30m`: if you expect hand-edited narrative sections to keep breaking the lint as more stories run at once.
- `5m`: if one instance in a week is all you expect.

While this is open, I am planning the touches, forecast, and tasks.

### 2026-10-08T00:06:24Z orchestrator
Set `time_lost_per_cycle: 15m` on S-0318, the figure you recommend. Neither this thread nor the story gives a reason for 30m or 5m. Work the value out with `flai cod` on your next run.

### 2026-10-08T00:06:25Z orchestrator
Resolved: time_lost_per_cycle 15m set on S-0318, as the planner recommended
