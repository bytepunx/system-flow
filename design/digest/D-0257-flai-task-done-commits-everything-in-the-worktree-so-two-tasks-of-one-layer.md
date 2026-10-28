---
id: D-0257
title: flai task done commits everything in the worktree, so two tasks of one layer cannot be closed apart
kind: cause
sources: [I-0104]
topics: [flai, cli, go, template, conventions, code, all]
paths: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, docs/users/flai-reference.md, design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, design/conventions/git.md, template/root/design/conventions/git.md, template/CHANGELOG.md, design/issues/I-0104-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md, design/issues/I-0108-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md, design/issues/summary.md, design/adrs/0128-flai-task-done-commits-the-paths-the-closing-task-covers-and-those-no-other.md, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md, design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md]
date: 2026-10-08
status: absorbed
absorbed_by: "design/adrs/0128-flai-task-done-commits-the-paths-the-closing-task-covers-and-those-no-other.md#decision"
asked_again: []
cited: []
---

# D-0257 flai task done commits everything in the worktree, so two tasks of one layer cannot be closed apart

I-0104 was closed on 2026-10-08.

The cause:

```text
flai task done commits everything in the worktree, so two tasks of one layer cannot be closed apart
```

The fix:

```text
Fixed by S-0322 (ADR-0128): flai task done commits only the changed paths the closing task's touches cover and those no other open task covers, and leaves and lists the paths only another open task covers, so two tasks of one layer close apart, each with its own commit and touches. The commit step is in flai/internal/taskdone/taskdone.go; TestRunClosesTwoTasksOfOneStoryApart reproduces I-0104's two tasks, and -m is needed only when there is something to commit.
```

Source: I-0104.
