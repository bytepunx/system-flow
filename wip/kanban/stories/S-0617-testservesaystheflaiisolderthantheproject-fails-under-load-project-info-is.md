---
id: S-0617
type: story
nature: remediation
title: "TestServeSaysTheFlaiIsOlderThanTheProject fails under load: project.info is asked before flai serve has worked out the outdated flai"
status: backlog
owner: alex
created: 2026-11-07T23:39:25Z
updated: 2026-11-07T23:39:25Z
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
    at: 2026-11-07T23:39:25Z
---
# S-0617 TestServeSaysTheFlaiIsOlderThanTheProject fails under load: project.info is asked before flai serve has worked out the outdated flai

## Goal

This story remediates [I-0252](../../../design/issues/I-0252-testservesaystheflaiisolderthantheproject-fails-under-load-project-info-is.md), "TestServeSaysTheFlaiIsOlderThanTheProject fails under load: project.info is asked before flai serve has worked out the outdated flai". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0252 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0252 is closed with `flai issue close I-0252 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0252. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-02T06:04:39Z, 5.7 days before this story; under one cycle counts as one).
