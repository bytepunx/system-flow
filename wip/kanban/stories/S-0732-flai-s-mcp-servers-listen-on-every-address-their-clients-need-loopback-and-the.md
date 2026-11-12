---
id: S-0732
type: story
nature: improvement
title: flai's MCP servers listen on every address their clients need, loopback and the container engine's gateway together, set as a list or found by flai for the host's engine
status: ready
owner: alex
created: 2026-11-12T08:00:53Z
updated: 2026-11-12T16:04:53Z
transitions:
  - to: ready
    at: 2026-11-12T15:59:55Z
    by: orchestrator
tags: [cli]
topics: [mcp, docker]
touches: [flai/internal/config/config.go, flai/internal/config/config_test.go, flai/internal/serve/mcpbind.go, flai/internal/serve/mcpbind_test.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go, flai/internal/serve/agents_start.go, flai/cmd/serve_mcp.go, flai/cmd/serve_mcp_test.go, flai/cmd/serve_settings.go, flai/cmd/serve_settings_test.go, flai/cmd/serve_view.go, flai/cmd/serve_actions_test.go, flai/cmd/mcp_http.go, flai/cmd/mcp_http_test.go, flai/cmd/host.go, flai/cmd/host_test.go, docs/users/flai-reference.md, docs/users/flai.md, docs/operators/settings.md, docs/operators/index.md, design/system/flai-cli.md, design/adrs/README.md, design/adrs]
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
      seconds: 14
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 10
          output: 1677
          cache_read: 792413
          cache_write: 11511
          cost: 0.2776
    - kind: orchestrator
      seconds: 38
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 47
          cache_read: 426213
          cache_write: 12890
          cost: 0.125
cost_of_delay:
  inputs:
    penalty_per_week: 20
    by: orchestrator
    at: 2026-11-12T15:58:04Z
  value: 20
  by: planner-S-0732
  at: 2026-11-12T15:58:18Z
forecast:
  duration: 1h30m
  delivery: 2026-11-12T18:50:00Z
  basis: "Its own forecast of 1h30m; 2nd in the pull order with an in-progress limit of 3, behind S-0659, S-0683, S-0709 and S-0723."
  by: flai
  at: 2026-11-12T16:04:53Z
finalized:
  by: orchestrator
  at: 2026-11-12T15:59:19Z
---
# S-0732 flai's MCP servers listen on every address their clients need, loopback and the container engine's gateway together, set as a list or found by flai for the host's engine

## Goal

`mcp.bind` (`flai serve mcp bind`) takes one address for every project's MCP server on the host, 127.0.0.1 by default. Which address an agent's container reaches the host on depends on the Docker engine: Docker Desktop routes `host.docker.internal` to the host's loopback, and a native engine routes it to the gateway of the container's network, 172.17.0.1 for the default bridge. On 2026-11-12 the operator's host ran the native engine after a restart without Docker Desktop, and the start of a planner was refused with "routes host.docker.internal to the gateway of the network "bridge", 172.17.0.1, which a server on 127.0.0.1 does not answer". Each way out with one address costs something: binding 172.17.0.1 takes the servers off loopback for clients on the host, binding 0.0.0.0 puts them on the LAN with a token flai does not encrypt, and `--network host` gives up the container's network isolation.

After this story the servers listen on more than one address at once, so that loopback clients and containers are both served, and the operator either names the addresses or lets flai find them for the engine it runs on.

## Acceptance criteria
- [ ] `mcp.bind` holds one address or several, and `flai host` starts each project's MCP server listening on all of them; `flai serve mcp bind` takes several, `flai serve mcp show` prints them, and a configuration with one address reads as before
- [ ] A value that leaves the choice to flai, such as `auto`, listens on 127.0.0.1 and, when the host's Docker engine routes `host.docker.internal` to a network's gateway rather than to loopback, on the gateway of the network the container harness runs on; on an engine that reaches loopback, or with no Docker, it listens on 127.0.0.1 alone
- [ ] With loopback and the gateway both listened on, a container's agent on a native engine starts and reaches the server with neither `flai serve mcp bind 172.17.0.1` nor `--network host`, and a client on the host still reaches it on 127.0.0.1, shown by a test
- [ ] The refusal in `flai/internal/serve/container.go` names the new way first, and is still given when the addresses listened on do not include one the container reaches
- [ ] The warning for listening beyond this machine is still logged for an address that is neither loopback nor a container network's gateway
- [ ] `docs/operators/` (the settings table and the MCP server section), `design/system/flai-cli.md`, and the command help say what `mcp.bind` holds and what the automatic value does; an ADR records whether the automatic value becomes the default

