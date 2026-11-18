---
id: S-0812
type: story
nature: feature
title: flai test --at ci runs the suites verify_at.ci lists, and the template's and this repository's CI call it
status: review
parent: E-0034
owner: alex
created: 2026-11-16T23:42:08Z
updated: 2026-11-18T09:15:42Z
transitions:
  - to: ready
    at: 2026-11-17T08:27:41Z
    by: alex
  - to: in-progress
    at: 2026-11-18T08:30:35Z
    by: agent-S-0812
  - to: review
    at: 2026-11-18T09:15:42Z
    by: agent-S-0812
tags: [flai, template]
topics: [verify, workflow]
touches: [flai/cmd/test.go, flai/cmd/test_test.go, flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/internal/verify/verify.go, docs/users/flai-reference.md, docs/operators/settings.md, template/root/.github/workflows/system-flow-check.yml, template/CHANGELOG.md, design/system/template.md, flai/cmd/new_test.go, scripts/template-test.sh, ".github/workflows/system-flow-check.yml", ".github/workflows/flai.yml", ".github/workflows/flaiover.yml", design/tech/ci.md, design/system/devex.md, design/conventions/code-quality.md, docs/users/flai.md, design/system/flai-cli.md, design/system/workflow.md, flai/internal/workflow/suites.go, flai/internal/workflow/suites_test.go, docs/operators/runbooks/change-the-board.md]
after: [S-0806]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2743
  priced_by: harness
  turns:
    - day: 2026-11-18
      ceremony: 1
      test_runs: 1
      work: 59
  models:
    - model: claude-opus-5-5
      input: 330
      output: 93918
      cache_read: 23094915
      cache_write: 823632
      cost: 11.6272
  strategic:
    - kind: planner
      seconds: 15
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 22
          output: 1125
          cache_read: 2354106
          cache_write: 29002
          cost: 0.7049
    - kind: orchestrator
      seconds: 118
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 22
          output: 212
          cache_read: 1032589
          cache_write: 111616
          cost: 0.3574
        - model: claude-sonnet-5-5
          input: 4
          output: 30
          cache_read: 14191
          cache_write: 34169
          cost: 0.0392
cost_of_delay:
  value: 115.74
  by: planner-S-0812
  at: 2026-11-16T23:48:32Z
forecast:
  duration: 50m
  delivery: 2026-11-18T09:16:00Z
  basis: "Its own forecast of 50m; 2nd in the pull order with an in-progress limit of 3, behind S-0811, S-0834, S-0842 and S-0809."
  by: flai
  at: 2026-11-18T08:13:57Z
finalized:
  by: orchestrator
  at: 2026-11-16T23:51:54Z
---
# S-0812 flai test --at ci runs the suites verify_at.ci lists, and the template's and this repository's CI call it

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §6. Today the CI workflow files hold commands of their own, apart from the manifest's tiers.

## Acceptance criteria

- [x] `flai test --at <point>` runs the tiers of the suites the file lists for `ci`, `batch`, or `publish`, over all paths.
- [x] The template's CI workflow calls `flai test --at ci` once flai is installed in it, and a project made from the template runs its suites in CI with no command of its own.
- [x] This repository's CI workflows call it where a suite stands for what the job ran, and a note in the design says which jobs keep commands of their own and why.
- [x] The user guide says how a project's CI runs its suites.

## Tasks

- T-3437 flai test --at ci, batch, or publish runs the tiers of the suites the workflow's verify_at lists, over all paths
- T-3438 The template's CI workflow installs flai and runs flai test --at ci in place of commands of its own
- T-3439 This repository's CI jobs call flai test --at ci where one suite stands for what the job ran, and the CI design says which keep their own commands and why
- T-3440 The user guide says how a project's CI runs its suites, and the CLI and workflow design describe flai test --at

## Notes

- After S-0806: the command reads verify_at, which S-0806 parses.

### Planning

After: S-0806 stays. This story reads the `verify_at:` key that S-0806 parses into the workflow package, and the baseline's `ci: all` that S-0806 writes. S-0806 itself waits for S-0805, so the tier suites and `flai test --suite` that this story narrows with are on main before it starts. That is why S-0805 is not named here.

Touches. The story declared none. Each one below is a file, and no folder touch is kept.

| Touch | Source | Why |
|-------|--------|-----|
| `flai/cmd/test.go`, `flai/cmd/test_test.go` | design, co-change | The `--at` flag (criterion 1) |
| `flai/internal/verify/select.go`, `select_test.go`, `verify.go` | co-change | The tier selection by suite over all paths |
| `flai/internal/workflow/workflow.go`, `workflow_test.go` | layout | A lookup of `verify_at` by point, if S-0806 leaves none |
| `docs/users/flai-reference.md`, `docs/operators/settings.md` | co-change, layout | Generated flag reference and flag index for a new flag |
| `template/root/.github/workflows/system-flow-check.yml` | design | The template's CI (criterion 2) |
| `template/CHANGELOG.md`, `design/system/template.md` | layout | A template change is logged and described |
| `flai/cmd/new_test.go`, `scripts/template-test.sh` | layout | Tests that the rendered project's CI runs its suites |
| `.github/workflows/system-flow-check.yml`, `flai.yml`, `flaiover.yml` | co-change, design | This repository's CI (criterion 3) |
| `design/tech/ci.md`, `design/system/devex.md` | co-change, design | The note on which jobs keep their own commands |
| `design/conventions/code-quality.md` | layout | Its project additions say what CI runs |
| `docs/users/flai.md`, `design/system/flai-cli.md`, `design/system/workflow.md` | co-change, design | The user guide (criterion 4) and the design of the command and the key |

Forecast: 50m, delivery 2026-11-17T03:28:00Z, as `flai forecast` gives it with the touches above: size 27 at the median 110 s per unit. Before the touches, it gave 8m at size 4. It stands. About half the size is docs and CI files, which go fast, and the close-out runs integration and smoke because the story touches `flai/` and `template/`. The delivery assumes the 14 stories ahead are pulled in order.

Cost of delay: 115.74 USD a week, as `flai cod` gives it. That is this story's share of E-0034's 1500 USD a week (the operator's 10h lost per cycle), by its 50m of the epic's 10h48m forecast over its 20 open stories without inputs. It stands: the story has no inputs of its own.
