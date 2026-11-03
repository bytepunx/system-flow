---
id: S-0573
type: story
nature: remediation
title: Stories in new projects ignore default agent/harness settings
status: backlog
owner: alex
created: 2026-11-03T17:24:00Z
updated: 2026-11-03T17:24:00Z
transitions: []
tags: [cli]
touches: [flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    penalty_per_week: 50
    by: alex
    at: 2026-11-03T17:24:00Z
---
# S-0573 Stories in new projects ignore default agent/harness settings

## Goal

In new projects, after an agent has been set in system-flow.yaml and the harness and default agent are set in the configuration screens, new stories fail to start with the error that there's missing agent/harness configuration.

Diagnose the cause and correct this behavior. All configuration should be set via the dashboard without additional steps in the CLI required.

## Acceptance criteria
- [ ] In new projects, where the harness and agent settings are configured, there are no additional steps required (i.e. editing the story's agent directly)
- [ ] Stories in projects where the harness and agent configuration are provided start successfully without additional CLI commands or file edits required

## Tasks

## Notes
