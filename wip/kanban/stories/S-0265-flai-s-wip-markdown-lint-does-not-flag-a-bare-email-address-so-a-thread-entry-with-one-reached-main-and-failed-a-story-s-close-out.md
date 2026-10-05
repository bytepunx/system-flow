---
id: S-0265
type: story
nature: remediation
title: flai's wip markdown lint does not flag a bare email address, so a thread entry with one reached main and failed a story's close-out
status: backlog
owner: alex
created: 2026-10-05T00:03:14Z
updated: 2026-10-05T00:03:14Z
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
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-05T00:03:14Z
---
# S-0265 flai's wip markdown lint does not flag a bare email address, so a thread entry with one reached main and failed a story's close-out

## Goal

This story remediates [I-0056](../../../design/issues/I-0056-flai-s-wip-markdown-lint-does-not-flag-a-bare-email-address-so-a-thread-entry-with-one-reached-main-and-failed-a-story-s-close-out.md), "flai's wip markdown lint does not flag a bare email address, so a thread entry with one reached main and failed a story's close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0056 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0056 is closed with `flai issue close I-0056 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0056. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-02T16:11:21Z, 2.3 days before this story; under one cycle counts as one).
