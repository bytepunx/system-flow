---
id: ADR-0221
title: "While the merge queue can batch, a hand acceptance is refused until the queue tiers ran at the story's head, flai verify --all runs only those tiers after a close-out at the same patch, and one flag lands the story without them, recorded in its Notes"
status: accepted
date: 2026-10-30
supersedes: []
superseded_by: []
refines: [ADR-0173, ADR-0174]
---

# ADR-0221 While the merge queue can batch, a hand acceptance is refused until the queue tiers ran at the story's head, flai verify --all runs only those tiers after a close-out at the same patch, and one flag lands the story without them, recorded in its Notes

## Context

While the merge queue can batch, `flai verify` and the close-out leave the tiers the manifest marks `queue` to the batch ([ADR-0174](0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md)). In this repository those are `integration` and `smoke`. ADR-0174's consequences say a hand acceptance is unchanged and runs no tier, so a story accepted by hand lands without them unless someone runs `flai verify --all` on it first.

[I-0203](../issues/I-0203-a-story-the-operator-accepts-outside-the-merge-queue-lands-without-the-queue-tiers-so-a-stale-generated-file-on-main-fails-every-later-batch-at-integration.md) records two instances:

- S-0491 changed `flai manifest set`'s help without regenerating `docs/users/flai-reference.md`. The operator accepted it by hand, so `integration` never ran on it. The next batch failed S-0490 at `TestReferenceIsCurrent`, though S-0490 changes no command.
- S-0416, accepted by hand, left a guard test expecting wording S-0535 had changed. The next batch failed S-0538 and S-0519 at `integration`, though neither touched the guard.

A batch that fails is bisected, and a batch of one blames its story. So main left red by a hand acceptance sends back every innocent story after it.

`flai verify --all` was no cheap remedy either: the report's fingerprint held the names of the tiers left to the queue, so a `--all` run after a close-out at the same patch reused nothing and ran every tier again.

## Decision

While the merge queue can batch, a hand acceptance is refused until the queue tiers ran at the story's head, `flai verify --all` runs only those tiers after a close-out at the same patch, and one flag lands the story without them, recorded in its Notes.

- **When.** The queue can batch by ADR-0174's three conditions, read from the configuration the running flai reads: the `merge-queue` host action is on for the project, `orchestration.permissions.accept_reviews` is on, and acceptance is local. With any of them off, a hand acceptance is unchanged.
- **Which acceptances.** The operator's, from `flai accept`, `flai move <S-nnnn> done`, the board, and the dashboard: every acceptance but the orchestrator's and the queue's landing (`--batch`). The orchestrator already enqueues instead while the queue batches (S-0403), and the batch runs every tier.
- **The blocker.** The acceptance reads the story's last `flai verify` report. It is a blocker, `queue_tiers`, when there is no report, the report is not at the story branch's head, it did not pass, or it did not run a tier marked `queue` that the story's paths select: one it left `queued`, or one that did not pass. Its text names the tiers and `flai verify <S-nnnn> --all`. `flai accept --dry-run` and the dashboard's preview list it, and the acceptance refuses it before anything is merged.
- **Reuse.** The report's fingerprint hashes the whole selection, the tiers run and those left to the queue alike. A run that leaves tiers to the queue and one that runs them all, at the same patch, share it. `flai verify --all` then reuses each tier the close-out passed and runs only the tiers it left `queued`. A change to a tier's command, or to the tiers as declared, still changes the fingerprint and reuses nothing.
- **The override.** `flai accept` and `flai move <S-nnnn> done` take `--without-queue-tiers`, which lands the story anyway. The acceptance writes, with its commit, a `### Accepted without the queue tiers` section in the story's Notes, naming the tiers that did not run, who accepted, and when. The flag is the CLI's alone: the board and the dashboard have no override. It is refused with `--by orchestrator` and with `--batch`.

This refines ADR-0174: its consequence that `flai accept` by hand runs no tier and lands without the `queue` tiers no longer holds while the queue can batch; a hand acceptance still runs no tier, but waits for a report that ran them. It refines ADR-0173 on the same point: a story accepted by hand still leaves the queue.

## Consequences

- A story accepted by hand while the queue can batch lands only once the `queue` tiers passed at its head, or with the operator's recorded override. A stale generated file or test expectation that only `integration` or `smoke` finds is found before it reaches main.
- Accepting by hand a story whose close-out left the tiers to the queue costs one `flai verify <S-nnnn> --all`, which runs only the queue tiers, about twelve minutes here. Enqueueing the story instead costs the operator nothing.
- The first verify after the upgrade reuses nothing, since the fingerprint changes once.
- The dashboard's acceptance shows the blocker and has no override: an operator who judges the tiers irrelevant lands the story from a shell.
- The check reads the story's last report from the main checkout's `.flai-cache`, as the orchestrator's acceptance and the queue's eligibility do.

## Alternatives considered

- **`flai accept` runs the queue tiers itself before it merges.** Rejected: about twelve minutes here, past the dashboard's ten-minute timeout on `accept.run` and `item.move`, so the host API and flaiover would need changes too.
- **A hand acceptance enqueues the story, and the queue lands it as the operator.** Rejected: done from the board would no longer mean done, and a red batch would undo the operator's acceptance.
- **A batch runs its failing tier on main alone before it blames a story.** A second defence, against main red for any reason. It is a change of scope apart from I-0203, proposed on TH-0701 and not built here.
- **No override.** Rejected: the operator keeps the call when the tiers cannot bear on the change, such as a documentation-only story.
