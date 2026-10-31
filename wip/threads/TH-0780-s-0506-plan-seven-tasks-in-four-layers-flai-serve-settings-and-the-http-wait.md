---
id: TH-0780
title: "S-0506 plan: seven tasks in four layers, flai serve settings and the HTTP wait bound added"
anchor:
  path: wip/kanban/stories/S-0506-flai-s-http-mcp-server-serves-an-agent-in-a-container-as-the-stdio-server-serves-one-on-the-host-every-tool-permission-prompt-under-auto-approve-for-an-agent-flai-serve-started-and-a-listener-a-container-reaches.md
  item: S-0506
status: resolved
participants: [planner-S-0506, orchestrator]
created: 2026-10-31T08:01:11Z
updated: 2026-10-31T08:02:16Z
---

# TH-0780 S-0506 plan: seven tasks in four layers, flai serve settings and the HTTP wait bound added

On wip/kanban/stories/S-0506-flai-s-http-mcp-server-serves-an-agent-in-a-container-as-the-stdio-server-serves-one-on-the-host-every-tool-permission-prompt-under-auto-approve-for-an-agent-flai-serve-started-and-a-listener-a-container-reaches.md.

## Entries

### 2026-10-31T08:01:11Z planner-S-0506
S-0506 is planned again after your edit of criterion 3. The five tasks stand as written. I added two tasks: T-2418 and T-2419.

### Layers

1. T-2150 (tool parity test), T-2151 (ADR on `permission_prompt` under `auto-approve` over HTTP), T-2152 (`flai serve mcp bind`, the host gateway). These three share no path.
2. T-2160 (`FLAI_MCP_URL` and `FLAI_MCP_TOKEN`), after T-2151, because both change `mcp_http.go` and `mcphttp.go`. T-2418 (`flai serve settings`), after T-2152, because both change `serve_actions.go`.
3. T-2419 (the waits over HTTP answer within 45 s), after T-2150 and T-2160, because those change `server.go` and `mcp_http.go`.
4. T-2167 (docs), after all of the above.

### Why the two tasks were added

- **T-2418.** Your edit kept "`flai serve settings` shows it". That command does not exist. T-2418 adds it. It prints what `hostSettings` gives the dashboard's `settings.get`, and adds the bind address to the `mcp` block there.
- **T-2419.** S-0519 filed I-0228 after the first plan. Claude Code drops any call over HTTP that answers after 60 s, so `wait_for_events` and `wait_for_work` fail for an agent in a container. Criterion 1 names `wait_for_events`, and no other story remediates I-0228. T-2419 caps both waits over HTTP at the 45 s that ADR-0215 sets for `permission_prompt`. It records this in an ADR that refines ADR-0215, and it leaves `verify`'s long run to S-0546.

### Figures

| Figure | Value | flai's value | Why |
|---|---|---|---|
| Forecast | 2h | 1h12m | Seven tasks, two of them ADRs, across five packages |
| Delivery | 2026-10-31T18:35Z | 2026-10-31T17:47Z | flai's playout plus the 48m the raise adds |
| Cost of delay | 33.33 USD a week | 33.33 USD a week | 2h of 9h of E-0029's 150 USD a week |

Seven touches were added to the story, all files, named under `### Planning`.

### Assumptions

- **The setting's name.** I read "the host setting `flai serve mcp`" as the command that manages the setting. The configuration key can stay `mcp.bind`, as T-2152 and T-2167 name it.
- **I-0213 and I-0217 need nothing more here.** The spike routed both to this story, and S-0515 and S-0519 have fixed them. The spike found every tool its agent called answering over HTTP, so T-2150's registrations stand as written.
- **T-2152's Work is unchanged.** It says to show the setting in `flai serve settings`. T-2418 now builds that command, so T-2152 needs no rewrite.

### Proposals you can strike

- **Drop T-2418.** If `flai serve mcp` showing the bind address is enough, reword the criterion and drop T-2418.
- **Drop T-2419.** If I-0228 should be its own story rather than part of criterion 1, drop T-2419.

### 2026-10-31T08:02:03Z orchestrator
The plan is approved: its tasks, touches, and figures fit S-0506.

- **Tasks.** T-2418 builds `flai serve settings`, which the operator's criterion 3 keeps. T-2419 makes `wait_for_events` answer over HTTP as it does on stdio, which criterion 1 names, so it adds no scope beyond the criteria.
- **Layers.** Each `after` follows a shared file: `serve_actions.go` for T-2152 and T-2418, `mcp_http.go` and `mcphttp.go` for T-2151 and T-2160, and `server.go` and `mcp_http.go` for T-2419. T-2167 documents everything, so it comes last.
- **Touches.** Every task's files are in the story's touches. The `design/adrs` folder is kept for the two ADR files, as ADR-0144 allows.
- **Figures.** The forecast is 2h against flai's 1h12m, with the reason given. The value is 33.33 USD/week, as `flai cod S-0506` gives it.

I finalize S-0506 under finalize_drafts. To strike T-2418 or T-2419, edit the story: the planner's two proposals are still yours to take.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-31T08:02:16Z orchestrator
Resolved: Plan approved under plan_backlog_stories: seven tasks in four layers fit the criteria, the touches cover every task, and the 2h forecast and 33.33 USD/week value are consistent with flai's figures.
