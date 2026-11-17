---
id: S-0827
type: story
nature: improvement
title: TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier
status: done
owner: alex
created: 2026-11-17T08:17:17Z
updated: 2026-11-17T09:04:30Z
transitions:
  - to: ready
    at: 2026-11-17T08:30:02Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-17T08:30:34Z
    by: agent-S-0827
  - to: review
    at: 2026-11-17T09:02:27Z
    by: agent-S-0827
  - to: done
    at: 2026-11-17T09:04:30Z
    by: orchestrator
tags: [flai, tests]
touches: [flai/internal/verify/select_test.go, design/issues/I-0327-testrepositorymanifestlintschangedflaioverfiles-skips-under-short-so-a-new-tier.md, design/issues/summary.md, design/issues/I-0401-main-s-go-test-tier-fails-testpublishformulaunderjudgement-since-the-publish.md, flai/cmd/publish_judgement_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1923
  priced_by: harness
  turns:
    - day: 2026-11-17
      ceremony: 3
      work: 33
  models:
    - model: claude-opus-5-5
      input: 80
      output: 12346
      cache_read: 7699054
      cache_write: 250359
      cost: 3.7899
  strategic:
    - kind: planner
      seconds: 107
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 172
          output: 758
          cache_read: 6649291
          cache_write: 497399
          cost: 0.4881
        - model: claude-opus-5-5
          input: 412
          output: 12193
          cache_read: 25095393
          cache_write: 752613
          cost: 7.9731
    - kind: orchestrator
      seconds: 214
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 60
          output: 479
          cache_read: 2749330
          cache_write: 230553
          cost: 0.9246
        - model: claude-sonnet-5-5
          input: 6
          output: 48
          cache_read: 45252
          cache_write: 37605
          cost: 0.0665
cost_of_delay:
  inputs:
    time_lost_per_cycle: 9m
    by: flai
    at: 2026-11-17T08:17:17Z
  value: 22.5
  by: planner-S-0827
  at: 2026-11-17T08:28:44Z
forecast:
  duration: 30m
  delivery: 2026-11-17T09:42:00Z
  basis: "Its own forecast of 30m; 2nd in the pull order with an in-progress limit of 3, behind S-0752, S-0806 and S-0754."
  by: flai
  at: 2026-11-17T08:30:09Z
finalized:
  by: orchestrator
  at: 2026-11-17T08:29:57Z
---
# S-0827 TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier

## Goal

This story remediates [I-0327](../../../design/issues/I-0327-testrepositorymanifestlintschangedflaioverfiles-skips-under-short-so-a-new-tier.md), "TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0327 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0327 is closed with `flai issue close I-0327 --reason` saying what fixed it

## Tasks
- T-3523 Run TestRepositoryManifestLintsChangedFlaioverFiles under -short
- T-3524 Close I-0327 with what fixed it
- T-3557 TestPublishFormulaUnderJudgement reads the publish formula the template ships, which main's checkout no longer holds

## Notes

Cost of delay inputs set by flai from I-0327. time_lost_per_cycle 9m: 10m per occurrence × 1 occurrence ÷ 1.1 cycles of 168h (first reported 2026-11-09T23:09:47Z, 7.4 days before this story).

### Planning

Proposed solution, from I-0327's one instance: drop the `testing.Short()` gate from `TestRepositoryManifestLintsChangedFlaioverFiles` in `flai/internal/verify/select_test.go`. The test reads only `system-flow.yaml` and the file tree through `os.DirFS`, needs no git, and keeps its "monorepo not present" skip. The repository-manifest tests beside it already run under `-short` for this reason. The go-test tier, which `flai test` runs with `-short`, then catches a tier the test must list. The test is its own reproduction: under `-short` it runs, and fails while its `want` differs from what the manifest selects.

This is S-0693's fix again. S-0693 made it and passed its close-out, but the operator cancelled it on 2026-11-11 because its branch undid other work on main (TH-1066), so the gate is still on main. The only other `-short` gate in the package, `realRepo` in `paths_test.go`, runs real git and stays.

Touches, and where each came from:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/verify/select_test.go` | design (I-0327) | The test and its gate |
| `design/issues/I-0327-…-tier.md` | layout | `flai issue close` writes the issue file |
| `design/issues/summary.md` | layout | `flai issue close` regenerates the summary |

`flai touches suggest S-0827` listed co-changes led by `docs/users/flai.md` (4%) and `design/system/flai-cli.md` (4%), mostly through `summary.md`. None is kept: the fix changes no tier, no selection code, and no documentation. No folder touch. No `after`: the story reads nothing another open story writes.

Figures:

- Forecast 30m, up from flai's 4m. flai's figure is 48 s per unit of size times size 5, and it leaves out the close-out, whose integration and smoke tiers a change under `flai/` selects. S-0693 made the same change in 21m of agent work, with a 6m47s close-out. Delivery 2026-11-17T14:56Z is flai's queue-based 14:30Z moved by the 26m added. The story is 16th in the pull order.
- Cost of delay value 22.50 USD a week, as `flai cod` gives it from the input of 9m per 168h cycle at 150 USD an hour. It stands: I-0327 has one instance.

### Accepted by the orchestrator

- Verified: dd24cb478b1b4b98142e16b659cdd6d9c506a3f5
- At: 2026-11-17T09:04:30Z

Verdict: pass
- 1: flai/internal/verify/select_test.go, flai/cmd/publish_judgement_test.go
- 2: design/issues/I-0327-testrepositorymanifestlintschangedflaioverfiles-skips-under-short-so-a-new-tier.md, design/issues/summary.md
