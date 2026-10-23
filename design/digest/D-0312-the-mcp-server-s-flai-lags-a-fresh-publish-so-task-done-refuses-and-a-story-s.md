---
id: D-0312
title: The MCP server's flai lags a fresh publish, so task_done refuses and a story's agent falls back to the tree's CLI
kind: cause
sources: [I-0147]
topics: [mcp, release, host, flai, cli, go, code, all]
paths: [design/adrs/README.md, design/system/flai-cli.md, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, flai/internal/mcpserver/installed.go, flai/internal/mcpserver/installed_test.go, flai/internal/mcpserver/server.go, docs/users/flai.md, docs/operators/index.md, design/system/agent-narrative.md, design/issues/I-0147-the-mcp-server-s-flai-lags-a-fresh-publish-so-task-done-refuses-and-a-story-s-agent-falls-back-to-the-tree-s-cli.md, design/issues/summary.md, design/adrs/0160-an-outdated-mcp-server-s-task-done-closes-the-task-with-the-installed-flai-once.md, flai/cmd/mcp.go, flai/cmd/mcp_http.go]
date: 2026-10-15
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0312 The MCP server's flai lags a fresh publish, so task_done refuses and a story's agent falls back to the tree's CLI

I-0147 was closed on 2026-10-15.

The cause:

```text
The MCP server's flai lags a fresh publish, so task_done refuses and a story's agent falls back to the tree's CLI
```

The fix:

```text
S-0390 (ADR-0160): when ADR-0142 refuses an outdated MCP server's task_done, it now closes the task with the flai installed at the server's own path once that flai has caught up, answering closed_by; TestTaskDoneClosesWithTheInstalledFlaiOnceItHasCaughtUp reproduces the instances.
```

Source: I-0147.
