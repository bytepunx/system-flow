---
id: TH-0887
title: "S-0606 draft lacks a cost of delay value; E-0033 has no inputs to give it one"
anchor:
  path: wip/kanban/stories/S-0606-item-get-answers-from-a-kept-item-list-and-the-file-watcher-stops-statting.md
  item: S-0606
status: resolved
participants: [orchestrator]
created: 2026-11-07T22:50:12Z
updated: 2026-11-07T22:54:40Z
---

# TH-0887 S-0606 draft lacks a cost of delay value; E-0033 has no inputs to give it one

On wip/kanban/stories/S-0606-item-get-answers-from-a-kept-item-list-and-the-file-watcher-stops-statting.md.

## Entries

### 2026-11-07T22:50:12Z orchestrator
**Recommended:** set a cost of delay input on E-0033, such as `time_lost_per_cycle`, so that the planner can work out a value for S-0606 and the other 13 drafts under it (S-0598 to S-0611). If you want one figure for the whole epic, 2h a cycle matches what I recommended for S-0596 on TH-0884: the time lost is your wait on slow flai serve calls.

```bash
flai edit E-0033 --time-lost-per-cycle 2h
```

**What is missing:** `flai promote --drafts` lists S-0606 at position 1 under `throughput`, not complete, for one reason only:

```text
lacks: no cost of delay value
```

Everything else is consistent. The four tasks (T-2615 to T-2618) cover the three criteria. The touches name every file the tasks change, with no folder. The forecast is 45m, its basis stated. `after: [S-0596]` is explained in the Notes: S-0596 writes the section of `design/system/server-performance.md` that criterion 3 records into.

**Why I do not set it:** the planner's Notes say neither S-0606 nor E-0033 has inputs, so `flai cod` refuses, and that the inputs are yours. It opened no thread asking for them, so I have no figure of its to take under `plan_backlog_stories`. A figure of my own is not mine to set.

Once an input is set and the planner gives S-0606 a value, I finalize it under `finalize_drafts`.

### 2026-11-07T22:54:40Z orchestrator
Resolved: S-0606 now has a cost of delay value of 30 a week, from penalty_per_week 30 set on TH-0886. flai promote --drafts listed it as complete, and I finalized it under finalize_drafts.
