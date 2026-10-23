---
id: S-0478
type: story
nature: improvement
title: flai components list, and the components examples check, take about 2.3s for flaiover since props.mjs resolves json props' types with typescript, up from 0.4s
status: done
owner: alex
created: 2026-10-21T11:59:35Z
updated: 2026-10-23T07:15:41Z
transitions:
  - to: ready
    at: 2026-10-22T20:55:35Z
    by: alex
  - to: in-progress
    at: 2026-10-23T06:44:07Z
    by: agent-S-0478
  - to: review
    at: 2026-10-23T07:03:21Z
    by: agent-S-0478
  - to: done
    at: 2026-10-23T07:15:41Z
    by: orchestrator
tags: []
touches: [design/system/component-library.md, docs/operators/runbooks/backup.md, flai/internal/components/props.go, flai/internal/components/props.mjs, flai/internal/components/props_test.go, design/issues/I-0166-flai-components-list-and-the-components-examples-check-take-about-2-3s-for-flaiover-since-props-mjs-resolves-json-props-types-with-typescript-up-from-0-4s.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1176
  priced_by: harness
  turns:
    - day: 2026-10-23
      ceremony: 1
      hand_edits: 2
      work: 52
  models:
    - model: claude-opus-5-5
      input: 112
      output: 44463
      cache_read: 8727573
      cache_write: 195249
      cost: 4.1972
  strategic:
    - kind: orchestrator
      seconds: 1953
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 181
          output: 3073
          cache_read: 66943086
          cache_write: 36181
          cost: 14.7894
        - model: claude-sonnet-5-5
          input: 12
          output: 92
          cache_read: 156247
          cache_write: 45180
          cost: 0.1394
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
- [x] The cause I-0166 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0166 is closed with `flai issue close I-0166 --reason` saying what fixed it

## Tasks
- T-2068 props.mjs keeps resolved skeletons in a cache stamped by every path typescript looked at, so an unchanged run skips typescript
- T-2069 I-0166 is closed as fixed by the skeleton cache

## Notes

Cost of delay inputs set by flai from I-0166. time_lost_per_cycle 2s: 2s per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-19T14:20:51Z, 1.9 days before this story; under one cycle counts as one).

### Accepted by the orchestrator

- Verified: fdf83cfcb3106377c3c34c5c51767680258647b1
- Batch: 20261023-070453
- Enqueued verified: 886afe032c4e492d3ae998f9ccdb3700cc684ab8
- At: 2026-10-23T07:15:41Z

Verdict: pass at 886afe032c4e492d3ae998f9ccdb3700cc684ab8
- 1: flai/internal/components/props.mjs, flai/internal/components/props.go, flai/internal/components/props_test.go, design/system/component-library.md, docs/operators/runbooks/backup.md
- 2: design/issues/I-0166-flai-components-list-and-the-components-examples-check-take-about-2-3s-for-flaiover-since-props-mjs-resolves-json-props-types-with-typescript-up-from-0-4s.md, design/issues/summary.md
