---
id: TH-1275
title: "S-0822 stays a draft: no touches, forecast, or cost of delay value yet"
anchor:
  path: wip/kanban/stories/S-0822-a-rebase-and-resolve-job-that-passes-records-stayed-in-review-without-reading.md
  item: S-0822
status: resolved
participants: [orchestrator]
created: 2026-11-17T08:22:54Z
updated: 2026-11-28T00:03:06Z
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

### 2026-11-19T08:40:24Z orchestrator
**Recommended:** run `flai plan S-0822` so its planner adds what TH-1274's answer needs. Then I finalize S-0822 under `finalize_drafts` once `flai promote --drafts` lists it complete.

```sh
flai plan S-0822
```

**What changed:** TH-1274 is settled as `stop`, confirmed after 48h. S-0822 must now stop a running rebase and resolve job when the story leaves review, and write an ADR that refines ADR-0313 §6.

**What the plan lacks:**

- A task that stops the job and ends it `stopped`, outcome `left review`. T-3520 to T-3522 do not cover this.
- A task for the ADR, with a `design/adrs` touch.
- A forecast duration and delivery. `flai promote --drafts` lists it incomplete for this.

**Why I did not start the planner:** `flai plan --candidates` leaves S-0822 out: "its planner ran until 2026-11-17T17:24:03Z and it has not changed since". I start the planner only for the stories it lists.

**Alternative:** finalize S-0822 as it is, once it has a forecast, and let its story agent add the stop task and the ADR when it reads TH-1274.

### 2026-11-28T00:03:06Z orchestrator
Resolved: S-0822 was replanned for TH-1274's stop answer (TH-1324), finalized, and moved to ready
