---
id: S-0519
type: story
nature: remediation
title: permission_prompt over the HTTP MCP server outlives Claude Code's tool-call timeout, so an agent in a container sees its protected write time out while the thread to the operator stays open
status: done
owner: alex
created: 2026-10-28T12:53:05Z
updated: 2026-10-29T03:54:23Z
transitions:
  - to: ready
    at: 2026-10-29T00:38:32Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-29T00:38:38Z
    by: agent-S-0519
  - to: review
    at: 2026-10-29T02:53:51Z
    by: agent-S-0519
  - to: in-progress
    at: 2026-10-29T03:24:27Z
    by: flai
  - to: review
    at: 2026-10-29T03:37:41Z
    by: agent-S-0519
  - to: done
    at: 2026-10-29T03:54:23Z
    by: orchestrator
blocked:
  - from: 2026-10-29T00:57:08Z
    until: 2026-10-29T00:57:09Z
    reason: awaiting the operator on TH-0726
  - from: 2026-10-29T00:57:26Z
    until: 2026-10-29T00:57:31Z
    reason: awaiting the operator on TH-0726
  - from: 2026-10-29T00:57:37Z
    until: 2026-10-29T01:02:14Z
    reason: awaiting the operator on TH-0726
  - from: 2026-10-29T01:02:23Z
    until: 2026-10-29T01:25:53Z
    reason: awaiting the operator on TH-0726
  - from: 2026-10-29T01:26:04Z
    until: 2026-10-29T01:56:28Z
    reason: awaiting the operator on TH-0726
  - from: 2026-10-29T01:56:39Z
    until: 2026-10-29T01:59:13Z
    reason: awaiting the operator on TH-0726
tags: [cli, mcp]
topics: [mcp, harness]
touches: [design/adrs, design/adrs/README.md, flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, flai/internal/mcpserver/server.go, flai/cmd/mcp_http.go, flai/cmd/mcp_http_test.go, flai/internal/ask/ask.go, flai/internal/ask/ask_test.go, docs/users/flai.md, docs/operators/index.md, design/system/agent-adapters.md, design/system/flai-cli.md, design/conventions/delegation.md, design/issues/I-0217-permission-prompt-over-the-http-mcp-server-outlives-claude-code-s-tool-call-timeout-so-an-agent-in-a-container-sees-its-protected-write-time-out-while-the-thread-to-the-operator-stays-open.md, design/issues/summary.md, design/issues/I-0228-claude-code-drops-any-call-to-flai-s-http-mcp-server-that-answers-after-60-seconds-so-wait-for-events-wait-for-work-verify-and-a-slow-task-done-fail-for-an-agent-in-a-container.md, flai/internal/mcpserver/folder.go, design/issues/I-0229-two-stories-accepted-apart-leave-main-red-s-0535-reworded-to-backlog-s-condition-and-s-0416-s-testadefinedroleisheldtoitsfile-still-expects-the-old-words.md, flai/internal/guard/roles_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 5230
  priced_by: harness
  turns:
    - day: 2026-10-29
      ceremony: 6
      hand_edits: 3
      work: 71
  models:
    - model: claude-opus-5-5
      input: 364
      output: 116788
      cache_read: 25579039
      cache_write: 1041727
      cost: 13.6158
  strategic:
    - kind: orchestrator
      seconds: 6070
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 693
          output: 11550
          cache_read: 94563621
          cache_write: 219671
          cost: 20.8653
        - model: claude-sonnet-5-5
          input: 16
          output: 95
          cache_read: 163272
          cache_write: 84373
          cost: 0.1707
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
  delivery: 2026-10-29T01:49:00Z
  basis: "Its own forecast of 1h; 5th in the pull order with an in-progress limit of 5, behind S-0532, S-0537, S-0534 and S-0536."
  by: flai
  at: 2026-10-29T00:38:58Z
finalized:
  by: orchestrator
  at: 2026-10-28T19:17:28Z
---
# S-0519 permission_prompt over the HTTP MCP server outlives Claude Code's tool-call timeout, so an agent in a container sees its protected write time out while the thread to the operator stays open

## Goal

This story remediates [I-0217](../../../design/issues/I-0217-permission-prompt-over-the-http-mcp-server-outlives-claude-code-s-tool-call-timeout-so-an-agent-in-a-container-sees-its-protected-write-time-out-while-the-thread-to-the-operator-stays-open.md), "permission_prompt over the HTTP MCP server outlives Claude Code's tool-call timeout, so an agent in a container sees its protected write time out while the thread to the operator stays open". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0217 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0217 is closed with `flai issue close I-0217 --reason` saying what fixed it

## Tasks
- T-2227 Claude Code's tool-call timeout over flai's HTTP MCP server is measured on the host and in the container, with MCP_TOOL_TIMEOUT unset and set, and the remedy is proposed on a thread
- T-2228 An ADR refining ADR-0124 gives permission_prompt a bound per transport: four minutes over stdio and, over HTTP, the bound the measurement settles below Claude Code's tool-call timeout
- T-2229 The HTTP MCP server's permission_prompt holds a write at most the HTTP bound and answers before Claude Code's tool-call timeout, with a test that reproduces I-0217
- T-2230 docs/users/flai.md, docs/operators/index.md, agent-adapters.md, flai-cli.md, and delegation.md state permission_prompt's bound per transport and what a client receives during a held call over HTTP
- T-2231 I-0217 is closed with flai issue close saying what fixed it, and the story's criteria are ticked
- T-2323 TestADefinedRoleIsHeldToItsFile expects to-backlog's wording as S-0535 changed it, so main's integration tier passes again

## Notes

Cost of delay inputs set by flai from I-0217. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-28T12:30:56Z, 0 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0519 on 2026-10-28. The story declared no touches, so `flai touches suggest S-0519` refused without a seed; seeded with the three files the hold lives in, `flai/internal/mcpserver/permission.go` (`permissionWait`, one value for both transports), `flai/internal/ask/ask.go` (`Asker.Hold`, the refusal naming the thread), and `flai/cmd/mcp_http.go` (where the HTTP server builds `mcpserver.Options`), it ranks 24 commits.

Touches. From co-change: `flai/internal/mcpserver/server.go` (54 percent; `Options` gains the bound), `flai/internal/mcpserver/permission_test.go` (29), `design/system/flai-cli.md` (21; the `flai mcp` row states the bound), `docs/users/flai.md` (21; § Writes to protected paths and the tools table), `flai/cmd/mcp_http_test.go` (12; the test that reproduces I-0217), and `docs/operators/index.md` (8; the "Long requests" row of § MCP over HTTP). Left out of the co-changes: `flai/cmd/mcp.go` and `flai/internal/mcpserver/folder.go`, which keep the stdio default and the tool's registration; `docs/operators/settings.md` and `docs/users/flai-reference.md`, since the bound is a constant, not a setting, and no command changes; the hostapi, serve, and root files. From the design: `design/system/agent-adapters.md` (the bound table of § Asking the owner before a protected write and the spike's paragraph of § A container harness of flai's own), `design/conventions/delegation.md` (the project addition stating the hold), `design/adrs/README.md` (the ADR's index row), and the folder `design/adrs`, kept because `flai adr new` in T-2228 names the ADR's file by its number and the title the story's agent writes, so no task can name it; the folder is kept for that file alone, and `flai task done` widens the touches with the committed path (ADR-0144). From the layout: `flai/internal/ask/ask.go` and `ask_test.go` (a bound of zero refuses at once), and `design/issues/I-0217-….md` with `design/issues/summary.md` (the close). Not touched: `flai/internal/mcphttp/mcphttp.go`, since the remedy recommended changes the bound and not the JSON response mode; `scripts/agent-container.sh`, which S-0521 is changing, since setting `MCP_TOOL_TIMEOUT` there is the alternative T-2227 argues against; `flai/internal/harness/harness.go`, whose start prompt says "at most four minutes", which stays true of a shorter bound.

