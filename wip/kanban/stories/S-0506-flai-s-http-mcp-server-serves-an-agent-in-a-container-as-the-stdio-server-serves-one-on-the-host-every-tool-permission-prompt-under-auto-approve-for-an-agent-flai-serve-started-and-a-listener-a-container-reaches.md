---
id: S-0506
type: story
nature: feature
title: "flai's HTTP MCP server serves an agent in a container as the stdio server serves one on the host: every tool, permission_prompt under auto-approve for an agent flai serve started, and a listener a container reaches"
status: backlog
parent: E-0029
owner: alex
created: 2026-10-26T04:42:51Z
updated: 2026-10-31T08:02:32Z
transitions: []
tags: [cli, mcp]
topics: [harness, docker, mcp]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, flai/internal/mcphttp/mcphttp.go, flai/internal/mcphttp/mcphttp_test.go, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/cmd/mcp_http_test.go, flai/cmd/host.go, flai/cmd/host_test.go, flai/cmd/serve_mcp.go, flai/cmd/serve_mcp_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/config/config.go, flai/internal/config/config_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, design/adrs/README.md, design/adrs, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, docs/operators/index.md, design/system/flai-cli.md, flai/cmd/serve.go, flai/cmd/serve_settings.go, flai/cmd/serve_settings_test.go, flai/internal/mcpserver/work.go, flai/internal/mcpserver/timing.go, design/issues/I-0228-claude-code-drops-any-call-to-flai-s-http-mcp-server-that-answers-after-60-seconds-so-wait-for-events-wait-for-work-verify-and-a-slow-task-done-fail-for-an-agent-in-a-container.md, design/issues/summary.md]
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
    - kind: planner
      seconds: 296
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 204
          output: 25920
          cache_read: 22262125
          cache_write: 138160
          cost: 6.5749
    - kind: orchestrator
      seconds: 646
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 150
          output: 1109
          cache_read: 12986216
          cache_write: 397875
          cost: 2.9612
cost_of_delay:
  value: 33.33
  by: planner-S-0506
  at: 2026-10-31T08:00:25Z
forecast:
  duration: 2h
  delivery: 2026-10-31T18:35:00Z
  basis: "flai forecast's 1h12m (110 s per unit over 101 done feature stories, size 39) raised to 2h: seven tasks, two of them ADRs, across five packages, with a by-hand check that a container reaches the listener; delivery is flai's playout plus the 48m the raise adds."
  by: planner-S-0506
  at: 2026-10-31T07:59:45Z
finalized:
  by: orchestrator
  at: 2026-10-31T08:02:32Z
---
# S-0506 flai's HTTP MCP server serves an agent in a container as the stdio server serves one on the host: every tool, permission_prompt under auto-approve for an agent flai serve started, and a listener a container reaches

## Goal

An agent in a container reaches flai only through the HTTP MCP server `flai host` keeps for the project ([ADR-0030](../../../design/adrs/0030-mcp-is-served-by-flai-on-the-host-over-stdio-and-http-and-the-dashboard-s-api.md), [ADR-0040](../../../design/adrs/0040-one-flai-host-per-machine-runs-flai-serve-and-each-project-s-mcp-server-as-its.md)), so that server must give it what `flai mcp` on stdio gives an agent on the host today. Three things stand in the way, and the spike S-0505 measures each: the HTTP server is started without the `auto-approve` host action, so `permission_prompt` over HTTP always asks (`flai/cmd/mcp_http.go` leaves `AutoApprove` out where `flai/cmd/mcp.go` passes it); the listener binds `127.0.0.1:4243`, which a container on Docker's bridge network cannot reach; and the token is one file under `.flai-cache`, read by hand. This story closes them, and whatever else the spike's issues name on the HTTP server, so that the container adapter can hand a run the URL and a token by environment and the agent inside works a story as one on the host does.

## Acceptance criteria

- [ ] A test compares the tool names the stdio server and the HTTP server register for the same agent and fails on any difference, and every tool the spike found missing or failing over HTTP, from the issues S-0505 filed, answers over HTTP as it does on stdio, `task_done`, `stream_sync`, `permission_prompt`, and `wait_for_events` among them
- [ ] An ADR settles whether `permission_prompt` over HTTP honours the `auto-approve` host action for an agent `flai serve` started and named, as the stdio server does, while the dashboard token still cannot grant it, and the HTTP server behaves as the ADR says, with a test for an agent flai serve started and one it did not
- [ ] `flai host` starts each project's `flai mcp http` on an address containers reach when the host setting `flai serve mcp` says so, loopback by default: `flai serve mcp bind <address>` sets it, `flai serve settings` shows it, a container on the default bridge reaches the listener through Docker's host gateway, and the listener still refuses a request with an `Origin` header and one without the bearer token
- [ ] `flai mcp token` and the server accept a token `flai serve` hands a run by environment as `FLAI_MCP_TOKEN`, with the URL as `FLAI_MCP_URL`, so that the container adapter can pass both without writing either into an image, a start request, or argv, and the server's log names the agent from `X-Flai-Agent`, never the token
- [ ] `docs/users/flai.md` § MCP over HTTP, `docs/operators/settings.md` (the new setting and the two variables, so `settings_doc_test` passes), `docs/operators/index.md`, and `design/system/flai-cli.md` describe the setting, the variables, and the behaviour the ADR settles

