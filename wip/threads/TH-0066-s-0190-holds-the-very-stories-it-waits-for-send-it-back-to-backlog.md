---
id: TH-0066
title: "S-0190 holds the very stories it waits for: send it back to backlog?"
anchor:
  path: wip/kanban/stories/S-0190-measure-what-delegating-costs-with-cheaper-sub-agents-and-one-verifier-run-before-review.md
  item: S-0190
status: resolved
participants: [agent-S-0190, alex]
created: 2026-10-02T12:40:01Z
updated: 2026-10-02T12:41:38Z
---

# TH-0066 S-0190 holds the very stories it waits for: send it back to backlog?

On wip/kanban/stories/S-0190-measure-what-delegating-costs-with-cheaper-sub-agents-and-one-verifier-run-before-review.md.

## Entries

### 2026-10-02T12:40:01Z agent-S-0190
Recommendation: I move S-0190 back to ready and then to backlog. You move it to ready again once an Opus story that changes code has ended.

The deadlock: S-0231, S-0191, S-0192, and S-0178 declare no touches, so flai serve holds them while S-0190 (and S-0193) is in progress. The board says each "starts when ... S-0190 and S-0193 are accepted, cancelled, or sent back". S-0190 is blocked waiting for one of them to run (TH-0064, option a). Neither side can move.

Nothing is lost by sending it back. The worktree has no commits and S-0194's numbers are in the narrative's Decisions. The measuring scripts are in .flai-cache/s0190/. The story's Notes already say to pull it once the runs exist.

Alternatives:
(b) You declare touches on one of the four stories, say S-0231 (LICENSE.md). Then it is no longer held by S-0190 and runs once S-0193 frees room. S-0190 stays blocked in progress until then.
(c) I leave S-0190 in progress and blocked. Nothing moves until S-0193 ends and something else changes.

### 2026-10-02T12:41:13Z alex
a

### 2026-10-02T12:41:38Z agent-S-0190
Resolved: S-0190 is back in backlog and no longer holds S-0231, S-0191, S-0192, or S-0178. They now wait only for S-0193.
