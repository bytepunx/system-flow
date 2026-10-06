---
id: S-0298
type: story
nature: improvement
title: The Host Updates Tab Allows an Operator to Rollback to A Previous Version
status: backlog
owner: alex
created: 2026-10-06T20:53:15Z
updated: 2026-10-06T20:53:15Z
transitions: []
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: alex
    at: 2026-10-06T20:53:15Z
---
# S-0298 The Host Updates Tab Allows an Operator to Rollback to A Previous Version

## Goal

Presently, the host Updates page allows the operator to deploy new versions. There isn't presently any tooling to allow an operator to roll back to a prior or specific version.

Introduce flai commands for the cli, http, mcp protocols so that the operator can discover available older versions and deploy one for both the flai CLI and the dashboard.

## Acceptance criteria
- [ ] New commands are available to allow the operator to view prior versions
- [ ] New commands are available to allow the operator to select a specific version of the CLI and dashboard to deploy
- [ ] Deploying a specific version works like updating to the latest

## Tasks

## Notes
