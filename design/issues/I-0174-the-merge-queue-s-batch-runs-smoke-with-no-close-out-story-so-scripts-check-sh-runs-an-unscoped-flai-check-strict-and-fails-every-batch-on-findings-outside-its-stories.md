---
id: I-0174
title: The merge queue's batch runs smoke with no CLOSE_OUT_STORY, so scripts/check.sh runs an unscoped flai check --strict and fails every batch on findings outside its stories
class: blocker
status: closed
count: 2
cost: 30m
first_reported: 2026-10-21T06:29:55Z
last_reported: 2026-10-21T07:20:28Z
updated: 2026-10-21T08:00:37Z
---

# I-0174 The merge queue's batch runs smoke with no CLOSE_OUT_STORY, so scripts/check.sh runs an unscoped flai check --strict and fails every batch on findings outside its stories

## Description
The merge queue's batch runs smoke with no CLOSE_OUT_STORY, so scripts/check.sh runs an unscoped flai check --strict and fails every batch on findings outside its stories

## Instances

### 2026-10-21T06:29:55Z
Story: S-0455.
Merge queue batch 20261021-061843 sent S-0455 back at smoke (TH-0551), at the commit b1d7174e whose own close-out passed smoke at 06:07Z; the thread quoted only "smoke: building bin/flai" (I-0172). Cause, reproduced in S-0455's worktree: `scripts/check.sh` passes `--story "$CLOSE_OUT_STORY" --record-issues` to `flai check --strict` only when CLOSE_OUT_STORY is set, and the batch's verification "sets no CLOSE_OUT_STORY" (design/system/flai-cli.md, The merge queue pass). `env -u CLOSE_OUT_STORY scripts/check.sh` exits 1 on six warnings, none in S-0455: item.archive S-0452 (cancelled, not archived), threads.archived TH-0526 (open on archived S-0404), issues.no-story I-0143, I-0144, I-0145, and the components.examples note for flaiover. With CLOSE_OUT_STORY=S-0455 the same check passes, and the rest of smoke (template, markdown lint, installer, LiteLLM gateway) passed with flai serve's PATH and keys. So every batch fails smoke, whatever its stories, while main carries any finding a scoped check passes over, and the bisection sends a sound story back each time. Fix options: the batch verifier sets CLOSE_OUT_STORY (or a list of its stories) for the tiers, or check.sh scopes to a batch's stories; until then, clearing the six findings on main (flai archive, resolving TH-0526, stories or closes for I-0143..I-0145) lets a batch pass.

### 2026-10-21T07:20:28Z
Story: S-0453.
Merge queue batches 20261021-065925, 20261021-070546, and 20261021-071126 each sent S-0453 back at smoke (TH-0555), at 83a35607, whose own close-out passed smoke at 06:05Z. The thread quoted only 'smoke: template render and check' (I-0172). Reproduced with scripts/smoke.sh in a scratch worktree of story/S-0453 rebased onto main d53c0c13: the template step passes and 'smoke: repository check' fails, an unscoped flai check --strict counting five warnings no story owns: item.archive on cancelled S-0452, threads.archived on TH-0526 (S-0404 archived), and issues.no-story on I-0143, I-0144, and I-0145. Every system-flow batch since 01:44Z (S-0454, S-0455, S-0453 three times) stopped at smoke.

## Remediation

Story S-0464 remediates this issue, created from it at 2026-10-21T07:19:58Z.
Closed 2026-10-21T08:00:37Z: S-0464 (ADR-0183): the merge queue's batch and each bisection step run their tiers with MERGE_QUEUE_STORIES naming their stories (T-1974), flai check --story takes several stories (T-1973), and scripts/check.sh, scripts/lint-md.sh, TestMonorepoIsClean, and TestRepositoryLintsClean scope to them, recording no issue (T-1975). TestTheBatchsTiersAreGivenItsStoriesAndNoCloseOut and TestEachBisectionStepIsVerifiedGivenItsStories reproduce the cause, a batch's tiers given no scope; TestCheckStoryTakesSeveralStories shows a finding outside both stories no longer fails a scoped --strict run.
