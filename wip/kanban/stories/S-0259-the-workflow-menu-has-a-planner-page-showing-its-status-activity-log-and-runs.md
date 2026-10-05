---
id: S-0259
type: story
nature: feature
title: The Workflow menu has a Planner page showing its status, activity log, and runs
status: ready
parent: E-0016
owner: alex
created: 2026-10-04T04:53:03Z
updated: 2026-10-05T00:03:14Z
transitions:
  - to: ready
    at: 2026-10-04T23:40:49Z
    by: alex
tags: [dashboard]
topics: [planning]
touches: [flaiover/src/routes/workflow, flaiover/src/routes/api, flaiover/src/lib/components, flaiover/src/lib/sitemenu.ts, flai/internal/hostapi, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0208, S-0225]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 36.59
  by: planner-E-0016
  at: 2026-10-04T04:53:06Z
forecast:
  duration: 45m
  delivery: 2026-10-05T01:21:00Z
  basis: "Its own forecast of 45m; 7th in the pull order with an in-progress limit of 3, behind S-0252, S-0253, S-0257, S-0244, S-0262, S-0258 and S-0260."
  by: flai
  at: 2026-10-05T00:03:14Z
finalized:
  by: alex
  at: 2026-10-04T23:40:46Z
---
# S-0259 The Workflow menu has a Planner page showing its status, activity log, and runs

## Goal

The planner runs today (S-0208) but has no place in the dashboard. Its page under the Workflow menu shows whether the `plan` action is enabled and a run is going, its activity document newest first, its past runs with their cost, and a form to run it on an epic or a story. It was split from S-0228 on TH-0107 so that it ships without waiting for the orchestrator and the analyzer.

## Acceptance criteria

- [ ] `/workflow/planner` is in the Workflow tier of the site menu, and shows the `plan` host action's state with the command that enables it
- [ ] The page shows the current run, if any, in the agent pane that story pages have
- [ ] The page shows the entries of `wip/agents/planner.md`, newest first, each with its items, duration, and cost, and the document's front matter totals
- [ ] The page lists the past runs with their cost and outcome
- [ ] A form runs the planner on an epic or a story through the `plan` host action, and says why when the action is off
- [ ] The page updates live from `/api/events` when `wip/agents/planner.md` changes
- [ ] `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` describe the page, and tests cover its data and the form

## Tasks

## Notes

- Split from S-0228 on TH-0107. S-0228 keeps the Orchestrator and Analyzer pages, and reuses this page's layout.
