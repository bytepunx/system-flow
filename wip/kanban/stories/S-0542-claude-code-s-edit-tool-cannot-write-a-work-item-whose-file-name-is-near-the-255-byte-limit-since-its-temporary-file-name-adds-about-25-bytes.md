---
id: S-0542
type: story
nature: improvement
title: Claude Code's Edit tool cannot write a work item whose file name is near the 255-byte limit, since its temporary file name adds about 25 bytes
status: backlog
owner: alex
created: 2026-10-29T15:14:37Z
updated: 2026-10-29T15:14:37Z
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
    time_lost_per_cycle: 3m
    by: flai
    at: 2026-10-29T15:14:37Z
---
# S-0542 Claude Code's Edit tool cannot write a work item whose file name is near the 255-byte limit, since its temporary file name adds about 25 bytes

## Goal

This story remediates [I-0231](../../../design/issues/I-0231-claude-code-s-edit-tool-cannot-write-a-work-item-whose-file-name-is-near-the-255-byte-limit-since-its-temporary-file-name-adds-about-25-bytes.md), "Claude Code's Edit tool cannot write a work item whose file name is near the 255-byte limit, since its temporary file name adds about 25 bytes". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0231 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0231 is closed with `flai issue close I-0231 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0231. time_lost_per_cycle 3m: 3m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-29T14:15:00Z, 0 days before this story; under one cycle counts as one).
