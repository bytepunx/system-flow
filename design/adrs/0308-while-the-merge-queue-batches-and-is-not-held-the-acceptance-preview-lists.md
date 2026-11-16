---
id: ADR-0308
title: "While the merge queue batches and is not held, the acceptance preview lists verify_behind_main only for an acceptance that lands by hand, and for one that would enqueue the story it names the stale tiers the batch runs"
status: accepted
date: 2026-11-16
supersedes: []
superseded_by: []
refines: [ADR-0291]
---

# ADR-0308 While the merge queue batches and is not held, the acceptance preview lists verify_behind_main only for an acceptance that lands by hand, and for one that would enqueue the story it names the stale tiers the batch runs

## Context

[ADR-0291](0291-every-acceptance-but-the-merge-queue-s-landing-waits-for-a-flai-verify-over.md) §1 says that while the merge queue batches, the operator's acceptance enqueues the story ([ADR-0225](0225-while-the-merge-queue-batches-and-is-not-held-the-operator-s-acceptance-of-a.md)). So its gate, `verify_behind_main`, meets only what still lands by hand. The batch runs every tier over the main it lands on.

The preview did not follow it. `flai accept --dry-run` and the dashboard's `accept.preview` read `verify_behind_main` for every acceptance but the queue's landing (`--batch`). A dry run is left out of `queuesInstead` in `flai/cmd/accept.go`, so the preview always previewed a landing by hand. The review page then disabled Accept and offered Rebase and resolve for stories the acceptance would have enqueued.

[I-0381](../issues/I-0381-the-acceptance-preview-lists-verify-behind-main-while-the-merge-queue-batches.md) records it. On 2026-11-15, with the `merge-queue` host action on, S-0742, S-0747, and S-0739 each showed the blocker, though `flai accept` would have enqueued each. A send-back costs an agent session, a sync, and a close-out of about twelve minutes.

## Decision

While the merge queue batches and is not held, the acceptance preview lists `verify_behind_main` only for an acceptance that lands by hand. For one that would enqueue the story, it names the stale tiers the batch runs.

1. **One predicate.** `preview.AcceptOptions` carries the queue's hold, `QueueHeld`: the merge queue state's `Held`, read by the caller. A state that cannot be read counts as held, so the gate never opens unread. One predicate, `AcceptOptions.Enqueues`, decides whether the operator's acceptance enqueues the story. It holds when the queue batches, it is not held, the acceptance is not the orchestrator's, it is not the queue's landing (`--batch`), and none of `--without-queue-tiers`, `--without-kluster`, and `--keep-removal` is given. `flai accept`'s `queuesInstead` uses the same predicate, so the preview and the acceptance agree.
2. **An acceptance that enqueues.** While `Enqueues` holds, the preview lists no `verify_behind_main` blocker. It still names the stale tiers (`stale_tiers`) and says, in `stale_tiers_queued`, that the merge queue's batch runs them again over the main it lands on. A report or a main that cannot be read is named there the same way, not as a blocker: the enqueue does not read it, and the batch runs every tier.
3. **An acceptance that lands by hand.** Every acceptance that lands by hand keeps the gate: one while the queue is held, one with any of the three flags, the orchestrator's acceptance, and any acceptance while the queue does not batch. The orchestrator's acceptance keeps it as ADR-0291 §5 says: the blocker is among the acceptance's blockers. The queue's landing stays ungated.
4. **No override.** No flag is added. ADR-0291 §6 stands.

This refines ADR-0291: §1's rule that the gate meets only what lands by hand now holds in the preview too.

## Consequences

- For a story the acceptance would enqueue, the review page offers Accept, which enqueues it, in place of Rebase and resolve.
- No agent session, sync, or twelve-minute close-out is spent on a send-back the batch makes needless.
- A story landed by hand still waits for a `flai verify` over main's new paths.
- `queue_tiers` ([ADR-0221](0221-while-the-merge-queue-can-batch-a-hand-acceptance-is-refused-until-the-queue.md)) is unchanged.

## Alternatives considered

- **The dashboard ignores the blocker while the queue batches.** Rejected: the dashboard would decide what flai decides, and it would miss the hold and the flags.
- **Drop the gate while the queue batches, whatever the flags.** Rejected: the flags and a held queue land by hand, which is what left main red in [I-0229](../issues/I-0229-two-stories-accepted-apart-leave-main-red-s-0535-reworded-to-backlog-s-condition-and-s-0416-s-testadefinedroleisheldtoitsfile-still-expects-the-old-words.md).
