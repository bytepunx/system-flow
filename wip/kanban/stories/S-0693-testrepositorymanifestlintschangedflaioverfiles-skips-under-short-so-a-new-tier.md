---
id: S-0693
type: story
nature: improvement
title: TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier
status: backlog
owner: alex
created: 2026-11-10T08:49:04Z
updated: 2026-11-10T08:49:04Z
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
    at: 2026-11-10T08:49:04Z
---
# S-0693 TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier

## Goal

This story remediates [I-0327](../../../design/issues/I-0327-testrepositorymanifestlintschangedflaioverfiles-skips-under-short-so-a-new-tier.md), "TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0327 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0327 is closed with `flai issue close I-0327 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0327. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-09T23:09:47Z, 0.4 days before this story; under one cycle counts as one).
