---
id: S-0490
type: story
nature: remediation
title: flai test runs go test with -short, so a real-git test a task adds is skipped until the integration tier runs in the merge queue's batch
status: backlog
owner: alex
created: 2026-10-23T20:20:26Z
updated: 2026-10-23T20:20:26Z
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
    time_lost_per_cycle: 40m
    by: flai
    at: 2026-10-23T20:20:26Z
---
# S-0490 flai test runs go test with -short, so a real-git test a task adds is skipped until the integration tier runs in the merge queue's batch

## Goal

This story remediates [I-0201](../../../design/issues/I-0201-flai-test-runs-go-test-with-short-so-a-real-git-test-a-task-adds-is-skipped-until-the-integration-tier-runs-in-the-merge-queue-s-batch.md), "flai test runs go test with -short, so a real-git test a task adds is skipped until the integration tier runs in the merge queue's batch". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0201 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0201 is closed with `flai issue close I-0201 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0201. time_lost_per_cycle 40m: 20m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-23T18:27:27Z, 0.1 days before this story; under one cycle counts as one).
