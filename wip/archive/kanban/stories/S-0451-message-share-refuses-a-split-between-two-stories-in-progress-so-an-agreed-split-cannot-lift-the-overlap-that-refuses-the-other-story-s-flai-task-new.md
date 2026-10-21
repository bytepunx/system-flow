---
id: S-0451
type: story
nature: improvement
title: message_share refuses a split between two stories in progress, so an agreed split cannot lift the overlap that refuses the other story's flai task new
status: done
owner: alex
created: 2026-10-20T18:59:50Z
updated: 2026-10-21T02:08:39Z
transitions:
  - to: ready
    at: 2026-10-20T19:01:47Z
    by: alex
  - to: in-progress
    at: 2026-10-20T19:02:10Z
    by: agent-S-0451
  - to: review
    at: 2026-10-20T22:31:30Z
    by: agent-S-0451
  - to: done
    at: 2026-10-21T02:08:39Z
    by: alex
tags: []
touches: [flai/internal/itemnew/itemnew.go, flai/internal/itemnew/itemnew_test.go, design/issues/I-0169-message-share-refuses-a-split-between-two-stories-in-progress-so-an-agreed-split-cannot-lift-the-overlap-that-refuses-the-other-story-s-flai-task-new.md, design/issues/summary.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0170-the-litellm-gateway-smoke-check-s-default-free-model-is-refused-by-claude-code-as-unrecognized-model-so-every-close-out-on-this-host-fails-smoke.md, docs/contributors/index.md, docs/operators/index.md, flai/internal/serve/gateway_test.go, scripts/litellm.sh]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1886
  priced_by: harness
  turns:
    - day: 2026-10-20
      ceremony: 1
      test_runs: 1
      hand_edits: 2
      work: 92
  models:
    - model: claude-opus-5-5
      input: 220
      output: 44925
      cache_read: 13582808
      cache_write: 538019
      cost: 7.7707
  strategic:
    - kind: orchestrator
      seconds: 1211
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 125
          output: 2115
          cache_read: 37128967
          cache_write: 46062
          cost: 9.3195
        - model: claude-sonnet-5-5
          input: 10
          output: 74
          cache_read: 108657
          cache_write: 41372
          cost: 0.1228
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-20T18:59:50Z
finalized:
  by: alex
  at: 2026-10-20T19:01:44Z
---
# S-0451 message_share refuses a split between two stories in progress, so an agreed split cannot lift the overlap that refuses the other story's flai task new

## Goal

This story remediates [I-0169](../../../design/issues/I-0169-message-share-refuses-a-split-between-two-stories-in-progress-so-an-agreed-split-cannot-lift-the-overlap-that-refuses-the-other-story-s-flai-task-new.md), "message_share refuses a split between two stories in progress, so an agreed split cannot lift the overlap that refuses the other story's flai task new". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0169 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0169 is closed with `flai issue close I-0169 --reason` saying what fixed it

## Tasks
- T-1912 itemnew keeps a new wip.overlap as a warning, so flai task new and item_new no longer refuse a task whose touches reach another story in progress
- T-1913 The CLI design and the user docs say a new item keeps a wip.overlap as a warning, and I-0169 is closed
- T-1928 The LiteLLM gateway check's default free model is nemotron-3-ultra, since OpenRouter retired nemotron-3-super

## Notes

Cost of delay inputs set by flai from I-0169. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-20T00:34:27Z, 0.8 days before this story; under one cycle counts as one).
