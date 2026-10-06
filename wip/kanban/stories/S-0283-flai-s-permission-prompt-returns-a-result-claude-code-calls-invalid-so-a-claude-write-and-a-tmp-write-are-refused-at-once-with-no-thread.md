---
id: S-0283
type: story
nature: remediation
title: flai's permission_prompt returns a result Claude Code calls invalid, so a .claude/ write and a /tmp write are refused at once with no thread
status: backlog
owner: alex
created: 2026-10-06T03:45:18Z
updated: 2026-10-06T03:45:18Z
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
    at: 2026-10-06T03:45:18Z
---
# S-0283 flai's permission_prompt returns a result Claude Code calls invalid, so a .claude/ write and a /tmp write are refused at once with no thread

## Goal

This story remediates [I-0082](../../../design/issues/I-0082-flai-s-permission-prompt-returns-a-result-claude-code-calls-invalid-so-a-claude-write-and-a-tmp-write-are-refused-at-once-with-no-thread.md), "flai's permission_prompt returns a result Claude Code calls invalid, so a .claude/ write and a /tmp write are refused at once with no thread". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0082 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0082 is closed with `flai issue close I-0082 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0082. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T03:19:21Z, 0 days before this story; under one cycle counts as one).
