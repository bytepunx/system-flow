---
id: S-0329
type: story
nature: improvement
title: Display all time stamps across the site in local time
status: backlog
owner: alex
created: 2026-10-07T19:44:24Z
updated: 2026-10-07T19:44:24Z
transitions: []
tags: [dashboard]
touches: [flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    penalty_per_week: 0
    by: alex
    at: 2026-10-07T19:44:24Z
---
# S-0329 Display all time stamps across the site in local time

## Goal

Every view in the site currently displays time stamps in UTC. It's important that we _record_ timestamps in UTC, but they should be displayed in the user's local time zone.

## Acceptance criteria
- [ ] No times display in the dashboard in UTC
- [ ] All times display in the dashboard in the local time zone

## Tasks

## Notes
