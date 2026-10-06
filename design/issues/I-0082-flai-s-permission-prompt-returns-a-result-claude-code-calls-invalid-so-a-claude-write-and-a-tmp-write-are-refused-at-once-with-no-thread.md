---
id: I-0082
title: flai's permission_prompt returns a result Claude Code calls invalid, so a .claude/ write and a /tmp write are refused at once with no thread
class: defect
status: closed
count: 3
cost: 6m
first_reported: 2026-10-06T03:19:21Z
last_reported: 2026-10-06T06:06:28Z
updated: 2026-10-06T10:10:31Z
---

# I-0082 flai's permission_prompt returns a result Claude Code calls invalid, so a .claude/ write and a /tmp write are refused at once with no thread

## Description
flai's permission_prompt returns a result Claude Code calls invalid, so a .claude/ write and a /tmp write are refused at once with no thread

## Instances

### 2026-10-06T03:19:21Z
Story: S-0219.
S-0219: the story agent's Write of .claude/agents/orchestrator.md and a task sub-agent's Write to /tmp were refused with 'Permission prompt tool returned an invalid result. Expected a single text block param with type="text" and a string text value.'; no Allow thread was opened, so the operator had to paste the file

### 2026-10-06T04:56:08Z
Story: S-0282.
T-0996's task sub-agent had Write to /tmp refused with 'Permission prompt tool returned an invalid result' and no thread; it edited in the worktree instead.

### 2026-10-06T06:06:28Z
Story: S-0220.
T-0894: the story's agent's Write to /tmp/s0220 and the task sub-agent's write to /tmp were refused at once with 'Permission prompt tool returned an invalid result'; the ADR body went through a heredoc instead.

## Remediation

Story S-0283 remediates this issue, created from it at 2026-10-06T03:45:18Z.
Closed 2026-10-06T10:10:31Z: Fixed by S-0283 (T-1001). The cause: flai registered permission_prompt with the go-sdk's typed mcp.AddTool, which declared an output schema and answered structuredContent beside the text block, and Claude Code accepts only a single text block, so it read every allow and every refusal as an invalid result. permission_prompt is now a raw handler that answers one text block holding the JSON decision, with no output schema and no structured content; arguments it cannot read are a deny in the same shape. TestPermissionPromptAnswersOneTextBlockAsClaudeCodeRequires in flai/internal/mcpserver/permission_test.go reproduces I-0082, and every permission test now checks the shape. The S-0219 .claude/ write was refused at once because S-0219 was still ready, a refusal ADR-0086 intends; the /tmp writes were refused as outside any worktree, also intended; with the fix the agent reads those reasons. Running agents get the fix once a flai release with S-0283 is installed, since their MCP server is the installed flai.
