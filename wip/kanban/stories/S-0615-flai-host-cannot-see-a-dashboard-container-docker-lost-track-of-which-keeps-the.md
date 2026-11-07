---
id: S-0615
type: story
nature: remediation
title: flai host cannot see a dashboard container Docker lost track of, which keeps the published port, so every container it starts is not answering and it restarts the dashboard every five minutes
status: backlog
owner: alex
created: 2026-11-07T23:39:11Z
updated: 2026-11-07T23:39:11Z
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
    time_lost_per_cycle: 144h
    by: flai
    at: 2026-11-07T23:39:11Z
---
# S-0615 flai host cannot see a dashboard container Docker lost track of, which keeps the published port, so every container it starts is not answering and it restarts the dashboard every five minutes

## Goal

This story remediates [I-0243](../../../design/issues/I-0243-flai-host-cannot-see-a-dashboard-container-docker-lost-track-of-which-keeps-the.md), "flai host cannot see a dashboard container Docker lost track of, which keeps the published port, so every container it starts is not answering and it restarts the dashboard every five minutes". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0243 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0243 is closed with `flai issue close I-0243 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0243. time_lost_per_cycle 144h: 144h per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-31T18:32:12Z, 7.2 days before this story).
