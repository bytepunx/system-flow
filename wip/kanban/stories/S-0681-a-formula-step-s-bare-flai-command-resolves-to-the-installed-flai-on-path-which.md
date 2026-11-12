---
id: S-0681
type: story
nature: improvement
title: A formula step's bare flai command resolves to the installed flai on PATH, which can be older than the flags the formula's steps use
status: review
owner: alex
created: 2026-11-10T08:37:05Z
updated: 2026-11-12T00:56:30Z
transitions:
  - to: ready
    at: 2026-11-10T10:58:14Z
    by: alex
  - to: in-progress
    at: 2026-11-12T00:41:15Z
    by: agent-S-0681
  - to: review
    at: 2026-11-12T00:56:30Z
    by: agent-S-0681
tags: []
touches: [design/system/project-manifest.md, design/system/strategic-agents.md, docs/users/flai.md, flai/cmd/formula.go, flai/cmd/formula_test.go, design/issues/I-0332-a-formula-step-s-bare-flai-command-resolves-to-the-installed-flai-on-path-which.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 935
  priced_by: harness
  turns:
    - day: 2026-11-12
      test_runs: 1
      work: 36
  models:
    - model: claude-opus-5-5
      input: 78
      output: 15833
      cache_read: 6146088
      cache_write: 193457
      cost: 3.0938
  strategic:
    - kind: orchestrator
      seconds: 123
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 47
          output: 325
          cache_read: 1778020
          cache_write: 223067
          cost: 0.5648
cost_of_delay:
  inputs:
    time_lost_per_cycle: 3m
    by: flai
    at: 2026-11-10T08:37:05Z
finalized:
  by: alex
  at: 2026-11-10T10:58:10Z
---
# S-0681 A formula step's bare flai command resolves to the installed flai on PATH, which can be older than the flags the formula's steps use

## Goal

This story remediates [I-0332](../../../design/issues/I-0332-a-formula-step-s-bare-flai-command-resolves-to-the-installed-flai-on-path-which.md), "A formula step's bare flai command resolves to the installed flai on PATH, which can be older than the flags the formula's steps use". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0332 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0332 is closed with `flai issue close I-0332 --reason` saying what fixed it

## Tasks
- T-3130 flai formula run writes a flai step's command with the manifest's flai.command when it names one
- T-3131 I-0332 is closed with the reason that a formula step's task names the manifest's flai.command

## Notes

Cost of delay inputs set by flai from I-0332. time_lost_per_cycle 3m: 3m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T08:28:53Z, 0 days before this story; under one cycle counts as one).
