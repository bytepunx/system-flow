---
id: D-0284
title: An outdated MCP flai's task_done commits every changed path under the closing task, so a layer's tasks cannot close apart
kind: cause
sources: [I-0127]
topics: [flai, cli, go, code, all]
paths: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/task_test.go, flai/cmd/task_done.go, docs/users/flai-reference.md, design/adrs/README.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, design/issues/I-0127-an-outdated-mcp-flai-s-task-done-commits-every-changed-path-under-the-closing-task-so-a-layer-s-tasks-cannot-close-apart.md, design/issues/summary.md, design/adrs/0142-flai-task-done-refuses-to-close-a-task-while-the-running-flai-s-major-or-minor.md]
date: 2026-10-12
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0284 An outdated MCP flai's task_done commits every changed path under the closing task, so a layer's tasks cannot close apart

I-0127 was closed on 2026-10-12.

The cause:

```text
An outdated MCP flai's task_done commits every changed path under the closing task, so a layer's tasks cannot close apart
```

The fix:

```text
S-0363: flai task done, the MCP tool task_done, and the host's task.done refuse, changing nothing, while the running flai's major or minor version is behind the newest flai release in the project's history (ADR-0142), so an outdated flai no longer closes a task by older rules and commits another open task's paths; a patch behind still warns through flai_outdated. The specific commit defect was fixed by S-0322 (ADR-0128).
```

Source: I-0127.
