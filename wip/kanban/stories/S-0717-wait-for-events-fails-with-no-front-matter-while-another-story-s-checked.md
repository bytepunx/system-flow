---
id: S-0717
type: story
nature: remediation
title: wait_for_events fails with no front matter while another story's checked creation holds a PLACEHOLDER task file
status: backlog
owner: alex
created: 2026-11-11T01:43:38Z
updated: 2026-11-11T01:43:38Z
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
    at: 2026-11-11T01:43:38Z
---
# S-0717 wait_for_events fails with no front matter while another story's checked creation holds a PLACEHOLDER task file

## Goal

This story remediates [I-0345](../../../design/issues/I-0345-wait-for-events-fails-with-no-front-matter-while-another-story-s-checked.md), "wait_for_events fails with no front matter while another story's checked creation holds a PLACEHOLDER task file". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0345 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0345 is closed with `flai issue close I-0345 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0345. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T11:10:23Z, 0.6 days before this story; under one cycle counts as one).
