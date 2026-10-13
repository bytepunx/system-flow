---
id: T-1387
type: task
nature: improvement
title: usage.Event and usage.Reader are defined, and a Claude Code reader turns stream-json into them
status: done
parent: S-0352
owner: alex
created: 2026-10-08T08:45:48Z
updated: 2026-10-12T20:49:03Z
transitions:
  - to: ready
    at: 2026-10-12T20:40:52Z
    by: agent-S-0352
  - to: in-progress
    at: 2026-10-12T20:40:53Z
    by: agent-S-0352
  - to: done
    at: 2026-10-12T20:49:03Z
    by: agent-S-0352
stream: S-0352
tags: [cli]
touches: [flai/internal/usage/event.go, flai/internal/usage/claudecode.go, flai/internal/usage/claudecode_test.go]
usage:
  source: log
  seconds: 490
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 51
      output: 20143
      cache_read: 3256894
      cache_write: 91963
      cost: 1.6122
---
# T-1387 usage.Event and usage.Reader are defined, and a Claude Code reader turns stream-json into them

## Work

- `event.go`: `Reader interface{ Read(line []byte) ([]Event, bool) }`, since one stream-json line can hold several tool calls, and `Event` with `Kind` (`SessionStart`, `Call`, `ToolCall`, `ToolResult`, `End`), `At`, `Session`, `Model`, `MessageID`, `Tokens` (`Input`, `Output`, `CacheRead`, `CacheWrite`), `CostUSD *float64`, `Models` (per-model totals with cost, for `End`), `ParentCall`, `Tool`, `ToolID`, `Description`, `Prompt`, `IsError`, and `Result` (the raw result for the tools whose result flai reads, `wait_for_events` among them).
- `claudecode.go`: the reader for stream-json, mapping `system`/`init`, `assistant` (`message.id`, `model`, `usage`, `content[]` `tool_use` with `name`, `id`, `input.description`, `input.prompt`), `user` (`tool_result` with `tool_use_id`, `is_error`, `task_description`), and `result` (`total_cost_usd`, `usage`, `modelUsage`), with `parent_tool_use_id` on each. A line that is not an event yields none.
- Tests from lines of real logs, kept as testdata or inline: one of each event, a repeated message ID, a sub-agent's call, a `wait_for_events` result.

First layer: it runs together with T-1388's capabilities, which touch other files.

## Done when

- `claudecode_test.go` covers every field `log.go` reads today.
- `flai test flai/internal/usage/` passes.

## Notes

Layer 1 of S-0352.
