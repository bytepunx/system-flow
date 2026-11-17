---
id: S-0811
type: story
nature: feature
title: The release's publish runs the suites verify_at.publish lists at the main branch's head before it tags, and publishes nothing when one fails
status: ready
parent: E-0034
owner: alex
created: 2026-11-16T23:42:07Z
updated: 2026-11-17T10:20:28Z
transitions:
  - to: ready
    at: 2026-11-17T08:27:38Z
    by: alex
tags: [flai, template]
topics: [verify, workflow, release]
touches: [flai/cmd/release.go, flai/cmd/release_publish_suites.go, flai/cmd/release_publish_suites_test.go, docs/users/flai-reference.md, flai/internal/hostapi/writes_outcome.go, flai/internal/hostapi/writes_test.go, design/workflows/formulas/publish.md, template/root/design/workflows/formulas/publish.md, flai/cmd/publish_formula_suites_test.go, template/CHANGELOG.md, design/system/pushing-from-the-board.md, design/system/workflow.md, docs/operators/runbooks/change-the-board.md, docs/users/flai.md]
after: [S-0806]
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
      seconds: 40
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 4
          output: 17
          cache_read: 285452
          cache_write: 5848
          cost: 0.0198
        - model: claude-opus-5-5
          input: 92
          output: 5256
          cache_read: 9658598
          cache_write: 94362
          cost: 2.9176
    - kind: orchestrator
      seconds: 63
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 15
          output: 175
          cache_read: 821199
          cache_write: 49799
          cost: 0.2669
cost_of_delay:
  value: 90.11
  by: planner-S-0811
  at: 2026-11-16T23:47:39Z
forecast:
  duration: 34m
  delivery: 2026-11-17T11:01:00Z
  basis: "Its own forecast of 34m; 1st in the pull order with an in-progress limit of 3, behind S-0806, S-0823 and S-0832."
  by: flai
  at: 2026-11-17T10:20:28Z
finalized:
  by: orchestrator
  at: 2026-11-16T23:48:10Z
---
# S-0811 The release's publish runs the suites verify_at.publish lists at the main branch's head before it tags, and publishes nothing when one fails

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §6. Today publishing runs no tier.

## Acceptance criteria

- [ ] With `verify_at.publish` set, the publish of what is pending runs those suites at the main branch's head before any tag or version file is written, with the output kept, and exits with its own code, naming the tier, when one fails.
- [ ] The publish formula stops on that code and its story says which suite failed.
- [ ] With no `publish` in the file, as on the baseline, publishing runs no tier, as before this story.
- [ ] The release design and the runbook describe it.

## Tasks
- T-3441 flai release --pending runs the suites verify_at.publish lists at the main branch's head before it applies or tags, and exits with a code of its own naming the failed tier
- T-3442 The dashboard's Publish and release_publish answer a publish a failed suite refused with the suite, the tier, and the log
- T-3443 The publish formula's publish step stops on a failed suite's exit, and its story records which suite and tier failed
- T-3445 The release design, the workflow design, the board runbook, and the user guide describe the suites a publish runs

## Notes

- After S-0806: the publish reads verify_at.publish, which S-0806 parses.

### Planning

After:

- S-0806, declared: this story reads `verify_at.publish` from the workflow file, which S-0806 parses into the workflow package. It also reads each tier's `suite` and the selection of a suite's tiers that S-0805 writes; S-0806 is itself after S-0805, so S-0805 is not named again here.

Touches, all files, no folder:

| Touch | From | Why |
|-------|------|-----|
| `flai/cmd/release.go` | co-change, layout | `publishPending` and `computeApplyAndTagPending` are where `flai release --pending` applies version files and tags; the suites run before them, and the command's help and a new exit code live here. |
| `flai/cmd/release_publish_suites.go` | layout | New file for reading `verify_at.publish`, selecting its suites' tiers, running them in the main checkout, and keeping the output. |
| `flai/cmd/release_publish_suites_test.go` | layout | New tests: a failing suite writes no tag or version file and exits with the new code naming the tier; no `publish` runs no tier. |
| `docs/users/flai-reference.md` | co-change (28%) | Generated from the command's help; `TestReferenceIsCurrent` fails until it is regenerated. |
| `flai/internal/hostapi/writes_outcome.go` | layout | The dashboard's Publish and the orchestrator's `release_publish` run `flai release --pending` and map its exits (`publishExits`); the new exit needs a payload that names the suite. |
| `flai/internal/hostapi/writes_test.go` | layout | Its test. |
| `design/workflows/formulas/publish.md` | criterion 2, design | The publish step's `stop` and the body's instruction to record the failed suite. |
| `template/root/design/workflows/formulas/publish.md` | design | The template ships the same formula. |
| `flai/cmd/publish_formula_suites_test.go` | layout | Beside `publish_judgement_test.go`: the publish step of an instance stops on the new exit. |
| `template/CHANGELOG.md` | co-change (15%) | The template's formula changes. |
| `design/system/pushing-from-the-board.md` | design, criterion 4 | Its opening current-state paragraph is the publishing design. |
| `design/system/workflow.md` | co-change (12%), design | The workflow file's section describes `verify_at`; `publish` gains what it does. |
| `docs/operators/runbooks/change-the-board.md` | criterion 4 | The runbook for changing the workflow file; it gains how to run suites at publish. |
| `docs/users/flai.md` | co-change (20%) | The publish section and the publish formula paragraph, which names the steps' `stop` exits. |

`flai touches suggest` was seeded with `flai/cmd/release.go`, `flai/internal/release/release.go`, and `design/workflows/formulas/publish.md`, since the story declared none. Left out of its list: `flai/internal/release/release.go` and its tests (the plan, apply, and tag stay as they are), `design/system/flai-cli.md`, `docs/operators/settings.md`, `design/conventions/git.md`, and the dashboard's files: none holds what this story changes.

Forecast 34m, delivery 2026-11-17T03:02:00Z, as `flai forecast` gives it with the 14 touches. The first figure, 8m, counted no touches. Kept as computed: the work is one command's path, one exit mapping, one formula in two copies, and four documents, the size of the large band's stories.

Cost of delay 90.11 USD a week, as `flai cod` gives it: S-0811's share of E-0034's 1500.00 USD a week, 34m of 9h26m over the epic's 20 open stories without inputs. Kept: this repository's own workflow file names no `publish` (ADR-0312 §8), so nothing here runs sooner for it, and no reason argues for more or less than its share.
