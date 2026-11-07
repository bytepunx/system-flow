---
id: S-0614
type: story
nature: remediation
title: The merge queue's send-back thread quotes go test's summary line with its tabs, the markdown lint refuses the thread (MD010), and the queue retries forever with the batch wedged
status: backlog
owner: alex
created: 2026-11-07T23:39:04Z
updated: 2026-11-07T23:39:04Z
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
    time_lost_per_cycle: 5h
    by: flai
    at: 2026-11-07T23:39:04Z
---
# S-0614 The merge queue's send-back thread quotes go test's summary line with its tabs, the markdown lint refuses the thread (MD010), and the queue retries forever with the batch wedged

## Goal

This story remediates [I-0242](../../../design/issues/I-0242-the-merge-queue-s-send-back-thread-quotes-go-test-s-summary-line-with-its-tabs.md), "The merge queue's send-back thread quotes go test's summary line with its tabs, the markdown lint refuses the thread (MD010), and the queue retries forever with the batch wedged". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0242 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0242 is closed with `flai issue close I-0242 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0242. time_lost_per_cycle 5h: 5h per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-31T16:43:40Z, 7.3 days before this story).
