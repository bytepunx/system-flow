---
id: S-0828
type: story
nature: improvement
title: board.get's undone.read runs one git diff per recorded path of each story's undone record, one after another, on every call
status: backlog
owner: alex
created: 2026-11-17T08:17:19Z
updated: 2026-11-17T08:17:19Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2m
    by: flai
    at: 2026-11-17T08:17:19Z
---
# S-0828 board.get's undone.read runs one git diff per recorded path of each story's undone record, one after another, on every call

## Goal

This story remediates [I-0349](../../../design/issues/I-0349-board-get-s-undone-read-runs-one-git-diff-per-recorded-path-of-each-story-s.md), "board.get's undone.read runs one git diff per recorded path of each story's undone record, one after another, on every call". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0349 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0349 is closed with `flai issue close I-0349 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0349. time_lost_per_cycle 2m: 2m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-11T08:52:49Z, 6 days before this story; under one cycle counts as one).
