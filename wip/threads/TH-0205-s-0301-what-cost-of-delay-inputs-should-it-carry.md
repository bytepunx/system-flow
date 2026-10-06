---
id: TH-0205
title: "S-0301: what cost of delay inputs should it carry?"
anchor:
  path: wip/kanban/stories/S-0301-flai-upgrade-consistently-pulls-an-old-template-no-matter-what.md
  item: S-0301
status: answered
participants: [planner-S-0301, alex]
created: 2026-10-06T22:46:55Z
updated: 2026-10-06T22:53:05Z
---

# TH-0205 S-0301: what cost of delay inputs should it carry?

On wip/kanban/stories/S-0301-flai-upgrade-consistently-pulls-an-old-template-no-matter-what.md.

## Entries

### 2026-10-06T22:46:55Z planner-S-0301
S-0301 has no cost of delay inputs and no epic, so `flai cod S-0301` has nothing to work from. The inputs are yours to set.

**My recommendation:** `time_lost_per_cycle: 30m`, with no revenue and no penalty. Until this is fixed, each upgrade of a project made from the template either leaves it on 1.0.18 or needs the lock file edited by hand. I reckon that costs about half an hour of your time per cycle. At `planning.hour_rate` 150, that would make the value about 75 a week, if flai counts one cycle a week.

Reply `yes` to set that, or give your own figures, for example `revenue_per_week`, `penalty_per_week` or `time_lost_per_cycle`. I'm drafting the touches, forecast and tasks in the meantime.

### 2026-10-06T22:53:05Z alex
yes
