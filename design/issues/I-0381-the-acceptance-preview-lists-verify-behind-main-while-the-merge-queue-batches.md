---
id: I-0381
title: The acceptance preview lists verify_behind_main while the merge queue batches, so the review page disables Accept and offers Rebase and resolve for a story the acceptance would enqueue
class: defect
status: closed
count: 1
cost: 1h
first_reported: 2026-11-16T00:41:30Z
last_reported: 2026-11-16T00:41:30Z
updated: 2026-11-16T01:31:05Z
---

# I-0381 The acceptance preview lists verify_behind_main while the merge queue batches, so the review page disables Accept and offers Rebase and resolve for a story the acceptance would enqueue

## Description
The acceptance preview lists verify_behind_main while the merge queue batches, so the review page disables Accept and offers Rebase and resolve for a story the acceptance would enqueue

## Instances

### 2026-11-16T00:41:30Z
With the merge-queue host action on for this project and the queue batching, flai accept S-0742 --dry-run still listed verify_behind_main, and the review page (flaiover/src/lib/components/Review.svelte) disabled Accept and offered Rebase and resolve for every story in review but S-0789. The acceptance itself enqueues the story before the gate is read (queuesInstead in flai/cmd/accept.go, which a dry run is left out of), and the batch runs every tier over the main it lands on, so ADR-0291 §1 says the gate meets only what lands by hand. verifyBehindMainGate in flai/internal/preview/accept.go skips only for the queue's landing (opts.Batch), not while the queue batches (opts.QueueBatches).

## Impact
On 2026-11-15, with the merge-queue host action on, S-0742, S-0747, and S-0739 each showed the blocker and Rebase and resolve, though flai accept would have enqueued each; a send-back costs an agent session, a sync, and a close-out of about twelve minutes

- time_lost_per_cycle: 2h

## Remediation

While the merge queue batches and is not held, the operator's acceptance preview does not list `verify_behind_main` as a blocker: `verifyBehindMainGate` in `flai/internal/preview/accept.go` skips when `opts.QueueBatches` is set and the acceptance would enqueue the story, as `queuesInstead` decides, and still reports the stale tiers so the preview can say the batch will run them. The review page then offers Accept, which enqueues, in place of Rebase and resolve. An acceptance that lands by hand (`--without-queue-tiers`, `--without-kluster`, `--keep-removal`, or a held queue) keeps the gate. Refine ADR-0291 with an ADR that says so.

Story S-0796 remediates this issue, created from it at 2026-11-16T00:41:47Z.
Closed 2026-11-16T01:31:05Z: S-0796 (ADR-0308): preview.AcceptOptions carries the merge queue's hold and one predicate, Enqueues, which flai accept's queuesInstead shares; while it holds, verifyBehindMainGate names the stale tiers in stale_tiers_queued in place of the verify_behind_main blocker, and flai accept --dry-run and the host's accept.preview read the hold, so the review page offers Accept, which enqueues, for such a story. Tests in flai/internal/preview, flai/cmd/accept_queue_test.go, and flai/internal/hostapi/reads_test.go.
