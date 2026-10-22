---
id: S-0478
type: story
nature: improvement
title: flai components list, and the components examples check, take about 2.3s for flaiover since props.mjs resolves json props' types with typescript, up from 0.4s
status: ready
owner: alex
created: 2026-10-21T11:59:35Z
updated: 2026-10-22T20:55:35Z
transitions:
  - to: ready
    at: 2026-10-22T20:55:35Z
    by: alex
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 43
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 17
          output: 307
          cache_read: 3107006
          cache_write: 2859
          cost: 0.6867
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2s
    by: flai
    at: 2026-10-21T11:59:35Z
finalized:
  by: alex
  at: 2026-10-22T20:52:09Z
---
# S-0478 flai components list, and the components examples check, take about 2.3s for flaiover since props.mjs resolves json props' types with typescript, up from 0.4s

## Goal

This story remediates [I-0166](../../../design/issues/I-0166-flai-components-list-and-the-components-examples-check-take-about-2-3s-for-flaiover-since-props-mjs-resolves-json-props-types-with-typescript-up-from-0-4s.md), "flai components list, and the components examples check, take about 2.3s for flaiover since props.mjs resolves json props' types with typescript, up from 0.4s". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0166 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0166 is closed with `flai issue close I-0166 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0166. time_lost_per_cycle 2s: 2s per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-19T14:20:51Z, 1.9 days before this story; under one cycle counts as one).
