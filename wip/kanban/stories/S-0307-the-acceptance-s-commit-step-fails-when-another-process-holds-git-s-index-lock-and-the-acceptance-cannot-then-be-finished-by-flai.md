---
id: S-0307
type: story
nature: remediation
title: The acceptance's commit step fails when another process holds git's index lock, and the acceptance cannot then be finished by flai
status: backlog
owner: alex
created: 2026-10-07T01:07:13Z
updated: 2026-10-07T01:10:39Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2m
    by: flai
    at: 2026-10-07T01:07:13Z
finalized:
  by: alex
  at: 2026-10-07T01:10:39Z
---
# S-0307 The acceptance's commit step fails when another process holds git's index lock, and the acceptance cannot then be finished by flai

## Goal

This story remediates [I-0100](../../../design/issues/I-0100-the-acceptance-s-commit-step-fails-when-another-process-holds-git-s-index-lock-and-the-acceptance-cannot-then-be-finished-by-flai.md), "The acceptance's commit step fails when another process holds git's index lock, and the acceptance cannot then be finished by flai". The issue recommends this solution:

Directions to weigh: retry the commit a few times when git reports the index lock held, since the other writer is brief; or make an acceptance that is done and archived but whose commit is missing resumable, as one that is done but not archived already is ("completed from step 0 without a second transition"); and record the commit's failure in the journal.

## Acceptance criteria
- [ ] The cause I-0100 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0100 is closed with `flai issue close I-0100 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0100. time_lost_per_cycle 2m: 2m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T23:31:35Z, 0.1 days before this story; under one cycle counts as one).
