---
id: ADR-0190
title: "An item's queue time is the seconds it spent in ready before it started, so a stretch back in backlog is not counted as waiting"
status: accepted
date: 2026-10-23
supersedes: []
superseded_by: []
---

# ADR-0190 An item's queue time is the seconds it spent in ready before it started, so a stretch back in backlog is not counted as waiting

## Context

Queue time was `started - committed`, with `committed` the first move to `ready` ([metrics.md](../system/metrics.md)). A story the operator sends back from `ready` to `backlog`, and makes ready again later, counted the stretch in backlog as queue time. The bottlenecks analysis of 2026-10-15 read queue time minus `held_seconds`, the time the hold rules held a story in `ready`, as time a story waited in `ready` unheld, and found five E-0016 stories that waited about 75 h so ([I-0153](../issues/I-0153-ready-stories-waited-up-to-3-5-days-unheld-with-in-progress-lanes-free-and-review-under-its-limit.md)). They had not: the operator sent all five back to backlog on 2026-10-04 and made them ready again on 2026-10-07, and their time in `ready` was held almost whole. S-0216 had 305,070 s of queue time, 35,789 s in `ready`, and 35,693 s held.

## Decision

An item's queue time is the sum of its `ready` intervals that end at or before `started`: the seconds it spent in `ready` before its first move to `in-progress`. `committed` stays the first move to `ready`.

- An item never moved out of `ready` but to `in-progress` has the queue time it had, `started - committed`.
- An item sent back to `backlog`, or cancelled, and made ready again counts only its time in `ready`.
- Queue time exists, as before, only when the item has both `committed` and `started`, with `started` not before `committed`.

Rejected: moving `committed` to the last move to `ready` before `started`. That drops the first stretch in `ready` from queue time while `held_seconds` keeps its holds, so held time could exceed queue time.

## Consequences

- Queue time is the time an item waited in `ready`, held or not, and queue time minus `held_seconds` is the time it waited in `ready` unheld, up to the holds of a stretch in `ready` after a move back from `in-progress`.
- `flai stats` and the dashboard's queue time distribution drop the backlog stretches of stories moved back out of `ready`; no other figure changes.
- Time in backlog after a move back from `ready` is still in `time_in_state` under `backlog`, and in lead time.
