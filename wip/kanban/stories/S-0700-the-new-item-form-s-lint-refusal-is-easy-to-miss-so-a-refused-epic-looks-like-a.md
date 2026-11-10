---
id: S-0700
type: story
nature: remediation
title: The new item form's lint refusal is easy to miss, so a refused epic looks like a hang
status: backlog
owner: alex
created: 2026-11-10T08:55:31Z
updated: 2026-11-10T08:55:31Z
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
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-11-10T08:55:31Z
---
# S-0700 The new item form's lint refusal is easy to miss, so a refused epic looks like a hang

## Goal

This story remediates [I-0263](../../../design/issues/I-0263-the-new-item-form-s-lint-refusal-is-easy-to-miss-so-a-refused-epic-looks-like-a.md), "The new item form's lint refusal is easy to miss, so a refused epic looks like a hang". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0263 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0263 is closed with `flai issue close I-0263 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0263. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-06T21:31:26Z, 3.5 days before this story; under one cycle counts as one).
