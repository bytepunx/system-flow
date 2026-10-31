---
id: TH-0761
title: S-0527 passes verify and meets both criteria at 7aff11bb, but the queue will not take it while TH-0697 is open
anchor:
  path: wip/kanban/stories/S-0527-metrics-md-charts-says-burn-up-draws-a-forecast-line-from-throughput-and-cycle-time-filters-by-nature-but-neither-chart-does.md
  item: S-0527
status: resolved
participants: [orchestrator, alex]
created: 2026-10-30T00:20:48Z
updated: 2026-10-31T07:46:20Z
---

# TH-0761 S-0527 passes verify and meets both criteria at 7aff11bb, but the queue will not take it while TH-0697 is open

On wip/kanban/stories/S-0527-metrics-md-charts-says-burn-up-draws-a-forecast-line-from-throughput-and-cycle-time-filters-by-nature-but-neither-chart-does.md.

## Entries

### 2026-10-30T00:20:48Z orchestrator
**Recommended:** answer TH-0697, my recommendation to correct `design/system/metrics.md` § Charts rather than build the Burn-up forecast line and the Cycle Time nature filter. Confirm it or answer otherwise. Once TH-0697 is answered and resolved, I enqueue S-0527.

## Where S-0527 stands

- **Verify:** passed at `7aff11bb39f017e3a201793b89f9dda8cb42eeb1`, the head of `story/S-0527`, against main `fae791a8`. Every step passed, the `flaiover` tier included.
- **Criterion 1 met:**
  - ADR-0219 (accepted) refines ADR-0054 and ADR-0057.
  - `design/system/metrics.md` § Charts no longer claims a forecast line or a nature filter, and it adds rows for Agent Waiting, Strategic Cost, and Strategic Use.
  - `flaiover/src/lib/viz/explain.test.ts` holds every catalog chart to a row.
- **Criterion 2 met:** I-0193 is closed with a reason, and `design/issues/summary.md` is updated.
- **Touches:** every changed file is under the story's touches.

## The one blocker

- The queue refuses `thread_open`: TH-0697, the planner's plan thread, is still open.
- It is open because it asks a change of scope: it takes two features out of the contract instead of building them. That decision is yours.
- I cannot resolve it while my recommendation waits on you.
- With `confirm_after: 48h`, flai serve confirms that recommendation on 2026-10-30T15:11:48Z if you do not answer first.

I resolve this thread once TH-0697 is settled.

### 2026-10-31T07:46:20Z alex
Resolved: S-0527 was accepted
