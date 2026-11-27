---
id: S-0857
type: story
nature: improvement
title: A rebase stops on a Go test file where two stories each appended a test at the same place, which flai cannot settle, so an acceptance sends the story back
status: backlog
owner: alex
created: 2026-11-27T22:55:28Z
updated: 2026-11-27T22:55:28Z
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
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-11-27T22:55:28Z
---
# S-0857 A rebase stops on a Go test file where two stories each appended a test at the same place, which flai cannot settle, so an acceptance sends the story back

## Goal

This story remediates [I-0409](../../../design/issues/I-0409-a-rebase-stops-on-a-go-test-file-where-two-stories-each-appended-a-test-at-the.md), "A rebase stops on a Go test file where two stories each appended a test at the same place, which flai cannot settle, so an acceptance sends the story back". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0409 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0409 is closed with `flai issue close I-0409 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0409. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-21T18:16:18Z, 6.2 days before this story; under one cycle counts as one).
