---
id: S-0295
type: story
nature: improvement
title: One in-progress story holds every ready story by overlap, through folder-wide touches and docs/users/flai.md, so the board runs one story at a time under a limit of three
status: backlog
owner: alex
created: 2026-10-06T11:44:50Z
updated: 2026-10-06T11:47:30Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5h12m
    by: flai
    at: 2026-10-06T11:44:50Z
finalized:
  by: alex
  at: 2026-10-06T11:47:30Z
---
# S-0295 One in-progress story holds every ready story by overlap, through folder-wide touches and docs/users/flai.md, so the board runs one story at a time under a limit of three

## Goal

This story remediates [I-0087](../../../design/issues/I-0087-one-in-progress-story-holds-every-ready-story-by-overlap-through-folder-wide-touches-and-docs-users-flai-md-so-the-board-runs-one-story-at-a-time-under-a-limit-of-three.md), "One in-progress story holds every ready story by overlap, through folder-wide touches and docs/users/flai.md, so the board runs one story at a time under a limit of three". The issue recommends this solution:

Not designed yet. Directions to weigh:

- The planner narrows a folder claim to the files a story will change, where its tasks already name them.
- A document every story edits is split by section or by command, or is generated, so that two stories change different files.
- Whether a story in review still needs to hold others is a question for ADR-0046: its branch is finished and synced, and a later story syncs onto main when it is accepted.

## Acceptance criteria
- [ ] The cause I-0087 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0087 is closed with `flai issue close I-0087 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0087. time_lost_per_cycle 5h12m: 1h44m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-10-06T10:32:27Z, 0.1 days before this story; under one cycle counts as one).
