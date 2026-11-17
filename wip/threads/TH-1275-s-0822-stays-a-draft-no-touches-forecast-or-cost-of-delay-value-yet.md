---
id: TH-1275
title: "S-0822 stays a draft: no touches, forecast, or cost of delay value yet"
anchor:
  path: wip/kanban/stories/S-0822-a-rebase-and-resolve-job-that-passes-records-stayed-in-review-without-reading.md
  item: S-0822
status: open
participants: [orchestrator]
created: 2026-11-17T08:22:54Z
updated: 2026-11-17T08:22:54Z
---

# TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet

On wip/kanban/stories/S-0822-a-rebase-and-resolve-job-that-passes-records-stayed-in-review-without-reading.md.

## Entries

### 2026-11-17T08:22:54Z orchestrator
**Recommended:** leave S-0822 a draft until its planner finishes, after you answer TH-1274. I then finalize it under `finalize_drafts`.

`flai promote --drafts` lists it incomplete. It lacks:

- touches
- a forecast duration and delivery
- a cost of delay value (its input, `time_lost_per_cycle: 30m`, is set)

The planner's plan waits on TH-1274 (`stop` or `run`), since `stop` adds a task and an ADR. Nothing is needed from you on this thread beyond that answer.
