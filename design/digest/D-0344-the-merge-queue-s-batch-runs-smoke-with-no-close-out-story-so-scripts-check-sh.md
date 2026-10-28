---
id: D-0344
title: The merge queue's batch runs smoke with no CLOSE_OUT_STORY, so scripts/check.sh runs an unscoped flai check --strict and fails every batch on findings outside its stories
kind: cause
sources: [I-0174]
topics: [cli, template, orchestration, flai, go, conventions, code, all]
paths: [flai/cmd/check.go, flai/cmd/check_test.go, flai/internal/check/story.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/internal/check/check_test.go, flai/internal/mdlint/repo_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/bisect_test.go, flai/internal/verify/story.go, scripts/check.sh, scripts/lint-md.sh, scripts/README.md, template/root/scripts/check.sh, template/root/scripts/README.md, template/CHANGELOG.md, docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md, design/system/devex.md, design/adrs/README.md, design/issues/I-0174-the-merge-queue-s-batch-runs-smoke-with-no-close-out-story-so-scripts-check-sh-runs-an-unscoped-flai-check-strict-and-fails-every-batch-on-findings-outside-its-stories.md, design/issues/summary.md, flai/internal/mergequeue/land_test.go, flai/internal/serve/mergequeue_test.go, design/adrs/0183-the-merge-queue-s-batch-and-each-bisection-step-run-their-tiers-with-merge.md]
date: 2026-10-21
status: absorbed
absorbed_by: "design/adrs/0183-the-merge-queue-s-batch-and-each-bisection-step-run-their-tiers-with-merge.md#decision"
asked_again: []
cited: []
---

# D-0344 The merge queue's batch runs smoke with no CLOSE_OUT_STORY, so scripts/check.sh runs an unscoped flai check --strict and fails every batch on findings outside its stories

I-0174 was closed on 2026-10-21.

The cause:

```text
The merge queue's batch runs smoke with no CLOSE_OUT_STORY, so scripts/check.sh runs an unscoped flai check --strict and fails every batch on findings outside its stories
```

The fix:

```text
S-0464 (ADR-0183): the merge queue's batch and each bisection step run their tiers with MERGE_QUEUE_STORIES naming their stories (T-1974), flai check --story takes several stories (T-1973), and scripts/check.sh, scripts/lint-md.sh, TestMonorepoIsClean, and TestRepositoryLintsClean scope to them, recording no issue (T-1975). TestTheBatchsTiersAreGivenItsStoriesAndNoCloseOut and TestEachBisectionStepIsVerifiedGivenItsStories reproduce the cause, a batch's tiers given no scope; TestCheckStoryTakesSeveralStories shows a finding outside both stories no longer fails a scoped --strict run.
```

Source: I-0174.
