---
id: S-0653
type: story
nature: remediation
title: TestAFormulaRoleStepsEndIsToldToTheInstancesAgent fails on its timing under load, in the go-test tier of any story that changes a package serve imports
status: backlog
owner: alex
created: 2026-11-10T05:45:21Z
updated: 2026-11-10T05:45:21Z
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
    time_lost_per_cycle: 12m
    by: flai
    at: 2026-11-10T05:45:21Z
---
# S-0653 TestAFormulaRoleStepsEndIsToldToTheInstancesAgent fails on its timing under load, in the go-test tier of any story that changes a package serve imports

## Goal

This story remediates [I-0314](../../../design/issues/I-0314-testaformularolestepsendistoldtotheinstancesagent-fails-on-its-timing-under.md), "TestAFormulaRoleStepsEndIsToldToTheInstancesAgent fails on its timing under load, in the go-test tier of any story that changes a package serve imports". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0314 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0314 is closed with `flai issue close I-0314 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0314. time_lost_per_cycle 12m: 6m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-11-08T02:50:07Z, 2.1 days before this story; under one cycle counts as one).
