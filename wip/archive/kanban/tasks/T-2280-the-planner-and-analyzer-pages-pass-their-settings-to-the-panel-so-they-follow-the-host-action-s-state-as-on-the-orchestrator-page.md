---
id: T-2280
type: task
nature: improvement
title: The Planner and Analyzer pages pass their settings to the panel, so they follow the host action's state as on the Orchestrator page
status: done
parent: S-0533
owner: alex
created: 2026-10-28T21:14:26Z
updated: 2026-10-28T21:23:15Z
transitions:
  - to: ready
    at: 2026-10-28T21:21:02Z
    by: system-flow
  - to: in-progress
    at: 2026-10-28T21:21:03Z
    by: system-flow
  - to: done
    at: 2026-10-28T21:23:15Z
    by: agent-S-0533
stream: S-0533
tags: [flaiover, settings]
touches: [flaiover/src/routes/workflow/planner/+page.svelte, flaiover/src/routes/workflow/planner/planner.svelte.test.ts, flaiover/src/routes/workflow/analyzer/+page.svelte, flaiover/src/routes/workflow/analyzer/analyzer.svelte.test.ts]
after: [T-2278]
usage:
  source: log
  seconds: 132
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 20
      output: 4849
      cache_read: 1695255
      cache_write: 45102
      cost: 0.7969
---
# T-2280 The Planner and Analyzer pages pass their settings to the panel, so they follow the host action's state as on the Orchestrator page

## Work

- Do on `flaiover/src/routes/workflow/planner/+page.svelte` (`StrategicSettings` told `planning`) and `flaiover/src/routes/workflow/analyzer/+page.svelte` (told `analysis`) what T-2279 does on the Orchestrator page. Move the settings block, its error, and its `strategic-absent` notice from the foot of the page into a snippet given to `PlannerPanel` or `AnalyzerPanel` as `settings`. Render it under the heading while the panel cannot show.
- In `planner.svelte.test.ts` and `analyzer.svelte.test.ts`, assert the document order: `<kind>-action`, the settings panel, then `<kind>-current` and `<kind>-activity`. Assert the `strategic-absent` notice sits at the same place.

It waits for T-2278, which gives the panels their `settings` snippet. It runs beside T-2279, with no path in common.

## Done when

- [ ] On `/workflow/planner` and `/workflow/analyzer` the settings panel follows the host action's section and precedes Running now and the activity log
- [ ] The `strategic-absent` notice takes the panel's place on both pages
- [ ] The order tests in both page tests pass, and `flai test` on the changed files passes

## Notes