## Tasks
- T-2150 A test compares the tool names the stdio and HTTP servers register for one agent, and every tool the spike found missing over HTTP is registered there
- T-2151 An ADR settles permission_prompt over HTTP under the auto-approve host action for an agent flai serve started, and the HTTP server follows it
- T-2152 flai serve mcp bind sets the address flai host gives each project's MCP server, loopback by default, and a container on the bridge reaches it through the host gateway
- T-2160 FLAI_MCP_URL and FLAI_MCP_TOKEN hand a run the server and its token by environment, flai mcp token prints the one in force, and the log names the agent and never the token
- T-2167 docs/users/flai.md, docs/operators/settings.md, docs/operators/index.md, and flai-cli.md describe mcp.bind, the two variables, and permission_prompt over HTTP
- T-2418 flai serve settings prints the host settings the dashboard's settings.get gives, its mcp block naming the bind address flai serve mcp sets
- T-2419 wait_for_events and wait_for_work over HTTP answer within 45 seconds, below Claude Code's 60-second tool-call timeout, as an ADR refining ADR-0215 settles, and I-0228 says what is left

## Notes

What this story does not do: run anything in a container itself; the adapter is S-0508. It waits on S-0505 because the spike's issues say which tools the HTTP server lacks and which network mode reached the listener on this WSL2 host; the known gaps above stand whatever the spike finds.

### Planning

Planned again by planner-S-0506 on 2026-10-31, after the operator edited criterion 3 and resolved TH-0664. The first plan was planner-E-0029's, revisited by planner-S-0506 on 2026-10-26. The five tasks stand as written. Two tasks are added, and T-2167 now waits for both.

What changed since 2026-10-26:

- **Criterion 3.** The operator's edit names the host setting `flai serve mcp` and keeps `flai serve settings` as what shows it. That command does not exist. T-2418 adds it, printing what `hostSettings` gives the dashboard's `settings.get`, with the bind address in its `mcp` block. T-2152's Work already says to show the setting in `flai serve settings`, so T-2418 builds that command after T-2152.
- **The spike.** S-0505 is done. Of the issues it routed to this story, I-0213 (403 on a non-loopback `Host`) is fixed by S-0515 (ADR-0216) and I-0217 (`permission_prompt` outliving Claude Code's timeout) by S-0519 (ADR-0215). The spike found every tool its agent called answering over HTTP, so T-2150's parity test and its registrations stand as written. It reached the loopback listener through `host.docker.internal`, Docker Desktop's host gateway, on this WSL2 host, so `mcp.bind` matters for a Docker engine without that gateway. T-2152 stands.
- **I-0228.** S-0519 filed it on 2026-10-29: Claude Code drops any HTTP call that answers after 60 seconds, so `wait_for_events` and `wait_for_work` fail for an agent in a container. Criterion 1 names `wait_for_events`, and no other story remediates I-0228, so T-2419 bounds both waits over HTTP at 45 seconds, as ADR-0215 bounds `permission_prompt`. `verify`'s long run is S-0546's (the close-out as a `flai serve` job), and T-2419 records that on the issue.

Touches. `flai touches suggest S-0506` over the 27 touches declared (949 of 2692 commits) adds nothing. Its co-changes (`flaiover-dashboard.md`, `docs/users/flaiover.md`, `workflow.md`, `hostapi/writes.go`, `template/CHANGELOG.md`, `project-manifest.md`, `cmd/root.go`, `mcpserver/folder_test.go`, `harness.go`) are left out. The story adds no top-level command, no hostapi method, no harness, and no manifest key, and it changes no dashboard. Every declared touch is kept with the source recorded before:

- **Co-change with the three seeds:** the `mcpserver`, `mcphttp`, and `cmd/mcp*` files.
- **Layout:**
  - `cmd/host.go`, `config.go`, `serve_actions.go`, `serve/agents.go`, `docs/operators/index.md`, and their tests.
  - `cmd/serve_mcp.go` and its test, where `mcpAddrFor` sets each project's address.
- **Added from the layout on 2026-10-31:**
  - T-2418: `flai/cmd/serve.go`, which registers the `flai serve` subcommands, and the new `flai/cmd/serve_settings.go` and its test.
  - T-2419: `flai/internal/mcpserver/timing.go`, where `LongestWait` lives, and `flai/internal/mcpserver/work.go`, which holds `wait_for_work` for one project.
  - T-2419 also touches I-0228's file and `design/issues/summary.md`, through `flai issue`.

The folder `design/adrs` is kept from the layout. `flai adr new` in T-2151 and T-2419 names each ADR's file by its number and the title the story's agent writes, so no task can name it. The folder is kept for those two files alone, and `flai task done` widens the touches with each committed path (ADR-0144).

Forecast. `flai forecast S-0506` gives 1h12m: a median of 110 s per unit over 101 done feature stories on claude-opus-5-5 in the large band, times size 39 (5 criteria, 34 touches). That puts delivery at 2026-10-31T17:47Z, 15th in the pull order. The forecast is raised to 2h, up from the 1h30m kept on 2026-10-26, because the story now has seven tasks. Two of them write an ADR, the changes span five packages (cmd, config, mcpserver, mcphttp, serve), and one step checks by hand that a container reaches the listener. The band's median comes from stories whose touches are mostly documents. The delivery written, 2026-10-31T18:35Z, is flai's playout plus the 48m the raise adds.

Cost of delay. `flai cod S-0506` gives 33.33 USD a week, written as flai gives it. That is the story's share of E-0029's `penalty_per_week` of 150 USD, which the operator set on TH-0657. The share goes by forecast duration: 2h of 9h over the epic's four open stories without inputs. The story has no inputs of its own and needs none, since it delivers for the epic alone.
