---
id: S-0781
type: story
nature: improvement
title: "flai check finds `story.unaccepted` outside the story at close-out"
status: backlog
owner: alex
created: 2026-11-13T22:46:13Z
updated: 2026-11-13T22:46:13Z
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
      seconds: 17
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 38
          cache_read: 183457
          cache_write: 22473
          cost: 0.0607
draft: true
---
# S-0781 flai check finds `story.unaccepted` outside the story at close-out

## Goal

This story remediates [I-0371](../../../design/issues/I-0371-flai-check-finds-story-unaccepted-outside-the-story-at-close-out.md), "flai check finds `story.unaccepted` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0371 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0371 is closed with `flai issue close I-0371 --reason` saying what fixed it

## Tasks

## Notes
