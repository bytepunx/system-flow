---
id: S-0563
type: story
nature: remediation
title: The orchestrator's conventions sit within about 1.5KB of its 80KB pack budget, so a few more lines in work-management.md or strategic-agents.md drop every design brief from its pack
status: backlog
owner: alex
created: 2026-11-02T04:21:24Z
updated: 2026-11-02T04:21:53Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 36m
    by: flai
    at: 2026-11-02T04:21:24Z
finalized:
  by: alex
  at: 2026-11-02T04:21:53Z
---
# S-0563 The orchestrator's conventions sit within about 1.5KB of its 80KB pack budget, so a few more lines in work-management.md or strategic-agents.md drop every design brief from its pack

## Goal

This story remediates [I-0250](../../../design/issues/I-0250-the-orchestrator-s-conventions-sit-within-about-1-5kb-of-its-80kb-pack-budget.md), "The orchestrator's conventions sit within about 1.5KB of its 80KB pack budget, so a few more lines in work-management.md or strategic-agents.md drop every design brief from its pack". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0250 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0250 is closed with `flai issue close I-0250 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0250. time_lost_per_cycle 36m: 18m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-11-02T03:46:07Z, 0 days before this story; under one cycle counts as one).

Made by agent-S-0518 on the operator's word on TH-0804 (option 1, confirmed 2026-11-02T04:21:05Z). S-0518 waits in progress until this fix is on main, since its close-out and the merge queue's batch both stop at `TestAThreadPackIsUnderAQuarterOfTheGeneralPack`.

The orchestrator measured the test on TH-0804. In the main checkout the thread pack is 22.6% of the general pack, and the test passes. In S-0518's worktree it is 28.3%, and the test fails. The thread pack is 43,932 bytes in both. The general pack drops its design briefs in the worktree (155,062 against 194,234 bytes), and the only difference there is S-0518's issue records under `design/issues`. So the result depends on what the repository holds, wip and issues included. The fix should make the test independent of that, or give the orchestrator's budget a real margin, as well as trimming its conventions.
