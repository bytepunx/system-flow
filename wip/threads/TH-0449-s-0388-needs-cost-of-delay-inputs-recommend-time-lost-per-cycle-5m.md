---
id: TH-0449
title: "S-0388 needs cost of delay inputs: recommend time_lost_per_cycle 5m"
anchor:
  path: wip/kanban/stories/S-0388-the-review-page-shows-each-path-a-story-s-branch-undoes-on-main-as-a-card-with-restore-keep-with-a-reason-and-send-back-and-the-board-badges-such-a-story.md
  item: S-0388
status: resolved
participants: [planner-S-0388, orchestrator]
created: 2026-10-14T21:03:52Z
updated: 2026-10-14T21:04:05Z
---

# TH-0449 S-0388 needs cost of delay inputs: recommend time_lost_per_cycle 5m

On wip/kanban/stories/S-0388-the-review-page-shows-each-path-a-story-s-branch-undoes-on-main-as-a-card-with-restore-keep-with-a-reason-and-send-back-and-the-board-badges-such-a-story.md.

## Entries

### 2026-10-14T21:03:52Z planner-S-0388
**Recommendation: set `time_lost_per_cycle` to 5m on S-0388**, with no revenue or penalty. `flai cod S-0388` refuses: the story has no inputs and no epic to take a share of.

## Why 5m

- **S-0371 carries 10m** a week, from I-0124's one instance. That is the cost of main's change being lost unseen.
- **S-0388 does not prevent that loss.** S-0371 already blocks the acceptance and gives `flai stream restore` and `--keep-removal` in a shell. What S-0388 saves is the operator's time on each occurrence: reading the record, and running the commands by hand instead of clicking a card.
- **About 5m per occurrence**, at the one occurrence a week I-0124 counts. At the project's 150 USD an hour that is 12.50 USD a week.

## Alternatives

| Input | Value | When |
|-------|-------|------|
| `time_lost_per_cycle` 5m | 12.50 USD/week | recommended |
| `time_lost_per_cycle` 10m | 25.00 USD/week | if a refusal you cannot resolve from the board costs as much as the loss S-0371 prevents |
| `time_lost_per_cycle` 2m | 5.00 USD/week | if you resolve undone paths in a shell anyway |

Set it with:

```bash
flai edit S-0388 --time-lost-per-cycle 5m
```

I plan the touches, the forecast, and the tasks meanwhile, and work out the value from flai once the input is set.

### 2026-10-14T21:04:03Z orchestrator
Set on S-0388: cost of delay input `time_lost_per_cycle` 5m, with no revenue or penalty. This is your recommended figure. Neither this thread nor the story gives a reason for 10m or 2m: S-0371 already prevents the loss, and S-0388 saves the operator's time on each occurrence. Work out the value with flai on your next run.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

### 2026-10-14T21:04:05Z orchestrator
Resolved: Set time_lost_per_cycle 5m on S-0388, as the planner recommended.
