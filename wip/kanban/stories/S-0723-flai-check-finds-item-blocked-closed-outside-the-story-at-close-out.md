---
id: S-0723
type: story
nature: improvement
title: "flai check finds `item.blocked-closed` outside the story at close-out"
status: backlog
owner: alex
created: 2026-11-11T07:23:14Z
updated: 2026-11-12T07:06:11Z
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
      seconds: 7
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 10
          cache_read: 62150
          cache_write: 9071
          cost: 0.0202
finalized:
  by: alex
  at: 2026-11-12T07:06:11Z
---
# S-0723 flai check finds `item.blocked-closed` outside the story at close-out

## Goal

This story remediates [I-0347](../../../design/issues/I-0347-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md), "flai check finds `item.blocked-closed` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0347 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0347 is closed with `flai issue close I-0347 --reason` saying what fixed it

## Tasks

## Notes
