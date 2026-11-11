---
id: TH-1065
title: "S-0692: which analysis.schedule should T-3037 set? (cadence from TH-1038)"
anchor:
  path: wip/kanban/stories/S-0692-no-analysis-schedule-or-analysis-formula-is-set-so-the-analyzer-ran-once-in.md
  item: S-0692
status: resolved
participants: [agent-S-0692, orchestrator, alex]
created: 2026-11-11T02:22:36Z
updated: 2026-11-11T02:26:30Z
---

# TH-1065 S-0692: which analysis.schedule should T-3037 set? (cadence from TH-1038)

On wip/kanban/stories/S-0692-no-analysis-schedule-or-analysis-formula-is-set-so-the-analyzer-ran-once-in.md.

## Entries

### 2026-11-11T02:22:36Z agent-S-0692
S-0692's code and docs are done: `flai manifest set` now writes `analysis.formula` (T-3031, T-3034). Only T-3037 is left before the close, and it needs your cadence. TH-1038 asks the same question; answer either thread.

**Recommended:** weekly, `0 6 * * 1`, Mondays at 06:00 UTC. That is about 3.5 USD a month at the last run's 0.82 USD.

**Alternatives:**

| Cadence | Spec | About a month |
|---------|------|---------------|
| Daily | `daily` | 25 USD |
| Twice a month | `0 6 1,15 * *` | 1.6 USD |

T-3037 sets the spec you choose with `analysis.formula: analysis-cadence`. Then T-3040 closes I-0292. The schedule starts nothing until you turn on the `analyze` host action for this project.

### 2026-11-11T02:23:55Z orchestrator (recommendation)
**Recommended:** weekly, `0 6 * * 1`, Mondays at 06:00 UTC, with `analysis.formula: analysis-cadence`. That is about 3.5 USD a month at the last run's 0.82 USD.

**Why:**

- It is the same answer I recommend on TH-1038, which is still open. One answer from you settles both threads.
- A week is `planning.cycle`'s default, 168h. This project's manifest does not set it.
- The 2026-11-07 analysis found the W44 regressions a week late. A weekly run closes that gap.
- The schedule is a cron spec in UTC, and `analysis.formula` makes the `analysis-cadence` instance each time it comes round. Nothing starts until you turn on the `analyze` host action.

**Alternatives:**

| Cadence | Spec | About a month |
|---------|------|---------------|
| Daily | `daily` | 25 USD |
| Twice a month | `0 6 1,15 * *` | 1.6 USD |

**Why I do not answer it:** the cadence sets recurring spend. No ADR, design section, or digest entry records a cadence for this project, so the choice is yours alone. T-3037 sets the spec you choose, then T-3040 closes I-0292.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/system/strategic-agents.md § The analysis schedule

### 2026-11-11T02:25:28Z alex
Confirmed the recommendation of 2026-11-11T02:23:55Z orchestrator.

Source: design/system/strategic-agents.md § The analysis schedule

### 2026-11-11T02:26:30Z agent-S-0692
Resolved: Weekly confirmed. T-3037 set analysis.schedule `0 6 * * 1` and analysis.formula analysis-cadence (4d052e4), and T-3040 closed I-0292 (4e59db8).
