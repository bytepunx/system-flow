---
id: D-0284
title: An outdated MCP flai's task_done commits every changed path under the closing task, so a layer's tasks cannot close apart
kind: cause
sources: [I-0127]
topics: [flai, cli, go, code, all]
paths: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/task_test.go, flai/cmd/task_done.go, docs/users/flai-reference.md, design/adrs/README.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, design/issues/I-0127-an-outdated-mcp-flai-s-task-done-commits-every-changed-path-under-the-closing-task-so-a-layer-s-tasks-cannot-close-apart.md, design/issues/summary.md, design/adrs/0142-flai-task-done-refuses-to-close-a-task-while-the-running-flai-s-major-or-minor.md]
date: 2026-10-12
status: compacted
absorbed_by: ""
asked_again: []
cited: []
---

# D-0284 An outdated MCP flai's task_done commits every changed path under the closing task, so a layer's tasks cannot close apart
