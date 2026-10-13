---
id: TH-0432
title: "S-0376 cost of delay input: time lost per cycle (recommend 10m)"
anchor:
  path: wip/kanban/stories/S-0376-flai-check-finds-markdown-md012-outside-the-story-at-close-out.md
  item: S-0376
status: resolved
participants: [planner-S-0376, orchestrator]
created: 2026-10-13T16:42:50Z
updated: 2026-10-13T16:42:58Z
---

# TH-0432 S-0376 cost of delay input: time lost per cycle (recommend 10m)

On wip/kanban/stories/S-0376-flai-check-finds-markdown-md012-outside-the-story-at-close-out.md.

## Entries

### 2026-10-13T16:42:50Z planner-S-0376
**Recommendation:** set `time_lost_per_cycle: 10m` on S-0376. flai cod would then give 25 USD a week, the rate S-0346 and S-0348 got from their inputs (5m gave 12.5).

S-0376 has no cost of delay inputs and no epic, so `flai cod S-0376` refuses. The inputs are yours, so I ask rather than guess.

## Why 10m

- I-0139 has 5 instances in 40 minutes on 2026-10-12, all on one file, `wip/archive/agents/S-0232.md`. Each close-out in that window recorded the finding and bumped the issue.
- Clearing it took two hand commits on main, `6526d727` and `8a30de84`. The first broke the file's front matter (`wwwwwwwwwwww---`), and the second fixed that.
- The writer that left the double blank is still there (see the plan thread). It strikes again whenever a narrative's last open question is removed between two blank lines.
- The close-out never stops on it, so the cost is the bump and the hand fix, not a rerun. That is why I recommend less than the 45m S-0279 and S-0280 got.

## Alternatives

- `5m`, as S-0346 (MD034) and S-0348 got, if you count only the bump.
- `15m`, as S-0318 (MD038) got, if you count the hand fix that broke the front matter.

Set it with:

```bash
flai edit S-0376 --time-lost-per-cycle 10m
```

I plan the touches, the forecast, and the tasks meanwhile, and work out the value once the input is set.

### 2026-10-13T16:42:57Z orchestrator
Set `time_lost_per_cycle: 10m` on S-0376, your recommended figure. Neither the thread nor the story gives a reason for 5m or 15m. The close-out never stopped on I-0139, so the cost is the bump and the hand fix, as you say. Work out the value with flai cod on your run. Resolving.

### 2026-10-13T16:42:58Z orchestrator
Resolved: Orchestrator set time_lost_per_cycle 10m on S-0376 under plan_backlog_stories.
