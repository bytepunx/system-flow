---
id: S-0506
type: story
nature: feature
title: "flai's HTTP MCP server serves an agent in a container as the stdio server serves one on the host: every tool, permission_prompt under auto-approve for an agent flai serve started, and a listener a container reaches"
status: backlog
parent: E-0029
owner: alex
created: 2026-10-26T04:42:51Z
updated: 2026-10-31T06:24:12Z
transitions: []
tags: [cli, mcp]
topics: [harness, docker, mcp]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, flai/internal/mcphttp/mcphttp.go, flai/internal/mcphttp/mcphttp_test.go, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/cmd/mcp_http_test.go, flai/cmd/host.go, flai/cmd/host_test.go, flai/cmd/serve_mcp.go, flai/cmd/serve_mcp_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/config/config.go, flai/internal/config/config_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, design/adrs/README.md, design/adrs, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, docs/operators/index.md, design/system/flai-cli.md]
after: [S-0505]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 376
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 66
          output: 515
          cache_read: 9969521
          cache_write: 180030
          cost: 2.2277
draft: true
cost_of_delay:
  value: 16.67
  by: planner-S-0506
  at: 2026-10-26T12:30:19Z
forecast:
  duration: 1h30m
  delivery: 2026-10-31T11:32:00Z
  basis: "Its own forecast of 1h30m; 11th in the pull order with an in-progress limit of 5, behind S-0424, S-0425, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411 and S-0412."
  by: flai
  at: 2026-10-31T06:24:12Z
---
# S-0506 flai's HTTP MCP server serves an agent in a container as the stdio server serves one on the host: every tool, permission_prompt under auto-approve for an agent flai serve started, and a listener a container reaches

## Goal

An agent in a container reaches flai only through the HTTP MCP server `flai host` keeps for the project ([ADR-0030](../../../design/adrs/0030-mcp-is-served-by-flai-on-the-host-over-stdio-and-http-and-the-dashboard-s-api.md), [ADR-0040](../../../design/adrs/0040-one-flai-host-per-machine-runs-flai-serve-and-each-project-s-mcp-server-as-its.md)), so that server must give it what `flai mcp` on stdio gives an agent on the host today. Three things stand in the way, and the spike S-0505 measures each: the HTTP server is started without the `auto-approve` host action, so `permission_prompt` over HTTP always asks (`flai/cmd/mcp_http.go` leaves `AutoApprove` out where `flai/cmd/mcp.go` passes it); the listener binds `127.0.0.1:4243`, which a container on Docker's bridge network cannot reach; and the token is one file under `.flai-cache`, read by hand. This story closes them, and whatever else the spike's issues name on the HTTP server, so that the container adapter can hand a run the URL and a token by environment and the agent inside works a story as one on the host does.

## Acceptance criteria

- [ ] A test compares the tool names the stdio server and the HTTP server register for the same agent and fails on any difference, and every tool the spike found missing or failing over HTTP, from the issues S-0505 filed, answers over HTTP as it does on stdio, `task_done`, `stream_sync`, `permission_prompt`, and `wait_for_events` among them
- [ ] An ADR settles whether `permission_prompt` over HTTP honours the `auto-approve` host action for an agent `flai serve` started and named, as the stdio server does, while the dashboard token still cannot grant it, and the HTTP server behaves as the ADR says, with a test for an agent flai serve started and one it did not
- [ ] `flai host` starts each project's `flai mcp http` on an address containers reach when the host setting `mcp.bind` says so, loopback by default: `flai serve mcp bind <address>` sets it, `flai serve settings` shows it, a container on the default bridge reaches the listener through Docker's host gateway, and the listener still refuses a request with an `Origin` header and one without the bearer token
- [ ] `flai mcp token` and the server accept a token `flai serve` hands a run by environment as `FLAI_MCP_TOKEN`, with the URL as `FLAI_MCP_URL`, so that the container adapter can pass both without writing either into an image, a start request, or argv, and the server's log names the agent from `X-Flai-Agent`, never the token
- [ ] `docs/users/flai.md` § MCP over HTTP, `docs/operators/settings.md` (the new setting and the two variables, so `settings_doc_test` passes), `docs/operators/index.md`, and `design/system/flai-cli.md` describe the setting, the variables, and the behaviour the ADR settles

