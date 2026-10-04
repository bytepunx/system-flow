---
id: S-0228
type: story
nature: feature
title: The Workflow menu has Planner, Orchestrator, and Analyzer pages showing their status, activity log, and runs
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:18Z
updated: 2026-10-04T04:48:23Z
transitions: []
tags: [dashboard]
touches: [flaiover/src/routes, flaiover/src/lib/components, flaiover/src/routes/api, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/sitemenu.ts, flai/internal/hostapi]
after: [S-0208, S-0218, S-0223]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 97.56
  by: planner-E-0016
  at: 2026-10-04T04:48:23Z
forecast:
  duration: 2h
  delivery: 2026-10-05T00:30:00Z
  basis: "Three pages with live updates and run forms, each larger than S-0202 (1197 s), like S-0204 (4322 s) in all; waits for S-0218 and S-0223."
  by: planner-E-0016
  at: 2026-10-04T04:45:04Z
---
# S-0228 The Workflow menu has Planner, Orchestrator, and Analyzer pages showing their status, activity log, and runs

## Goal

The strategic agents are invisible without a place in the dashboard. Each gets a page under the Workflow menu with whether it is enabled and running, its activity document rendered newest first, its runs with cost and outcome, and the actions the operator can take (run the planner on an item, run the analyzer, start or stop the orchestrator).

## Acceptance criteria
- [ ] `/workflow/planner`, `/workflow/orchestrator`, `/workflow/analyzer` in the two-tier site menu's Workflow tier; each shows: the host action's state with the command that enables it, the current run if any (the agent pane as story pages have it), the activity log's entries with duration and cost, the front matter totals, and a list of past runs
- [ ] The planner page has a form to run it on an epic or story; the analyzer page a Run button with the focus; the orchestrator page a Stop button and the last decisions with their reasons
- [ ] Pages update live from `/api/events` as the activity documents change
- [ ] `design/system/flaiover-dashboard.md` and the user guide describe them; tests cover each page's data and actions

## Tasks

## Notes
