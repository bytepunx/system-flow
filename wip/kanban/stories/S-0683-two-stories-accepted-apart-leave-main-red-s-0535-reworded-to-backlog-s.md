---
id: S-0683
type: story
nature: remediation
title: "Two stories accepted apart leave main red: S-0535 reworded to-backlog's condition, and S-0416's TestADefinedRoleIsHeldToItsFile still expects the old words"
status: backlog
owner: alex
created: 2026-11-10T08:37:08Z
updated: 2026-11-10T08:37:08Z
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
    time_lost_per_cycle: 18m
    by: flai
    at: 2026-11-10T08:37:08Z
---
# S-0683 Two stories accepted apart leave main red: S-0535 reworded to-backlog's condition, and S-0416's TestADefinedRoleIsHeldToItsFile still expects the old words

## Goal

This story remediates [I-0229](../../../design/issues/I-0229-two-stories-accepted-apart-leave-main-red-s-0535-reworded-to-backlog-s-condition-and-s-0416-s-testadefinedroleisheldtoitsfile-still-expects-the-old-words.md), "Two stories accepted apart leave main red: S-0535 reworded to-backlog's condition, and S-0416's TestADefinedRoleIsHeldToItsFile still expects the old words". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0229 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0229 is closed with `flai issue close I-0229 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0229. time_lost_per_cycle 18m: 30m per occurrence × 1 occurrence ÷ 1.7 cycles of 168h (first reported 2026-10-29T03:27:16Z, 12.2 days before this story).
