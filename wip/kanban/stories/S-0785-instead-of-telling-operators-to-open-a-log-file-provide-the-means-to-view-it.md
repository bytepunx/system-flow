---
id: S-0785
type: story
nature: feature
title: Instead of telling operators to open a log file, provide the means to view it
status: backlog
owner: alex
created: 2026-11-14T10:23:44Z
updated: 2026-11-14T10:23:44Z
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
    time_lost_per_cycle: 5m
    by: alex
    at: 2026-11-14T10:23:44Z
---
# S-0785 Instead of telling operators to open a log file, provide the means to view it

## Goal

Right now errors in several areas (agents status, deploy status, issues with review, etc.) will instruct the operator to view the contents of a log file. The issue is that this requires the operator to have shell access to the host machine, which is not always the case.

These log paths should work as hyperlinks that can be vied within flaiover instead.

## Acceptance criteria
- [ ] No status/error messages should present a path to a log file as a plain string
- [ ] There should be a log viewer component that flai displays log contents in
- [ ] The path to the log file in error and status messages is displayed as a hyperlink that takes the operator to page with the log viewer set to the log file's contents

## Tasks

## Notes
