---
id: S-0692
type: story
nature: improvement
title: No analysis.schedule or analysis.formula is set, so the analyzer ran once in thirty days
status: backlog
owner: alex
created: 2026-11-10T08:49:03Z
updated: 2026-11-10T10:55:48Z
transitions: []
tags: [flai, analyzer, manifest]
topics: [analysis, cli]
touches: [flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/check/check_roles.go, flai/cmd/manifest.go, flai/cmd/manifest_test.go, flai/cmd/serve_actions_test.go, docs/users/flai-reference.md, docs/users/flai.md, docs/operators/settings.md, docs/operators/index.md, design/system/strategic-agents.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, template/CHANGELOG.md, system-flow.yaml, design/issues/I-0292-no-analysis-schedule-or-analysis-formula-is-set-so-the-analyzer-ran-once-in.md, design/issues/summary.md]
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
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3102
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 14
          output: 2840
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2h
    by: orchestrator
    at: 2026-11-10T10:55:48Z
forecast:
  duration: 38m
  delivery: 2026-11-11T00:57:00Z
  basis: "Its own forecast of 38m; 23rd in the pull order with an in-progress limit of 3, behind S-0665, S-0708, S-0713, S-0667, S-0674, S-0704, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663, S-0666, S-0675, S-0676, S-0677, S-0678, S-0679, S-0681, S-0683, S-0686, S-0688, S-0689, S-0690 and S-0691."
  by: flai
  at: 2026-11-10T10:54:40Z
---
# S-0692 No analysis.schedule or analysis.formula is set, so the analyzer ran once in thirty days

## Goal

This story remediates [I-0292](../../../design/issues/I-0292-no-analysis-schedule-or-analysis-formula-is-set-so-the-analyzer-ran-once-in.md), "No analysis.schedule or analysis.formula is set, so the analyzer ran once in thirty days", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0292 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0292 is closed with `flai issue close I-0292 --reason` saying what fixed it

## Tasks
- T-3031 flai manifest set writes analysis.formula and refuses a formula no file defines or one with an input that has no default
- T-3034 Document that flai manifest set and the Analyzer page write analysis.formula, and regenerate the command reference
- T-3037 Set this project's analysis.schedule and analysis.formula analysis-cadence through flai manifest set
- T-3040 Close I-0292 with the reason that the analysis schedule now starts analysis-cadence

## Notes

The analyzer found I-0292: the analysis the goal links gives its evidence.

Strategic usage carried over from I-0292: analyzer 42s, 0.8156 USD.

### Planning

Proposed solution, on TH-1038: set this project's `analysis.schedule` and `analysis.formula: analysis-cadence`, after making `flai manifest set` write `analysis.formula`. Today the formula is "changed by hand" (`design/system/strategic-agents.md` § The analysis schedule), and `CLAUDE.md` gives that key to flai. The new key is the part a test reproduces (T-3031).

Touches, all files, no folder:

- Layout, from the explorer's read of the code: `flai/internal/manifest/settings.go` and `settings_test.go` (the settings catalog), `flai/internal/formulas/formulas.go` and `formulas_test.go` (the shared check), `flai/internal/check/check_roles.go` (`manifest.analysis` today), `flai/cmd/manifest.go` and `manifest_test.go`, and `flai/cmd/serve_actions_test.go`, which lists the analysis keys.
- Design and docs that say the formula is set by hand or list the keys: `design/system/strategic-agents.md`, `design/system/project-manifest.md`, `docs/operators/settings.md`, `docs/operators/index.md`, `docs/users/flai.md`, and `docs/users/flai-reference.md`, generated from the command's help.
- Co-change (`flai touches suggest`): `design/system/flai-cli.md` (23 percent, its `flai manifest set` row lists the keys) and `design/system/flaiover-dashboard.md` (7 percent, the Analyzer page's settings panel names its keys). The rest it listed, such as `flai/internal/hostapi/writes.go`, change with these files for other reasons: the host API takes the catalog through `manifest.SettingFor`, with no list of its own.
- Goal and criteria: `system-flow.yaml`, written by `flai manifest set`; the issue file and `design/issues/summary.md`, for the close; `template/CHANGELOG.md`, for the flai change.

No `after`: the story reads no other story's output.

Forecast: 38m, flai's figure, kept. Its size counts 20 touches. Nine of them are one-line document edits, but T-3031 moves a check into a shared function and T-3034 regenerates the command reference, so the figure is not lowered.

Cost of delay: no value yet. The story has no inputs and no epic, and I-0292's Impact gives no figure flai parses. TH-1038 asks the operator. The recommendation there is `time_lost_per_cycle: 2h`. The value is worked out with `flai cod` once the inputs are set.

The cadence is money, and the operator's. TH-1038 recommends weekly, `0 6 * * 1`, about 3.5 USD a month at the last run's 0.82 USD. T-3037 sets what the operator chooses.
