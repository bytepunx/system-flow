---
id: ADR-0131
title: "A harness without a guard hook runs a story's agent but none of its roles unless the host sets guard none for it, and a protected write on it is refused unless auto-approve is on"
status: accepted
date: 2026-10-08
supersedes: []
superseded_by: []
refines: [ADR-0060, ADR-0086, ADR-0102]
---

# ADR-0131 A harness without a guard hook runs a story's agent but none of its roles unless the host sets guard none for it, and a protected write on it is refused unless auto-approve is on

## Context

`flai guard` refuses a sub-agent's call that would change a work item, a thread, or the repository's history ([ADR-0060](0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md)) and, while `auto-approve` is off, its write to a file under `.claude/` ([ADR-0102](0102-while-auto-approve-is-off-flai-guard-refuses-a-story-s-sub-agent-a-write-to-a.md)); `permission_prompt` asks the owner before a protected write ([ADR-0086](0086-flai-serve-gives-a-claude-code-agent-flai-s-permission-prompt-as-its-permission.md)). Both run only because Claude Code calls them before each tool. The finding of S-0339, [agent-adapters.md](../system/agent-adapters.md), shows harnesses whose hook cannot tell a sub-agent's call from the agent's own (Goose), whose hook needs trust the operator grants (Codex CLI), or that have no hook at all (aider): on them flai cannot tell whose write it is, and a protected write is either denied outright by the harness's own permission map or allowed outright under its auto mode.

## Decision

A harness without a guard hook runs a story's agent but none of its roles unless the host sets guard none for it, and a protected write on it is refused unless auto-approve is on.

- An adapter's capabilities ([ADR-0130](0130-every-harness-meets-flai-through-neutral-contracts-an-adapter-s-capabilities-a.md)) say whether the harness runs the guard and whether it can hold a call to ask. `flai serve` reads them before a start.
- On a harness that cannot run the guard, or whose guard cannot tell a sub-agent's call, the story's agent may run, since its own calls are judged by the harness's permissions as ADR-0038 has it, but a story whose `roles` name a sub-agent is refused at start, as a role on another harness is refused today, with the reason and the setting that lifts it.
- The operator lifts it per harness on the host, `agent.harnesses.<name>.guard: none`, taking the risk that a sub-agent writes what only the story's agent should; the refusal names the setting.
- On a harness that cannot hold a call to ask, a protected write is refused: `flai serve` refuses to start a story's agent there unless the project's `auto-approve` host action is on, or the host's arguments for that harness deny the protected paths in the harness's own vocabulary, and the refusal says which.

## Consequences

- No harness is ruled out: Codex CLI, OpenCode, and Goose can each run a story's agent before their hooks are adopted, with roles and protected writes held back, and gain them as their adapters add a guard reader and hold-and-ask.
- The dashboard's story page and `flai serve`'s start log say why a story was refused and what lifts it.
- `docs/operators/settings.md` gains `agent.harnesses.<name>.guard`, and `docs/users/flai.md` says what a harness without a guard loses.

## Alternatives considered

- Refusing such a harness outright: rules out every harness until it has a hook flai reads, which is most of them today, and makes the first adapter of each a large story rather than a small one.
- Letting roles run unguarded by default: a sub-agent could move items, write threads, or commit, which ADR-0060 exists to prevent; the risk is the operator's to take, not the default.
- Enforcing the guard in `flai mcp` instead: a sub-agent shares its parent's MCP connection in every harness, and shell and file writes never reach flai, so the server cannot tell the callers apart.
