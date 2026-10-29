---
id: T-2227
type: task
nature: research
title: Claude Code's tool-call timeout over flai's HTTP MCP server is measured on the host and in the container, with MCP_TOOL_TIMEOUT unset and set, and the remedy is proposed on a thread
status: in-progress
parent: S-0519
owner: alex
created: 2026-10-28T13:44:42Z
updated: 2026-10-29T00:40:07Z
transitions:
  - to: ready
    at: 2026-10-29T00:40:05Z
    by: agent-S-0519
  - to: in-progress
    at: 2026-10-29T00:40:07Z
    by: agent-S-0519
stream: S-0519
tags: [cli, mcp, docker]
usage:
  source: log
  seconds: 1047
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 111
      output: 30705
      cache_read: 5113339
      cache_write: 324463
      cost: 3.7012
---
# T-2227 Claude Code's tool-call timeout over flai's HTTP MCP server is measured on the host and in the container, with MCP_TOOL_TIMEOUT unset and set, and the remedy is proposed on a thread

## Work

Measure how long Claude Code waits for a tool call over flai's HTTP MCP server before it gives up with `Error calling tool: The operation timed out.`: ADR-0124 took Claude Code's HTTP timeout for five minutes and set the hold at four, and I-0217 shows the timeout shorter than the hold, with no figure recorded. The probe is a headless `claude -p --output-format stream-json --verbose` on the host with `--mcp-config` naming a file in the shape `flai mcp status` prints for the project's HTTP server, the bearer from `flai mcp token`, `X-Flai-Agent` set to `agent-S-0519`, your own name, since `flai guard` refuses a call naming another author ([ADR-0143](../../../design/adrs/0143-flai-guard-refuses-a-story-s-agent-the-planner-and-the-analyzer-a-call-that.md)), and `--strict-mcp-config`, asked to call `wait_for_events` with `timeout_seconds` 290 and nothing else. That call holds the request open with nothing the client receives until it ends: `flai/internal/mcphttp/mcphttp.go` answers as `application/json`, so the progress notifications `keepAlive` in `flai/internal/mcpserver/timing.go` sends are not delivered while the call runs. Time from the call to the error in the stream. Run it with `MCP_TOOL_TIMEOUT` unset, then set above four minutes (such as `600000`), to learn whether Claude Code honours it over HTTP. Repeat the unset run inside `flai-agent-claude-code:local`, which is present, through `scripts/agent-container.sh S-0519` with a prompt on standard input that makes the same call rather than works the story: the instance was Claude Code 2.1.312 there, the host runs 2.1.311. Log the figures and both versions in the narrative with `flai stream log`.

Then propose the remedy on a thread on S-0519 with `thread_open`, your recommendation first and the figures beside it. Recommended: the HTTP server's `permission_prompt` holds a write at most a bound of its own, below the measured timeout, and zero, the refusal at once naming the thread as ADR-0124's refusal does, when the timeout leaves no useful hold; stdio keeps four minutes; the thread stays open and the same write made again takes the answer, as today. Name the alternatives and why not: `MCP_TOOL_TIMEOUT` set in the container's environment where Claude Code honours it, which rests on a harness variable rather than flai's bound and leaves a hand-run container as it was; SSE answers with progress notifications on the HTTP server, which keep the agent blocked while the operator is away, as ADR-0124 found. When a long call other than `permission_prompt` is cut short by the same timeout, `wait_for_events` up to its 30 minutes, `task_done`, or `verify`, record it as an issue with `flai issue new` from the worktree; fixing that is not this story's.

This task waits for nothing. When nothing is left but the operator's answer, write the narrative's Current state and Next steps naming the thread and end: `flai serve` starts you again when it is answered.

## Done when

The timeout over HTTP is a measured number of seconds on the host and in the container, with `MCP_TOOL_TIMEOUT` unset and set, logged in the narrative with the Claude Code versions; the thread on S-0519 proposes the remedy with the recommendation first and the figures; and the operator has answered it.

## Notes

This task changes no file of the repository: its findings go to the narrative and the thread, and the ADR of the next task records them.
