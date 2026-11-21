---
id: S-0811
type: story
nature: feature
title: The release's publish runs the suites verify_at.publish lists at the main branch's head before it tags, and publishes nothing when one fails
status: in-progress
parent: E-0034
owner: alex
created: 2026-11-16T23:42:07Z
updated: 2026-11-21T18:36:48Z
transitions:
  - to: ready
    at: 2026-11-17T08:27:38Z
    by: alex
  - to: in-progress
    at: 2026-11-18T08:12:44Z
    by: agent-S-0811
  - to: review
    at: 2026-11-18T09:07:52Z
    by: agent-S-0811
  - to: in-progress
    at: 2026-11-21T18:13:45Z
    by: flai
  - to: review
    at: 2026-11-21T18:32:45Z
    by: agent-S-0811
  - to: in-progress
    at: 2026-11-21T18:36:48Z
    by: flai
tags: [flai, template]
topics: [verify, workflow, release]
touches: [flai/cmd/release.go, flai/cmd/release_publish_suites.go, flai/cmd/release_publish_suites_test.go, docs/users/flai-reference.md, flai/internal/hostapi/writes_outcome.go, flai/internal/hostapi/writes_test.go, design/workflows/formulas/publish.md, template/root/design/workflows/formulas/publish.md, flai/cmd/publish_formula_suites_test.go, template/CHANGELOG.md, design/system/pushing-from-the-board.md, design/system/workflow.md, docs/operators/runbooks/change-the-board.md, docs/users/flai.md, flai/internal/hostapi/writes_actions.go, design/issues/I-0408-sibling-stories-of-one-adr-each-rewrite-the-same-nothing-else-reads-the-keys.md, design/issues/summary.md]
after: [S-0806]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4551
  priced_by: harness
  turns:
    - day: 2026-11-18
      ceremony: 1
      work: 66
    - day: 2026-11-21
      ceremony: 7
      test_runs: 4
      work: 36
  models:
    - model: claude-opus-5-5
      input: 472
      output: 111254
      cache_read: 29115974
      cache_write: 814611
      cost: 13.5219
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
      seconds: 377
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 99
          output: 702
          cache_read: 3777394
          cache_write: 487634
          cost: 1.3398
        - model: claude-sonnet-5-5
          input: 12
          output: 139
          cache_read: 83096
          cache_write: 64073
          cost: 0.1193
cost_of_delay:
  value: 90.11
  by: planner-S-0811
  at: 2026-11-16T23:47:39Z
forecast:
  duration: 34m
  delivery: 2026-11-18T08:14:00Z
  basis: "Its own forecast of 34m; 1st in the pull order with an in-progress limit of 3, behind S-0834, S-0840 and S-0842."
  by: flai
  at: 2026-11-18T07:32:43Z
finalized:
  by: orchestrator
  at: 2026-11-16T23:48:10Z
---
# S-0811 The release's publish runs the suites verify_at.publish lists at the main branch's head before it tags, and publishes nothing when one fails

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §6. Today publishing runs no tier.

## Acceptance criteria

- [x] With `verify_at.publish` set, the publish of what is pending runs those suites at the main branch's head before any tag or version file is written, with the output kept, and exits with its own code, naming the tier, when one fails.
- [x] The publish formula stops on that code and its story says which suite failed.
- [x] With no `publish` in the file, as on the baseline, publishing runs no tier, as before this story.
- [x] The release design and the runbook describe it.

## Tasks
- T-3441 flai release --pending runs the suites verify_at.publish lists at the main branch's head before it applies or tags, and exits with a code of its own naming the failed tier
- T-3442 The dashboard's Publish and release_publish answer a publish a failed suite refused with the suite, the tier, and the log
- T-3443 The publish formula's publish step fails on a failed suite's exit, so its instance blocks and asks, naming the suite and tier
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
- 2026-11-21T18:13:45Z: moved to in-progress: rebase and resolve job S-0811-20261121T181345Z: the rebase of story/S-0811 onto main stopped on conflicts in design/system/workflow.md, docs/operators/runbooks/change-the-board.md, which still waits in the worktree: a rebase of story/S-0811 was already in progress in /home/alex/git/bytepunx/system-flow/.flai-cache/worktrees/S-0811, which flai left as it found it; its agent settles it in the worktree and syncs it with flai stream sync S-0811
- 2026-11-21T18:36:48Z: moved to in-progress: rebase and resolve job S-0811-20261121T183647Z: the rebase of story/S-0811 onto main stopped on conflicts in design/system/workflow.md, docs/operators/runbooks/change-the-board.md, which flai aborted, so the branch is as it was; its agent settles it in the worktree and syncs it with flai stream sync S-0811
