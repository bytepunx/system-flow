---
id: S-0324
type: story
nature: remediation
title: "flai's MD034 does not report a bare `www.` literal, which markdownlint reports as a bare URL"
status: backlog
owner: alex
created: 2026-10-07T18:59:55Z
updated: 2026-10-07T18:59:55Z
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
    at: 2026-10-07T18:59:55Z
---
# S-0324 flai's MD034 does not report a bare `www.` literal, which markdownlint reports as a bare URL

## Goal

This story remediates [I-0110](../../../design/issues/I-0110-flai-s-md034-does-not-report-a-bare-www-literal-which-markdownlint-reports-as-a-bare-url.md), "flai's MD034 does not report a bare `www.` literal, which markdownlint reports as a bare URL". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0110 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0110 is closed with `flai issue close I-0110 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0110. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T08:55:00Z, 0.4 days before this story; under one cycle counts as one).
