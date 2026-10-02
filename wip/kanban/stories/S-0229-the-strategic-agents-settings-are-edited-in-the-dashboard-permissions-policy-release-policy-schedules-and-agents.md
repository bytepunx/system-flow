---
id: S-0229
type: story
nature: feature
title: "The strategic agents' settings are edited in the dashboard: permissions, policy, release policy, schedules, and agents"
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:18Z
updated: 2026-10-02T11:54:45Z
transitions: []
tags: [dashboard, flai]
touches: [flaiover/src/routes, flai/internal/hostapi, flai/internal/manifest, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0218, S-0211, S-0223]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
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
