---
id: S-0364
type: story
nature: improvement
title: flai test vets and tests only the changed Go packages, so a change that breaks a package importing them is found only by the integration tier
status: backlog
owner: alex
created: 2026-10-12T14:12:31Z
updated: 2026-10-12T14:12:31Z
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
    at: 2026-10-12T14:12:31Z
---
# S-0364 flai test vets and tests only the changed Go packages, so a change that breaks a package importing them is found only by the integration tier

## Goal

This story remediates [I-0128](../../../design/issues/I-0128-flai-test-vets-and-tests-only-the-changed-go-packages-so-a-change-that-breaks-a-package-importing-them-is-found-only-by-the-integration-tier.md), "flai test vets and tests only the changed Go packages, so a change that breaks a package importing them is found only by the integration tier". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0128 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0128 is closed with `flai issue close I-0128 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0128. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-08T10:50:48Z, 4.1 days before this story; under one cycle counts as one).
