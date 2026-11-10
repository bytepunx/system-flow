---
id: S-0698
type: story
nature: improvement
title: This repository's board.get takes 300 ms or more in half its calls with 13 boards at once, from the merge queue's git while stories are in review
status: backlog
owner: alex
created: 2026-11-10T08:54:37Z
updated: 2026-11-10T08:54:37Z
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
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-10T08:54:37Z
---
# S-0698 This repository's board.get takes 300 ms or more in half its calls with 13 boards at once, from the merge queue's git while stories are in review

## Goal

This story remediates [I-0333](../../../design/issues/I-0333-this-repository-s-board-get-takes-300-ms-or-more-in-half-its-calls-with-13.md), "This repository's board.get takes 300 ms or more in half its calls with 13 boards at once, from the merge queue's git while stories are in review". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0333 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0333 is closed with `flai issue close I-0333 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0333. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T08:41:34Z, 0 days before this story; under one cycle counts as one).
