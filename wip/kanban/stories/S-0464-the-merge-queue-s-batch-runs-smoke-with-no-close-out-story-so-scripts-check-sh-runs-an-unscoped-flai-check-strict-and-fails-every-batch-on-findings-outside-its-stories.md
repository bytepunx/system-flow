---
id: S-0464
type: story
nature: remediation
title: The merge queue's batch runs smoke with no CLOSE_OUT_STORY, so scripts/check.sh runs an unscoped flai check --strict and fails every batch on findings outside its stories
status: in-progress
owner: alex
created: 2026-10-21T07:19:58Z
updated: 2026-10-21T07:31:41Z
transitions:
  - to: ready
    at: 2026-10-21T07:27:18Z
    by: alex
  - to: in-progress
    at: 2026-10-21T07:31:41Z
    by: agent-S-0464
tags: [flai, template]
topics: [cli, template, orchestration]
touches: [flai/cmd/check.go, flai/cmd/check_test.go, flai/internal/check/story.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/internal/check/check_test.go, flai/internal/mdlint/repo_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/bisect_test.go, flai/internal/verify/story.go, scripts/check.sh, scripts/lint-md.sh, scripts/README.md, template/root/scripts/check.sh, template/root/scripts/README.md, template/CHANGELOG.md, docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md, design/system/devex.md, design/adrs/README.md, design/issues/I-0174-the-merge-queue-s-batch-runs-smoke-with-no-close-out-story-so-scripts-check-sh-runs-an-unscoped-flai-check-strict-and-fails-every-batch-on-findings-outside-its-stories.md, design/issues/summary.md]
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
      seconds: 197
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 28
          output: 402
          cache_read: 5278247
          cache_write: 14660
          cost: 1.3271
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-10-21T07:19:58Z
  value: 50
  by: planner-S-0464
  at: 2026-10-21T07:31:24Z
forecast:
  duration: 1h30m
  delivery: 2026-10-21T09:15:00Z
  basis: "flai forecast's 41m (89 s per unit over 31 large remediation stories on claude-opus-5-5, size 27) raised to 1h30m for four tasks in three layers across two Go packages, the scripts, the template, and a new ADR; first in the pull order behind S-0453."
  by: planner-S-0464
  at: 2026-10-21T07:31:24Z
finalized:
  by: alex
  at: 2026-10-21T07:26:42Z
---
# S-0464 The merge queue's batch runs smoke with no CLOSE_OUT_STORY, so scripts/check.sh runs an unscoped flai check --strict and fails every batch on findings outside its stories

## Goal

This story remediates [I-0174](../../../design/issues/I-0174-the-merge-queue-s-batch-runs-smoke-with-no-close-out-story-so-scripts-check-sh-runs-an-unscoped-flai-check-strict-and-fails-every-batch-on-findings-outside-its-stories.md), "The merge queue's batch runs smoke with no CLOSE_OUT_STORY, so scripts/check.sh runs an unscoped flai check --strict and fails every batch on findings outside its stories". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0174 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0174 is closed with `flai issue close I-0174 --reason` saying what fixed it

## Tasks
- T-1973 flai check --story takes several stories, a finding outside every one of them a note
- T-1974 The merge queue's batch and bisection run their tiers with MERGE_QUEUE_STORIES naming the stories they verify
- T-1975 scripts/check.sh, scripts/lint-md.sh, and the integration checks scope to MERGE_QUEUE_STORIES in a batch, recording no issue
- T-1976 Record the batch's scoped checks in an ADR, the design, and the docs, and close I-0174

## Notes

Cost of delay inputs set by flai from I-0174. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-21T06:29:55Z, 0 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from I-0174's instance and TH-0555: the batch's verifier names the stories it verifies in a new variable, `MERGE_QUEUE_STORIES` (comma-separated), and the checks inside the tiers scope to them as a close-out's scope to `CLOSE_OUT_STORY`, but record no issue, since the batch's worktree is never committed. `CLOSE_OUT_STORY` stays a close-out's alone. `flai check --story` takes several stories so that one call scopes to the batch.

Touches, file by file; no folder touch is kept:

| Paths | Source |
|-------|--------|
| `scripts/check.sh` | declared by the title and I-0174 |
| `flai/internal/mergequeue/batch.go`, `bisect.go`, their tests; `flai/internal/verify/story.go` | design: `flai-cli.md` § The merge queue pass says `verify.Tiers` sets no `CLOSE_OUT_STORY`; `CheckoutVerifier` and `Tiers` live there |
| `flai/cmd/check.go`, `check_test.go`; `flai/internal/check/story.go`, `scope.go`, `scope_test.go` | layout: `--story`, `RunScoped`, and `ScopeToStory` take one story today |
| `flai/internal/check/check_test.go`, `flai/internal/mdlint/repo_test.go`, `scripts/lint-md.sh` | layout: the other readers of `CLOSE_OUT_STORY` inside the queue's tiers (integration and smoke) |
| `template/root/scripts/check.sh`, `template/root/scripts/README.md`, `template/CHANGELOG.md`, `scripts/README.md` | layout: the template's copy of `check.sh` and the READMEs that describe the variable |
| `docs/users/flai.md`, `docs/users/flai-reference.md`, `design/system/flai-cli.md`, `design/system/devex.md` | layout and `documentation.md`: the pages that describe `CLOSE_OUT_STORY` and `flai check --story` |
| `design/adrs/README.md` | design: a new ADR refining ADR-0085 and ADR-0174; `flai task done` widens the touches with its file |
| `design/issues/I-0174-….md`, `design/issues/summary.md` | criterion 2, closing the issue |

`flai touches suggest` found nothing to start from, then, from these paths, only broad co-change (`docs/operators/index.md`, `design/system/flaiover-dashboard.md`, `design/conventions/work-management.md`, each at 11% or less): none is added.

Forecast: 1h30m, raised from `flai forecast`'s 41m. Its median is over remediation stories at large size, but this one has four tasks in three layers, changes the merge queue's `Verifier` signature that its tests fake, and adds an ADR and a template change.

Cost of delay: 50 USD a week, as `flai cod` gives it from the inputs (20m lost per 168h cycle at 150 USD an hour). It stands. TH-0555 counts five failed batches, not I-0174's one, so the inputs may undercount; they are the operator's, and the plan's thread says so.
