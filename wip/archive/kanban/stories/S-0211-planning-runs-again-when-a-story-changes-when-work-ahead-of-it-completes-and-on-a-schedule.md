---
id: S-0211
type: story
nature: improvement
title: Planning runs again when a story changes, when work ahead of it completes, and on a schedule
status: done
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:14Z
updated: 2026-10-04T21:23:05Z
transitions:
  - to: ready
    at: 2026-10-03T20:33:48Z
    by: alex
  - to: in-progress
    at: 2026-10-04T20:35:20Z
    by: agent-S-0211
  - to: review
    at: 2026-10-04T21:22:55Z
    by: agent-S-0211
  - to: done
    at: 2026-10-04T21:23:05Z
    by: alex
tags: [flai]
touches: [flai/internal/serve, flai/internal/hostapi, flai/internal/mcpserver, flai/internal/manifest, flai/internal/cron, flai/internal/planning, flai/internal/workitem/activity.go, flai/internal/workitem/activity_test.go, flai/internal/workitem/changes.go, flai/internal/workitem/changes_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/cmd/touches.go, flai/cmd/touches_test.go, flaiover/src/routes/settings, flaiover/src/lib/settings.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, design/adrs/0084-flai-serve-plans-again-on-its-own-behind-the-plan-host-action-on-an-edit-when.md, design/adrs/README.md, design/issues/I-0068-the-mcp-prime-pack-for-a-story-is-larger-than-claude-code-s-tool-result-limit-so-the-agent-reads-it-back-from-a-saved-file.md, design/issues/summary.md, design/system/agent-narrative.md, design/system/strategic-agents.md, design/system/project-manifest.md, design/system/flai-cli.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flaiover.md, design/system/flaiover-dashboard.md]
after: [S-0209, S-0210]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2890
  models:
    - model: claude-haiku-4-5-20251001
      input: 274
      output: 10001
      cache_read: 1467818
      cache_write: 65985
      cost: 0.2795
    - model: claude-opus-5-5
      input: 560
      output: 218534
      cache_read: 35096714
      cache_write: 887014
      cost: 16.7118
    - model: claude-sonnet-5
      input: 84
      output: 16220
      cache_read: 3709656
      cache_write: 134410
      cost: 1.2403
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
- [x] With `plan` enabled, `flai serve` queues a planner run for a story when its goal, criteria, touches, or cost of delay inputs change (an `edited` change) and its forecast or value is older than the edit
- [x] When a story is accepted, cancelled, or reordered, the ready and backlog stories after it get their forecasts recomputed deterministically (`flai forecast`, no agent) and the planner is queued only when the manifest's `planning.replan` allows it (`never`, `deterministic`, `agent`)
- [x] A schedule in the manifest (`planning.schedule`, a cron expression or `daily`) runs the planner over the ready column; runs are coalesced so one story is planned once per trigger window
- [x] Each run is logged in the planner's activity document with what triggered it; the settings page shows the triggers
- [x] Tests cover each trigger and the coalescing

## Tasks
- T-0820 flai parses a five-field cron expression or daily and finds its next time
- T-0821 planning replays a story's delivery with its own forecast duration
- T-0822 A planner run records its trigger, and its activity entry says it
- T-0823 The manifest takes planning.replan and planning.schedule, and flai check validates them
- T-0824 flai serve queues the planner on an edit, replans when work ahead completes or the order changes, and runs the schedule, coalescing runs
- T-0825 The settings page shows the planning triggers
- T-0826 The design and the user guides describe when the planner runs again

## Notes

Built as ADR-0084 decides; `design/system/strategic-agents.md` § Planning again describes it.

- Every trigger needs the `plan` host action on. An edit replans only a story already planned: one with a forecast or a cost of delay value. A cost of delay edit counts when its inputs changed after the value. A planner's own edits never trigger.
- `planning.replan`: `never` does nothing; `deterministic`, the default, replays the forecasts of ready and backlog stories with `planning.Replay`, keeping each story's own duration, and writes the deliveries that moved in one commit by `flai`; `agent` also queues the planner for those stories.
- `planning.schedule` is a five-field cron in UTC or `daily`, parsed by `flai/internal/cron`.
- The queue holds a story once, with its triggers merged, and runs one planner at a time per project. Each run records its triggers in `serve/agents.json` and in a `- Trigger:` line on its `planner.md` entry.
- `flai touches` now leaves an edit notice, so a touches change triggers a replan as `flai edit`'s does.
- Tests: `flai/internal/serve/replan_test.go` (each trigger, the policy, the gate, the queue), `flai/internal/cron`, `flai/internal/planning/replay_test.go`, `flai/internal/manifest`, `flai/cmd` (settings view, touches notice), and the Settings panel's component test.
