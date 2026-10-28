---
id: D-0286
title: flai task done widens touches with a path the commit reverted to main, re-growing a claim the story narrowed
kind: cause
sources: [I-0121]
topics: [touches, flai, cli, go, code, all]
paths: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/task_done.go, flai/internal/mcpserver/task.go, docs/users/flai-reference.md, design/adrs/README.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0121-flai-task-done-widens-touches-with-a-path-the-commit-reverted-to-main-re-growing-a-claim-the-story-narrowed.md, design/issues/summary.md, design/adrs/0144-flai-task-done-widens-touches-only-with-the-committed-paths-the-story-s-branch.md]
date: 2026-10-12
status: absorbed
absorbed_by: "design/adrs/0144-flai-task-done-widens-touches-only-with-the-committed-paths-the-story-s-branch.md#decision"
asked_again: []
cited: []
---

# D-0286 flai task done widens touches with a path the commit reverted to main, re-growing a claim the story narrowed

I-0121 was closed on 2026-10-12.

The cause:

```text
flai task done widens touches with a path the commit reverted to main, re-growing a claim the story narrowed
```

The fix:

```text
S-0370: flai task done now widens touches only with the committed paths the story's branch still changes against its merge base with main, read after the sync, and lists the others as unchanged (ADR-0144); TestRunWidensTouchesOnlyWithPathsTheBranchStillChanges reproduces the revert.
```

Source: I-0121.
