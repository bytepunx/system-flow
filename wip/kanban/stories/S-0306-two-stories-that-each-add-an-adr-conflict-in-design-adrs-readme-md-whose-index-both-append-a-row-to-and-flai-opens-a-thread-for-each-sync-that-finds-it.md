---
id: S-0306
type: story
nature: improvement
title: Two stories that each add an ADR conflict in design/adrs/README.md, whose index both append a row to, and flai opens a thread for each sync that finds it
status: backlog
owner: alex
created: 2026-10-07T01:07:12Z
updated: 2026-10-07T02:17:33Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-07T01:07:12Z
finalized:
  by: alex
  at: 2026-10-07T02:17:33Z
---
# S-0306 Two stories that each add an ADR conflict in design/adrs/README.md, whose index both append a row to, and flai opens a thread for each sync that finds it

## Goal

This story remediates [I-0099](../../../design/issues/I-0099-two-stories-that-each-add-an-adr-conflict-in-design-adrs-readme-md-whose-index-both-append-a-row-to-and-flai-opens-a-thread-for-each-sync-that-finds-it.md), "Two stories that each add an ADR conflict in design/adrs/README.md, whose index both append a row to, and flai opens a thread for each sync that finds it". The issue recommends this solution:

The index is derived from the ADR files, as `summary.md` is from the issues: `flai adr new` could regenerate it whole, and `flai stream sync` and `flai accept` could treat it as a generated file (ADR-0098), writing it again from the ADR files when a rebase stops on it alone, and leaving it out of a pair's conflicts. `design/adrs/README.md` would then need to carry nothing but the generated table, or hold its hand-written part apart from it.

## Acceptance criteria
- [ ] The cause I-0099 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0099 is closed with `flai issue close I-0099 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0099. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T23:31:35Z, 0.1 days before this story; under one cycle counts as one).
