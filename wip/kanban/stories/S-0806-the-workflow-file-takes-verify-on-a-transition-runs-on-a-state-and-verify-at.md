---
id: S-0806
type: story
nature: feature
title: The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour
status: backlog
parent: E-0034
owner: alex
created: 2026-11-16T23:42:03Z
updated: 2026-11-17T07:46:44Z
transitions: []
tags: [flai, template]
topics: [verify, workflow]
touches: [flai/internal/workflow/workflow.go, flai/internal/workflow/workflow_test.go, flai/internal/workflow/baseline.go, flai/internal/workflow/suites.go, flai/internal/workflow/suites_test.go, flai/internal/workflow/baseline/board.md, flai/internal/workflow/baseline_test.go, template/root/design/workflows/board.md, template/CHANGELOG.md, flai/internal/workitem/front-matter-fields.txt, flai/internal/workitem/fields_test.go, flai/internal/upgrade/upgrade_test.go, flai/internal/check/check_workflow.go, flai/internal/check/workflow_test.go, design/system/workflow.md, design/system/work-hierarchy.md, design/system/flai-cli.md, docs/operators/runbooks/change-the-board.md, docs/users/flai.md]
after: [S-0805]
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
      seconds: 55
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 49
          output: 6855
          cache_read: 5304789
          cache_write: 106543
          cost: 1.9709
    - kind: orchestrator
      seconds: 14
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 37
          cache_read: 211181
          cache_write: 1913
          cost: 0.0653
cost_of_delay:
  value: 110.09
  by: planner-S-0806
  at: 2026-11-16T23:50:08Z
forecast:
  duration: 1h20m
  delivery: 2026-11-17T11:16:00Z
  basis: "Its own forecast of 1h20m; 6th in the pull order with an in-progress limit of 3, behind S-0750, S-0805, S-0819, S-0241, S-0752, S-0754, S-0756 and S-0758."
  by: flai
  at: 2026-11-17T07:46:44Z
finalized:
  by: orchestrator
  at: 2026-11-16T23:56:26Z
---
# S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §2, §4, §6, §8, §9: the keys and their checks only. Nothing reads them to decide a run yet; the stories after this one do.

## Acceptance criteria

- [ ] The workflow package parses `verify:` on a transition (suites or `all`, with a severity), `runs:` on a state (`suites:` and `kluster:` of `deploy` and `checks`), and top-level `verify_at:` with `batch`, `publish`, and `ci`, and still refuses any other key.
- [ ] The baseline and the template's file give `verify: all` at `warn` on the move to review, `verify: all` on the acceptance, `batch: all`, `ci: all`, no `publish`, and `runs: {kluster: [deploy, checks]}` on `review`, and the two files stay equal.
- [ ] `flai check` reports a suite the file names that no tier has as an error, and a suite of a tier that is not `manual` which no transition, state, or point names as a warning.
- [ ] A flai from before this story refuses a file that uses the keys with words that name the flai that reads them, and `flai upgrade` does not write the keys into a project's own file.
- [ ] The workflow design and the reference describe the keys.

## Tasks
- T-3469 The workflow package parses verify on a transition, runs on a state, and verify_at, and lists the suites a file names
- T-3471 The baseline and the template's board give verify, runs, and verify_at as today's behaviour, and stay equal
- T-3473 The workflow file's keys raise flai.minimum on publish, and flai upgrade keeps a project's file without them
- T-3474 flai check reports a suite the workflow file names that no tier has, and a non-manual tier's suite nothing names
- T-3477 The workflow design, the board runbook, and the check codes describe verify, runs, verify_at, and the suite findings

## Notes

- After S-0805: the checks of the file's suites read the suite S-0805 gives each tier.

### Planning

Touches, file by file, no folder:

| Touch | From | Why |
|-------|------|-----|
| `flai/internal/workflow/workflow.go`, `workflow_test.go` | layout | The YAML shapes of a state, a transition, and the file, decoded strictly, and their parse-time checks (criterion 1) |
| `flai/internal/workflow/baseline.go` | layout | `Clone` deep-copies the new lists, which `TestBaselineIsACopy` checks |
| `flai/internal/workflow/suites.go`, `suites_test.go` | layout | New: the suites a file names, by transition, state, and point, with `all`, for the check and the stories after this one |
| `flai/internal/workflow/baseline/board.md`, `template/root/design/workflows/board.md`, `flai/internal/workflow/baseline_test.go` | design | ADR-0312 §8; `TestBaselineIsTheTemplatesBoard` keeps the two equal (criterion 2) |
| `template/CHANGELOG.md` | layout | The template's board changes |
| `flai/internal/workitem/front-matter-fields.txt`, `fields_test.go` | design | What a publish reads to raise `flai.minimum`, so an older flai stops with `manifest.TooOldError` naming the release ([work-hierarchy.md](../../../design/system/work-hierarchy.md), criterion 4) |
| `flai/internal/upgrade/upgrade_test.go` | layout | `upgrade.kept` already never writes a project's file (ADR-0307); a test pins it for a file without the keys (criterion 4) |
| `flai/internal/check/check_workflow.go`, `flai/internal/check/workflow_test.go` | layout | The two suite findings, against the manifest's tiers (criterion 3) |
| `design/system/workflow.md`, `docs/operators/runbooks/change-the-board.md` | design | The workflow design and the operator's reference for the file (criterion 5) |
| `docs/users/flai.md`, `design/system/flai-cli.md` | layout | The tables of `flai check`'s workflow codes |
| `design/system/work-hierarchy.md` | design | What the fields file lists, once it lists the workflow file's keys |

`flai touches suggest S-0806` gave nothing: the story declared no touches. The co-change of the paths above was not run.

`after`: S-0805. This story reads the `suite` that S-0805 adds to a tier in `flai/internal/manifest/manifest_tests.go` (`TestTier`, defaulting to the tier's name), which the check of criterion 3 compares the file's suites with. S-0805 also accepts ADR-0312, which this story implements.

Forecast: 1h20m, delivery 2026-11-17T06:15:00Z. `flai forecast` gave 45m with the 19 touches, from 110 s per unit of size times 24, and delivery 05:40Z. Raised to 1h20m: the five workflow stories of E-0034 most like this one, S-0737, S-0738, S-0741, S-0748, and S-0755, took 34m to 1h30m of agent time each, about 1h19m at the median, against forecasts of 29m to 53m. This one has as many parts: a parser, two equal files, a check, the minimum, and four documents. The delivery is flai's with the 35m added.

Cost of delay: 110.09 USD a week, as `flai cod` gives it once the forecast is 1h20m: S-0806's share of E-0034's 1500 USD a week, from the operator's input there of 10h lost per cycle, weighted by its 1h20m against the 18h10m forecast over the epic's 20 open stories without inputs. It stands as flai gives it.
