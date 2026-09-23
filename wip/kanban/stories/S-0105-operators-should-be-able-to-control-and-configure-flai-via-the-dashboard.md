---
id: S-0105
type: story
nature: feature
title: Operators should be able to control and configure flai via the dashboard
status: review
parent: E-0008
owner: alex
created: 2026-09-23T20:21:50Z
updated: 2026-09-23T21:29:02Z
transitions:
  - to: ready
    at: 2026-09-23T20:21:55Z
    by: alex
  - to: in-progress
    at: 2026-09-23T20:36:06Z
    by: system-flow
  - to: review
    at: 2026-09-23T21:29:02Z
    by: system-flow
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd]
---
# S-0105 Operators should be able to control and configure flai via the dashboard

## Goal

There are a number of things that currently require the operator to have shell access on the host to enable or disable. These should be manageable by project via the dashboard.

## Acceptance criteria
- [x] operations that presently require an operator to enter commands in the shell are exposed per project in the dashboard

## Tasks
- T-0371 A settings host action, enabled only in a shell, lets a dashboard change the host's settings through flai
- T-0372 The dashboard's settings page shows and changes each project's host settings
- T-0373 Tried end to end from the dashboard; ADR and docs

## Notes
