---
id: S-0288
type: story
nature: remediation
title: flaiover's notify.test.ts fails now and then under the full vitest run because project.info reads a system-flow.yaml with no version
status: backlog
owner: alex
created: 2026-10-06T09:56:51Z
updated: 2026-10-06T09:56:51Z
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
    at: 2026-10-06T09:56:51Z
---
# S-0288 flaiover's notify.test.ts fails now and then under the full vitest run because project.info reads a system-flow.yaml with no version

## Goal

This story remediates [I-0085](../../../design/issues/I-0085-flaiover-s-notify-test-ts-fails-now-and-then-under-the-full-vitest-run-because-project-info-reads-a-system-flow-yaml-with-no-version.md), "flaiover's notify.test.ts fails now and then under the full vitest run because project.info reads a system-flow.yaml with no version". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0085 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0085 is closed with `flai issue close I-0085 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0085. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T07:08:02Z, 0.1 days before this story; under one cycle counts as one).
