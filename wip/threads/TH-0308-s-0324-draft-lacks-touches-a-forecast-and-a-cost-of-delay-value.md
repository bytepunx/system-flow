---
id: TH-0308
title: S-0324 draft lacks touches, a forecast, and a cost of delay value
anchor:
  path: wip/kanban/stories/S-0324-flai-s-md034-does-not-report-a-bare-www-literal-which-markdownlint-reports-as-a-bare-url.md
  item: S-0324
status: resolved
participants: [orchestrator]
created: 2026-10-07T19:08:09Z
updated: 2026-10-08T00:29:53Z
---

# TH-0308 S-0324 draft lacks touches, a forecast, and a cost of delay value

On wip/kanban/stories/S-0324-flai-s-md034-does-not-report-a-bare-www-literal-which-markdownlint-reports-as-a-bare-url.md.

## Entries

### 2026-10-07T19:08:09Z orchestrator
Recommendation: ask the planner to plan S-0324 (dashboard Plan, or `flai plan S-0324`), then set its cost of delay inputs so `flai cod` can give a value.

I am leaving S-0324 as a draft. `flai promote --drafts` lists it as incomplete: no touches, no forecast duration, no forecast delivery, no cost of delay value.

It is in the same area as S-0289 (flai's markdown lint and table pipes in code spans). You may want them planned together.

### 2026-10-08T00:29:53Z orchestrator
Resolved: Settled: planner-S-0324 set touches, a 12m forecast, and a 12.50 USD/week value, and the orchestrator finalized S-0324; S-0289 shares no file with it
