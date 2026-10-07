---
id: T-0928
type: task
nature: feature
title: The Workflow menu's Orchestrator page shows its state, current run, decisions, and runs, and stops and starts it
status: done
parent: S-0228
owner: alex
created: 2026-10-05T05:45:04Z
updated: 2026-10-07T00:31:15Z
transitions:
  - to: ready
    at: 2026-10-07T00:23:19Z
    by: agent-S-0228
  - to: in-progress
    at: 2026-10-07T00:23:19Z
    by: agent-S-0228
  - to: done
    at: 2026-10-07T00:31:15Z
    by: agent-S-0228
stream: S-0228
tags: [dashboard]
touches: [flaiover/src/routes/workflow/orchestrator, flaiover/src/lib/components/OrchestratorPanel.svelte, flaiover/src/lib/components/OrchestratorPanel.svelte.test.ts, flaiover/src/lib/components/StrategicAgentPanel.svelte, flaiover/src/lib/components/PlannerPanel.svelte, flaiover/src/lib/components/PlannerPanel.svelte.test.ts, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/activity.ts, flaiover/src/lib/components/AgentStream.svelte, flaiover/src/lib/components/AgentStream.svelte.test.ts]
after: [T-0921]
usage:
  source: log
  seconds: 476
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 76
      output: 442
      cache_read: 4115085
      cache_write: 167277
      cost: 1.8975
---
# T-0928 The Workflow menu's Orchestrator page shows its state, current run, decisions, and runs, and stops and starts it

## Work

`/workflow/orchestrator`, beside Planner in the Workflow tier of the site menu (`flaiover/src/lib/sitemenu.ts`), laid out as S-0259's Planner page is and built on what `/api/orchestrator` answers:

- Take what the Planner page and the Orchestrator page share out of `PlannerPanel.svelte` into `StrategicAgentPanel.svelte`: the host action's state with the command that turns it on, the current run in the agent pane story pages have (`AgentStream.svelte`, told to read a run by role), the activity document's front matter totals, its entries newest first with their items, seconds, and cost, and the past runs with their cost and outcome. The Planner page keeps what it shows and its tests.
- The orchestrator's state: on, off with `flai serve enable orchestrate`, or held after a Stop.
- Its last decisions with their reasons, the newest entries of `wip/agents/orchestrator.md`, above the whole log.
- A Stop button while a run goes, and a Start button while it is held, through `POST /api/orchestrator`, each saying what flai answered; both disabled, saying why, while the action is off.
- It loads again when a narrative changes, which `orchestrator.md` is (`follow(['narrative'])`), and when flai serve says an orchestrator run started, ended, or was held.

It waits for T-0921, whose routes and `strategic.ts` it reads.

## Done when

- [ ] The page and its component show each part above, with component tests for the state (on, off, held), the decisions, the runs, the live reload, and Stop and Start
- [ ] The Planner page's tests still pass on the shared panel
- [ ] The site menu lists Orchestrator in the Workflow tier, with its test
- [ ] `npx vitest run` passes for the files it changed, and `npm run check` reports nothing in them

## Notes
