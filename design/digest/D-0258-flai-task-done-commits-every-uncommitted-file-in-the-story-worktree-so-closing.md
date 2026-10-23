---
id: D-0258
title: flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit
kind: cause
sources: [I-0108]
topics: [flai, cli, go, template, conventions, code, all]
paths: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, docs/users/flai-reference.md, design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, design/conventions/git.md, template/root/design/conventions/git.md, template/CHANGELOG.md, design/issues/I-0104-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md, design/issues/I-0108-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md, design/issues/summary.md, design/adrs/0128-flai-task-done-commits-the-paths-the-closing-task-covers-and-those-no-other.md, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md, design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md]
date: 2026-10-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0258 flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit

I-0108 was closed on 2026-10-08.

The cause:

```text
flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit
```

The fix:

```text
Fixed by S-0322 (ADR-0128): flai task done no longer runs git add -A; it commits the changed paths the closing task's touches cover and those no other open task covers, and leaves the paths only another open task covers for that task's close, so closing one task of a layer no longer sweeps its siblings' work into its commit. The commit step is in flai/internal/taskdone/taskdone.go; TestRunClosesALayerOfThreeTasksInACommitEach reproduces S-0274's layer of three tasks.
```

Source: I-0108.
