---
id: TH-0302
title: S-0318 draft lacks touches, a forecast, and a cost of delay value
anchor:
  path: wip/kanban/stories/S-0318-flai-check-finds-markdown-md038-outside-the-story-at-close-out.md
  item: S-0318
status: resolved
participants: [orchestrator]
created: 2026-10-07T19:07:58Z
updated: 2026-10-08T00:11:42Z
---

# TH-0302 S-0318 draft lacks touches, a forecast, and a cost of delay value

On wip/kanban/stories/S-0318-flai-check-finds-markdown-md038-outside-the-story-at-close-out.md.

## Entries

### 2026-10-07T19:07:58Z orchestrator
Recommendation: ask the planner to plan S-0318 (dashboard Plan, or `flai plan S-0318`), then set its cost of delay inputs so `flai cod` can give a value.

I am leaving S-0318 as a draft. `flai promote --drafts` lists it as incomplete: no touches, no forecast duration, no forecast delivery, no cost of delay value.

It is the same kind of problem as S-0280 (`item.archive`), S-0323, and S-0325 (`narrative.state`): a check finding outside the story at close-out. One remedy may cover all four. You may want them planned together.

### 2026-10-08T00:11:42Z orchestrator
Resolved: Settled: planner-S-0318 set touches, a 25m forecast, and a 37.5 USD/week value (from the 15m input), and the orchestrator finalized S-0318; the S-0323/S-0325 overlap is raised on TH-0346
