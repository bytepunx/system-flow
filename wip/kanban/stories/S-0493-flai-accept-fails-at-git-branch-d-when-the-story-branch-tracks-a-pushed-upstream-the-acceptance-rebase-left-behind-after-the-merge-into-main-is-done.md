---
id: S-0493
type: story
nature: remediation
title: flai accept fails at git branch -d when the story branch tracks a pushed upstream the acceptance rebase left behind, after the merge into main is done
status: backlog
owner: alex
created: 2026-10-24T13:15:31Z
updated: 2026-10-24T13:15:31Z
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
    at: 2026-10-24T13:15:31Z
---
# S-0493 flai accept fails at git branch -d when the story branch tracks a pushed upstream the acceptance rebase left behind, after the merge into main is done

## Goal

This story remediates [I-0204](../../../design/issues/I-0204-flai-accept-fails-at-git-branch-d-when-the-story-branch-tracks-a-pushed-upstream-the-acceptance-rebase-left-behind-after-the-merge-into-main-is-done.md), "flai accept fails at git branch -d when the story branch tracks a pushed upstream the acceptance rebase left behind, after the merge into main is done". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0204 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0204 is closed with `flai issue close I-0204 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0204. time_lost_per_cycle 10m carried over from I-0204's Impact section. Its Impact time_lost_per_cycle was taken rather than the 10m derived from its cost and count.
