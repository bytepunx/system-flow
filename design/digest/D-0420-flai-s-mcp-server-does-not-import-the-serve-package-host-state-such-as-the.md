---
id: D-0420
title: "flai's MCP server does not import the serve package: host state such as the budget reaches it through an Options function cmd wires to the serve folder, as AutoApprove does"
kind: decision
sources: [S-0534]
topics: [cli, go]
paths: [flai/internal/mcpserver, flai/cmd/mcp.go, flai/internal/serve/budget.go]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0420 flai's MCP server does not import the serve package: host state such as the budget reaches it through an Options function cmd wires to the serve folder, as AutoApprove does

S-0534 put the host's budget state in every inbox the MCP server answers. Its narrative records how that state reaches the server without a new dependency (wip/archive/agents/S-0534.md, Decisions):

> `mcpserver` does not import `serve`: the inbox's `budget` reaches the MCP server through an `Options` function `cmd` wires to the serve folder, as `AutoApprove` does (assumption 5).

The same option carried it into the inboxes inside `task_done` and `story_start` answers (a `HostBudget` option on `taskdone` and `storystart`), and the HTTP MCP server is wired as the stdio one is. `design/system/flai-cli.md` § The host's budget state says who reads `budget.json`, but not this layering. A later story that shows other serve-folder state through MCP takes the same route: an option function set in `cmd`, never an import of `serve` from `mcpserver`.

Sources: S-0534, TH-0711 (assumption 5).
