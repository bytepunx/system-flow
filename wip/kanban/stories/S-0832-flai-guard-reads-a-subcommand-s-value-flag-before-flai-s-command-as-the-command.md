---
id: S-0832
type: story
nature: remediation
title: flai guard reads a subcommand's value flag before flai's command as the command, so a rule can miss the command flai runs
status: backlog
owner: alex
created: 2026-11-17T08:17:26Z
updated: 2026-11-17T08:17:26Z
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
    at: 2026-11-17T08:17:26Z
---
# S-0832 flai guard reads a subcommand's value flag before flai's command as the command, so a rule can miss the command flai runs

## Goal

This story remediates [I-0358](../../../design/issues/I-0358-flai-guard-reads-a-subcommand-s-value-flag-before-flai-s-command-as-the-command.md), "flai guard reads a subcommand's value flag before flai's command as the command, so a rule can miss the command flai runs". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0358 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0358 is closed with `flai issue close I-0358 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0358. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T08:42:15Z, 4 days before this story; under one cycle counts as one).
