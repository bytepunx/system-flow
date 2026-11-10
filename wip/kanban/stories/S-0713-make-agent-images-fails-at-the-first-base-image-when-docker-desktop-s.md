---
id: S-0713
type: story
nature: improvement
title: make agent-images fails at the first base image when Docker Desktop's credential helper desktop.exe cannot start from WSL
status: backlog
owner: alex
created: 2026-11-10T10:11:03Z
updated: 2026-11-10T10:11:03Z
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
    at: 2026-11-10T10:11:03Z
---
# S-0713 make agent-images fails at the first base image when Docker Desktop's credential helper desktop.exe cannot start from WSL

## Goal

This story remediates [I-0339](../../../design/issues/I-0339-make-agent-images-fails-at-the-first-base-image-when-docker-desktop-s.md), "make agent-images fails at the first base image when Docker Desktop's credential helper desktop.exe cannot start from WSL". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0339 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0339 is closed with `flai issue close I-0339 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0339. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T09:13:13Z, 0 days before this story; under one cycle counts as one).
