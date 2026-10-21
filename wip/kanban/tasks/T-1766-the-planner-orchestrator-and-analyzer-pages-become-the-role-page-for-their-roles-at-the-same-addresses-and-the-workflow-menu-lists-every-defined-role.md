---
id: T-1766
type: task
nature: feature
title: The Planner, Orchestrator, and Analyzer pages become the role page for their roles at the same addresses, and the Workflow menu lists every defined role
status: backlog
parent: S-0420
owner: alex
created: 2026-10-16T05:47:00Z
updated: 2026-10-21T12:29:23Z
transitions: []
stream: S-0420
tags: [dashboard]
touches: [flaiover/src/routes/workflow/planner/+page.svelte, flaiover/src/routes/workflow/planner/planner.svelte.test.ts, flaiover/src/routes/workflow/orchestrator/+page.svelte, flaiover/src/routes/workflow/orchestrator/orchestrator.svelte.test.ts, flaiover/src/routes/workflow/analyzer/+page.svelte, flaiover/src/routes/workflow/analyzer/analyzer.svelte.test.ts, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/components/NavRail.svelte, flaiover/src/lib/components/NavRail.svelte.test.ts, flaiover/src/routes/+layout.svelte, flaiover/src/routes/layout.svelte.test.ts, flaiover/src/lib/components/StrategicAgentPanel.svelte]
after: [T-1765]
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