`after`: none. The fix needs nothing of S-0506 (permission_prompt under auto-approve over HTTP) or S-0515 (the Host header) on main; the overlap on `permission.go` and `mcp_http.go` is the board's to hold while either is in progress.

Forecast. `flai forecast S-0519` gives 29m (median 95 s per unit over 38 done remediation stories on claude-opus-5-5 in the large band; size 18: 2 criteria, 16 touches) and delivery 2026-10-29T02:15Z, 32nd in the pull order. Raised to 1h: T-2227 runs three or four probes that each wait up to the timeout, one of them in the container, which the band's remediation stories did not, and T-2228 writes an ADR; S-0309, the stdio remedy of the same hold, was forecast 35m without a measurement. The delivery written, 2026-10-29T02:46Z, is flai's play-out start (its delivery less its 29m) plus the hour. T-2227 ends on a thread to the operator, whose wait is not agent time.

Cost of delay. `flai cod S-0519` gives 25.00 USD a week: 10m of time lost per 168h cycle, the input flai set from I-0217, at 150 USD an hour and one cycle a week. Written as flai gives it: the issue occurred once, in a spike, and the figure will rise with each container agent that meets it until the fix lands.
- 2026-10-29T03:24:27Z: moved to in-progress: merge queue batch 20261029-031456 fails with S-0519 at 2ec7500ae00e: integration failed: flai/internal/guard/roles_test.go:139 TestADefinedRoleIsHeldToItsFile: {Session: Subagent:false Agent: Role:review Story: Kind:3 Name: Command: Path: Tool:item_move Args:{ID:S-0001 Story: To:review Type: Draft:false Commit:false Recommendation:false Source: Fields:[] CostOfDelay:map[]}}: {Outcome:1 Reason:the role reviewer cannot move S-0001 to review: it may make it only under to-backlog (on a call, a move to backlog and nowhere else; on a trigger, only a move to backlog), or to-ready (a move to ready and nowhere else). Ask the operator with thread_open on the item if it …

### Accepted by the orchestrator

- Verified: fa8edce27028ee4105946131cd349aa4aef201a7
- Batch: 20261029-033842
- Enqueued verified: fa8edce27028ee4105946131cd349aa4aef201a7
- At: 2026-10-29T03:54:23Z

Verdict: S-0519 meets both criteria at fa8edce; flai verify passed with integration and smoke left to the queue, every changed file is under the story's touches, and T-2323 fixes main's failing guard test (I-0229).
- 1: flai/internal/mcpserver/permission.go, flai/internal/mcpserver/server.go, flai/cmd/mcp_http.go, flai/cmd/mcp_http_test.go, flai/internal/mcpserver/folder.go, flai/internal/ask/ask.go, flai/internal/ask/ask_test.go, flai/internal/mcpserver/permission_test.go, design/adrs/0215-permission-prompt-holds-a-write-at-most-four-minutes-over-stdio-and-at-most-45.md, design/adrs/README.md, design/system/agent-adapters.md, design/system/flai-cli.md, design/conventions/delegation.md, docs/users/flai.md, docs/operators/index.md, flai/internal/guard/roles_test.go
- 2: design/issues/I-0217-permission-prompt-over-the-http-mcp-server-outlives-claude-code-s-tool-call-timeout-so-an-agent-in-a-container-sees-its-protected-write-time-out-while-the-thread-to-the-operator-stays-open.md, design/issues/summary.md
