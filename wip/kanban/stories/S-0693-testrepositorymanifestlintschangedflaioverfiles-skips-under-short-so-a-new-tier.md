---
id: S-0693
type: story
nature: improvement
title: TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier
status: backlog
owner: alex
created: 2026-11-10T08:49:04Z
updated: 2026-11-10T11:05:34Z
transitions: []
tags: [flai, tests]
topics: [code]
touches: [flai/internal/verify/select_test.go, design/issues/I-0327-testrepositorymanifestlintschangedflaioverfiles-skips-under-short-so-a-new-tier.md, design/issues/summary.md]
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
    - kind: planner
      seconds: 23
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 90
          output: 391
          cache_read: 5254190
          cache_write: 165340
          cost: 0.3648
        - model: claude-opus-5-5
          input: 115
          output: 10572
          cache_read: 8571742
          cache_write: 163919
          cost: 3.0848
    - kind: orchestrator
      seconds: 23
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 4
          cache_read: 362301
          cache_write: 10560
          cost: 0.102
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-10T08:49:04Z
  value: 25
  by: planner-S-0693
  at: 2026-11-10T10:29:13Z
forecast:
  duration: 45m
  delivery: 2026-11-11T01:02:00Z
  basis: "Its own forecast of 45m; 23rd in the pull order with an in-progress limit of 3, behind S-0665, S-0667, S-0708, S-0666, S-0688, S-0677, S-0679, S-0675, S-0676, S-0678, S-0674, S-0704, S-0681, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663, S-0683, S-0686, S-0689, S-0690, S-0691 and S-0692."
  by: flai
  at: 2026-11-10T11:05:34Z
finalized:
  by: orchestrator
  at: 2026-11-10T11:04:43Z
---
# S-0693 TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier

## Goal

This story remediates [I-0327](../../../design/issues/I-0327-testrepositorymanifestlintschangedflaioverfiles-skips-under-short-so-a-new-tier.md), "TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0327 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0327 is closed with `flai issue close I-0327 --reason` saying what fixed it

## Tasks
- T-2982 Run TestRepositoryManifestLintsChangedFlaioverFiles under -short
- T-2985 Close I-0327 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0327. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-09T23:09:47Z, 0.4 days before this story; under one cycle counts as one).

### Planning

Proposed solution, from I-0327's one instance: drop the `testing.Short()` gate from `TestRepositoryManifestLintsChangedFlaioverFiles` in `flai/internal/verify/select_test.go`. The test reads only `system-flow.yaml` and the file tree through `os.DirFS`, needs no git, and keeps its "monorepo not present" skip. The repository-manifest tests beside it (`TestRepositoryManifestRunsTheGatewayChecksOnlyByName`, `TestRepositoryManifestTypeChecksChangedFlaioverFiles`, `TestRepositoryManifestRunsFlaioverTestsOfTheBuiltFlaiForHostAPIChanges`) already run under `-short` for this reason. Then the go-test tier, which `flai test` runs with `-short` on a change to `system-flow.yaml` or `select_test.go`, catches a new tier the test must list. The test itself is the reproduction: under `-short` it now runs, and fails while its `want` lacks a tier the manifest selects.

Touches, and where each came from:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/verify/select_test.go` | design (I-0327) | The test and its gate, at line 474 |
| `design/issues/I-0327-…-tier.md` | layout | `flai issue close` writes the issue file |
| `design/issues/summary.md` | layout | `flai issue close` regenerates the summary |

`flai touches suggest S-0693 flai/internal/verify/select_test.go` listed by co-change `system-flow.yaml` (55%), `flai/internal/verify/select.go`, `verify.go`, and others. None is kept: the fix changes no tier, no selection code, and no documentation. `design/conventions/tooling.md` lists the tests that skip under `-short`, and this test is not among them, so it needs no change. No folder touch. No `after`: the story reads nothing another open story writes.

Figures:

- Forecast 45m, adjusted from flai's 2m. flai's figure is 59 s per unit of size times size 2, counted before any touch was declared. It leaves out the close-out, whose integration and smoke tiers a change under `flai/` selects. Close-outs on this host took 8m14s (S-0672) and 20m44s (S-0669) today. Delivery 2026-11-10T14:51Z is flai's queue-based 14:08Z moved by the 43m added; the story is 25th in the pull order.
- Cost of delay value 25 USD a week, as `flai cod` gives it from the operator's input of 10m per 168h cycle at 150 USD an hour. It stands: I-0327 has one instance, and the input is the operator's.
