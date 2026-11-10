---
id: S-0705
type: story
nature: improvement
title: A role a project defines in a role file of its own cannot run on the container harness, since the start request does not carry its file to pilot
status: backlog
owner: alex
created: 2026-11-10T09:11:33Z
updated: 2026-11-10T09:11:33Z
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
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-10T09:11:33Z
---
# S-0705 A role a project defines in a role file of its own cannot run on the container harness, since the start request does not carry its file to pilot

## Goal

This story remediates [I-0258](../../../design/issues/I-0258-a-role-a-project-defines-in-a-role-file-of-its-own-cannot-run-on-the-container.md), "A role a project defines in a role file of its own cannot run on the container harness, since the start request does not carry its file to pilot". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0258 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0258 is closed with `flai issue close I-0258 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0258. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1.1 cycles of 168h (first reported 2026-11-02T19:06:43Z, 7.6 days before this story).
