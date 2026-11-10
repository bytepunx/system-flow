---
id: S-0650
type: story
nature: remediation
title: "TestComponentsServeEndsWhenViteExits fails under load: the fake vite exits before it answers"
status: backlog
owner: alex
created: 2026-11-10T05:12:13Z
updated: 2026-11-10T05:12:13Z
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
    time_lost_per_cycle: 8m
    by: flai
    at: 2026-11-10T05:12:13Z
---
# S-0650 TestComponentsServeEndsWhenViteExits fails under load: the fake vite exits before it answers

## Goal

This story remediates [I-0325](../../../design/issues/I-0325-testcomponentsserveendswhenviteexits-fails-under-load-the-fake-vite-exits.md), "TestComponentsServeEndsWhenViteExits fails under load: the fake vite exits before it answers". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0325 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0325 is closed with `flai issue close I-0325 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0325. time_lost_per_cycle 8m: 8m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-09T22:51:04Z, 0.3 days before this story; under one cycle counts as one).
