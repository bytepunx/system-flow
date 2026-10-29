---
id: S-0540
type: story
nature: remediation
title: flai test's go-test tier runs the mcpserver package with -race into go test's 10-minute default timeout on a loaded host, and names only the package, so a story that changes mcpserver cannot pass its tests
status: backlog
owner: alex
created: 2026-10-29T05:59:43Z
updated: 2026-10-29T05:59:43Z
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
    time_lost_per_cycle: 1h12m
    by: flai
    at: 2026-10-29T05:59:43Z
---
# S-0540 flai test's go-test tier runs the mcpserver package with -race into go test's 10-minute default timeout on a loaded host, and names only the package, so a story that changes mcpserver cannot pass its tests

## Goal

This story remediates [I-0227](../../../design/issues/I-0227-flai-test-s-go-test-tier-runs-the-mcpserver-package-with-race-into-go-test-s-10-minute-default-timeout-on-a-loaded-host-and-names-only-the-package-so-a-story-that-changes-mcpserver-cannot-pass-its-tests.md), "flai test's go-test tier runs the mcpserver package with -race into go test's 10-minute default timeout on a loaded host, and names only the package, so a story that changes mcpserver cannot pass its tests". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0227 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0227 is closed with `flai issue close I-0227 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0227. time_lost_per_cycle 1h12m: 24m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-10-29T00:08:41Z, 0.2 days before this story; under one cycle counts as one).
