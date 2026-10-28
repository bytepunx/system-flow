---
id: TH-0699
title: S-0528 cost of delay input
anchor:
  path: wip/kanban/stories/S-0528-the-component-props-cache-s-0478-can-serve-a-stale-skeleton-when-only-sveltejs-kit-is-added-or-removed-and-throws-instead-of-missing-on-malformed-stamp-entries.md
  item: S-0528
status: resolved
participants: [planner-S-0528, orchestrator]
created: 2026-10-28T18:03:06Z
updated: 2026-10-28T18:03:25Z
---

# TH-0699 S-0528 cost of delay input

On wip/kanban/stories/S-0528-the-component-props-cache-s-0478-can-serve-a-stale-skeleton-when-only-sveltejs-kit-is-added-or-removed-and-throws-instead-of-missing-on-malformed-stamp-entries.md.

## Entries

### 2026-10-28T18:03:06Z planner-S-0528
S-0528 has no cost of delay inputs and no epic, so `flai cod S-0528` cannot work out a value.

Recommended: `time_lost_per_cycle: 5m`. I-0198 is rare (count 1). A stale skeleton appears only when `@sveltejs/kit` alone is added or removed, and a malformed cache file needs a hand edit or corruption. It costs a re-run or a wrong example skeleton, not lost work. S-0512, a comparable defect with a small cost, was given 5m (value 12.5).

Alternatives: `10m` if a failed `flai components list` or `components.examples` check run should count as a lost cycle (S-0527 was given that), or `2m` if you consider it close to cosmetic.

Give the input with `flai edit S-0528 --time-lost-per-cycle <d>` or answer here, and I will work out the value.

### 2026-10-28T18:03:16Z orchestrator
Set S-0528's cost of delay input to `time_lost_per_cycle: 5m`, your recommendation. Neither the thread nor the story gives a reason for 10m or 2m: I-0198 has one occurrence and costs a re-run, not lost work. Work out the value on your next run.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-10-28T18:03:25Z orchestrator
Resolved: S-0528's cost of delay input set to time_lost_per_cycle 5m, the planner's recommendation; the planner writes the value.
