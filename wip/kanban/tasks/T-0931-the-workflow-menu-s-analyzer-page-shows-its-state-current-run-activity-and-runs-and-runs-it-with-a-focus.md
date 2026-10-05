---
id: T-0931
type: task
nature: feature
title: The Workflow menu's Analyzer page shows its state, current run, activity, and runs, and runs it with a focus
status: backlog
parent: S-0228
owner: alex
created: 2026-10-05T05:45:13Z
updated: 2026-10-05T05:45:13Z
transitions: []
stream: S-0228
tags: [dashboard]
touches: [flaiover/src/routes/workflow/analyzer, flaiover/src/lib/components/AnalyzerPanel.svelte, flaiover/src/lib/components/AnalyzerPanel.svelte.test.ts, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts]
after: [T-0928]
---
# T-0931 The Workflow menu's Analyzer page shows its state, current run, activity, and runs, and runs it with a focus

## Work

`/workflow/analyzer`, beside Planner and Orchestrator in the Workflow tier of the site menu, on the shared `StrategicAgentPanel.svelte` and what `/api/analyzer` answers:

- The `analyze` host action's state: on, or off with `flai serve enable analyze`.
- The current run in the agent pane, the activity document's front matter totals, its entries newest first with their items, seconds, and cost, each naming the report it wrote with a link to it on the documents page, and the past runs with their cost and outcome.
- A Run button with a focus to choose, `bottlenecks`, `intent`, `risk`, or none, through `POST /api/analyzer`, saying what flai answered; disabled, saying why, while the action is off or a run goes.
- It loads again when a narrative changes, which `wip/agents/analyzer.md` is (`follow(['narrative'])`), and when flai serve says an analyzer run started or ended.

It waits for T-0928, which takes the shared panel out of the Planner page and adds Orchestrator to the site menu, which this task also changes.

## Done when

- [ ] The page and its component show each part above, with component tests for the state, the runs, the entries, the live reload, and Run with and without a focus, on and off
- [ ] The site menu lists Analyzer in the Workflow tier, with its test
- [ ] `npx vitest run` passes for the files it changed, and `npm run check` reports nothing in them

## Notes
