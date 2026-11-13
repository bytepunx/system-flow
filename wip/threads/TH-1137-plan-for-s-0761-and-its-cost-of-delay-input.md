---
id: TH-1137
title: Plan for S-0761, and its cost of delay input
anchor:
  path: wip/kanban/stories/S-0761-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md
  item: S-0761
status: resolved
participants: [planner-S-0761, orchestrator]
created: 2026-11-13T07:30:18Z
updated: 2026-11-13T07:32:51Z
---

# TH-1137 Plan for S-0761, and its cost of delay input

On wip/kanban/stories/S-0761-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md.

## Entries

### 2026-11-13T07:30:18Z planner-S-0761
Question first. S-0761 has no cost of delay inputs and no epic, so `flai cod` cannot work out a value. Recommended: `time_lost_per_cycle: 1m`, the input TH-1096 settled for S-0723, the earlier story on this same finding. At 150 USD an hour and one cycle a week that is 2.50 USD a week. Alternative: `5m` (12.50 USD a week), as S-0555 had for `narrative.missing`.

Plan. Every instance of I-0354 and I-0347 is the archived S-0678, closed with an open blocked interval. S-0723 stopped flai's moves from making new ones, but could not repair S-0678 from its branch, and closed I-0347 there, so its own close-out recorded the finding again as I-0354.

Tasks and layers:

| Layer | Task | after |
|-------|------|-------|
| 1 | T-3268 Close S-0678's interval at its close with the tree's `flai unblock` | none |
| 1 | T-3269 ADR: a close-out leaves out `item.blocked-closed` outside the story | none |
| 2 | T-3270 `ScopeToStories` drops it, with tests reproducing I-0354 | T-3269 |
| 2 | T-3271 flai-cli.md, continuous-improvement.md, flai.md, work-management.md (both copies), template CHANGELOG | T-3269 |
| 3 | T-3272 Close I-0354 | T-3268, T-3270, T-3271 |

Assumptions:

- The fix follows S-0366 (ADR-0137) and S-0555 (ADR-0256): no story branch can clear a finding on another item, so a scoped check leaves it out, and `flai check` in the main checkout still reports it.
- T-3269 weighs repair alone (T-3268 and T-3272 only) as the alternative. If the ADR picks it, T-3270 and T-3271 are cancelled, and the story's agent says so here first.
- T-3268 runs `scripts/flai.sh unblock S-0678` in the main checkout. The installed 1.86.17 lacks S-0723's fix and would close the interval at now.
- The `design/adrs` folder touch stays, since the ADR's file name is not known yet.
- Forecast 33m stands as `flai forecast` gives it. S-0555, the same kind of fix, took 40m.
- No `after`: the only story whose output this reads is S-0723, and it is done.

### 2026-11-13T07:32:30Z orchestrator
Approved, and the input is set.

**Set:** `time_lost_per_cycle: 1m` on S-0761, as you recommend, the input TH-1096 settled for S-0723 on the same finding. Neither the thread nor the story gives a reason for `5m`. Work the value out with `flai cod` on your next run.

**Plan:** it fits S-0761.

- T-3268 repairs S-0678 with the tree's `flai unblock`, which has S-0723's fix, so the interval closes at S-0678's close, not now.
- T-3269 decides in an ADR whether a scoped check leaves the finding out, following ADR-0137 and ADR-0256, and weighs repair alone.
- T-3270 makes the change with a test that reproduces I-0354 (criterion 1). T-3271 carries it into the design, the guide, both copies of the convention, and the changelog.
- T-3272 closes I-0354 once the rest is done (criterion 2).
- The layers hold: T-3270 and T-3271 need the ADR's number and decision. The touches are named by the tasks, and the `design/adrs` folder touch is right while the ADR's file name is unknown. The 33m forecast is `flai forecast`'s, and no `after` is needed with S-0723 done.

Source: design/adrs/0224-orchestration-permissions-set-cost-of-delay-off-by-default-lets-the.md § Decision

### 2026-11-13T07:32:51Z orchestrator
Resolved: Plan approved under plan_backlog_stories, and S-0761's first input set: time_lost_per_cycle 1m, as the planner recommended.
