---
id: E-0016
type: epic
nature: feature
title: Introduce Orchestration, Planning, and Analysis agents
status: in-progress
owner: alex
created: 2026-10-02T10:51:49Z
updated: 2026-10-04T04:53:11Z
transitions:
  - to: ready
    at: 2026-10-03T06:23:24Z
    by: alex
  - to: in-progress
    at: 2026-10-03T06:23:28Z
    by: alex
tags: [dashboard, cli]
topics: [orchestration, planning, analysis]
usage:
  source: sum
  seconds: 59596
  models:
    - model: claude-haiku-4-5-20251001
      input: 2858
      output: 114564
      cache_read: 18338620
      cache_write: 1038611
      cost: 3.7079
    - model: claude-opus-5-5
      input: 7834
      output: 2482564
      cache_read: 514663492
      cache_write: 10721939
      cost: 219.9562
    - model: claude-sonnet-5
      input: 1426
      output: 362160
      cache_read: 49806550
      cache_write: 2695429
      cost: 20.3244
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10h
    by: alex
    at: 2026-10-04T04:42:30Z
  value: 1500
  by: planner-E-0016
  at: 2026-10-04T04:48:12Z
---
# E-0016 Introduce Orchestration, Planning, and Analysis agents

## Outcome

Three new agents get introduced to system-flow:

- The planner:
  - capable of taking an epic and creating stories based on the initial understanding of the epic
  - examines stories to determine and enrich them with:
    - cost of delay estimates (if key information is supplied by the operator)
    - estimated time to completion and estimated delivery date (real world wall clock estimate given agent/model/configuration and story complexity, past timing data)
    - an accurate list of files anticipated in touch to help improve parallelism with safety
- The orchestrator:
  - an agentic operator that manages projects with backlog and ready items
  - answers threads to keep work moving
  - decides when to publish releases
  - can optimize work through the system based on either maximized throughput or cost of delay
- The analyzer:
  - an agentic analysis tool that:
    - identifies bottlenecks
    - identifies gaps between implementation and intent
    - identifies technical and security risks
    - creates stories for improvements and remediations with cost of delay pre-calculated based on impact

## Stories
- S-0199 Work items carry planning data: draft, cost of delay inputs and value, and a forecast with who set it
- S-0200 An epic follows its stories: to ready and in-progress with the first, to review and done with the last
- S-0201 The story page shows [Draft] and finalizes a draft story, and cards mark drafts
- S-0202 A card's right-click menu offers Finalize on a draft story and the card's other actions
- S-0203 A story created from an issue is a draft and carries the issue's cost of delay inputs
- S-0204 The new-item and edit forms take cost of delay inputs in a collapsed panel
- S-0205 flai stats and metrics.md gain the planning, waiting, and strategic-agent metrics
- S-0206 Planner, orchestrator, and analyzer runs are logged in an activity document, with cost, seconds, and work completed in its front matter
- S-0207 The conventions and prime know the planner, orchestrator, and analyzer as roles
- S-0208 The planner is an agent flai serve starts for an epic or a story, behind the plan host action
- S-0209 The planner drafts an epic's stories into the backlog and revisits the children it already has
- S-0210 The planner enriches a story with predicted touches, a forecast, and a cost of delay value
- S-0211 Planning runs again when a story changes, when work ahead of it completes, and on a schedule
- S-0212 Charts compare forecasts and estimates with what happened
- S-0213 Charts show cost of delay outstanding, incurred, and what the pull order costs
- S-0214 Charts show parallelism, holds, and touches drift
- S-0215 A chart shows how long agents spend waiting on threads and on review
- S-0216 Strategic Cost and Strategic Use charts show what the planner, orchestrator, and analyzer add against delivery
- S-0217 flai exposes the orchestrator's deterministic operations as commands: ordering by policy, promotion candidates, and release evaluation
- S-0218 The orchestrator is a long-running agent per project behind the orchestrate host action, with permissions the operator sets and a decision log
- S-0219 The orchestrator moves and orders work by its policy within its permissions
- S-0220 The orchestrator answers threads, or recommends an answer, as its permission allows
- S-0221 The orchestrator accepts stories in review when permitted
- S-0222 The orchestrator publishes by the release policy: judgement, theme, or cost of delay threshold
- S-0223 The analyzer runs on demand or on a schedule and writes a report under design/analysis
- S-0224 The analyzer files its actionable findings as issues with their impact, and the issue step turns them into draft stories
- S-0225 The planner's cost is recorded on the epic or story it planned
- S-0226 The orchestrator's cost is recorded on the story or epic each decision concerned
- S-0227 The analyzer's cost is recorded on the issues it filed and the stories made from them
- S-0228 The Workflow menu has Orchestrator and Analyzer pages showing their status, activity log, and runs
- S-0229 The strategic agents' settings are edited in the dashboard: permissions, policy, release policy, schedules, and agents
- S-0255 The planner drafts a story's tasks into the backlog and revisits the children it already has
- S-0259 The Workflow menu has a Planner page showing its status, activity log, and runs

## Notes
