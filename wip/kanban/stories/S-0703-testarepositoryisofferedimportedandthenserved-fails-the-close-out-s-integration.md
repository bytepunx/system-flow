---
id: S-0703
type: story
nature: remediation
title: "TestARepositoryIsOfferedImportedAndThenServed fails the close-out's integration tier with \"failed to read frame header: EOF\" under host load, on a branch that does not touch flai/internal/serve"
status: backlog
owner: alex
created: 2026-11-10T08:55:40Z
updated: 2026-11-10T08:55:40Z
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
    time_lost_per_cycle: 9m
    by: flai
    at: 2026-11-10T08:55:40Z
---
# S-0703 TestARepositoryIsOfferedImportedAndThenServed fails the close-out's integration tier with "failed to read frame header: EOF" under host load, on a branch that does not touch flai/internal/serve

## Goal

This story remediates [I-0267](../../../design/issues/I-0267-testarepositoryisofferedimportedandthenserved-fails-the-close-out-s-integration.md), "TestARepositoryIsOfferedImportedAndThenServed fails the close-out's integration tier with "failed to read frame header: EOF" under host load, on a branch that does not touch flai/internal/serve". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0267 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0267 is closed with `flai issue close I-0267 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0267. time_lost_per_cycle 9m: 9m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T02:16:58Z, 3.3 days before this story; under one cycle counts as one).
