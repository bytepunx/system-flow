---
id: S-0519
type: story
nature: remediation
title: permission_prompt over the HTTP MCP server outlives Claude Code's tool-call timeout, so an agent in a container sees its protected write time out while the thread to the operator stays open
status: backlog
owner: alex
created: 2026-10-28T12:53:05Z
updated: 2026-10-28T20:51:46Z
transitions: []
tags: [cli, mcp]
topics: [mcp, harness]
touches: [design/adrs, design/adrs/README.md, flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, flai/internal/mcpserver/server.go, flai/cmd/mcp_http.go, flai/cmd/mcp_http_test.go, flai/internal/ask/ask.go, flai/internal/ask/ask_test.go, docs/users/flai.md, docs/operators/index.md, design/system/agent-adapters.md, design/system/flai-cli.md, design/conventions/delegation.md, design/issues/I-0217-permission-prompt-over-the-http-mcp-server-outlives-claude-code-s-tool-call-timeout-so-an-agent-in-a-container-sees-its-protected-write-time-out-while-the-thread-to-the-operator-stays-open.md, design/issues/summary.md]
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
      seconds: 759
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 102
          output: 1517
          cache_read: 27145069
          cache_write: 41203
          cost: 5.9686
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-28T12:53:05Z
  value: 25
  by: planner-S-0519
  at: 2026-10-28T13:46:59Z
forecast:
  duration: 1h
  delivery: 2026-10-29T09:12:00Z
  basis: "Its own forecast of 1h; 31st in the pull order with an in-progress limit of 5, behind S-0416, S-0530, S-0419, S-0423, S-0418, S-0421, S-0422, S-0420, S-0417, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0424, S-0425, S-0506, S-0507, S-0508, S-0509, S-0510, S-0511, S-0512, S-0514, S-0515, S-0516, S-0517 and S-0518."
  by: flai
  at: 2026-10-28T20:51:46Z
finalized:
  by: orchestrator
  at: 2026-10-28T19:17:28Z
---
# S-0519 permission_prompt over the HTTP MCP server outlives Claude Code's tool-call timeout, so an agent in a container sees its protected write time out while the thread to the operator stays open

## Goal

This story remediates [I-0217](../../../design/issues/I-0217-permission-prompt-over-the-http-mcp-server-outlives-claude-code-s-tool-call-timeout-so-an-agent-in-a-container-sees-its-protected-write-time-out-while-the-thread-to-the-operator-stays-open.md), "permission_prompt over the HTTP MCP server outlives Claude Code's tool-call timeout, so an agent in a container sees its protected write time out while the thread to the operator stays open". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0217 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0217 is closed with `flai issue close I-0217 --reason` saying what fixed it

## Tasks
- T-2227 Claude Code's tool-call timeout over flai's HTTP MCP server is measured on the host and in the container, with MCP_TOOL_TIMEOUT unset and set, and the remedy is proposed on a thread
- T-2228 An ADR refining ADR-0124 gives permission_prompt a bound per transport: four minutes over stdio and, over HTTP, the bound the measurement settles below Claude Code's tool-call timeout
- T-2229 The HTTP MCP server's permission_prompt holds a write at most the HTTP bound and answers before Claude Code's tool-call timeout, with a test that reproduces I-0217
- T-2230 docs/users/flai.md, docs/operators/index.md, agent-adapters.md, flai-cli.md, and delegation.md state permission_prompt's bound per transport and what a client receives during a held call over HTTP
- T-2231 I-0217 is closed with flai issue close saying what fixed it, and the story's criteria are ticked

## Notes

Cost of delay inputs set by flai from I-0217. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-28T12:30:56Z, 0 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0519 on 2026-10-28. The story declared no touches, so `flai touches suggest S-0519` refused without a seed; seeded with the three files the hold lives in, `flai/internal/mcpserver/permission.go` (`permissionWait`, one value for both transports), `flai/internal/ask/ask.go` (`Asker.Hold`, the refusal naming the thread), and `flai/cmd/mcp_http.go` (where the HTTP server builds `mcpserver.Options`), it ranks 24 commits.

Touches. From co-change: `flai/internal/mcpserver/server.go` (54 percent; `Options` gains the bound), `flai/internal/mcpserver/permission_test.go` (29), `design/system/flai-cli.md` (21; the `flai mcp` row states the bound), `docs/users/flai.md` (21; § Writes to protected paths and the tools table), `flai/cmd/mcp_http_test.go` (12; the test that reproduces I-0217), and `docs/operators/index.md` (8; the "Long requests" row of § MCP over HTTP). Left out of the co-changes: `flai/cmd/mcp.go` and `flai/internal/mcpserver/folder.go`, which keep the stdio default and the tool's registration; `docs/operators/settings.md` and `docs/users/flai-reference.md`, since the bound is a constant, not a setting, and no command changes; the hostapi, serve, and root files. From the design: `design/system/agent-adapters.md` (the bound table of § Asking the owner before a protected write and the spike's paragraph of § A container harness of flai's own), `design/conventions/delegation.md` (the project addition stating the hold), `design/adrs/README.md` (the ADR's index row), and the folder `design/adrs`, kept because `flai adr new` in T-2228 names the ADR's file by its number and the title the story's agent writes, so no task can name it; the folder is kept for that file alone, and `flai task done` widens the touches with the committed path (ADR-0144). From the layout: `flai/internal/ask/ask.go` and `ask_test.go` (a bound of zero refuses at once), and `design/issues/I-0217-….md` with `design/issues/summary.md` (the close). Not touched: `flai/internal/mcphttp/mcphttp.go`, since the remedy recommended changes the bound and not the JSON response mode; `scripts/agent-container.sh`, which S-0521 is changing, since setting `MCP_TOOL_TIMEOUT` there is the alternative T-2227 argues against; `flai/internal/harness/harness.go`, whose start prompt says "at most four minutes", which stays true of a shorter bound.

`after`: none. The fix needs nothing of S-0506 (permission_prompt under auto-approve over HTTP) or S-0515 (the Host header) on main; the overlap on `permission.go` and `mcp_http.go` is the board's to hold while either is in progress.

Forecast. `flai forecast S-0519` gives 29m (median 95 s per unit over 38 done remediation stories on claude-opus-5-5 in the large band; size 18: 2 criteria, 16 touches) and delivery 2026-10-29T02:15Z, 32nd in the pull order. Raised to 1h: T-2227 runs three or four probes that each wait up to the timeout, one of them in the container, which the band's remediation stories did not, and T-2228 writes an ADR; S-0309, the stdio remedy of the same hold, was forecast 35m without a measurement. The delivery written, 2026-10-29T02:46Z, is flai's play-out start (its delivery less its 29m) plus the hour. T-2227 ends on a thread to the operator, whose wait is not agent time.

Cost of delay. `flai cod S-0519` gives 25.00 USD a week: 10m of time lost per 168h cycle, the input flai set from I-0217, at 150 USD an hour and one cycle a week. Written as flai gives it: the issue occurred once, in a spike, and the figure will rise with each container agent that meets it until the fix lands.
