---
id: S-0501
type: story
nature: remediation
title: search.test.ts 'finds a body phrase across wip and design' times out at 5 s on the monorepo while other stories' close-outs run
status: backlog
owner: alex
created: 2026-10-26T04:01:35Z
updated: 2026-10-26T04:01:35Z
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
    at: 2026-10-26T04:01:35Z
---
# S-0501 search.test.ts 'finds a body phrase across wip and design' times out at 5 s on the monorepo while other stories' close-outs run

## Goal

This story remediates [I-0185](../../../design/issues/I-0185-search-test-ts-finds-a-body-phrase-across-wip-and-design-times-out-at-5-s-on-the-monorepo-while-other-stories-close-outs-run.md), "search.test.ts 'finds a body phrase across wip and design' times out at 5 s on the monorepo while other stories' close-outs run". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0185 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0185 is closed with `flai issue close I-0185 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0185. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-21T12:58:18Z, 4.6 days before this story; under one cycle counts as one).
