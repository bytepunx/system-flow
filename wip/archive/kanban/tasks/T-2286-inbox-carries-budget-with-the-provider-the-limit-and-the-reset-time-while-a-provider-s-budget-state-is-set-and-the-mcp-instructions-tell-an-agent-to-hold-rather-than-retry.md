---
id: T-2286
type: task
nature: improvement
title: inbox carries budget with the provider, the limit, and the reset time while a provider's budget state is set, and the MCP instructions tell an agent to hold rather than retry
status: done
parent: S-0534
owner: alex
created: 2026-10-28T21:26:36Z
updated: 2026-10-29T06:33:15Z
transitions:
  - to: ready
    at: 2026-10-29T05:56:34Z
    by: agent-S-0534
  - to: in-progress
    at: 2026-10-29T05:56:35Z
    by: agent-S-0534
  - to: done
    at: 2026-10-29T06:33:15Z
    by: agent-S-0534
stream: S-0534
tags: [flai, mcp, inbox, budget]
touches: [flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/folder_test.go, flai/internal/mcpserver/inbox_test.go, flai/cmd/mcp.go, flai/cmd/mcp_test.go, flai/cmd/mcp_http.go, flai/internal/taskdone/taskdone.go, flai/internal/storystart/start.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/story_start.go]
after: [T-2282]
usage:
  source: log
  seconds: 2200
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 60
      output: 27039
      cache_read: 4030499
      cache_write: 135902
      cost: 2.1
---
# T-2286 inbox carries budget with the provider, the limit, and the reset time while a provider's budget state is set, and the MCP instructions tell an agent to hold rather than retry

## Work

Give `inbox.Inbox` a `budget` field, omitted when empty: one entry per provider whose state is set, each with `provider`, `limit`, `since`, `until`, and `probing` when a probe runs.

The MCP server reads the state, without writing it, from `budget.json` in the serve folder beside the host configuration. It finds that folder as `flai serve` does: the config path `--config`, then `FLAI_CONFIG`, then `~/.flai/config.json`, through `serve.DirFor`. Pass the reader in through `mcpserver`'s options as a callback, as `AutoApprove` is passed (`server.go`, `cmd/mcp.go`), so tests give a fixed state. A missing or unreadable file is no budget, never an error.

Every agent on the host gets the field: a story's agent, a sub-agent, and each role, whether its MCP server runs over stdio (`flai mcp`) or over HTTP. The folder form of the server (`folder.go`) carries it too.

Add a sentence to the MCP server's instructions, in both `server.go` and `folder.go`: when `inbox` carries `budget` for your provider, the host is out of budget until `until`. Do not retry or start work that calls the model in a loop. Write your narrative's state and end; `flai serve` starts you again once the limit resets.

This task waits for T-2282's state and reader. It runs beside the run-end task, which touches none of these files.

## Done when

- [ ] `inbox` returns `budget` with the provider, the limit, and `until` while the state is set, and leaves it out once it clears, tested over stdio and in the folder form.
- [ ] A missing serve folder or `budget.json` gives no `budget` and no error.
- [ ] The instructions name the field and say to hold rather than retry, and a test checks the sentence is in both forms.
- [ ] `flai test` passes on the changed paths.

## Notes

Planned by planner-S-0534. Reading the serve folder from the MCP server, rather than a copy in each project's cache, is an assumption on the plan's thread.