## Tasks
- T-2150 A test compares the tool names the stdio and HTTP servers register for one agent, and every tool the spike found missing over HTTP is registered there
- T-2151 An ADR settles permission_prompt over HTTP under the auto-approve host action for an agent flai serve started, and the HTTP server follows it
- T-2152 flai serve mcp bind sets the address flai host gives each project's MCP server, loopback by default, and a container on the bridge reaches it through the host gateway
- T-2160 FLAI_MCP_URL and FLAI_MCP_TOKEN hand a run the server and its token by environment, flai mcp token prints the one in force, and the log names the agent and never the token
- T-2167 docs/users/flai.md, docs/operators/settings.md, docs/operators/index.md, and flai-cli.md describe mcp.bind, the two variables, and permission_prompt over HTTP

## Notes

What this story does not do: run anything in a container itself; the adapter is S-0508. It waits on S-0505 because the spike's issues say which tools the HTTP server lacks and which network mode reached the listener on this WSL2 host; the known gaps above stand whatever the spike finds.

### Planning

Planned by planner-S-0506 on 2026-10-26, revisiting the plan planner-E-0029 wrote with the epic the same day. The five tasks stand as written; T-2152 gains two touches and the story the same two.

Touches. `flai touches suggest S-0506` over the 25 touches declared (846 of 2427 commits) adds nothing: its co-changes (`flaiover-dashboard.md`, `docs/users/flaiover.md`, `workflow.md`, `hostapi/writes.go`, `issues/summary.md`, `template/CHANGELOG.md`, `cmd/root.go`, `mcpserver/folder_test.go`, `harness.go`) are left out, since the story adds no command to `root.go`, no hostapi method, no harness, and changes no dashboard. Every declared touch is kept, each with the source planner-E-0029 recorded: the `mcpserver`, `mcphttp`, and `cmd/mcp*` files from co-change with the three seeds; `cmd/host.go`, `config.go`, `serve_actions.go`, `serve/agents.go`, `docs/operators/index.md`, and their tests from the layout, where `hostLauncher.mcpAddr` chooses each project's address, the host settings live, `flai serve` settings are set, and `spawn` sets a run's environment. Added from the layout: `flai/cmd/serve_mcp.go` and `flai/cmd/serve_mcp_test.go`, since `mcpAddrFor` there takes the host of every project's address from `mcpDefaultAddr` and the remembered address file (`mcpAddrFile`) keeps the host as well, so `mcp.bind` is applied there and a remembered loopback address is re-hosted when the setting changes (T-2152). The folder `design/adrs` is kept (layout): `flai adr new` in T-2151 names the ADR's file by its number and the title the story's agent writes, so no task can name it; the folder is kept for that file alone, and `flai task done` widens the touches with the committed path (ADR-0144).

Forecast. `flai forecast S-0506` gives 54m (median 107 s per unit over 93 done feature stories on claude-opus-5-5 in the large band; size 30: 5 criteria, 25 touches) and delivery 2026-10-28T00:21Z, 24th in the pull order by a cycle factor of 1.17. Kept at 1h30m for the reason planner-E-0029 gave: the story writes an ADR and changes five packages (cmd, config, mcpserver, mcphttp, serve) with a by-hand check that a container reaches the listener, where the band's median is set by stories whose touches are mostly documents. The delivery written, 2026-10-28T00:57Z, is flai's playout plus the 36m the raised duration adds.

Cost of delay. `flai cod S-0506` gives 16.67 USD a week: the story's share of E-0029's `penalty_per_week` of 150 USD, which the operator set on TH-0657, by forecast duration, 1h30m of 13h30m over the epic's six open stories without inputs. Written as flai gives it: the story has no inputs of its own and needs none, since it delivers for the epic alone.
