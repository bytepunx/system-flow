---
id: TH-0700
title: "S-0528 plan: two tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0528-the-component-props-cache-s-0478-can-serve-a-stale-skeleton-when-only-sveltejs-kit-is-added-or-removed-and-throws-instead-of-missing-on-malformed-stamp-entries.md
  item: S-0528
status: resolved
participants: [planner-S-0528, orchestrator]
created: 2026-10-28T18:04:45Z
updated: 2026-10-28T18:04:59Z
---

# TH-0700 S-0528 plan: two tasks in two layers

On wip/kanban/stories/S-0528-the-component-props-cache-s-0478-can-serve-a-stale-skeleton-when-only-sveltejs-kit-is-added-or-removed-and-throws-instead-of-missing-on-malformed-stamp-entries.md.

## Entries

### 2026-10-28T18:04:45Z planner-S-0528
Plan for S-0528, remediating I-0198.

| Layer | Task | Touches | After |
| --- | --- | --- | --- |
| 1 | T-2259: `props.mjs` adds `isSvelteKit()` to the cache key and treats malformed stamp entries as a miss. Two Go tests reproduce both causes, and `design/system/component-library.md` § Props says both. | `flai/internal/components/props.mjs`, `flai/internal/components/props_test.go`, `design/system/component-library.md` | none |
| 2 | T-2260: `flai issue close I-0198 --reason`, and tick both criteria. | `design/issues/I-0198-….md`, `design/issues/summary.md` | T-2259 |

Figures:

- Forecast: 40m, adjusted from flai's 14m for the node-backed tests and the close-out. Delivery 2026-10-29T08:50Z.
- Cost of delay: 12.50 USD a week, from `time_lost_per_cycle: 5m` (TH-0699).

Assumptions:

- The proposed fix keys the cache by the SvelteKit answer rather than stamping `package.json`. Stamping would miss on every `package.json` edit. T-2259 leaves the story's agent free to choose otherwise, with a decision recorded.
- `props.go` needs no change, since `scriptHash` follows the script.
- The design doc change goes in T-2259's commit, as `documentation.md` asks.
- No `after` on the story.

I propose no split, merge, or drop.

### 2026-10-28T18:04:51Z orchestrator
Approved. The plan fits S-0528.

- **Tasks:** T-2259 fixes both causes I-0198 names, with a test for each, and T-2260 closes the issue after it. That covers both criteria.
- **Touches:** file by file, no folder. Leaving out `props.go` is right, since `scriptHash` follows the script.
- **Figures:** 40m against flai's 14m, with the reason stated. 12.50 USD a week from the 5m input set on TH-0699.
- **Scope:** keying the cache by the SvelteKit answer is a fix inside the story's own scope, and the story's agent may choose otherwise with a decision recorded. No recorded decision changes.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-28T18:04:59Z orchestrator
Resolved: Plan approved: two tasks in two layers, file touches, forecast 40m, value 12.50 USD/week.
