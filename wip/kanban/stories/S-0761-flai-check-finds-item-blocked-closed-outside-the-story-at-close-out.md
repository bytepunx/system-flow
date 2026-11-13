---
id: S-0761
type: story
nature: improvement
title: "flai check finds `item.blocked-closed` outside the story at close-out"
status: backlog
owner: alex
created: 2026-11-13T07:23:41Z
updated: 2026-11-13T07:27:32Z
transitions: []
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
      seconds: 4
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 8
          cache_read: 32556
          cache_write: 104
          cost: 0.0094
finalized:
  by: alex
  at: 2026-11-13T07:27:32Z
---
# S-0761 flai check finds `item.blocked-closed` outside the story at close-out

## Goal

This story remediates [I-0354](../../../design/issues/I-0354-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md), "flai check finds `item.blocked-closed` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0354 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0354 is closed with `flai issue close I-0354 --reason` saying what fixed it

## Tasks

## Notes
