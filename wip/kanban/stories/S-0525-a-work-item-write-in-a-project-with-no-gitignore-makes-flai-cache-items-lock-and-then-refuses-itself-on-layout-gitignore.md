---
id: S-0525
type: story
nature: remediation
title: A work item write in a project with no .gitignore makes .flai-cache/items.lock and then refuses itself on layout.gitignore
status: backlog
owner: alex
created: 2026-10-28T12:53:57Z
updated: 2026-10-28T12:53:57Z
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
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-28T12:53:57Z
---
# S-0525 A work item write in a project with no .gitignore makes .flai-cache/items.lock and then refuses itself on layout.gitignore

## Goal

This story remediates [I-0191](../../../design/issues/I-0191-a-work-item-write-in-a-project-with-no-gitignore-makes-flai-cache-items-lock-and-then-refuses-itself-on-layout-gitignore.md), "A work item write in a project with no .gitignore makes .flai-cache/items.lock and then refuses itself on layout.gitignore". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0191 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0191 is closed with `flai issue close I-0191 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0191. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-22T20:40:04Z, 5.7 days before this story; under one cycle counts as one).
