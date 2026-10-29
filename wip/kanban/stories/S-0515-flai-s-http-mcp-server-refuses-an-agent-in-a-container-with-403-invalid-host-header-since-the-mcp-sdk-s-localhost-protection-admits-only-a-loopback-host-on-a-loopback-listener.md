---
id: S-0515
type: story
nature: remediation
title: flai's HTTP MCP server refuses an agent in a container with 403 invalid Host header, since the MCP SDK's localhost protection admits only a loopback Host on a loopback listener
status: backlog
owner: alex
created: 2026-10-28T12:52:34Z
updated: 2026-10-29T04:05:47Z
transitions: []
tags: [cli, mcp]
topics: [mcp, docker]
touches: [flai/internal/mcphttp/mcphttp.go, flai/internal/mcphttp/mcphttp_test.go, flai/cmd/mcp_http.go, docs/users/flai-reference.md, design/adrs/README.md, design/adrs, docs/users/flai.md, docs/operators/index.md, design/system/flai-cli.md, design/issues/I-0213-flai-s-http-mcp-server-refuses-an-agent-in-a-container-with-403-invalid-host-header-since-the-mcp-sdk-s-localhost-protection-admits-only-a-loopback-host-on-a-loopback-listener.md, design/issues/summary.md]
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
      seconds: 313
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 1099
          output: 52554
          cache_read: 4323131
          cache_write: 394452
          cost: 11.6085
    - kind: orchestrator
      seconds: 228
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 85
          output: 1164
          cache_read: 14186680
          cache_write: 69510
          cost: 3.1303
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-10-28T12:52:34Z
  value: 50
  by: planner-S-0515
  at: 2026-10-28T13:03:21Z
forecast:
  duration: 45m
  delivery: 2026-10-29T13:14:00Z
  basis: "Its own forecast of 45m; 20th in the pull order with an in-progress limit of 5, behind S-0417, S-0539, S-0537, S-0534, S-0536, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0424, S-0425, S-0506, S-0508, S-0509, S-0510, S-0511 and S-0512."
  by: flai
  at: 2026-10-29T04:05:47Z
finalized:
  by: alex
  at: 2026-10-28T12:57:10Z
---
# S-0515 flai's HTTP MCP server refuses an agent in a container with 403 invalid Host header, since the MCP SDK's localhost protection admits only a loopback Host on a loopback listener

## Goal

This story remediates [I-0213](../../../design/issues/I-0213-flai-s-http-mcp-server-refuses-an-agent-in-a-container-with-403-invalid-host-header-since-the-mcp-sdk-s-localhost-protection-admits-only-a-loopback-host-on-a-loopback-listener.md), "flai's HTTP MCP server refuses an agent in a container with 403 invalid Host header, since the MCP SDK's localhost protection admits only a loopback Host on a loopback listener". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0213 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0213 is closed with `flai issue close I-0213 --reason` saying what fixed it

## Tasks
- T-2206 mcphttp builds both SDK handlers with the localhost protection off, a test reproduces I-0213's 403 on a foreign Host, and flai mcp http's help says the token is the gate
- T-2208 An ADR records that the bearer token and the Origin refusal are the HTTP MCP server's gates and the SDK's localhost Host check is off, refining ADR-0030
- T-2210 The user, operator, and design documents say the token is the HTTP MCP server's gate whatever Host a request names, and I-0213 is closed

## Notes

Cost of delay inputs set by flai from I-0213. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-28T12:22:29Z, 0 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0515 on 2026-10-28. The solution proposed, from I-0213's one instance: the bearer token and the `Origin` refusal are the HTTP MCP server's gates, and the go-sdk's localhost protection, which refuses a loopback listener a request whose `Host` is not loopback, is switched off on both SDK handlers (`DisableLocalhostProtection`). flai refuses an `Origin` and checks the token before the SDK sees a request, so the DNS-rebinding page the SDK's check guards against, which has no token, is refused already; a container that dials `host.docker.internal:4243` carries the token and is admitted whatever `Host` it names. No allowed-hosts setting: it would ask the operator to list the names each container runtime gives the host, for no gate the token does not already give. An ADR records the decision, refining ADR-0030. The fix stands on its own, whatever address the listener binds, so the story names no `after`: S-0506 moves the listener and the token, and either may land first; the two overlap on `mcphttp.go`, `mcphttp_test.go`, `mcp_http.go`, and the three documents, so the pull order settles which holds the other.

Touches. The story declared none, and `flai touches suggest S-0515` seeded with `flai/internal/mcphttp/mcphttp.go` and `flai/cmd/mcp_http.go`, the two files the instance names (17 of 2473 commits), lists `mcpserver/server.go`, `cmd/mcp.go`, `mcpserver/folder.go`, and their tests at 47 to 71%: left out, since the tools and the stdio server do not change. Kept from that list (co-change): `docs/users/flai.md`, `docs/operators/index.md`, `design/system/flai-cli.md`, and `docs/users/flai-reference.md`, each of which describes the server's gates, and the reference is regenerated by `make flai-reference` when `mcp_http.go`'s help changes, since `reference_test` fails when it is stale. From the layout: `flai/internal/mcphttp/mcphttp.go`, where `Handler` builds the two SDK handlers; `flai/internal/mcphttp/mcphttp_test.go`, whose `lab` listens on loopback, so a request with a foreign `Host` reproduces the 403 there; `flai/cmd/mcp_http.go`, whose `Long` help says who is refused; `design/adrs/README.md`, the index `flai adr new` writes; and the issue's file and `design/issues/summary.md`, which `flai issue close` rewrites. The folder `design/adrs` is kept: `flai adr new` names the ADR's file by its number and title, which no task can name, and `flai task done` widens the touches with the committed path (ADR-0144).

Forecast. `flai forecast S-0515` gives 5m (median 131 s per unit over 12 done remediation stories on claude-opus-5-5 in the small band; size 2: 2 criteria, 0 touches) and delivery 2026-10-29T08:23Z, 29th in the pull order with an in-progress limit of 5. Raised to 45m: the size counted no touches, and the story makes a Go change with a test that reproduces the 403, writes an ADR, changes three documents, regenerates the reference, and closes the issue. The delivery written, 2026-10-29T09:03Z, is flai's plus the 40m the raised duration adds.

Cost of delay. `flai cod S-0515` gives 50.00 USD a week from the operator's input flai set from I-0213: 20m of time lost per 168h cycle at 150 USD an hour, 1.00 cycles a week. Written as flai gives it; the input stands as set.
