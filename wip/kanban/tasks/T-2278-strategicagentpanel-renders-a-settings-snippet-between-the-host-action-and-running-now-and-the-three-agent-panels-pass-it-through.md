---
id: T-2278
type: task
nature: improvement
title: StrategicAgentPanel renders a settings snippet between the host action and Running now, and the three agent panels pass it through
status: done
parent: S-0533
owner: alex
created: 2026-10-28T21:13:45Z
updated: 2026-10-28T21:19:48Z
transitions:
  - to: ready
    at: 2026-10-28T21:17:56Z
    by: agent-S-0533
  - to: in-progress
    at: 2026-10-28T21:17:57Z
    by: agent-S-0533
  - to: done
    at: 2026-10-28T21:19:48Z
    by: agent-S-0533
stream: S-0533
tags: [flaiover, settings]
touches: [flaiover/src/lib/components/StrategicAgentPanel.svelte, flaiover/src/lib/components/OrchestratorPanel.svelte, flaiover/src/lib/components/PlannerPanel.svelte, flaiover/src/lib/components/AnalyzerPanel.svelte, flaiover/src/lib/components/OrchestratorPanel.svelte.test.ts]
usage:
  source: log
  seconds: 111
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 8
      output: 1871
      cache_read: 654015
      cache_write: 17400
      cost: 0.3074
---
# T-2278 StrategicAgentPanel renders a settings snippet between the host action and Running now, and the three agent panels pass it through

## Work

- Give `StrategicAgentPanel.svelte` an optional `settings` snippet. Render it right after the `{kind}-action` section and before `{kind}-current` (Running now), so the order is: host action, settings, Running now, `beforeActivity` (the orchestrator's last decisions), Activity, Runs.
- Add a `settings` prop to `OrchestratorPanel.svelte`, `PlannerPanel.svelte`, and `AnalyzerPanel.svelte`, each passed straight to `StrategicAgentPanel`.
- Update the panel's header comment to name the new slot.
- In `OrchestratorPanel.svelte.test.ts`, render the panel with a settings snippet. Assert its element comes after `orchestrator-action` and before `orchestrator-current`, `orchestrator-decisions`, and `orchestrator-activity` in document order (`compareDocumentPosition`).
- If `flai check` reports `components.examples` for these components, update their files under `design/components/flaiover/src/lib/components/`.

It waits for nothing: the page tasks build on it.

## Done when

- [ ] `StrategicAgentPanel` takes `settings` and renders it between the host action section and Running now
- [ ] The three panels pass `settings` through
- [ ] The order test in `OrchestratorPanel.svelte.test.ts` passes, and `flai test` on the changed files passes

## Notes
