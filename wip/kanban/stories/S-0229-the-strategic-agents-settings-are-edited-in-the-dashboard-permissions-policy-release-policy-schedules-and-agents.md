---
id: S-0229
type: story
nature: feature
title: "The strategic agents' settings are edited in the dashboard: permissions, policy, release policy, schedules, and agents"
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:18Z
updated: 2026-10-05T03:10:20Z
transitions: []
tags: [dashboard, flai]
touches: [flaiover/src/routes, flai/internal/hostapi, flai/internal/manifest, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/routes/settings, flaiover/src/lib/components, design/system/project-manifest.md, docs/operators/settings.md]
after: [S-0218, S-0211, S-0223]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 97.56
  by: planner-E-0016
  at: 2026-10-04T04:48:24Z
forecast:
  duration: 2h
  delivery: 2026-10-05T21:25:00Z
  basis: "Its own forecast of 2h; 23rd in the pull order with an in-progress limit of 3, behind S-0253, S-0257, S-0244, S-0262, S-0258, S-0260, S-0212, S-0213, S-0214, S-0215, S-0216, S-0217, S-0218, S-0219, S-0220, S-0221, S-0222, S-0223, S-0224, S-0226, S-0227 and S-0228."
  by: flai
  at: 2026-10-05T03:10:20Z
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

## Notes
