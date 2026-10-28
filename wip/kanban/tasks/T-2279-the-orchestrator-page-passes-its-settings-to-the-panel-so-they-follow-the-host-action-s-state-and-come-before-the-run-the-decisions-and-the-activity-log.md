---
id: T-2279
type: task
nature: improvement
title: The Orchestrator page passes its settings to the panel, so they follow the host action's state and come before the run, the decisions, and the activity log
status: done
parent: S-0533
owner: alex
created: 2026-10-28T21:14:09Z
updated: 2026-10-28T21:22:00Z
transitions:
  - to: ready
    at: 2026-10-28T21:20:59Z
    by: system-flow
  - to: in-progress
    at: 2026-10-28T21:21:00Z
    by: system-flow
  - to: done
    at: 2026-10-28T21:22:00Z
    by: agent-S-0533
stream: S-0533
tags: [flaiover, orchestrator, settings]
touches: [flaiover/src/routes/workflow/orchestrator/+page.svelte, flaiover/src/routes/workflow/orchestrator/orchestrator.svelte.test.ts]
after: [T-2278]
usage:
  source: log
  seconds: 60
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 1149
      cache_read: 401620
      cache_write: 10685
      cost: 0.1888
---
# T-2279 The Orchestrator page passes its settings to the panel, so they follow the host action's state and come before the run, the decisions, and the activity log

## Work

- In `flaiover/src/routes/workflow/orchestrator/+page.svelte`, move the settings block out of the `mt-8 border-t` div at the page's foot. That block is the settings error, `StrategicSettings` told `orchestration`, and the `strategic-absent` notice. Put it in a snippet given to `OrchestratorPanel` as `settings`, and drop the divider.
- While the panel cannot show (`view` is null, loading or failed), render the same snippet under the heading and the error, so the settings stay reachable.
- The `strategic-absent` notice shown when `settings.get` has no `host.strategic` renders where `StrategicSettings` would.
- In `orchestrator.svelte.test.ts`, assert the document order: `orchestrator-action`, then the settings panel, then `orchestrator-current`, `orchestrator-decisions`, `orchestrator-activity`. Assert the `strategic-absent` notice sits at the same place when the settings have no `host.strategic`.

It waits for T-2278, which gives the panel its `settings` snippet.

## Done when

- [ ] On `/workflow/orchestrator` the settings panel follows the host action's section and precedes Running now, the last decisions, and the activity log
- [ ] The `strategic-absent` notice takes the panel's place
- [ ] The order tests in `orchestrator.svelte.test.ts` pass, and `flai test` on the changed files passes

## Notes
