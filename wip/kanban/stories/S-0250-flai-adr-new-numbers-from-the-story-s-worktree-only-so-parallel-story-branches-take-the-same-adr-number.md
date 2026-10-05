---
id: S-0250
type: story
nature: improvement
title: flai adr new numbers from the story's worktree only, so parallel story branches take the same ADR number
status: backlog
owner: alex
created: 2026-10-03T18:03:59Z
updated: 2026-10-03T18:03:59Z
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
    - kind: planner
      seconds: 145
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 624
          output: 126
          cache_read: 4433381
          cache_write: 227929
          cost: 0.9225
        - model: claude-opus-5-5
          input: 497
          output: 16195
          cache_read: 16618464
          cache_write: 470742
          cost: 8.9524
---
# S-0250 flai adr new numbers from the story's worktree only, so parallel story branches take the same ADR number

## Goal

This story remediates [I-0062](../../../design/issues/I-0062-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md), "flai adr new numbers from the story's worktree only, so parallel story branches take the same ADR number". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0062 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0062 is closed with `flai issue close I-0062 --reason` saying what fixed it

## Tasks

## Notes
