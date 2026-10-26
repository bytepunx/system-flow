---
id: TH-0664
title: "Plan for S-0506: five tasks in three layers, revisited"
anchor:
  path: wip/kanban/stories/S-0506-flai-s-http-mcp-server-serves-an-agent-in-a-container-as-the-stdio-server-serves-one-on-the-host-every-tool-permission-prompt-under-auto-approve-for-an-agent-flai-serve-started-and-a-listener-a-container-reaches.md
  item: S-0506
status: open
participants: [planner-S-0506, orchestrator]
created: 2026-10-26T12:31:43Z
updated: 2026-10-26T12:31:57Z
---

# TH-0664 Plan for S-0506: five tasks in three layers, revisited

On wip/kanban/stories/S-0506-flai-s-http-mcp-server-serves-an-agent-in-a-container-as-the-stdio-server-serves-one-on-the-host-every-tool-permission-prompt-under-auto-approve-for-an-agent-flai-serve-started-and-a-listener-a-container-reaches.md.

## Entries

### 2026-10-26T12:31:43Z planner-S-0506
S-0506 is planned. planner-E-0029's five tasks stand as written; no task is split, merged, or dropped, and none is missing against the criteria.

Layers:

1. T-2150 (tool parity test, missing tools registered over HTTP), T-2151 (ADR on `permission_prompt` under `auto-approve` over HTTP, server follows it), T-2152 (`mcp.bind`, `flai serve mcp bind`, the host gateway reachable). No path in common: T-2150 changes `mcpserver/server.go`, `folder.go`, and test files; T-2151 `cmd/mcp_http.go`, `mcpserver/permission.go`, `mcphttp/mcphttp.go`; T-2152 `cmd/host.go`, `cmd/serve_mcp.go`, `config.go`, `serve_actions.go`.
2. T-2160 (`FLAI_MCP_URL` and `FLAI_MCP_TOKEN` by environment, log names the agent, never the token), after T-2151, which changes `mcp_http.go` and `mcphttp.go` too.
3. T-2167 (docs), after T-2150, T-2152, and T-2160, whose behaviour it documents.

What changed:

- T-2152 and the story gain `flai/cmd/serve_mcp.go` and `serve_mcp_test.go`: `mcpAddrFor` there takes the host of every project's address from `mcpDefaultAddr` (`127.0.0.1:4243` in `mcp_http.go`), and the remembered address file keeps the host too, so the bind is applied there and a remembered loopback address is re-hosted when the setting changes. `mcp_http.go` is left to T-2151 and T-2160 so that layer 1 stays parallel.
- Forecast kept at 1h30m (flai gives 54m), delivery 2026-10-28T00:57Z; cost of delay 16.67 USD a week from E-0029's 150 USD `penalty_per_week`. Reasons under `### Planning`.

Assumptions:

- `flai serve settings`, which criterion 3 names as what shows `mcp.bind`, is no command today: `flai serve` has `actions`, `agent`, `harness`, `provider`, `forge`, and `components`, each showing its own settings. I read the criterion as "the `flai serve` command that shows the setting": T-2152 adds `flai serve mcp` showing the bind address with `bind <address>` as its subcommand, in the pattern of `flai serve agent`. If you want a `flai serve settings` command that shows every host setting, say so here and the story's agent adds it under T-2152; the criterion's wording is yours to amend when you finalize the draft.
- `settings_doc_test` (`flai/cmd/settings_doc_test.go`) fails once T-2152 adds the key and T-2160 the two variables until T-2167 gives them rows, so `flai test flai/cmd` on those tasks will report it; the story's agent closes T-2167 before the close-out, as the layers say.
- The story waits on S-0505 (`after`) because the spike's issues name the tools missing over HTTP and the network mode that reached the listener on this WSL2 host; the three known gaps in the goal stand whatever it finds, so T-2150, T-2151, and T-2152 need no change from the spike beyond the list of tools T-2150 registers.
- T-2150's parity test compares registered tool names, not behaviour; a tool that differs by transport for a reason records it in the test, as its Work says.

### 2026-10-26T12:31:57Z orchestrator (recommendation)
**Recommended for the operator:** approve the plan. Amend criterion 3 to name `flai serve mcp`, the command T-2152 builds, in place of `flai serve settings`, then finalize S-0506.

**The plan fits.**
- **Tasks:** five tasks in three layers. Layer 1 (T-2150, T-2151, T-2152) shares no path, and keeping `mcp_http.go` out of T-2152 keeps it parallel.
- **Ordering:** T-2160 waits on T-2151, which changes the same two files. T-2167, the docs, comes last, and it is also what makes `settings_doc_test` pass.
- **Figures:** the forecast is 1h30m, against flai's 54m. The value is 16.67 USD/week, E-0029's 150 shared by duration. `after: S-0505` is right, since the spike names the tools missing over HTTP.

**Why I do not finalize it or answer myself.** Criterion 3 says `flai serve settings` shows `mcp.bind`, but no such command exists and the plan does not add one. Finalized as written, the criterion could not be ticked against the diff.

There are two ways to fix it, and choosing between them is your call:
- **Reword the criterion (recommended).** `flai serve mcp` follows the per-area pattern of `flai serve agent`, `forge` and `provider`, and adds no scope.
- **Build `flai serve settings`.** A single command listing every host setting adds scope to the story.

The work-management rules keep an agent from inventing scope, and a finalized story's words are yours.

Source: design/conventions/work-management.md § Rules
