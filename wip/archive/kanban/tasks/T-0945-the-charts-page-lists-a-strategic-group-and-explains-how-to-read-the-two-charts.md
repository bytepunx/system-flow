---
id: T-0945
type: task
nature: feature
title: The charts page lists a Strategic group and explains how to read the two charts
status: done
parent: S-0216
owner: alex
created: 2026-10-05T05:45:44Z
updated: 2026-10-07T09:56:46Z
transitions:
  - to: ready
    at: 2026-10-07T09:43:43Z
    by: agent-S-0216
  - to: in-progress
    at: 2026-10-07T09:43:43Z
    by: agent-S-0216
  - to: done
    at: 2026-10-07T09:56:46Z
    by: agent-S-0216
stream: S-0216
tags: [dashboard]
touches: ["flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-0939, T-0941]
usage:
  source: log
  seconds: 783
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 66
      output: 27052
      cache_read: 3172540
      cache_write: 207242
      cost: 2.3521
---
# T-0945 The charts page lists a Strategic group and explains how to read the two charts

## Work

- Waits for T-0939 and T-0941: the page draws what they build, and its test renders both kinds.
- In `flaiover/src/routes/charts/[kind]/+page.svelte`, list the Strategic group after Flow and Usage in the charts navigation.
- Add a note under each Strategic chart, as the page does for `cost-spent`, saying how to read it:
  - Strategic Cost: the bars are what planning, orchestration, and analysis cost each day; the line is what a delivered story cost its agents; the ratio is the overhead they add per story.
  - Strategic Use: a fall in the waiting or cycle-time lines while the agents' hours rise is the return on their time.
- The table view lists the days with each figure, as every chart page does.

## Done when

- `charts.svelte.test.ts` renders `/charts/strategic-cost` and `/charts/strategic-use` from a fixture report. It finds the Strategic group, both notes, the ratio, and a table row per day.
- `scripts/flaiover-test.sh` passes.

## Notes
