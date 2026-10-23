---
id: D-0173
title: flai's permission_prompt returns a result Claude Code calls invalid, so a .claude/ write and a /tmp write are refused at once with no thread
kind: cause
sources: [I-0082]
topics: [cli, flai, go, code, all]
paths: [flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, flai/internal/mcpserver/folder.go, docs/users/flai.md, design/issues/I-0082-flai-s-permission-prompt-returns-a-result-claude-code-calls-invalid-so-a-claude-write-and-a-tmp-write-are-refused-at-once-with-no-thread.md, design/issues/summary.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md, design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md]
date: 2026-10-06
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0173 flai's permission_prompt returns a result Claude Code calls invalid, so a .claude/ write and a /tmp write are refused at once with no thread

I-0082 was closed on 2026-10-06.

The cause:

```text
flai's permission_prompt returns a result Claude Code calls invalid, so a .claude/ write and a /tmp write are refused at once with no thread
```

The fix:

```text
Fixed by S-0283 (T-1001). The cause: flai registered permission_prompt with the go-sdk's typed mcp.AddTool, which declared an output schema and answered structuredContent beside the text block, and Claude Code accepts only a single text block, so it read every allow and every refusal as an invalid result. permission_prompt is now a raw handler that answers one text block holding the JSON decision, with no output schema and no structured content; arguments it cannot read are a deny in the same shape. TestPermissionPromptAnswersOneTextBlockAsClaudeCodeRequires in flai/internal/mcpserver/permission_test.go reproduces I-0082, and every permission test now checks the shape. The S-0219 .claude/ write was refused at once because S-0219 was still ready, a refusal ADR-0086 intends; the /tmp writes were refused as outside any worktree, also intended; with the fix the agent reads those reasons. Running agents get the fix once a flai release with S-0283 is installed, since their MCP server is the installed flai.
```

Source: I-0082.
