---
id: TH-0428
title: "S-0370 and S-0351 have sat in ready 50 minutes without an agent: restart their starts"
anchor:
  path: wip/kanban/stories/S-0370-flai-task-done-widens-touches-with-a-path-the-commit-reverted-to-main-re-growing-a-claim-the-story-narrowed.md
  item: S-0370
status: open
participants: [orchestrator]
created: 2026-10-12T18:48:21Z
updated: 2026-10-12T18:48:21Z
---

# TH-0428 S-0370 and S-0351 have sat in ready 50 minutes without an agent: restart their starts

On wip/kanban/stories/S-0370-flai-task-done-widens-touches-with-a-path-the-commit-reverted-to-main-re-growing-a-claim-the-story-narrowed.md.

## Entries

### 2026-10-12T18:48:21Z orchestrator
**Recommendation:** start S-0370's and S-0351's agents again from the dashboard's board, or check `flai serve status` for a failed start on each. I cannot run `flai serve status` (flai guard).

## What I see

- S-0370 (rank 1, 21m) and S-0351 (rank 2, 45m) entered ready at 17:57, while flai serve was the stalled 1.40.1 (TH-0411).
- Neither is held, and only S-0236 is in progress, against a limit of 5.
- When flai serve came back at 18:13, it started S-0363, S-0373, S-0236, and S-0369, which were ranked below them, but not these two. So their starts probably failed during the stall and were not tried again. ADR-0043 leaves a restart of a failed agent to you.
- S-0351 needs `OPENROUTER_API_KEY` in its agent's environment, the planner's assumption 1 on TH-0386. Without it the agent asks on a thread, and the key-gated tests skip.

## Also waiting on you

- TH-0427: when to publish. Five stories are now unreleased behind S-0237's compatibility break: S-0237, S-0363, S-0368, S-0369, and S-0373.
- TH-0420, TH-0422, and TH-0424: plan decisions for S-0371, S-0372, and S-0374.
- TH-0421: ADR-0141.
