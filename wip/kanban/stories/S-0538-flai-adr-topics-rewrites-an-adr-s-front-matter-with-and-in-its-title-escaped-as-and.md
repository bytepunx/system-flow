---
id: S-0538
type: story
nature: remediation
title: "flai adr topics rewrites an ADR's front matter with < and > in its title escaped as < and >"
status: backlog
owner: alex
created: 2026-10-29T02:01:36Z
updated: 2026-10-29T02:01:36Z
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
    at: 2026-10-29T02:01:36Z
---
# S-0538 flai adr topics rewrites an ADR's front matter with < and > in its title escaped as < and >

## Goal

This story remediates [I-0222](../../../design/issues/I-0222-flai-adr-topics-rewrites-an-adr-s-front-matter-with-and-in-its-title-escaped-as-and.md), "flai adr topics rewrites an ADR's front matter with < and > in its title escaped as < and >". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0222 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0222 is closed with `flai issue close I-0222 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0222. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-28T21:41:05Z, 0.2 days before this story; under one cycle counts as one).
