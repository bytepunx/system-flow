---
id: ADR-0243
title: "Over HTTP, permission_prompt honours the auto-approve host action only for the agent flai serve runs for the story whose worktree the write is in, and asks every other agent"
status: accepted
date: 2026-11-02
supersedes: []
superseded_by: []
refines: [ADR-0086, ADR-0097]
topics: [mcp, cli]
---

# ADR-0243 Over HTTP, permission_prompt honours the auto-approve host action only for the agent flai serve runs for the story whose worktree the write is in, and asks every other agent

## Context

[ADR-0086](0086-flai-serve-gives-a-claude-code-agent-flai-s-permission-prompt-as-its-permission.md) has flai's `permission_prompt` allow a write to a path Claude Code protects in an in-progress story's worktree once the owner answers allow on a thread, or at once while the operator has the `auto-approve` host action on. [ADR-0097](0097-permission-prompt-takes-an-answer-from-the-story-s-owner-or-the-project-s-owner.md) says who may answer. Both were written for the server `flai mcp` runs over stdio, one process per agent session, which `flai serve` starts with the agent's name (`--agent`).

The server over HTTP (`flai mcp start`, [ADR-0030](0030-mcp-is-served-by-flai-on-the-host-over-stdio-and-http-and-the-dashboard-s-api.md)) serves an agent that cannot start a process on the host, such as one in a container. It never read `auto-approve`: `flai/cmd/mcp_http.go` left `Options.AutoApprove` unset, so `permission_prompt` asked every agent there on a thread, whatever the operator had turned on. An agent `flai serve` started in a container therefore waited on the operator for a write the same agent on the host would have made at once. S-0506 makes the HTTP server serve such an agent as the stdio server serves one on the host.

The two servers differ in who the caller is. Over stdio the server's process is the agent's own, started by `flai serve` or by the operator. Over HTTP one process serves every holder of its bearer token ([ADR-0216](0216-the-http-mcp-server-s-gates-are-the-bearer-token-and-the-origin-refusal.md)), and the caller names itself with the `X-Flai-Agent` header, so a name there is only a claim. Granting `auto-approve` to every caller over HTTP would grant it to whatever holds the token and names itself anything.

`flai serve` records every run it starts in `agents.json` in its folder beside the operator's configuration: for a story's agent, the story, the name it gave the agent (`agent-S-nnnn`, or the name of the run it resumes), and the process it started with when that process started. That record is the host's own, written by the process that started the agent; no caller over HTTP writes it.

## Decision

**Over HTTP, `permission_prompt` honours the `auto-approve` host action only for the agent `flai serve` runs for the story whose worktree the write is in, and asks every other agent on a thread.** This refines ADR-0086's decision 3 for the HTTP server. Who may answer the thread stays as ADR-0097 says, and the stdio server is unchanged.

1. The HTTP server reads `auto-approve` for the project at every request, as the stdio server does, from the configuration the server was started with.
2. With it on, a write that ADR-0086's decision 1 would approve is allowed at once, with no thread, only when `flai serve`'s record names a run for the story whose worktree holds the path, the run's agent has the name the caller gives in `X-Flai-Agent`, the run has not ended and was started, and the process it started is alive and is the one it started.
3. Any other caller over HTTP, an agent under another name, the agent of another story, or one whose run has ended, is asked on a thread on the story as if `auto-approve` were off, within ADR-0215's 45-second bound. The server logs that it asked, with the agent's name and the story, and never the token.
4. A record the server cannot read names no run, so the server asks.
5. The dashboard token ([ADR-0018](0018-dashboard-token.md)) grants nothing here. It is a different token, which the MCP server never accepts, and `auto-approve` stays the operator's shell setting that no dashboard sees or changes ([ADR-0067](0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md), ADR-0086).

`mcpserver` takes the check as an option, `Options.ServedAgent`, consulted only when `Options.AutoApprove` says yes. `flai/cmd` fills it for the HTTP server from `flai serve`'s record, so `mcpserver` does not import `flai serve`'s package. The stdio server leaves it unset.

## Consequences

- An agent `flai serve` started in a container, under `auto-approve`, writes its `.claude/` files and `.mcp.json` at once, as it would on the host. Its sub-agents send the same `X-Flai-Agent` from the same MCP configuration, so they are allowed too, as on the host; with `auto-approve` off, `flai guard` still refuses a sub-agent's write first ([ADR-0102](0102-while-auto-approve-is-off-flai-guard-refuses-a-story-s-sub-agent-a-write-to-a.md)).
- An agent of the operator's own, or any agent `flai serve` did not start, that reaches the HTTP server is still asked on a thread, even under `auto-approve`. Over stdio such an agent is not asked, since the stdio server cannot tell who started it and its process is the caller's own.
- The token remains the gate. A token holder that gives the name of a run `flai serve` has under way, for that run's story, is allowed as that agent is, while it runs. Holding the token already lets a caller do anything else an agent may, so `flai mcp token --rotate` remains the answer to a leaked token.
- Every allowed write is logged, as before, with the agent's name; each refusal to auto-approve is logged too, so the operator can see why a thread was opened under `auto-approve`.
- `docs/users/flai.md` and `design/system/flai-cli.md` said the shared HTTP server always asks. They now say it asks every agent but the one `flai serve` runs for the story.

## Alternatives considered

- **Keep asking every agent over HTTP.** It is the safest reading of a shared server, but an agent `flai serve` started in a container then waits up to 45 seconds and a retry on a write the operator already allowed by turning `auto-approve` on, and S-0506's aim, the same service in a container as on the host, fails.
- **Honour `auto-approve` for every caller over HTTP.** It trusts the `X-Flai-Agent` header, which any token holder sets, so it would widen the grant from the agents the operator's `flai serve` started to whatever holds the token.
- **A per-run token that `flai serve` hands each agent it starts.** It would make the name more than a claim, but it adds a token store, its rotation, and a second gate to the server for a grant the run record already bounds; it can follow if the shared token is ever handed to callers the operator does not run.
- **Honour it for any running agent `flai serve` started, whatever the story.** One story's agent could then write protected files in another story's worktree unasked; tying the name to the story whose worktree holds the path keeps each agent to its own.
