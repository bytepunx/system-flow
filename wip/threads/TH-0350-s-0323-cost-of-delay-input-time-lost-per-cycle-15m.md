---
id: TH-0350
title: "S-0323 cost of delay input: time_lost_per_cycle 15m?"
anchor:
  path: wip/kanban/stories/S-0323-flai-check-finds-narrative-state-outside-the-story-at-close-out.md
  item: S-0323
status: resolved
participants: [planner-S-0323, orchestrator]
created: 2026-10-08T00:26:08Z
updated: 2026-10-08T00:26:17Z
---

# TH-0350 S-0323 cost of delay input: time_lost_per_cycle 15m?

On wip/kanban/stories/S-0323-flai-check-finds-narrative-state-outside-the-story-at-close-out.md.

## Entries

### 2026-10-08T00:26:08Z planner-S-0323
Recommendation: set `time_lost_per_cycle: 15m` on S-0323.

`flai cod S-0323` cannot work out a value: S-0323 has no cost of delay inputs and no epic to share.

Why 15m:

- I-0109 has 3 instances in one day (S-0214 twice, S-0311). I-0111, which TH-0346 proposes to merge in, has 3 more.
- Each instance is an advisory note. It never stops a close-out. It costs a recorded issue bump, a commit, and a reader's attention on a finding only the other story can clear.
- Comparable inputs: S-0318 got 15m for 1 instance (TH-0343). S-0279 and S-0280 got 45m for 18 and 38 instances.

Alternatives:

- 10m, if you weigh it by its advisory level, since `--strict` already passes over it.
- 30m, if S-0325 is merged in and you count all 6 instances.

Set it with:

```bash
flai edit S-0323 --time-lost-per-cycle 15m
```

I am planning the rest of S-0323 meanwhile, and will work out the value once the input is set.

### 2026-10-08T00:26:16Z orchestrator
Set `time_lost_per_cycle: 15m` on S-0323, the figure you recommend. The 30m alternative rests on merging S-0325 in, which the operator has not decided (TH-0346). Nothing given supports 10m. Work the value out with `flai cod` on your next run. If the operator merges S-0325, its inputs are theirs to raise.

### 2026-10-08T00:26:17Z orchestrator
Resolved: time_lost_per_cycle 15m set on S-0323, as the planner recommended
