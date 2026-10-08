---
id: S-0325
type: story
nature: improvement
title: "flai check finds `narrative.state` outside the story at close-out"
status: cancelled
owner: alex
created: 2026-10-07T18:59:57Z
updated: 2026-10-08T04:25:40Z
transitions:
  - to: cancelled
    at: 2026-10-08T04:25:40Z
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
      seconds: 53
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 11
          output: 143
          cache_read: 3365716
          cache_write: 10641
          cost: 0.8319
draft: true
---
# S-0325 flai check finds `narrative.state` outside the story at close-out

## Goal

This story remediates [I-0111](../../../design/issues/I-0111-flai-check-finds-narrative-state-outside-the-story-at-close-out.md), "flai check finds `narrative.state` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0111 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0111 is closed with `flai issue close I-0111 --reason` saying what fixed it

## Tasks

## Notes
- 2026-10-08T04:25:40Z: moved to cancelled: duplicate of S-0323
