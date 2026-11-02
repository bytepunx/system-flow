---
id: T-1766
type: task
nature: feature
title: The Planner, Orchestrator, and Analyzer pages become the role page for their roles at the same addresses, and the Workflow menu lists every defined role
status: done
parent: S-0420
owner: alex
created: 2026-10-16T05:47:00Z
updated: 2026-11-02T01:24:04Z
transitions:
  - to: ready
    at: 2026-11-02T01:05:19Z
    by: agent-S-0420
  - to: in-progress
    at: 2026-11-02T01:05:20Z
    by: agent-S-0420
  - to: done
    at: 2026-11-02T01:24:04Z
    by: agent-S-0420
stream: S-0420
tags: [dashboard]
touches: [flaiover/src/routes/workflow/planner/+page.svelte, flaiover/src/routes/workflow/planner/planner.svelte.test.ts, flaiover/src/routes/workflow/orchestrator/+page.svelte, flaiover/src/routes/workflow/orchestrator/orchestrator.svelte.test.ts, flaiover/src/routes/workflow/analyzer/+page.svelte, flaiover/src/routes/workflow/analyzer/analyzer.svelte.test.ts, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/components/NavRail.svelte, flaiover/src/lib/components/NavRail.svelte.test.ts, flaiover/src/routes/+layout.svelte, flaiover/src/routes/layout.svelte.test.ts, flaiover/src/lib/components/StrategicAgentPanel.svelte, flaiover/src/lib/components/PlannerPanel.svelte, flaiover/src/lib/components/PlannerPanel.svelte.test.ts, flaiover/src/lib/components/OrchestratorPanel.svelte, flaiover/src/lib/components/OrchestratorPanel.svelte.test.ts, flaiover/src/lib/components/AnalyzerPanel.svelte, flaiover/src/lib/components/AnalyzerPanel.svelte.test.ts, flaiover/src/lib/components/RolePage.svelte]
after: [T-1765]
usage:
  source: log
  seconds: 1122
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 186
      output: 991
      cache_read: 12953392
      cache_write: 340309
      cost: 4.6618
---
# T-1766 The Planner, Orchestrator, and Analyzer pages become the role page for their roles at the same addresses, and the Workflow menu lists every defined role

## Work

- Make the three pages render `RolePage` for their role, with what each shows today and nothing less: the planner's runs by item, the orchestrator's last decisions and its Stop and Start, and the analyzer's focus and report links.
- Retire `StrategicAgentPanel.svelte`, or keep it as a part of `RolePage`.
- The rail's Agents group (`sitemenu.ts`, `NavRail.svelte`) lists each role `agent.status` reports, the three first in today's order.
- The existing page tests pass with no changed assertion; add a menu test with a fixture role.

Waits for the settings panel, so the pages move once, whole.

## Done when

The three page tests and `sitemenu.test.ts` pass. `flai test` on `flaiover/src` is clean.

## Notes

- Replanned 2026-10-21 (TH-0580): the third Work bullet named the Workflow group in `sitemenu.ts`. Since S-0458 the three pages sit in the rail's Agents group, which `NavRail.svelte` renders, so the bullet now names it, as the orchestrator approved.
