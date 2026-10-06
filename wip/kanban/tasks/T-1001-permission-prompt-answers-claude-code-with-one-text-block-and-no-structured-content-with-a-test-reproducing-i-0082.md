---
id: T-1001
type: task
nature: remediation
title: permission_prompt answers Claude Code with one text block and no structured content, with a test reproducing I-0082
status: backlog
parent: S-0283
owner: alex
created: 2026-10-06T06:26:49Z
updated: 2026-10-06T06:26:49Z
transitions: []
stream: S-0283
tags: [flai]
touches: [flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, flai/internal/mcpserver/folder.go]
---
# T-1001 permission_prompt answers Claude Code with one text block and no structured content, with a test reproducing I-0082

## Work

Find the cause before building the fix, and write the proposal in the narrative's `## Decisions` first, as the story's goal asks.

The likely cause: `folder.go` registers `permission_prompt` with `mcp.AddTool` and a typed handler, `mcp.ToolHandlerFor[PermissionIn, PermissionOut]`. The go-sdk (v1.8.0, `mcp/server.go`) then infers an output schema for the tool, validates `PermissionOut` against it, and sends the result as `structuredContent` beside a text block. Claude Code expects a permission tool to answer with exactly one content block of type `text` whose `text` is the JSON `{"behavior": ...}`. Confirm this over the wire: call the tool through an in-memory client session and read the raw `CallToolResult` and `tools/list` entry, not only `Content[0]`, which is all `askPermission` in `permission_test.go` reads, and why the tests passed.

Build the fix:

- Answer every outcome, the allow, the deny at once, the auto-approve, and the owner's answer on a thread, with a `CallToolResult` holding one `TextContent` of the JSON and no `StructuredContent`. Give the tool no output schema: register it with an untyped output, or with the server's raw `AddTool`, in `folder.go`.
- A deny keeps its reason in `message`, so an agent refused a write to `/tmp` or to the main checkout reads why.
- Check from I-0082's S-0219 instance which branch refused the story agent's `.claude/agents/orchestrator.md` write with no thread: the result's shape, or a deny before asking, such as a path outside the worktree. If it was a deny ADR-0086 does not intend, record it with `flai issue new` in the worktree rather than widen this task.

In `permission_test.go`, reproduce I-0082: for an allow and a deny, the result has exactly one content block, it is text, its text is a JSON object with `behavior`, and the result has no structured content; the tool lists no output schema.

## Done when

- the I-0082 test fails on the code before the change and passes after it
- the existing permission tests pass
- `scripts/flai-test.sh` passes

## Notes
