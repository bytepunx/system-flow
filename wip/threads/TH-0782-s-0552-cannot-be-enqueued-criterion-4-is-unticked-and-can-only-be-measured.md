---
id: TH-0782
title: "S-0552 cannot be enqueued: criterion 4 is unticked and can only be measured after landing"
anchor:
  path: wip/kanban/stories/S-0552-acceptances-land-only-through-the-merge-queue-and-the-queue-waits-a-window-to.md
  item: S-0552
status: resolved
participants: [orchestrator, alex]
created: 2026-10-31T08:39:50Z
updated: 2026-10-31T09:02:20Z
---

# TH-0782 S-0552 cannot be enqueued: criterion 4 is unticked and can only be measured after landing

On wip/kanban/stories/S-0552-acceptances-land-only-through-the-merge-queue-and-the-queue-waits-a-window-to.md.

## Entries

### 2026-10-31T08:39:50Z orchestrator
Recommendation: split criterion 4. Keep its first half, `flai stats` reports batch size and verifications per landed story, as the criterion, and tick it once checked. Move the second half, a mean batch size above one in the fortnight after landing, to a follow-up research story under E-0031. The queue can then take S-0552. You could also accept it yourself as it stands.

What blocks it:
- `criteria_unticked`: criterion 4 is unticked. The merge queue lists S-0552 as ineligible for this reason, and `queue_add` refuses a story whose criteria are not all ticked.
- The story's notes say why: the batch size after landing can only be measured once the story has landed and a fortnight has passed.

What passed:
- `flai verify` passed at the head of `story/S-0552`, 382c790baadfacf59949b1d8166afa34f6224e4b. The integration and smoke tiers are left to the queue.

I have not enqueued it. I will verify it and enqueue it once criterion 4 is ticked or rewritten.

### 2026-10-31T09:02:20Z alex
Resolved: S-0552 was accepted