## Tasks
- T-3144 mcp.bind holds one address, several, or auto, and a configuration with one address reads as before
- T-3145 One lookup finds the host's Docker engine and the gateway of the container harness's network, for the automatic value and the container check alike
- T-3146 flai serve mcp bind takes several addresses or auto, and flai serve mcp show and flai serve settings print them
- T-3147 flai mcp http listens on several addresses on one port, records each in its state, and warns only of an address beyond loopback and the container gateways
- T-3148 flai host starts each project's MCP server on every address mcp.bind gives, resolving auto at each start, on a port free on all of them
- T-3150 A container start passes when the server listens on an address the container reaches, and the refusal names flai serve mcp bind auto first
- T-3153 An ADR records the automatic value and whether it is the default, and the operator docs and flai-cli.md say what mcp.bind holds

## Notes

What exists: `config.MCPBind`, `WithMCPBind`, and `checkMCPBind` in `flai/internal/config/config.go`; the check of the engine's route and the refusal in `flai/internal/serve/container.go` (S-0619, I-0254), which already finds the gateway the automatic value needs.

Open for planning: whether `auto` is the default for a new host, since the engine can change under a host, as it did here between Docker Desktop and the native engine; and what the servers do when the engine or the gateway changes while they run, restart by themselves or report it.

Related: E-0036, which the same change of engine prompted, since the engine's image store changed with it.

### Planning

Touches, by where each came from (the story declared none):

- Design: `flai/internal/config/config.go` and its test (`MCPHost.Bind`, `MCPBind`, `WithMCPBind`, `checkMCPBind`), and `flai/internal/serve/container.go` and its test (`containerReachable`, `containerEngine`, the refusal), both named in the Notes.
- Layout: `flai/cmd/serve_mcp.go` and its test (`flai serve mcp bind`, `show`, `clear`, `mcpBind`, `mcpAddrFor`), `flai/cmd/serve_settings.go`, `flai/cmd/serve_view.go`, and their tests, with `flai/cmd/serve_actions_test.go` (`TestServeMCPBind`, the settings JSON's `mcp` block), which print and set `mcp.bind`; `flai/cmd/mcp_http.go` and its test, where `flai mcp http` listens on one `--addr`, writes its state, and logs the warning for listening beyond this machine; `flai/cmd/host.go` and its test (`hostLauncher.mcp`, `mcpAddr`), where `flai host` starts each server on `mcp.bind`; `flai/internal/serve/agents_start.go` (`mcpEnv`), which reads the server's state for the agent's URL.
- New file: `flai/internal/serve/mcpbind.go` and its test, the engine and gateway lookup the automatic value and the container check share, moved out of `container.go`.
- Co-change (`flai touches suggest` from the files above): `docs/users/flai-reference.md` (21%, generated from the command help), `docs/operators/settings.md` (25%), `design/system/flai-cli.md` (16%), `docs/users/flai.md` (16%), `docs/operators/index.md` (8%). The suggestions under `flai/internal/mcpserver/` and `flai/cmd/serve_actions.go` are left out: the tools and the host actions do not change.
- Folder `design/adrs/`: kept, because the ADR criterion 6 asks for is a new file whose number and slug `flai adr new` gives only when it is written; `design/adrs/README.md`, its index, is named too.

No `after`: the story reads nothing another open story writes.

Forecast: `flai forecast` gave 54m, 107 s per unit of size over 175 done improvement stories in the large band, times size 30 (6 criteria, 24 touches). Raised to 1h30m: seven tasks in four layers cross the configuration, the commands, the host, the HTTP server, and the container check, and end with an ADR and four documents. S-0506, which built the listener and the container route across the same files, took 80 minutes. The delivery moves by the same 36 minutes, to 2026-11-13T00:55:00Z.

Cost of delay: 20 USD a week, as `flai cod` works it out from the one input the orchestrator set on TH-1099, `penalty_per_week: 20`. It stands unadjusted. The penalty is the whole of it: the story brings no revenue, and the time each refused container start loses is already priced in the penalty.

Design assumptions, asked on TH-1099: `auto` is the default for an unset `mcp.bind`, and a running server reports a change of engine rather than restarting by itself. T-3144, T-3148, and T-3153 build on them.
