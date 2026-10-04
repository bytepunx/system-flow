---
id: S-0211
type: story
nature: improvement
title: Planning runs again when a story changes, when work ahead of it completes, and on a schedule
status: ready
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:14Z
updated: 2026-10-04T04:48:13Z
transitions:
  - to: ready
    at: 2026-10-03T20:33:48Z
    by: alex
tags: [flai]
touches: [flai/internal/serve, flai/internal/hostapi, flai/internal/mcpserver, flai/internal/manifest, flaiover/src/routes/settings, design/system/strategic-agents.md, design/system/project-manifest.md, design/system/flai-cli.md, docs/operators/settings.md, docs/users/flai.md]
after: [S-0209, S-0210]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 85.37
  by: planner-E-0016
  at: 2026-10-04T04:48:13Z
forecast:
  duration: 1h45m
  delivery: 2026-10-04T14:30:00Z
  basis: "Event triggers, a new scheduler in flai serve (none exists), coalescing and a settings readout, like S-0206 (4753 s); serial after S-0210 and S-0217, which share flai/cmd, hostapi and mcpserver."
  by: planner-E-0016
  at: 2026-10-04T04:44:04Z
---
# S-0211 Planning runs again when a story changes, when work ahead of it completes, and on a schedule

## Goal

A forecast and a cost of delay value go stale: the story's words change, a story ahead of it is accepted or cancelled so its delivery date moves, or the queue is reordered. The planner should run again when it should, within what the operator allows.

## Acceptance criteria
- [ ] With `plan` enabled, `flai serve` queues a planner run for a story when its goal, criteria, touches, or cost of delay inputs change (an `edited` change) and its forecast or value is older than the edit
- [ ] When a story is accepted, cancelled, or reordered, the ready and backlog stories after it get their forecasts recomputed deterministically (`flai forecast`, no agent) and the planner is queued only when the manifest's `planning.replan` allows it (`never`, `deterministic`, `agent`)
- [ ] A schedule in the manifest (`planning.schedule`, a cron expression or `daily`) runs the planner over the ready column; runs are coalesced so one story is planned once per trigger window
- [ ] Each run is logged in the planner's activity document with what triggered it; the settings page shows the triggers
- [ ] Tests cover each trigger and the coalescing

## Tasks

## Notes
