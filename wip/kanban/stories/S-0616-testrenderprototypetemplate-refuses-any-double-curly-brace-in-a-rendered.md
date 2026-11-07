---
id: S-0616
type: story
nature: remediation
title: TestRenderPrototypeTemplate refuses any double curly brace in a rendered project, so a baseline formula the template ships cannot use an input
status: backlog
owner: alex
created: 2026-11-07T23:39:18Z
updated: 2026-11-07T23:39:18Z
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
    at: 2026-11-07T23:39:18Z
---
# S-0616 TestRenderPrototypeTemplate refuses any double curly brace in a rendered project, so a baseline formula the template ships cannot use an input

## Goal

This story remediates [I-0249](../../../design/issues/I-0249-testrenderprototypetemplate-refuses-any-double-curly-brace-in-a-rendered.md), "TestRenderPrototypeTemplate refuses any double curly brace in a rendered project, so a baseline formula the template ships cannot use an input". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0249 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0249 is closed with `flai issue close I-0249 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0249. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-02T02:50:10Z, 5.9 days before this story; under one cycle counts as one).
