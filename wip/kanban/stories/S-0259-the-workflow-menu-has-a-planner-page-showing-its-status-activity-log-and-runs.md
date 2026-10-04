---
id: S-0259
type: story
nature: feature
title: The Workflow menu has a Planner page showing its status, activity log, and runs
status: backlog
parent: E-0016
owner: alex
created: 2026-10-04T04:53:03Z
updated: 2026-10-04T04:53:06Z
transitions: []
tags: [dashboard]
topics: [planning]
touches: [flaiover/src/routes/workflow, flaiover/src/routes/api, flaiover/src/lib/components, flaiover/src/lib/sitemenu.ts, flai/internal/hostapi, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0208, S-0225]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  value: 36.59
  by: planner-E-0016
  at: 2026-10-04T04:53:06Z
forecast:
  duration: 45m
  delivery: 2026-10-04T09:30:00Z
  basis: "One page over the activity document S-0206 writes and the plan action S-0208 serves, larger than S-0202 (1197 s); in the dashboard lane after S-0212, because both change flaiover/src/lib/sitemenu.ts."
  by: planner-E-0016
  at: 2026-10-04T04:53:06Z
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
