---
id: T-1793
type: task
nature: feature
title: The dashboard's Strategic Cost and Strategic Use charts and each role page draw a series per role from flai stats
status: in-progress
parent: S-0425
owner: alex
created: 2026-10-16T05:49:35Z
updated: 2026-11-02T19:02:12Z
transitions:
  - to: ready
    at: 2026-11-02T19:02:11Z
    by: agent-S-0425
  - to: in-progress
    at: 2026-11-02T19:02:12Z
    by: agent-S-0425
stream: S-0425
tags: [dashboard]
touches: [flaiover/src/lib/strategic.ts, flaiover/src/lib/strategic.test.ts, flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/usage.ts, flaiover/src/lib/usage.test.ts, flaiover/src/lib/components/RolePage.svelte, flaiover/src/lib/components/RolePage.svelte.test.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts"]
after: [T-1790]
---
# T-1793 The dashboard's Strategic Cost and Strategic Use charts and each role page draw a series per role from flai stats

## Work

- `strategic.ts` and `usage.ts` read the per-role figures `flai stats` now reports.
- The strategic charts in `viz/charts.ts` draw one series per role, the three built-in roles in today's colours and each other role in the next palette slot (`design/tech/charts.md` § Palette and rules).
- `RolePage` shows the role's own cost and use.
- Write tests with a fixture role.

Waits for the stats task, whose figures it draws.

## Done when

The tests pass, and `flai test` on `flaiover/src` is clean.

## Notes
