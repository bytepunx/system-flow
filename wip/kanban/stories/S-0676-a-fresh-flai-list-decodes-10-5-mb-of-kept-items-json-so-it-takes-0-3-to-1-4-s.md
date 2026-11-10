---
id: S-0676
type: story
nature: improvement
title: A fresh flai list decodes 10.5 MB of kept items JSON, so it takes 0.3 to 1.4 s, not the 0.3 s S-0600's title names
status: backlog
owner: alex
created: 2026-11-10T08:34:53Z
updated: 2026-11-10T08:34:53Z
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
    time_lost_per_cycle: 1s
    by: flai
    at: 2026-11-10T08:34:53Z
---
# S-0676 A fresh flai list decodes 10.5 MB of kept items JSON, so it takes 0.3 to 1.4 s, not the 0.3 s S-0600's title names

## Goal

This story remediates [I-0310](../../../design/issues/I-0310-a-fresh-flai-list-decodes-10-5-mb-of-kept-items-json-so-it-takes-0-3-to-1-4-s.md), "A fresh flai list decodes 10.5 MB of kept items JSON, so it takes 0.3 to 1.4 s, not the 0.3 s S-0600's title names". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0310 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0310 is closed with `flai issue close I-0310 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0310. time_lost_per_cycle 1s: 1s per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-08T01:17:08Z, 2.3 days before this story; under one cycle counts as one).
