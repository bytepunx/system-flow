---
id: S-0483
type: story
nature: remediation
title: Accepting a story that adds a flaiover dependency leaves the main checkout's flaiover/node_modules without it, so every close-out's integration tier fails TestEveryListedFlaioverComponentCompilesInTheHarness until someone installs
status: backlog
owner: alex
created: 2026-10-21T12:00:17Z
updated: 2026-10-21T12:03:49Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-21T12:00:17Z
finalized:
  by: alex
  at: 2026-10-21T12:03:49Z
---
# S-0483 Accepting a story that adds a flaiover dependency leaves the main checkout's flaiover/node_modules without it, so every close-out's integration tier fails TestEveryListedFlaioverComponentCompilesInTheHarness until someone installs

## Goal

This story remediates [I-0179](../../../design/issues/I-0179-accepting-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-flaiover-node-modules-without-it-so-every-close-out-s-integration-tier-fails-testeverylistedflaiovercomponentcompilesintheharness-until-someone-installs.md), "Accepting a story that adds a flaiover dependency leaves the main checkout's flaiover/node_modules without it, so every close-out's integration tier fails TestEveryListedFlaioverComponentCompilesInTheHarness until someone installs". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0179 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0179 is closed with `flai issue close I-0179 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0179. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-21T08:44:04Z, 0.1 days before this story; under one cycle counts as one).
