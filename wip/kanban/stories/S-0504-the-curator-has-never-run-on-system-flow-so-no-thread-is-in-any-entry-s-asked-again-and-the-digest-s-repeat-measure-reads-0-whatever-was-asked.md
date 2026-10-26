---
id: S-0504
type: story
nature: remediation
title: The curator has never run on system-flow, so no thread is in any entry's asked_again and the digest's repeat measure reads 0 whatever was asked
status: backlog
owner: alex
created: 2026-10-26T04:01:51Z
updated: 2026-10-26T04:01:51Z
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
    at: 2026-10-26T04:01:51Z
---
# S-0504 The curator has never run on system-flow, so no thread is in any entry's asked_again and the digest's repeat measure reads 0 whatever was asked

## Goal

This story remediates [I-0207](../../../design/issues/I-0207-the-curator-has-never-run-on-system-flow-so-no-thread-is-in-any-entry-s-asked-again-and-the-digest-s-repeat-measure-reads-0-whatever-was-asked.md), "The curator has never run on system-flow, so no thread is in any entry's asked_again and the digest's repeat measure reads 0 whatever was asked". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0207 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0207 is closed with `flai issue close I-0207 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0207. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-24T23:48:27Z, 1.2 days before this story; under one cycle counts as one).
