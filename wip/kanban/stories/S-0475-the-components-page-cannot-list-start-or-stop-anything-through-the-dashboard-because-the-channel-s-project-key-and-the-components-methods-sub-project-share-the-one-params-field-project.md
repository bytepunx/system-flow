---
id: S-0475
type: story
nature: remediation
title: The Components page cannot list, start, or stop anything through the dashboard, because the channel's project key and the components methods' sub-project share the one params field project
status: backlog
owner: alex
created: 2026-10-21T11:59:09Z
updated: 2026-10-21T11:59:09Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    penalty_per_week: 1000
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-10-21T11:59:09Z
---
# S-0475 The Components page cannot list, start, or stop anything through the dashboard, because the channel's project key and the components methods' sub-project share the one params field project

## Goal

This story remediates [I-0163](../../../design/issues/I-0163-the-components-page-cannot-list-start-or-stop-anything-through-the-dashboard-because-the-channel-s-project-key-and-the-components-methods-sub-project-share-the-one-params-field-project.md), "The Components page cannot list, start, or stop anything through the dashboard, because the channel's project key and the components methods' sub-project share the one params field project". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0163 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0163 is closed with `flai issue close I-0163 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0163. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-18T02:25:28Z, 3.4 days before this story; under one cycle counts as one). penalty_per_week 1000 carried over from I-0163's Impact section.
