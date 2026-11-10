---
id: S-0695
type: story
nature: remediation
title: repo.test.ts's whole-monorepo read times out at vitest's 5s default under load and stops a close-out that does not touch it
status: backlog
owner: alex
created: 2026-11-10T08:54:28Z
updated: 2026-11-10T19:26:06Z
transitions: []
tags: [flaiover, testing]
topics: [testing]
touches: [flaiover/src/lib/server/repo.test.ts, design/issues/I-0269-repo-test-ts-s-whole-monorepo-read-times-out-at-vitest-s-5s-default-under-load.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 21
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 28
          cache_read: 685977
          cache_write: 2228
          cost: 0.1881
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
  delivery: 2026-11-11T07:47:00Z
  basis: "Its own forecast of 14m; 24th in the pull order with an in-progress limit of 3, behind S-0666, S-0667, S-0708, S-0688, S-0677, S-0679, S-0675, S-0676, S-0678, S-0674, S-0704, S-0681, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663, S-0683, S-0686, S-0689, S-0690, S-0691, S-0692, S-0693 and S-0694."
  by: flai
  at: 2026-11-10T19:26:06Z
finalized:
  by: orchestrator
  at: 2026-11-10T11:04:42Z
---
# S-0695 repo.test.ts's whole-monorepo read times out at vitest's 5s default under load and stops a close-out that does not touch it

## Goal

This story remediates [I-0269](../../../design/issues/I-0269-repo-test-ts-s-whole-monorepo-read-times-out-at-vitest-s-5s-default-under-load.md), "repo.test.ts's whole-monorepo read times out at vitest's 5s default under load and stops a close-out that does not touch it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0269 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0269 is closed with `flai issue close I-0269 --reason` saying what fixed it

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
