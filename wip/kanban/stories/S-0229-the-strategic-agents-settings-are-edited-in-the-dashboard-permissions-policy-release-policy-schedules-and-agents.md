---
id: S-0229
type: story
nature: feature
title: "The strategic agents' settings are edited in the dashboard: permissions, policy, release policy, schedules, and agents"
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:18Z
updated: 2026-10-05T05:46:14Z
transitions: []
tags: [dashboard, flai]
touches: [flaiover/src/routes, flai/internal/hostapi, flai/internal/manifest, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/routes/settings, flaiover/src/lib/components, design/system/project-manifest.md, docs/operators/settings.md, flai/cmd/manifest.go, flai/cmd/manifest_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/guard, flai/internal/harness, flaiover/src/lib/settings.ts, flaiover/src/lib/server/agent.ts, design/system/strategic-agents.md, design/system/flai-cli.md, design/adrs, docs/operators/index.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [S-0218, S-0211, S-0223]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 84.52
  by: planner-S-0229
  at: 2026-10-05T05:45:18Z
forecast:
  duration: 1h15m
  delivery: 2026-10-05T21:25:00Z
  basis: "flai forecast's 1h (132 s per unit times size 27) raised to 1h15m, the agent time of S-0204 (72m) and S-0201 (63m), the closest done stories that likewise add a flai write and its dashboard form; 18th in the pull order with an in-progress limit of 3, delivered at its start plus 1h15m times the cycle factor 1.52."
  by: planner-S-0229
  at: 2026-10-05T05:45:12Z
---
# S-0229 The strategic agents' settings are edited in the dashboard: permissions, policy, release policy, schedules, and agents

## Goal

The operator decides how much freedom each strategic agent has. The manifest holds the settings; the dashboard should edit them where the agents' pages are, behind the `settings` host action, with each setting explained.

## Acceptance criteria
- [ ] The orchestrator page has a settings panel for `orchestration.permissions` (each with a sentence on what it allows and its risk), `orchestration.policy`, and `orchestration.release` (kind, thresholds, theme, whole epics); the planner page for `planning` (agent, replan, schedule, hour rate, cycle, default duration); the analyzer page for `analysis` (agent, schedule)
- [ ] Saving writes the manifest through a hostapi write (`settings.manifest`) that validates and refuses with the field and reason; a running orchestrator picks up changed permissions on its next decision
- [ ] The panels are read-only with the reason when the `settings` host action is off
- [ ] `design/system/flaiover-dashboard.md`, `project-manifest.md`, and the operator guide describe every setting; tests cover a save, a refusal, and read-only

## Tasks
- T-0952 flai manifest set writes the strategic agents' settings to system-flow.yaml, validating each and refusing with the field and the reason

## Notes
