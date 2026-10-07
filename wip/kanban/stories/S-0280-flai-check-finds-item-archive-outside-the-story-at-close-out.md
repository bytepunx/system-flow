---
id: S-0280
type: story
nature: improvement
title: "flai check finds `item.archive` outside the story at close-out"
status: backlog
owner: alex
created: 2026-10-05T07:09:05Z
updated: 2026-10-07T15:03:46Z
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
      seconds: 2
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 6
          cache_read: 24135
          cache_write: 161
          cost: 0.0064
draft: true
---
# S-0280 flai check finds `item.archive` outside the story at close-out

## Goal

This story remediates [I-0078](../../../design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md), "flai check finds `item.archive` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0078 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0078 is closed with `flai issue close I-0078 --reason` saying what fixed it

## Tasks
- T-1164 An ADR refining ADR-0085 and ADR-0115 records the remedy for I-0078, proposed from its instances
- T-1165 A close-out records no item.archive in an issue, and its scoped check leaves out one that does not name the story
- T-1166 I-0078 is closed with flai issue close, saying that a close-out records no item.archive

## Notes
