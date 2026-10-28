---
id: S-0514
type: story
nature: remediation
title: A curator run that ends on the harness's usage limit is recorded failed, logs no activity, and is not started again until the next acceptance or schedule
status: backlog
owner: alex
created: 2026-10-28T12:40:16Z
updated: 2026-10-28T12:40:16Z
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
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-10-28T12:40:16Z
---
# S-0514 A curator run that ends on the harness's usage limit is recorded failed, logs no activity, and is not started again until the next acceptance or schedule

## Goal

This story remediates [I-0212](../../../design/issues/I-0212-a-curator-run-that-ends-on-the-harness-s-usage-limit-is-recorded-failed-logs-no-activity-and-is-not-started-again-until-the-next-acceptance-or-schedule.md), "A curator run that ends on the harness's usage limit is recorded failed, logs no activity, and is not started again until the next acceptance or schedule". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0212 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0212 is closed with `flai issue close I-0212 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0212. time_lost_per_cycle 30m: 30m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-28T12:03:27Z, 0 days before this story; under one cycle counts as one).
