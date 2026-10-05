---
id: S-0259
type: story
nature: feature
title: The Workflow menu has a Planner page showing its status, activity log, and runs
status: done
parent: E-0016
owner: alex
created: 2026-10-04T04:53:03Z
updated: 2026-10-05T00:34:46Z
transitions:
  - to: ready
    at: 2026-10-04T23:40:49Z
    by: alex
  - to: in-progress
    at: 2026-10-05T00:03:35Z
    by: agent-S-0259
  - to: review
    at: 2026-10-05T00:32:56Z
    by: agent-S-0259
  - to: done
    at: 2026-10-05T00:34:46Z
    by: alex
tags: [dashboard]
topics: [planning]
touches: [flaiover/src/routes/workflow, flaiover/src/routes/api, flaiover/src/lib/components, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/planner.ts, flaiover/src/lib/planner.test.ts, flaiover/src/lib/activity.ts, flaiover/src/lib/usage.ts, flaiover/src/lib/server/agent.ts, flai/internal/hostapi, flai/internal/serve/stream.go, flai/internal/serve/stream_test.go, flai/cmd/serve_actions.go, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0208, S-0225]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2058
  models:
    - model: claude-haiku-4-5-20251001
      input: 122
      output: 7630
      cache_read: 535485
      cache_write: 62936
      cost: 0.1705
    - model: claude-opus-5-5
      input: 340
      output: 122677
      cache_read: 16750933
      cache_write: 513701
      cost: 9.1066
    - model: claude-sonnet-5
      input: 62
      output: 12937
      cache_read: 1766922
      cache_write: 121772
      cost: 0.7873
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

- [x] `/workflow/planner` is in the Workflow tier of the site menu, and shows the `plan` host action's state with the command that enables it
- [x] The page shows the current run, if any, in the agent pane that story pages have
- [x] The page shows the entries of `wip/agents/planner.md`, newest first, each with its items, duration, and cost, and the document's front matter totals
- [x] The page lists the past runs with their cost and outcome
- [x] A form runs the planner on an epic or a story through the `plan` host action, and says why when the action is off
- [x] The page updates live from `/api/events` when `wip/agents/planner.md` changes
- [x] `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` describe the page, and tests cover its data and the form

## Tasks
- T-0836 flai serves the planner's activity document and its run's stream to the dashboard
- T-0837 The dashboard's API answers the planner's state, activity, and runs, and its run's stream
- T-0838 The Workflow menu's Planner page shows the planner's state, current run, activity, and runs, and runs it
- T-0839 The design and the user guide describe the Planner page

## Notes

- Split from S-0228 on TH-0107. S-0228 keeps the Orchestrator and Analyzer pages, and reuses this page's layout.
