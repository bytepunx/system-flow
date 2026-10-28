---
id: D-0245
title: A story agent's .claude/ write waits thirty minutes on an unanswered permission thread, then fails on Claude Code's MCP idle timeout
kind: cause
sources: [I-0103]
topics: [cli, agents, flai, go, code, all]
paths: [flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/delegation.md, docs/users/flai.md, design/system/flai-cli.md, design/adrs, design/adrs/README.md, design/issues/I-0103-a-story-agent-s-claude-write-waits-thirty-minutes-on-an-unanswered-permission-thread-then-fails-on-claude-code-s-mcp-idle-timeout.md, design/issues/summary.md, design/issues/I-0108-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md]
date: 2026-10-08
status: absorbed
absorbed_by: "design/adrs/0124-permission-prompt-holds-a-write-at-most-four-minutes-then-refuses-it-and-leaves.md#decision"
asked_again: []
cited: []
---

# D-0245 A story agent's .claude/ write waits thirty minutes on an unanswered permission thread, then fails on Claude Code's MCP idle timeout

I-0103 was closed on 2026-10-08.

The cause:

```text
A story agent's .claude/ write waits thirty minutes on an unanswered permission thread, then fails on Claude Code's MCP idle timeout
```

The fix:

```text
Fixed by S-0309 (ADR-0124). permission_prompt holds a protected write at most four minutes, below Claude Code's MCP idle timeout. When no answer comes in that time, or the session ends, it refuses the write naming the thread and leaves the thread open. When the agent makes the same write again, the answer given since is taken from that thread, or it waits again. The start prompt and delegation.md now tell the agent to retry the write once the thread is answered, instead of staging files in .flai-cache/ with cp commands. TestPermissionPromptRefusesUnansweredWithinTheBoundAndLeavesTheThreadOpen reproduces the cause.
```

Source: I-0103.
