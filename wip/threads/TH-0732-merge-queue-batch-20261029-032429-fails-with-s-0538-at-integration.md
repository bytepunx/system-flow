---
id: TH-0732
title: Merge queue batch 20261029-032429 fails with S-0538 at integration
anchor:
  path: wip/kanban/stories/S-0538-flai-adr-topics-rewrites-an-adr-s-front-matter-with-and-in-its-title-escaped-as-and.md
  item: S-0538
status: resolved
participants: [flai serve, agent-S-0538]
created: 2026-10-29T03:35:46Z
updated: 2026-10-29T03:54:29Z
---

# TH-0732 Merge queue batch 20261029-032429 fails with S-0538 at integration

On wip/kanban/stories/S-0538-flai-adr-topics-rewrites-an-adr-s-front-matter-with-and-in-its-title-escaped-as-and.md.

## Entries

### 2026-10-29T03:35:46Z flai serve
Merge queue batch `20261029-032429` fails with S-0538 at `integration`. The queue verified S-0538 alone on main at `c24575c2eecb`, at `2a955340bc1e`, the copy of its verified head `2a955340bc1e` rebased in the batch, and `integration` failed:

```text
flai/internal/guard/roles_test.go:139 TestADefinedRoleIsHeldToItsFile: {Session: Subagent:false Agent: Role:review Story: Kind:3 Name: Command: Path: Tool:item_move Args:{ID:S-0001 Story: To:review Type: Draft:false Commit:false Recommendation:false Source: Fields:[] CostOfDelay:map[]}}: {Outcome:1 Reason:the role reviewer cannot move S-0001 to review: it may make it only under to-backlog (on a call, a move to backlog and nowhere else; on a trigger, only a move to backlog), or to-ready (a move to ready and nowhere else). Ask the operator with thread_open on the item if it needs doing (design/workflows/roles/reviewer.md, ADR-0207). Call:item_move S-0001 review Needs:}
```

The whole output of each tier it ran is in `.flai-cache/merge-queue/batches/20261029-032429`.

flai moved S-0538 back to in-progress and took it out of the merge queue. Fix it on `story/S-0538`, run flai verify S-0538, move it to review, and enqueue it again with flai queue add S-0538.

### 2026-10-29T03:47:16Z agent-S-0538
Fixed at e6c38cd. S-0538's own change did not cause this: main fails `TestADefinedRoleIsHeldToItsFile` in `flai/internal/guard/roles_test.go`. S-0416 was accepted by hand (afbc75b) after S-0535 had reworded the to-backlog condition, and its test still expected the old reason (I-0203, bumped). T-2324 updates the expectation. It is the same one-line edit as S-0519's T-2323, so the two rebase cleanly over each other. The close-out passed; S-0538 is back in review.

### 2026-10-29T03:54:29Z flai serve
Resolved: S-0538 is back in review at e6c38cda9850, a new commit; it is batched again once it is enqueued
