---
id: S-0592
type: story
nature: improvement
title: flai test on a change to flai's host API runs none of flaiover's server tests that drive the real flai, so a changed gate in hostapi first fails in the close-out
status: done
owner: alex
created: 2026-11-07T15:15:58Z
updated: 2026-11-07T22:31:02Z
transitions:
  - to: ready
    at: 2026-11-07T22:04:23Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-07T22:04:54Z
    by: agent-S-0592
  - to: review
    at: 2026-11-07T22:29:29Z
    by: agent-S-0592
  - to: done
    at: 2026-11-07T22:31:02Z
    by: orchestrator
tags: [flai, flaiover, tests]
touches: [system-flow.yaml, flai/internal/verify/select_test.go, flai/cmd/settings_doc_test.go, design/conventions/tooling.md, design/system/devex.md, design/issues/I-0304-flai-test-on-a-change-to-flai-s-host-api-runs-none-of-flaiover-s-server-tests.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1490
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 2
      test_runs: 1
      hand_edits: 2
      work: 35
  models:
    - model: claude-opus-5-5
      input: 136
      output: 25518
      cache_read: 15388230
      cache_write: 312251
      cost: 5.9839
  strategic:
    - kind: orchestrator
      seconds: 274
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 59
          output: 489
          cache_read: 1858580
          cache_write: 288383
          cost: 0.5436
        - model: claude-sonnet-5-5
          input: 6
          output: 43
          cache_read: 40474
          cache_write: 43899
          cost: 0.0613
cost_of_delay:
  inputs:
    time_lost_per_cycle: 16m
    by: flai
    at: 2026-11-07T15:15:58Z
  value: 40
  by: planner-S-0592
  at: 2026-11-07T22:01:21Z
forecast:
  duration: 45m
  delivery: 2026-11-07T22:58:00Z
  basis: "Its own forecast of 45m; 3rd in the pull order with an in-progress limit of 5, behind S-0595, S-0591 and S-0241."
  by: flai
  at: 2026-11-07T22:02:42Z
finalized:
  by: orchestrator
  at: 2026-11-07T22:03:33Z
---
# S-0592 flai test on a change to flai's host API runs none of flaiover's server tests that drive the real flai, so a changed gate in hostapi first fails in the close-out

## Goal

This story remediates [I-0304](../../../design/issues/I-0304-flai-test-on-a-change-to-flai-s-host-api-runs-none-of-flaiover-s-server-tests.md), "flai test on a change to flai's host API runs none of flaiover's server tests that drive the real flai, so a changed gate in hostapi first fails in the close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0304 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0304 is closed with `flai issue close I-0304 --reason` saying what fixed it

## Tasks
- T-2594 TestSettingsIndexIsComplete checks the settings index's keys under -short
- T-2595 A change to flai's host API selects a tier that runs flaiover's tests of the real flai
- T-2596 A change to the settings index selects a tier that runs TestSettingsIndex
- T-2597 The tier lists name the two new tiers, and I-0304 is closed

## Notes

Cost of delay inputs set by flai from I-0304. time_lost_per_cycle 16m: 8m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-11-07T11:41:27Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Proposed solution, from the issue's two instances:

1. A change to flai's host API selects flaiover's tests that run the built flai. Every such test file under `flaiover/src` imports `flaiover/src/lib/server/testing.ts`, which finds the binary (`FLAI_BIN`, else `bin/flai`). A new tier in `system-flow.yaml`, `flaiover-flai`, runs `scripts/flaiover-unit.sh related --run --reporter=json src/lib/server/testing.ts` in `flaiover`, a fixed command as `svelte-check` has, whenever a path under `flai/internal/hostapi/` or `flai/cmd/hostapi.go` changes. `flaiover-unit.sh` builds `bin/flai` first. The `flaiover` tier covers it, so `--all` runs those tests once. Widening the `vitest` tier's paths would not work: its `{files}` keeps only paths under `flaiover`, so a hostapi-only change leaves it out.
2. The settings index is checked where it changes. `TestSettingsIndexIsComplete` (`flai/cmd/settings_doc_test.go`) skips under `-short`, the flag the `go-test` tier passes. Its key checks are reflection over `config.Config` and `manifest.Manifest` and cost nothing; only the source walks for environment variables cost time. The key checks run under `-short`, so a change under `flai/internal/config` or `flai/internal/manifest` runs them through `go-test`'s importers. A second tier, `settings-index`, runs `go test -short -run '^TestSettingsIndex' ./cmd` with a timeout when `docs/operators/settings.md` changes, which selects no Go tier today.

Touches, and where each came from (`flai touches suggest` listed nothing, since the story declared no touches):

| Touch | Source |
|-------|--------|
| `system-flow.yaml` | design: the two tiers; layout: S-0530 added `svelte-check` there |
| `flai/internal/verify/select_test.go` | layout: the repository-manifest tests, such as `TestRepositoryManifestTypeChecksChangedFlaioverFiles`, reproduce a tier's selection |
| `flai/cmd/settings_doc_test.go` | design: the issue's second instance names it |
| `design/conventions/tooling.md` | layout: its project addition lists the tiers `flai test` runs |
| `design/system/devex.md` | layout: its close-out row lists the tiers a flai or flaiover story runs |
| `design/issues/I-0304-…md`, `design/issues/summary.md` | criteria: the issue's close |

No folder touch, and no `after`: the story reads no other story's output.

Figures:

- Forecast 45m, adjusted from flai's 3m. flai sized the story at 2 with 0 touches, so it counted only the criteria. S-0530, which added one selected tier with its repository-manifest test, took 1135 s; this story adds two tiers, splits a Go test, and closes the issue.
- Delivery 2026-11-07T22:50:00Z: the forecast from now, since the in-progress limit has room. It moves once the operator finalizes the draft.
- Cost of delay 40 USD a week, as `flai cod` gives it from the issue's input of 16m a cycle at 150 USD an hour. It stands: the issue's count and cost are what it rests on.

### Accepted by the orchestrator

- Verified: bd7a67504cce7445c4552a93b2414c55d93df0e9
- At: 2026-11-07T22:31:02Z

Verdict: pass at bd7a67504cce7445c4552a93b2414c55d93df0e9
- 1: system-flow.yaml, flai/cmd/settings_doc_test.go, flai/internal/verify/select_test.go, design/conventions/tooling.md, design/system/devex.md
- 2: design/issues/I-0304-flai-test-on-a-change-to-flai-s-host-api-runs-none-of-flaiover-s-server-tests.md, design/issues/summary.md
