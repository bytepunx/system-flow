---
id: S-0714
type: story
nature: remediation
title: "TestEveryListedFlaioverComponentCompilesInTheHarness fails a close-out with \"the setup src/routes/layout.css does not compile: 404, 404 page not found\" while the main checkout's flaiover node_modules is intact"
status: backlog
owner: alex
created: 2026-11-10T10:54:41Z
updated: 2026-11-10T10:54:41Z
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
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-10T10:54:41Z
---
# S-0714 TestEveryListedFlaioverComponentCompilesInTheHarness fails a close-out with "the setup src/routes/layout.css does not compile: 404, 404 page not found" while the main checkout's flaiover node_modules is intact

## Goal

This story remediates [I-0342](../../../design/issues/I-0342-testeverylistedflaiovercomponentcompilesintheharness-fails-a-close-out-with-the.md), "TestEveryListedFlaioverComponentCompilesInTheHarness fails a close-out with "the setup src/routes/layout.css does not compile: 404, 404 page not found" while the main checkout's flaiover node_modules is intact". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0342 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0342 is closed with `flai issue close I-0342 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0342. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T10:40:56Z, 0 days before this story; under one cycle counts as one).
