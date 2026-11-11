---
id: MS-0524
title: "S-0689 holds S-0712 on overlap alone: may S-0712 start beside it?"
from: S-0712
to: S-0689
about: [flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go]
status: closed
participants: [flai, agent-S-0689, alex]
created: 2026-11-11T18:17:00Z
updated: 2026-11-11T21:55:29Z
---

# MS-0524 S-0689 holds S-0712 on overlap alone: may S-0712 start beside it?

Between S-0712 and S-0689, about `flai/internal/storygit/sync.go`, `flai/internal/storygit/sync_test.go`.

## Entries

### 2026-11-11T18:17:00Z flai S-0712
S-0689 holds S-0712 on overlap alone: may S-0712 start beside it?

S-0712, "flai task done's sync reports undoes_main for paths the story never changed when main gains a commit during the sync", is ready, and its claim overlaps S-0689's on `flai/internal/storygit/sync.go`, `flai/internal/storygit/sync_test.go`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0338](../../../design/issues/I-0338-flai-task-done-s-sync-reports-undoes-main-for-paths-the-story-never-changed.md), "flai task done's sync reports undoes_main for paths the story never changed when main gains a commit during the sync". The issue recommends no solution yet: propose one from its instances before building it.
>
> It also remediates [I-0335](../../../design/issues/I-0335-flai-stream-sync-reports-a-path-as-undone-by-the-branch-when-main-gains-a.md), "flai stream sync reports a path as undone by the branch when main gains a commit to it while the sync runs", which has the same cause. The operator merged S-0709 into this story on TH-1017.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0689 --remove <path>`, or on the task that names the path, if S-0689 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0524 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0712 can start.
3. Reply with `message_reply` saying why the hold stands: S-0712 waits until S-0689 moves to review, is cancelled, or is sent back.

### 2026-11-11T18:19:18Z agent-S-0689 S-0689 (settled)
The hold stands, briefly. S-0689 changes storygit.Rebase's handling of a stopped rebase, stopped(), ContinueOverResolved, and Sync's checks in sync.go. S-0712 pins main's head inside the same storygit.Rebase, so the two would edit the same function at once. S-0689 was sent back only for a docs conflict, which is now settled. Its close-out is running, and it moves to review as soon as the close-out passes, which releases S-0712.

### 2026-11-11T21:55:29Z alex
Closed: S-0689 was accepted
