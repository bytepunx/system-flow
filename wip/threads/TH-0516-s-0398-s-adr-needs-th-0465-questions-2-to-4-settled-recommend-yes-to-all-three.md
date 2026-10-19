---
id: TH-0516
title: "S-0398's ADR needs TH-0465 questions 2 to 4 settled: recommend yes to all three"
anchor:
  path: wip/kanban/stories/S-0398-flai-queue-lists-the-stories-in-review-that-would-batch-in-the-order-they-entered-review-and-why-each-other-story-would-not-and-the-board-the-mcp-board-and-the-dashboard-show-the-queue.md
  item: S-0398
status: resolved
participants: [agent-S-0398, orchestrator, alex]
created: 2026-10-19T14:35:36Z
updated: 2026-10-19T18:43:51Z
---

# TH-0516 S-0398's ADR needs TH-0465 questions 2 to 4 settled: recommend yes to all three

On wip/kanban/stories/S-0398-flai-queue-lists-the-stories-in-review-that-would-batch-in-the-order-they-entered-review-and-why-each-other-story-would-not-and-the-board-the-mcp-board-and-the-dashboard-show-the-queue.md.

## Entries

### 2026-10-19T14:35:36Z agent-S-0398
S-0398's last criterion asks for an ADR that records the merge queue "as TH-0465 settles it". TH-0465 has no answer yet. I am writing the ADR now as `proposed`, on the planner's recommended answers, and building the read-only queue view, which does not depend on them. I accept the ADR once you answer.

1. **What runs the queue.** Recommended: a model-free loop in `flai serve`, behind a new host action `merge-queue`, off by default, running only while `orchestration.permissions.accept_reviews` is on. Stories enter it only when enqueued, by the orchestrator with ADR-0093's evidence or by the operator.
2. **Batch size.** Recommended: at most 4, set as `merge_queue.batch_max` in `system-flow.yaml`.
3. **A story whose branch undoes main's change (ADR-0155).** Recommended: never batched; shown as not eligible with the paths.

Also recorded as decided, from TH-0466: the queue runs under `acceptance.mode: local` only.

Answer "yes to all three", or name the one to change. A different answer changes only the ADR's text, not this story's code.

### 2026-10-19T18:36:21Z orchestrator (recommendation)
**Recommendation:** yes to all three. ADR-0173 can then be accepted as written, and T-1638 closes.

**Why this needs you.** Only question 3 is already settled by a recorded decision. Questions 1 and 2 are new decisions, and TH-0465, which asked the same questions, has no answer yet. Confirming this answers both threads' questions 2 to 4. Question 1 on TH-0465, the cost of delay input, stays open.

1. **What runs the queue: a model-free loop in `flai serve`, behind a `merge-queue` host action that is off by default, and only while `accept_reviews` is on.**
   - Yes. When you finalized and promoted E-0022's stories, you took this design: S-0399's title says "behind the merge-queue host action", and S-0403 has the orchestrator enqueue each story instead of accepting it.
   - Building, verifying and bisecting a batch is git and arithmetic, so a model adds cost and no judgement.
   - Enqueueing with ADR-0093's evidence keeps the acceptance gate where it is today.
2. **Batch size: at most 4, as `merge_queue.batch_max`.**
   - Yes. A red batch of 4 bisects in at most two more runs.
   - The cap stays under this project's review limit, and it is a manifest setting, so you can change it without a release.
3. **A story whose branch undoes main's change is never batched, and is shown as not eligible with the paths.**
   - Yes. This follows from ADR-0155: "flai and `accept.run` refuse a keep with `--by orchestrator`: an undone path waits for the operator".
   - A queue that lands as the orchestrator could not land such a story, so batching it would only waste the batch's run.

**Also recorded, from TH-0466:** the queue runs under `acceptance.mode: local` only. ADR-0169's `forge_stale` would stop every rebased batch under pull-request acceptance.

Source: design/adrs/0155-a-rebased-story-branch-that-undoes-main-s-change-to-a-path-is-recorded-at-sync.md § Decision

### 2026-10-19T18:43:16Z alex
Confirmed the recommendation of 2026-10-19T18:36:21Z orchestrator.

Source: design/adrs/0155-a-rebased-story-branch-that-undoes-main-s-change-to-a-path-is-recorded-at-sync.md § Decision

### 2026-10-19T18:43:51Z agent-S-0398
Resolved.
