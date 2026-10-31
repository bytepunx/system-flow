---
id: ADR-0225
title: "While the merge queue batches and is not held, the operator's acceptance of a story in review enqueues it in place of landing it, so stories land only through the queue and no hand acceptance drops a batch, and a held queue, --without-queue-tiers, and --keep-removal land by hand"
status: accepted
date: 2026-10-31
supersedes: []
superseded_by: []
refines: [ADR-0173, ADR-0221]
topics: [orchestration, flai, cli, dashboard]
---

# ADR-0225 While the merge queue batches and is not held, the operator's acceptance of a story in review enqueues it in place of landing it, so stories land only through the queue and no hand acceptance drops a batch, and a held queue, --without-queue-tiers, and --keep-removal land by hand

## Context

The merge queue verifies the stories enqueued as one batch, bisects a red batch, and lands the green stories ([ADR-0173](0173-a-merge-queue-in-flai-serve-with-no-model-runs-behind-the-merge-queue-host.md), [ADR-0176](0176-the-merge-queue-lands-each-story-of-a-passing-batch-through-flai-accept-batch.md)). The orchestrator enqueues in place of accepting while the queue batches (S-0403). The operator's acceptance did not: `flai accept`, `flai move <S-nnnn> done`, and the board's accept landed the story by hand. That acceptance took the story out of the queue and dropped the batch that held it, so the stories in that batch were built and verified again.

From 2026-10-15 to 2026-10-29 the queue built 111 batches of mean size 1.0 to land 49 stories (S-0552's evidence, from `flai stats`). Each verification takes about ten minutes here.

[ADR-0221](0221-while-the-merge-queue-can-batch-a-hand-acceptance-is-refused-until-the-queue.md) gated a hand acceptance on the tiers marked `queue`, after a hand acceptance left main red ([I-0203](../issues/I-0203-a-story-the-operator-accepts-outside-the-merge-queue-lands-without-the-queue-tiers-so-a-stale-generated-file-on-main-fails-every-later-batch-at-integration.md)). That gate costs a `flai verify <S-nnnn> --all` of about twelve minutes per story, outside any batch. ADR-0221 rejected having the hand acceptance enqueue the story, on two grounds: done on the board would no longer mean done, and a red batch would undo the operator's acceptance. Both hold only if the acceptance moves the story to done before the queue lands it.

## Decision

While the merge queue batches and is not held, the operator's acceptance of a story in review enqueues it in place of landing it. So stories land only through the queue, and no hand acceptance drops a batch. A held queue, `--without-queue-tiers`, and `--keep-removal` land by hand.

- **When.** The queue batches by ADR-0174's three conditions: the `merge-queue` host action is on for the project, `orchestration.permissions.accept_reviews` is on, and acceptance is local. The queue is not held (`flai queue hold`). With any of these not so, the acceptance lands the story as before.
- **Which acceptances.** The operator's acceptance of a story in review, by `flai accept <S-nnnn>` or `flai move <S-nnnn> done`. The board's accept and move to done run them, through the host writes `accept.run` and `item.move`. The orchestrator's acceptance (`--by orchestrator`) and the queue's landing (`--batch`) are unchanged.
- **What it does.** It enqueues the story as `flai queue add <S-nnnn> --by <by>` does, `<by>` being `--by` or the configured author. Nothing is merged, moved, archived, or dropped, and the story stays in review. The queue lands it as whoever enqueued it once its batch passes, through `flai accept --batch` ([ADR-0176](0176-the-merge-queue-lands-each-story-of-a-passing-batch-through-flai-accept-batch.md)). The answer says the story was enqueued, and at which position. `--json` adds `queue`: `enqueued`, `batch`, `position`, `by`, and `at`.
- **Enqueued already.** A story enqueued already, or in the batch under way, is answered so, and nothing changes. The batch under way is not dropped.
- **Refused.** The enqueue's refusal is the acceptance's: a story the queue finds not eligible is refused with each reason the queue gives.
- **Landing by hand.** `--without-queue-tiers` lands the story by hand without the queue tiers and records that in its Notes, as ADR-0221 decided. `--keep-removal` lands the story by hand under ADR-0221's gate: the queue never batches a story whose branch undoes main's change ([ADR-0173](0173-a-merge-queue-in-flai-serve-with-no-model-runs-behind-the-merge-queue-host.md)), so a kept removal lands only by hand. While the queue is held, an acceptance lands by hand under ADR-0221's gate. A landing by hand takes the story out of the queue and drops the batch that holds it, as before. `--dry-run` previews the landing by hand.

This refines ADR-0221: the alternative it rejected now holds, with the story kept in review until the queue lands it. Its gate now applies to an acceptance that lands by hand. It refines ADR-0173: an acceptance takes a story out of the queue and drops its batch only when it lands the story by hand.

## Consequences

- While the queue batches, each story lands verified with its batch, its `queue` tiers among the batch's tiers. An acceptance on the board can no longer land a story that fails `integration` or `smoke`.
- Accepting a story while the queue batches costs no `flai verify --all`. The story lands once its batch's window ends and its verification passes.
- No acceptance drops a batch under way. The stories in it are not verified again on its account.
- Done on the board still means landed. An accepted story stays in review until the queue lands it, so the operator sees it there for some minutes.
- A red batch sends the story it is bisected to back to in-progress, with a thread, as for any story enqueued. The operator's acceptance was never recorded as done, so nothing is undone.
- An operator who needs a story landed now holds the queue first, or uses `--without-queue-tiers`. A batch already running when the queue is held finishes.
- The preview, `--dry-run` and the dashboard's `accept.preview`, still previews a landing by hand. So it can list ADR-0221's `queue_tiers` blocker for a story the acceptance would enqueue.

## Alternatives considered

- **Keep ADR-0221's gate alone.** Each hand acceptance costs a `flai verify --all` and drops the batch under way.
- **Enqueue while the queue is held too.** A held queue lands nothing, so the story would wait until someone resumes it. Holding the queue is how the operator takes landing back.
- **Move the story to done and let the queue land it after.** Done would not mean merged, the objection ADR-0221 raised.
- **Have `--without-queue-tiers` and `--keep-removal` enqueue too.** The first exists to land without the queue tiers, and the queue never batches the second's story.
