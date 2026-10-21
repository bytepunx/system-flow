---
id: S-0465
type: story
nature: improvement
title: flai move to done refuses unfixed lint as text only, so the dashboard picks each file's checkout from the refusal's words
status: backlog
owner: alex
created: 2026-10-21T07:24:59Z
updated: 2026-10-21T07:24:59Z
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
    time_lost_per_cycle: 9m
    by: flai
    at: 2026-10-21T07:24:59Z
---
# S-0465 flai move to done refuses unfixed lint as text only, so the dashboard picks each file's checkout from the refusal's words

## Goal

This story remediates [I-0145](../../../design/issues/I-0145-flai-move-to-done-refuses-unfixed-lint-as-text-only-so-the-dashboard-picks-each-file-s-checkout-from-the-refusal-s-words.md), "flai move to done refuses unfixed lint as text only, so the dashboard picks each file's checkout from the refusal's words". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0145 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0145 is closed with `flai issue close I-0145 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0145. time_lost_per_cycle 9m: 10m per occurrence × 1 occurrence ÷ 1.1 cycles of 168h (first reported 2026-10-13T18:47:09Z, 7.5 days before this story).
