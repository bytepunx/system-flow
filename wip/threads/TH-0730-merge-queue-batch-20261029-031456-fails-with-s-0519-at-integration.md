---
id: TH-0730
title: Merge queue batch 20261029-031456 fails with S-0519 at integration
anchor:
  path: wip/kanban/stories/S-0519-permission-prompt-over-the-http-mcp-server-outlives-claude-code-s-tool-call-timeout-so-an-agent-in-a-container-sees-its-protected-write-time-out-while-the-thread-to-the-operator-stays-open.md
  item: S-0519
status: resolved
participants: [flai serve, agent-S-0519]
created: 2026-10-29T03:24:27Z
updated: 2026-10-29T03:37:45Z
---

# TH-0730 Merge queue batch 20261029-031456 fails with S-0519 at integration

On wip/kanban/stories/S-0519-permission-prompt-over-the-http-mcp-server-outlives-claude-code-s-tool-call-timeout-so-an-agent-in-a-container-sees-its-protected-write-time-out-while-the-thread-to-the-operator-stays-open.md.

## Entries

### 2026-10-29T03:24:27Z flai serve
Merge queue batch `20261029-031456` fails with S-0519 at `integration`. The queue verified S-0519 alone on main at `c24575c2eecb`, at `2ec7500ae00e`, the copy of its verified head `fdac90c3eb5b` rebased in the batch, and `integration` failed:

```text
flai/internal/guard/roles_test.go:139 TestADefinedRoleIsHeldToItsFile: {Session: Subagent:false Agent: Role:review Story: Kind:3 Name: Command: Path: Tool:item_move Args:{ID:S-0001 Story: To:review Type: Draft:false Commit:false Recommendation:false Source: Fields:[] CostOfDelay:map[]}}: {Outcome:1 Reason:the role reviewer cannot move S-0001 to review: it may make it only under to-backlog (on a call, a move to backlog and nowhere else; on a trigger, only a move to backlog), or to-ready (a move to ready and nowhere else). Ask the operator with thread_open on the item if it needs doing (design/workflows/roles/reviewer.md, ADR-0207). Call:item_move S-0001 review Needs:}
```

The whole output of each tier it ran is in `.flai-cache/merge-queue/batches/20261029-031456`.

flai moved S-0519 back to in-progress and took it out of the merge queue. Fix it on `story/S-0519`, run flai verify S-0519, move it to review, and enqueue it again with flai queue add S-0519.

### 2026-10-29T03:27:36Z agent-S-0519
The failure is on main, not in S-0519. S-0535 reworded the `to-backlog` condition, and S-0416's `TestADefinedRoleIsHeldToItsFile`, accepted after it, still expected the old words. `flai test flai/internal/guard` fails on main at c24575c2e without S-0519. S-0519's T-2323 (fa8edce) updates the test's expected text, so every batch can pass integration again. I recorded this as I-0229. S-0519 goes back to review after its close-out.

### 2026-10-29T03:37:45Z flai serve
Resolved: S-0519 is back in review at fa8edce27028, a new commit; it is batched again once it is enqueued
