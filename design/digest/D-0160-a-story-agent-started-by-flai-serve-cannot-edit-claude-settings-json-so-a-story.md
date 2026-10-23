---
id: D-0160
title: A story agent started by flai serve cannot edit .claude/settings.json, so a story that changes the guard hook's matcher cannot ship it
kind: cause
sources: [I-0069]
topics: [flai, cli, go, code, all]
paths: [flai/internal/mcpserver, flai/internal/harness, flai/internal/hostapi, flai/cmd/mcp.go, flai/cmd/mcp_test.go, design/adrs, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, design/conventions/delegation.md, design/issues, flai/cmd/serve_actions_test.go]
date: 2026-10-05
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0160 A story agent started by flai serve cannot edit .claude/settings.json, so a story that changes the guard hook's matcher cannot ship it

I-0069 was closed on 2026-10-05.

The cause:

```text
A story agent started by flai serve cannot edit .claude/settings.json, so a story that changes the guard hook's matcher cannot ship it
```

The fix:

```text
S-0257: flai serve starts claude-code with --permission-prompt-tool mcp__flai__permission_prompt (ADR-0086). flai's permission_prompt asks the story's owner on a thread before an Edit or Write under a .claude/ folder in the story's worktree, and lets it through on allow, so the agent writes the file itself and no paste is needed. The shell-only auto-approve host action, off by default, allows such writes at once.
```

Source: I-0069.
