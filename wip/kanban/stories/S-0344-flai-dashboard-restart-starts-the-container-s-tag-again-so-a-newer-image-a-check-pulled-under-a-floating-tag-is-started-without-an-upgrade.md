---
id: S-0344
type: story
nature: remediation
title: flai dashboard restart starts the container's tag again, so a newer image a check pulled under a floating tag is started without an upgrade
status: backlog
owner: alex
created: 2026-10-08T08:08:17Z
updated: 2026-10-08T08:08:17Z
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
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-08T08:08:17Z
---
# S-0344 flai dashboard restart starts the container's tag again, so a newer image a check pulled under a floating tag is started without an upgrade

## Goal

This story remediates [I-0116](../../../design/issues/I-0116-flai-dashboard-restart-starts-the-container-s-tag-again-so-a-newer-image-a-check-pulled-under-a-floating-tag-is-started-without-an-upgrade.md), "flai dashboard restart starts the container's tag again, so a newer image a check pulled under a floating tag is started without an upgrade". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0116 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0116 is closed with `flai issue close I-0116 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0116. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-08T00:09:35Z, 0.3 days before this story; under one cycle counts as one).
