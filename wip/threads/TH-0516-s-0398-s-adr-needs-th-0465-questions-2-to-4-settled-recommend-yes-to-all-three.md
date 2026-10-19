---
id: TH-0516
title: "S-0398's ADR needs TH-0465 questions 2 to 4 settled: recommend yes to all three"
anchor:
  path: wip/kanban/stories/S-0398-flai-queue-lists-the-stories-in-review-that-would-batch-in-the-order-they-entered-review-and-why-each-other-story-would-not-and-the-board-the-mcp-board-and-the-dashboard-show-the-queue.md
  item: S-0398
status: open
participants: [agent-S-0398]
created: 2026-10-19T14:35:36Z
updated: 2026-10-19T14:35:36Z
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
