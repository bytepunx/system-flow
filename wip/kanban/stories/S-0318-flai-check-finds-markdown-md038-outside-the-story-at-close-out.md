---
id: S-0318
type: story
nature: improvement
title: "flai check finds `markdown.MD038` outside the story at close-out"
status: backlog
owner: alex
created: 2026-10-07T18:59:49Z
updated: 2026-10-07T18:59:49Z
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
          input: 2
          output: 12
          cache_read: 64217
          cache_write: 5432
          cost: 0.0172
draft: true
---
# S-0318 flai check finds `markdown.MD038` outside the story at close-out

## Goal

This story remediates [I-0096](../../../design/issues/I-0096-flai-check-finds-markdown-md038-outside-the-story-at-close-out.md), "flai check finds `markdown.MD038` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0096 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0096 is closed with `flai issue close I-0096 --reason` saying what fixed it

## Tasks

## Notes
