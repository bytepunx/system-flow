---
id: T-0874
type: task
nature: remediation
title: flai's MCP server has a permission_prompt tool that asks the operator on a thread before a .claude/ write
status: done
parent: S-0257
owner: alex
created: 2026-10-05T04:14:59Z
updated: 2026-10-05T04:24:44Z
transitions:
  - to: ready
    at: 2026-10-05T04:15:52Z
    by: agent-S-0257
  - to: in-progress
    at: 2026-10-05T04:15:53Z
    by: agent-S-0257
  - to: done
    at: 2026-10-05T04:24:44Z
    by: agent-S-0257
stream: S-0257
tags: []
touches: [flai/internal/mcpserver]
usage:
  source: log
  seconds: 531
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 97
      output: 480
      cache_read: 4798640
      cache_write: 149034
      cost: 2.0427
---
# T-0874 flai's MCP server has a permission_prompt tool that asks the operator on a thread before a .claude/ write

## Work

Add the MCP tool `permission_prompt` to `flai/internal/mcpserver`, for Claude Code's `--permission-prompt-tool`.

- **Input:** `{tool_name, input, tool_use_id}`.
- **Output:** one text content holding JSON, either `{"behavior":"allow","updatedInput":<input>}` or `{"behavior":"deny","message":"..."}`.
- **What it may approve:** an Edit, Write, MultiEdit, or NotebookEdit whose `file_path` (or `notebook_path`) is an absolute path inside a story's worktree (`repo.WorktreePath(id)`), with a `.claude` segment below the worktree. The story must be in progress.
- **Everything else:** denied at once, with a message saying what the tool approves.
- **Auto-approve:** `Options` gains a func that says whether the operator auto-approves for the project's root. It is read at every call. When it says yes, the tool allows at once and logs that it did.
- **Otherwise, ask the operator.** Open a thread on the story as the agent, titled "Allow <Tool> <path relative to the worktree>?". The entry carries the path, plus the whole content (Write) or the old and new text (Edit) in fenced blocks. Then block until an entry by anyone but the agent arrives, or the context ends.
- **Reading the answer:** a first word of allow, yes, approve, or approved (any case) allows. Anything else denies, with the operator's text as the message. Resolve the thread with what was decided.
- **Context ends:** deny.

## Done when

- The tool is advertised; `TestToolsAreAdvertised` is updated.
- Tests cover:
  - a path outside any worktree, a non-`.claude` path, another tool, and a story not in progress: each denied at once;
  - auto-approve: allowed with no thread;
  - an operator's "allow" answer: allowed;
  - another answer: denied with that text.

## Notes

The tool is called by Claude Code, not by the model; the guard's sub-agent allowlists need no change.
