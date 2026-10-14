---
id: S-0306
type: story
nature: improvement
title: Two stories that each add an ADR conflict in design/adrs/README.md, whose index both append a row to, and flai opens a thread for each sync that finds it
status: in-progress
owner: alex
created: 2026-10-07T01:07:12Z
updated: 2026-10-14T23:07:38Z
transitions:
  - to: ready
    at: 2026-10-14T20:54:58Z
    by: alex
  - to: in-progress
    at: 2026-10-14T23:02:21Z
    by: agent-S-0306
tags: []
touches: [design/adrs/0157-flai-stream-sync-and-flai-accept-merge-design-adrs-readme-md-s-index-row-by-row.md, flai/internal/adrindex/adrindex.go, flai/internal/adrindex/adrindex_test.go, design/adrs/README.md, design/system/flai-cli.md, docs/users/flai.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 307
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-10-14
      hand_edits: 1
      work: 27
  models:
    - model: claude-opus-5-5
      input: 106
      output: 793
      cache_read: 4563560
      cache_write: 269751
      cost: 2.2305
  strategic:
    - kind: orchestrator
      seconds: 391
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 59
          output: 927
          cache_read: 7800827
          cache_write: 20732
          cost: 1.9609
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
- T-1583 The ADR index in design/adrs/README.md merges row by row from a stopped rebase's three versions, keeping each side's rows in number order
- T-1584 ADR records that a rebase merges design/adrs/README.md's index by rows, and the design and user docs say so
- T-1585 flai stream sync, flai task done, and flai accept merge the ADR index themselves, and the trial merge leaves it out of a pair's conflicts
- T-1586 Close I-0099 as fixed by the rebase's row merge of the ADR index

## Notes

Cost of delay inputs set by flai from I-0099. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T23:31:35Z, 0.1 days before this story; under one cycle counts as one).
