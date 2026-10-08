---
id: TH-0300
title: S-0316 draft lacks touches, a forecast, and a cost of delay value
anchor:
  path: wip/kanban/stories/S-0316-the-dashboard-upgrade-stops-the-old-container-and-cannot-start-the-new-one-because-the-name-flaiover-is-still-in-use.md
  item: S-0316
status: resolved
participants: [orchestrator]
created: 2026-10-07T19:07:54Z
updated: 2026-10-07T23:55:53Z
---

# TH-0300 S-0316 draft lacks touches, a forecast, and a cost of delay value

On wip/kanban/stories/S-0316-the-dashboard-upgrade-stops-the-old-container-and-cannot-start-the-new-one-because-the-name-flaiover-is-still-in-use.md.

## Entries

### 2026-10-07T19:07:54Z orchestrator
Recommendation: ask the planner to plan S-0316 (dashboard Plan, or `flai plan S-0316`), then set its cost of delay inputs so `flai cod` can give a value.

I am leaving S-0316 as a draft. `flai promote --drafts` lists it as incomplete: no touches, no forecast duration, no forecast delivery, no cost of delay value.

### 2026-10-07T23:55:53Z orchestrator
Resolved: Settled: planner-S-0316 set touches, a 30m forecast, and a 30 USD/week value, and the orchestrator finalized S-0316
