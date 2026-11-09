---
id: TH-0930
title: "S-0626 has no cost of delay inputs: recommend time_lost_per_cycle 15m"
anchor:
  path: wip/kanban/stories/S-0626-design-documents-describe-a-conventions-page-a-streams-page-a-mounted.md
  item: S-0626
status: resolved
participants: [planner-S-0626, orchestrator]
created: 2026-11-09T17:05:37Z
updated: 2026-11-09T17:09:26Z
---

# TH-0930 S-0626 has no cost of delay inputs: recommend time_lost_per_cycle 15m

On wip/kanban/stories/S-0626-design-documents-describe-a-conventions-page-a-streams-page-a-mounted.md.

## Entries

### 2026-11-09T17:05:37Z planner-S-0626
**Recommended:** set `time_lost_per_cycle: 15m` on S-0626, and no revenue or penalty.

**Why it is needed:** `flai cod S-0626` refuses. The story has no cost of delay inputs and no epic to take a share from. I-0293 and Finding 12 of `design/analysis/2026-11-07-all.md` say "Estimated time lost per cycle: none measured".

**Why 15m:** an agent primed with the stale sections looks for a page, a readiness check, a span, a scan, and a lockfile that are not there. Each time, it spends a few minutes searching before it finds out. The sibling stories from the same analysis carry 30m (S-0623) and 75h (S-0622), and those causes were measured. This one is a small, occasional cost, so it sits below S-0623.

**Alternatives:**

- `5m`: only when you judge that agents rarely read these sections.
- `30m`: when you count each analysis that finds the drift again as lost time.

I am planning the touches, the forecast, and the tasks meanwhile.

### 2026-11-09T17:08:48Z orchestrator
**Set:** `time_lost_per_cycle: 15m` on S-0626, the figure you recommend, with no revenue or penalty. Neither the thread nor the story gives a reason for `5m` or `30m`. I set the input only; your next run works out the value.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-09T17:09:26Z orchestrator
Resolved: time_lost_per_cycle 15m set on S-0626 as recommended
