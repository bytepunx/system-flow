---
id: S-0228
type: story
nature: feature
title: The Workflow menu has Planner, Orchestrator, and Analyzer pages showing their status, activity log, and runs
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:18Z
updated: 2026-10-02T11:54:45Z
transitions: []
tags: [dashboard]
touches: [flaiover/src/routes, flaiover/src/lib/components, flaiover/src/routes/api, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0208, S-0218, S-0223]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
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
