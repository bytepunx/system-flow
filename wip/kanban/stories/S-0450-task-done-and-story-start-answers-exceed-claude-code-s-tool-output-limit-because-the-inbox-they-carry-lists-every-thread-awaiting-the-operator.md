---
id: S-0450
type: story
nature: improvement
title: task_done and story_start answers exceed Claude Code's tool output limit, because the inbox they carry lists every thread awaiting the operator
status: backlog
owner: alex
created: 2026-10-20T18:53:36Z
updated: 2026-10-20T18:53:36Z
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
    time_lost_per_cycle: 12m
    by: flai
    at: 2026-10-20T18:53:36Z
---
# S-0450 task_done and story_start answers exceed Claude Code's tool output limit, because the inbox they carry lists every thread awaiting the operator

## Goal

This story remediates [I-0167](../../../design/issues/I-0167-task-done-and-story-start-answers-exceed-claude-code-s-tool-output-limit-because-the-inbox-they-carry-lists-every-thread-awaiting-the-operator.md), "task_done and story_start answers exceed Claude Code's tool output limit, because the inbox they carry lists every thread awaiting the operator". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0167 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0167 is closed with `flai issue close I-0167 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0167. time_lost_per_cycle 12m: 4m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-10-19T15:28:45Z, 1.1 days before this story; under one cycle counts as one).
