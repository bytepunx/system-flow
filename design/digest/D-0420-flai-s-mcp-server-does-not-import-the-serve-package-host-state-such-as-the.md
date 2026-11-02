---
id: D-0420
title: "flai's MCP server does not import the serve package: host state such as the budget reaches it through an Options function cmd wires to the serve folder, as AutoApprove does"
kind: decision
sources: [S-0534, S-0546, S-0506]
topics: [cli, go, mcp]
paths: [flai/internal/mcpserver, flai/cmd/mcp.go, flai/internal/serve/budget.go, flai/internal/closeout, flai/cmd/mcp_http.go, flai/internal/mcpserver/permission.go]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: [S-0506]
---

# D-0420 flai's MCP server does not import the serve package: host state such as the budget reaches it through an Options function cmd wires to the serve folder, as AutoApprove does

S-0534 put the host's budget state in every inbox the MCP server answers. Its narrative records how that state reaches the server without a new dependency (wip/archive/agents/S-0534.md, Decisions):

> `mcpserver` does not import `serve`: the inbox's `budget` reaches the MCP server through an `Options` function `cmd` wires to the serve folder, as `AutoApprove` does (assumption 5).

The same option carried it into the inboxes inside `task_done` and `story_start` answers (a `HostBudget` option on `taskdone` and `storystart`), and the HTTP MCP server is wired as the stdio one is. `design/system/flai-cli.md` § The host's budget state says who reads `budget.json`, but not this layering. A later story that shows other serve-folder state through MCP takes the same route: an option function set in `cmd`, never an import of `serve` from `mcpserver`.

S-0546 met the same seam from the other side. Its MCP tool `verify` hands the close-out to flai serve as a job (ADR-0229), and the request side could not live in `serve`: "T-2407 moved the request side of the hand-off out of `serve` into a new `flai/internal/closeout` package: `serve`'s internal tests import `mcpserver`, so `mcpserver` cannot import `serve`" (wip/archive/agents/S-0546.md, Decisions). The request needs no serve folder, since it is a file under the project's `.flai-cache/verify/jobs/`. So when `mcpserver` needs to ask serve for something, put what both sides share in a package below them, as `closeout` is, or wire it through an option set in `cmd`.

S-0506 took the option route again for `permission_prompt` over HTTP, whose `auto-approve` holds only for the run `flai serve` started for the story the path belongs to (ADR-0243): "The check lives in `cmd` as an option function, since `mcpserver` does not import `serve` (D-0420)" (wip/archive/agents/S-0506.md, Decisions).

Sources: S-0534, TH-0711 (assumption 5), S-0546, S-0506.
