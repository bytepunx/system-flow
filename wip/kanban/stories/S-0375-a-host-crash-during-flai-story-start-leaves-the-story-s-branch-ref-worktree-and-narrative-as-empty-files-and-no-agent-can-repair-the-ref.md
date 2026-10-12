---
id: S-0375
type: story
nature: remediation
title: A host crash during flai story start leaves the story's branch ref, worktree, and narrative as empty files, and no agent can repair the ref
status: backlog
owner: alex
created: 2026-10-12T16:52:56Z
updated: 2026-10-12T16:52:56Z
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
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-10-12T16:52:56Z
---
# S-0375 A host crash during flai story start leaves the story's branch ref, worktree, and narrative as empty files, and no agent can repair the ref

## Goal

This story remediates [I-0137](../../../design/issues/I-0137-a-host-crash-during-flai-story-start-leaves-the-story-s-branch-ref-worktree-and-narrative-as-empty-files-and-no-agent-can-repair-the-ref.md), "A host crash during flai story start leaves the story's branch ref, worktree, and narrative as empty files, and no agent can repair the ref". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0137 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0137 is closed with `flai issue close I-0137 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0137. time_lost_per_cycle 30m: 30m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-12T15:29:54Z, 0.1 days before this story; under one cycle counts as one).
