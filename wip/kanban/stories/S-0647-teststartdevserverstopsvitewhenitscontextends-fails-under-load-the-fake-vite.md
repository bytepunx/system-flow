---
id: S-0647
type: story
nature: remediation
title: "TestStartDevServerStopsViteWhenItsContextEnds fails under load: the fake vite wrote no pid"
status: backlog
owner: alex
created: 2026-11-10T05:12:04Z
updated: 2026-11-10T05:12:04Z
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
    time_lost_per_cycle: 51m
    by: flai
    at: 2026-11-10T05:12:04Z
---
# S-0647 TestStartDevServerStopsViteWhenItsContextEnds fails under load: the fake vite wrote no pid

## Goal

This story remediates [I-0311](../../../design/issues/I-0311-teststartdevserverstopsvitewhenitscontextends-fails-under-load-the-fake-vite.md), "TestStartDevServerStopsViteWhenItsContextEnds fails under load: the fake vite wrote no pid". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0311 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0311 is closed with `flai issue close I-0311 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0311. time_lost_per_cycle 51m: 17m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-11-08T01:43:19Z, 2.1 days before this story; under one cycle counts as one).
