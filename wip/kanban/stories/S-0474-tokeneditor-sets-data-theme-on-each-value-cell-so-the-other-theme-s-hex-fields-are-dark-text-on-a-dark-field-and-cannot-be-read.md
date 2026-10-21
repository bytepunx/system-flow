---
id: S-0474
type: story
nature: remediation
title: TokenEditor sets data-theme on each value cell, so the other theme's hex fields are dark text on a dark field and cannot be read
status: backlog
owner: alex
created: 2026-10-21T11:59:01Z
updated: 2026-10-21T12:04:23Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-21T11:59:01Z
finalized:
  by: alex
  at: 2026-10-21T12:04:23Z
---
# S-0474 TokenEditor sets data-theme on each value cell, so the other theme's hex fields are dark text on a dark field and cannot be read

## Goal

This story remediates [I-0159](../../../design/issues/I-0159-tokeneditor-sets-data-theme-on-each-value-cell-so-the-other-theme-s-hex-fields-are-dark-text-on-a-dark-field-and-cannot-be-read.md), "TokenEditor sets data-theme on each value cell, so the other theme's hex fields are dark text on a dark field and cannot be read". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0159 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0159 is closed with `flai issue close I-0159 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0159. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-17T05:32:26Z, 4.3 days before this story; under one cycle counts as one).
