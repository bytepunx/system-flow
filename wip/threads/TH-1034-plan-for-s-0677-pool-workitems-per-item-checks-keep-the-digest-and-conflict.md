---
id: TH-1034
title: "Plan for S-0677: pool workItems' per-item checks, keep the digest and conflict-marker results by stat, four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0677-a-whole-flai-check-takes-3-3-to-4-9-s-on-a-loaded-host-after-s-0601-not-under-2.md
  item: S-0677
status: open
participants: [planner-S-0677]
created: 2026-11-10T10:32:48Z
updated: 2026-11-10T10:32:48Z
---

# TH-1034 Plan for S-0677: pool workItems' per-item checks, keep the digest and conflict-marker results by stat, four tasks in three layers

On wip/kanban/stories/S-0677-a-whole-flai-check-takes-3-3-to-4-9-s-on-a-loaded-host-after-s-0601-not-under-2.md.

## Entries

### 2026-11-10T10:32:48Z planner-S-0677
The plan: run `oneItem` on check.Run's pool with a children index, keep the digest rule's and the conflict-marker scan's per-file results under `.flai-cache` keyed by size and modification time, then measure and close I-0312.

## Tasks and layers

| Layer | Task | Waits for | Why |
|-------|------|-----------|-----|
| 1 | T-3010 Run oneItem over the items on check.Run's pool, with only the byID fill first and alone | none | |
| 1 | T-3012 Keep the digest rule's per-entry results under .flai-cache | none | |
| 2 | T-3016 Keep the conflict-marker scan's per-file results under .flai-cache | T-3012 | both change `check_records.go`; it follows T-3012's kept file |
| 3 | T-3020 Time a warm flai check against main, record it, close I-0312 | T-3010, T-3016 | it times their combined change |

## Figures

- Forecast 1h, delivery 2026-11-10T15:00:00Z. flai said 33m from 17 touches. Two new kept-cache files with tests, and timed runs on a loaded host, take more.
- Cost of delay value 0.08 USD a week, from the inputs flai set from I-0312.

## Assumptions

- T-3016 is an addition to what I-0312 names. Without it, `conflictMarkers` (0.45 to 1.10 s) stays on the critical path once the digest rule is cached.
- A kept digest result covers only what the entry's file alone decides. The source-ID and `absorbed_by` lookups run every check, so a removed source is still reported.
- `Repo.List`'s decode (0.66 to 1.10 s) is S-0676's (I-0310). This story does not change it and has no `after` on S-0676. The 2 s target for a whole check may need both stories.
- The story has no topics of its own. I set `cli` and `server-side`.

## For the operator

The cost of delay input `time_lost_per_cycle` 2s counts I-0312's one report. Every `flai check`, `flai verify`, close-out check step, and `flai task done` pays the time this story removes, many times a week. If you want the value to reflect that, set a larger input, for example 2s times the checks run per week. I changed nothing of yours. The plan stands at either figure.
