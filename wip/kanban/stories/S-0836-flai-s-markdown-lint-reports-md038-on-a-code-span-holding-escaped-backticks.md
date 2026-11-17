---
id: S-0836
type: story
nature: remediation
title: flai's markdown lint reports MD038 on a code span holding escaped backticks that markdownlint passes, so TestRepositoryLintsClean fails on a curator's digest entry
status: backlog
owner: alex
created: 2026-11-17T08:17:31Z
updated: 2026-11-17T08:17:31Z
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
    at: 2026-11-17T08:17:31Z
---
# S-0836 flai's markdown lint reports MD038 on a code span holding escaped backticks that markdownlint passes, so TestRepositoryLintsClean fails on a curator's digest entry

## Goal

This story remediates [I-0363](../../../design/issues/I-0363-flai-s-markdown-lint-reports-md038-on-a-code-span-holding-escaped-backticks.md), "flai's markdown lint reports MD038 on a code span holding escaped backticks that markdownlint passes, so TestRepositoryLintsClean fails on a curator's digest entry". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0363 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0363 is closed with `flai issue close I-0363 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0363. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T16:00:01Z, 3.7 days before this story; under one cycle counts as one).
