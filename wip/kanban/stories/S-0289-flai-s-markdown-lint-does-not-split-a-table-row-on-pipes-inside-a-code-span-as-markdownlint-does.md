---
id: S-0289
type: story
nature: remediation
title: flai's markdown lint does not split a table row on pipes inside a code span, as markdownlint does
status: backlog
owner: alex
created: 2026-10-06T09:56:52Z
updated: 2026-10-06T09:56:52Z
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
    time_lost_per_cycle: 6m
    by: flai
    at: 2026-10-06T09:56:52Z
---
# S-0289 flai's markdown lint does not split a table row on pipes inside a code span, as markdownlint does

## Goal

This story remediates [I-0077](../../../design/issues/I-0077-flai-s-markdown-lint-does-not-split-a-table-row-on-pipes-inside-a-code-span-as-markdownlint-does.md), "flai's markdown lint does not split a table row on pipes inside a code span, as markdownlint does". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0077 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0077 is closed with `flai issue close I-0077 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0077. time_lost_per_cycle 6m: 6m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-05T03:58:53Z, 1.2 days before this story; under one cycle counts as one).
