---
id: S-0228
type: story
nature: feature
title: The Workflow menu has Orchestrator and Analyzer pages showing their status, activity log, and runs
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:18Z
updated: 2026-10-04T04:53:11Z
transitions: []
tags: [dashboard]
touches: [flaiover/src/routes, flaiover/src/lib/components, flaiover/src/routes/api, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/sitemenu.ts, flai/internal/hostapi]
after: [S-0208, S-0218, S-0223, S-0259]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 60.98
  by: planner-E-0016
  at: 2026-10-04T04:53:11Z
forecast:
  duration: 1h15m
  delivery: 2026-10-05T00:30:00Z
  basis: "Two pages laid out as S-0259's Planner page, with a Run button and a Stop button, like S-0201 (3771 s) less the flai side; waits for S-0218 and S-0223."
  by: planner-E-0016
  at: 2026-10-04T04:53:11Z
---
# S-0228 The Workflow menu has Orchestrator and Analyzer pages showing their status, activity log, and runs

## Goal

The orchestrator and the analyzer are invisible without a place in the dashboard. Each gets a page under the Workflow menu, laid out like S-0259's Planner page. The page shows whether the agent's host action is enabled and running, its activity document newest first, its runs with cost and outcome, and the actions the operator can take: run the analyzer, or start and stop the orchestrator.

## Acceptance criteria
- [ ] `/workflow/orchestrator` and `/workflow/analyzer` are in the two-tier site menu's Workflow tier, beside the Planner page. Each shows the host action's state with the command that enables it, the current run if any (the agent pane that story pages have), the activity log's entries with duration and cost, the front matter totals, and a list of past runs
- [ ] The analyzer page has a Run button with the focus. The orchestrator page has a Stop button and the last decisions with their reasons
- [ ] Pages update live from `/api/events` as the activity documents change
- [ ] `design/system/flaiover-dashboard.md` and the user guide describe them; tests cover each page's data and actions

## Tasks

## Notes

- Split on TH-0107: the Planner page is S-0259, which ships first.
