---
id: S-0695
type: story
nature: remediation
title: repo.test.ts's whole-monorepo read times out at vitest's 5s default under load and stops a close-out that does not touch it
status: review
owner: alex
created: 2026-11-10T08:54:28Z
updated: 2026-11-11T01:58:10Z
transitions:
  - to: ready
    at: 2026-11-11T01:45:52Z
    by: alex
  - to: in-progress
    at: 2026-11-11T01:50:19Z
    by: agent-S-0695
  - to: review
    at: 2026-11-11T01:58:10Z
    by: agent-S-0695
tags: [flaiover, testing]
topics: [testing]
touches: [flaiover/src/lib/server/repo.test.ts, design/issues/I-0269-repo-test-ts-s-whole-monorepo-read-times-out-at-vitest-s-5s-default-under-load.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 479
  priced_by: harness
  turns:
    - day: 2026-11-11
      ceremony: 1
      work: 29
  models:
    - model: claude-opus-5-5
      input: 66
      output: 10222
      cache_read: 4542885
      cache_write: 171911
      cost: 2.4886
  strategic:
    - kind: orchestrator
      seconds: 72
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 22
          output: 152
          cache_read: 1352868
          cache_write: 96425
          cost: 0.3983
        - model: claude-sonnet-5-5
          input: 4
          output: 48
          cache_read: 12602
          cache_write: 31468
          cost: 0.0334
cost_of_delay:
  inputs:
    time_lost_per_cycle: 18m
    by: flai
    at: 2026-11-10T08:54:28Z
  value: 45
  by: planner-S-0695
  at: 2026-11-10T10:29:26Z
forecast:
  duration: 14m
  delivery: 2026-11-11T02:05:00Z
  basis: "Its own forecast of 14m; 1st in the pull order with an in-progress limit of 3, behind S-0674, S-0676 and S-0688."
  by: flai
  at: 2026-11-11T01:46:11Z
finalized:
  by: orchestrator
  at: 2026-11-10T11:04:42Z
---
# S-0695 repo.test.ts's whole-monorepo read times out at vitest's 5s default under load and stops a close-out that does not touch it

## Goal

This story remediates [I-0269](../../../design/issues/I-0269-repo-test-ts-s-whole-monorepo-read-times-out-at-vitest-s-5s-default-under-load.md), "repo.test.ts's whole-monorepo read times out at vitest's 5s default under load and stops a close-out that does not touch it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0269 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0269 is closed with `flai issue close I-0269 --reason` saying what fixed it

## Tasks
- T-2969 Give repo.test.ts's whole-monorepo describe a 30 s timeout, as search.test.ts's monorepo describe has
- T-2972 Close I-0269 with the timeout as what fixed it

## Notes

Cost of delay inputs set by flai from I-0269. time_lost_per_cycle 18m: 6m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-11-07T04:07:16Z, 3.2 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from I-0269's three instances: give the `Repo on the monorepo` describe in `flaiover/src/lib/server/repo.test.ts` `{ timeout: 30000 }`, as `flaiover/src/lib/server/search.test.ts`'s `adrs on the monorepo` describe has. The planner timed the read, `flai hostapi items.list` with the archive, at 0.75 to 0.90 s on the host at a load average of 130 on 24 cores. All three failures came on a loaded host, so 5 s leaves too little room when vitest's workers share the host. A timeout that depends on load has no test that reproduces it reliably. The first criterion allows for that ("where one fits"). D-0524 also names concurrent writes to the main checkout as a cause. T-2969 checks that cause and leaves it to its own issue if it is real.

Touches:

- `flaiover/src/lib/server/repo.test.ts`: design. The issue names it, and the fix is there.
- `design/issues/I-0269-repo-test-ts-s-whole-monorepo-read-times-out-at-vitest-s-5s-default-under-load.md`: layout. `flai issue close` writes it (criterion 2).
- `design/issues/summary.md`: layout. flai regenerates it as the issue closes.
- Co-change: `flai touches suggest` listed nothing above 5 %. Its top entries came from `summary.md`'s many co-changes and are unrelated, so none was added. No folder touch.

Tasks: T-2969 is the first layer (the timeout). T-2972 is the second (closing I-0269), after T-2969, because its reason names what T-2969 changed and measured.

Figures:

- Forecast: 14m, flai's figure for size 5, kept. It has no ADR or docs change, and its one code change is a single line.
- Cost of delay: 45 USD a week, from the 18m `time_lost_per_cycle` input at 150 USD an hour, as `flai cod` gives it. Kept: the input is the operator's, and it is the issue's measured cost.
- No `after`: the story reads no other story's output.
