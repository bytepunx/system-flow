---
id: S-0665
type: story
nature: improvement
title: flai serve records no image, container name, or mounts for a container harness run, so a run cannot be shown afterwards to have used the project's layer image
status: backlog
owner: alex
created: 2026-11-10T08:34:37Z
updated: 2026-11-10T08:34:37Z
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
    at: 2026-11-10T08:34:37Z
---
# S-0665 flai serve records no image, container name, or mounts for a container harness run, so a run cannot be shown afterwards to have used the project's layer image

## Goal

This story remediates [I-0274](../../../design/issues/I-0274-flai-serve-records-no-image-container-name-or-mounts-for-a-container-harness.md), "flai serve records no image, container name, or mounts for a container harness run, so a run cannot be shown afterwards to have used the project's layer image". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0274 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0274 is closed with `flai issue close I-0274 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0274. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T05:41:15Z, 3.1 days before this story; under one cycle counts as one).
