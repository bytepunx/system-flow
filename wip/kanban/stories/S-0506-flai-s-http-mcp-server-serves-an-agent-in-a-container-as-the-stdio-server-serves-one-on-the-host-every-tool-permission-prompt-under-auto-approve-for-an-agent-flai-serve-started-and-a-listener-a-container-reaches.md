---
id: S-0506
type: story
nature: feature
title: "flai's HTTP MCP server serves an agent in a container as the stdio server serves one on the host: every tool, permission_prompt under auto-approve for an agent flai serve started, and a listener a container reaches"
status: backlog
parent: E-0029
owner: alex
created: 2026-10-26T04:42:51Z
updated: 2026-10-26T04:57:14Z
transitions: []
tags: [cli, mcp]
topics: [harness, docker, mcp]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, flai/internal/mcphttp/mcphttp.go, flai/internal/mcphttp/mcphttp_test.go, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/cmd/mcp_http_test.go, flai/cmd/host.go, flai/cmd/host_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/config/config.go, flai/internal/config/config_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, design/adrs/README.md, design/adrs, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, docs/operators/index.md, design/system/flai-cli.md]
after: [S-0505]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 1h30m
  delivery: 2026-10-27T05:35:00Z
  basis: "flai forecast gave 54m (median 107 s per unit over 93 done feature stories on claude-opus-5-5 in the large band, size 30: 5 criteria, 25 touches) and delivery 2026-10-27T04:52Z, 25th in the pull order by a cycle factor of 1.17; raised to 1h30m because the story writes an ADR and changes five packages (cmd/host, config, mcpserver, mcphttp, serve) with a by-hand check that a container reaches the listener, where the band's median is set by stories whose touches are mostly documents; delivery is flai's start of 2026-10-27T03:49Z plus 1h30m times the same cycle factor."
  by: planner-E-0029
  at: 2026-10-26T04:57:14Z
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

Planned by planner-E-0029 on 2026-10-26 with the epic.

Touches. `flai touches suggest S-0506`, seeded with `flai/cmd/mcp_http.go`, `flai/internal/mcphttp/mcphttp.go`, and `flai/internal/mcpserver/permission.go` (24 commits), names `mcpserver/server.go`, `cmd/mcp.go`, `mcpserver/folder.go`, `folder_test.go`, `server_test.go`, `permission_test.go`, `flai-cli.md`, `docs/users/flai.md`, `settings.md`, `flai-reference.md`, `cmd/mcp_http_test.go`, `cmd/root.go`, `hostapi/writes.go`, and `serve/serve.go`; every one but `root.go`, `hostapi/writes.go`, `folder_test.go`, and `serve.go` is kept (co-change), since the story adds no command and no hostapi method. The rest are the tasks' (layout):

- `flai/cmd/host.go`, `host_test.go`, `flai/internal/config/config.go`, `config_test.go`, `flai/cmd/serve_actions.go`, `serve_actions_test.go`: `hostLauncher.mcpAddr` chooses each project's address, the host setting lives in `config.go`, and `flai serve` settings are in `serve_actions.go`. T-2152.
- `flai/internal/serve/agents.go`, `agents_test.go`: `spawn` sets the run's environment. T-2160.
- `docs/operators/index.md`: the operator's view of the HTTP server. T-2167.
- `design/adrs/README.md` and the folder `design/adrs` (layout, folder kept): `flai adr new` in T-2151 names the ADR's file by its number and the title the story's agent writes, so no task can name it; the folder is kept for that file alone, and `flai task done` widens the touches with the committed path (ADR-0144).

Forecast and cost of delay: see the second Planning paragraph, written once the figures are set.
