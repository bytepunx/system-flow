---
id: D-0271
title: flai task done commits the changed paths of another open task of the same layer
kind: cause
sources: [I-0126]
topics: [flai, cli, go, template, conventions, code, all]
paths: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, docs/users/flai-reference.md, design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, design/conventions/git.md, template/root/design/conventions/git.md, template/CHANGELOG.md, design/issues/I-0104-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md, design/issues/I-0108-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md, design/issues/summary.md, design/adrs/0128-flai-task-done-commits-the-paths-the-closing-task-covers-and-those-no-other.md, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md, design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md]
date: 2026-10-12
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0271 flai task done commits the changed paths of another open task of the same layer

I-0126 was closed on 2026-10-12.

The cause:

```text
flai task done commits the changed paths of another open task of the same layer
```

The fix:

```text
Fixed by S-0322, which implemented ADR-0128 and was published in flai 1.39.14 at 2026-10-08T09:30Z. The instance closed T-1377 at 10:35:55Z through an MCP server still running the installed flai 1.39.13, which committed every changed path. S-0362's tests show the current flai task done and MCP task_done leave an open sibling's paths uncommitted, list them in left, and keep them out of the closing task's and the story's touches: TestTaskDoneLeavesASiblingsConventionTemplateAndChangelog (flai/internal/mcpserver/task_test.go) reproduces the instance's layer, and TestRunLeavesOnlyAnOpenSiblingsPaths (flai/internal/taskdone/taskdone_test.go) covers backlog, ready, done, and cancelled siblings. Both fail when commit stages every changed path and pass on the current code, so taskdone.go needed no fix. The guard against an MCP flai older than the project's history is S-0363's, from I-0127.
```

Source: I-0126.
