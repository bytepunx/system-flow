---
id: ADR-0183
title: "The merge queue's batch and each bisection step run their tiers with MERGE_QUEUE_STORIES naming their stories, and the repository check, the markdown lint, and the integration checks scope to them as a close-out's scope to its story, recording no issue"
status: accepted
date: 2026-10-21
supersedes: []
superseded_by: []
refines: [ADR-0085, ADR-0174]
---

# ADR-0183 The merge queue's batch and each bisection step run their tiers with MERGE_QUEUE_STORIES naming their stories, and the repository check, the markdown lint, and the integration checks scope to them as a close-out's scope to its story, recording no issue

## Context

A close-out runs its tiers with `CLOSE_OUT_STORY` set to the story, and the checks inside them scope to it: `scripts/check.sh` runs `flai check --strict --story "$CLOSE_OUT_STORY" --record-issues`, so a finding outside the story is a note recorded in an issue and fails nothing ([ADR-0085](0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md)). `scripts/lint-md.sh`, `TestMonorepoIsClean`, and the repository lint test scope the same way.

The merge queue's batch runs the tiers, those marked `queue` among them, in a checkout that is no story's ([ADR-0174](0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md)). `verify.Tiers` set no `CLOSE_OUT_STORY`, so every check in the batch ran unscoped, and any finding on main outside the batch's stories, such as a cancelled story not yet archived or an open issue no story links, failed every batch. The bisection then sent a sound story back each time (I-0174).

## Decision

The merge queue's batch and each bisection step run their tiers with `MERGE_QUEUE_STORIES` set to the stories they verify, and the checks inside the tiers scope to those stories as a close-out's scope to its story, recording no issue.

- The batch gives its verifier the stories it kept, in batch order, and each bisection step the stories that step verifies. `CheckoutVerifier` sets `MERGE_QUEUE_STORIES` to them, comma-separated (`verify.MergeQueueStoriesEnv`), and still sets no `CLOSE_OUT_STORY`.
- `flai check --story` takes several stories, comma-separated or repeated. A finding is inside when it is inside any of them and outside only when it is outside every one, with each rule a scoped run passes over treated the same way. `--record-issues` is refused with more than one story.
- With `CLOSE_OUT_STORY` unset and `MERGE_QUEUE_STORIES` set, `scripts/check.sh` (and the template's) runs `flai check --strict --story "$MERGE_QUEUE_STORIES"` without `--record-issues`. `scripts/lint-md.sh` lints `wip/` only where the batch's checkout changes it. `TestMonorepoIsClean` and `TestRepositoryLintsClean` scope to the union of the stories' changes.
- `CLOSE_OUT_STORY` stays the close-out's alone, since it leads the scripts to record issues.

## Consequences

- A batch fails only on what its stories bring, as a close-out does, so a finding on main no longer sends a sound story back.
- A batch records no issue: its worktree is never committed, and files written there would stop the bisection's checkouts. The findings outside its stories are printed as notes in the tier's log.
- A project made from the template needs a flai that has this change. With an older flai the batch sets no `MERGE_QUEUE_STORIES`, and the check is unscoped as before.

## Alternatives considered

- Setting `CLOSE_OUT_STORY` in the batch, as I-0174 offered: it names one story, and it records issues in a worktree that is never committed.
- Clearing main's findings by hand before each batch: it fixes one batch, not the next finding that reaches main.
