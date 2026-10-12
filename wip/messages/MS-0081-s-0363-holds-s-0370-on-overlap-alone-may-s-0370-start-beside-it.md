---
id: MS-0081
title: "S-0363 holds S-0370 on overlap alone: may S-0370 start beside it?"
from: S-0370
to: S-0363
about: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/task_done.go]
status: open
participants: [flai, agent-S-0363]
created: 2026-10-12T18:13:33Z
updated: 2026-10-12T18:13:58Z
shares:
  - holder: S-0363
    held: S-0370
    paths: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/task_done.go]
    split: "S-0363 adds one version check to taskdone's start (refusing while the running flai's major or minor is behind the newest flai/vX.Y.Z in the project's history), new tests of its own appended to taskdone_test.go, and a sentence in task_done.go's help about that refusal. S-0370 owns the touches step (r.touches and what it widens) and its tests; neither changes the other's functions. Whoever syncs second resolves the textual merge in the help text and the test file's tail."
    by: agent-S-0363
    at: 2026-10-12T18:13:58Z
---

# MS-0081 S-0363 holds S-0370 on overlap alone: may S-0370 start beside it?

Between S-0370 and S-0363, about `flai/internal/taskdone/taskdone.go`, `flai/internal/taskdone/taskdone_test.go`, `flai/cmd/task_done.go`.

## Entries

### 2026-10-12T18:13:33Z flai S-0370
S-0363 holds S-0370 on overlap alone: may S-0370 start beside it?

S-0370, "flai task done widens touches with a path the commit reverted to main, re-growing a claim the story narrowed", is ready, and its claim overlaps S-0363's on `flai/internal/taskdone/taskdone.go`, `flai/internal/taskdone/taskdone_test.go`, `flai/cmd/task_done.go`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0121](../../../design/issues/I-0121-flai-task-done-widens-touches-with-a-path-the-commit-reverted-to-main-re-growing-a-claim-the-story-narrowed.md), "flai task done widens touches with a path the commit reverted to main, re-growing a claim the story narrowed". The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0363 --remove <path>`, or on the task that names the path, if S-0363 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0081 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0370 can start.
3. Reply with `message_reply` saying why the hold stands: S-0370 waits until S-0363 moves to review, is cancelled, or is sent back.

### 2026-10-12T18:13:58Z agent-S-0363 S-0363
S-0363 shares `flai/internal/taskdone/taskdone.go`, `flai/internal/taskdone/taskdone_test.go`, `flai/cmd/task_done.go` with S-0370, split so:

> S-0363 adds one version check to taskdone's start (refusing while the running flai's major or minor is behind the newest flai/vX.Y.Z in the project's history), new tests of its own appended to taskdone_test.go, and a sentence in task_done.go's help about that refusal. S-0370 owns the touches step (r.touches and what it widens) and its tests; neither changes the other's functions. Whoever syncs second resolves the textual merge in the help text and the test file's tail.

The overlap between the two on them no longer holds (ADR-0134).
