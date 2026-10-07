---
id: S-0328
type: story
nature: feature
title: Add Permission and Ability to Orchestrate to Trigger Planner
status: backlog
owner: alex
created: 2026-10-07T19:35:32Z
updated: 2026-10-07T19:35:32Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    penalty_per_week: 800
    by: alex
    at: 2026-10-07T19:35:32Z
---
# S-0328 Add Permission and Ability to Orchestrate to Trigger Planner

## Goal

The Orchestrator needs the ability to dispatch the planner on stories in the backlog that are missing a plan. The orchestrator settings should introduce a permission called `plan_backlog_stories` to enable this behavior.

## Acceptance criteria
- [ ] The orchestrator settings add a new permission called `plan_backlog_stories` that, when enabled, gives it the permission to delegate planning to the planner agent.
- [ ] The orchestrator will identify stories in the backlog that are missing a plan (touches, CoD, and tasks) and, when enabled, will delegate planning of the story to the planning agent.
- [ ] The orchestrator will monitor threads from the planner and approve them, answer outstanding questions, and choose a CoD or accept the planner's recommendation.

## Tasks

## Notes
