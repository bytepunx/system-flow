---
id: S-0299
type: story
nature: remediation
title: A story's sub-agent that writes under .claude/ blocks its layer for thirty minutes on a permission thread nobody answers, and the start prompt does not warn the agent beforehand
status: backlog
owner: alex
created: 2026-10-06T21:00:33Z
updated: 2026-10-06T21:00:33Z
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
    at: 2026-10-06T21:00:33Z
---
# S-0299 A story's sub-agent that writes under .claude/ blocks its layer for thirty minutes on a permission thread nobody answers, and the start prompt does not warn the agent beforehand

## Goal

This story remediates [I-0093](../../../design/issues/I-0093-a-story-s-sub-agent-that-writes-under-claude-blocks-its-layer-for-thirty-minutes-on-a-permission-thread-nobody-answers-and-the-start-prompt-does-not-warn-the-agent-beforehand.md), "A story's sub-agent that writes under .claude/ blocks its layer for thirty minutes on a permission thread nobody answers, and the start prompt does not warn the agent beforehand". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0093 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0093 is closed with `flai issue close I-0093 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0093. time_lost_per_cycle 6m: 6m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T20:53:14Z, 0 days before this story; under one cycle counts as one).
