---
id: ADR-0130
title: "Every harness meets flai through neutral contracts, an adapter's capabilities, a usage reader, the guard's call and verdict with hold-and-ask, and a protected list per harness, and Claude Code's hooks and permission tool are one adapter's form of them"
status: proposed
date: 2026-10-08
supersedes: []
superseded_by: []
refines: [ADR-0060, ADR-0065, ADR-0086, ADR-0106, ADR-0124]
---

# ADR-0130 Every harness meets flai through neutral contracts, an adapter's capabilities, a usage reader, the guard's call and verdict with hold-and-ask, and a protected list per harness, and Claude Code's hooks and permission tool are one adapter's form of them

## Context

Everything under `flai serve` that reads an agent's output, guards its calls, or answers its permission prompts reads Claude Code's formats: the `PreToolUse`, `SubagentStart`, and `SubagentStop` hooks and their JSON ([ADR-0060](0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md)), `--permission-prompt-tool` with its `behavior` and `updatedInput` shape and its four-minute bound ([ADR-0086](0086-flai-serve-gives-a-claude-code-agent-flai-s-permission-prompt-as-its-permission.md), [ADR-0124](0124-permission-prompt-holds-a-write-at-most-four-minutes-then-refuses-it-and-leaves.md)), `--agents` JSON for a role's model ([ADR-0065](0065-a-story-s-agent-carries-a-model-per-sub-agent-role-and-claude-code-runs-each.md)), Claude Code's list of protected paths ([ADR-0106](0106-a-story-whose-branch-changes-a-path-claude-code-protects-is-accepted-by-the.md)), and stream-json for usage. The finding of S-0339, [agent-adapters.md](../system/agent-adapters.md), shows that Codex CLI, OpenCode, and Goose each offer a hook, a permission mode, sub-agent definitions, and an event stream of their own shape, and that the policy under each of flai's mechanisms is the same for all of them.

## Decision

Every harness meets flai through neutral contracts, an adapter's capabilities, a usage reader, the guard's call and verdict with hold-and-ask, and a protected list per harness, and Claude Code's hooks and permission tool are one adapter's form of them.

- The `Adapter` interface gains `Capabilities()`, stating whether the harness can resume a session, run the guard, answer a permission question, give each role its own model, run a strategic agent from a definition, and report cost, so that `flai serve` refuses or degrades before a start rather than after a failure; and `Reader()`, which turns the harness's log into neutral events: session start, a model call with its tokens, model, cost when reported, and the call that started its sub-agent, a tool result, and the run's end. `flai/internal/usage` measures from those events, not from stream-json.
- `flai guard` keeps one decision, `Decide(Call) Verdict`, where a `Call` says who calls (the session, a sub-agent, a strategic role), what (a shell line, a file write, a flai tool), and with which arguments, and gains one input reader per harness that maps the harness's hook input, tool names, and sub-agent sign onto it. Each harness registers the guard its own way: `.claude/settings.json` today, `.codex/hooks.json`, an OpenCode plugin, a Goose `hooks.json`.
- Asking the owner before a protected write is the guard's hold-and-ask: when the policy says ask, the guard holds the call, puts the question on a thread as `permission_prompt` does, and answers allow or deny within the harness's bound. Claude Code's `--permission-prompt-tool` is that contract in Claude Code's form; ADR-0124's four minutes is Claude Code's bound, and each harness has its own.
- A role's model is laid over the role's definition in each harness's own format (`--agents` JSON, a Codex TOML definition, an OpenCode `agent` entry); `.claude/agents/` stays the source the adapters translate from, with a tool-name table per harness.
- The protected list in `flai/internal/protected` is flai's policy: Claude Code's paths today, plus each adopted harness's own files that change what its agent may do, such as `.codex/`, `.opencode/`, `.agents/plugins/`, and `AGENTS.md`.

## Consequences

- Adding a harness is an adapter: a `Start`, a `Capabilities`, a `Reader`, a guard input reader, a registration file in the template, and its protected paths. Nothing in the policies, the measurement, or the threads changes.
- The contracts are built inside the Claude Code adapter first, where one harness exercises them, so they cost little until a second adapter arrives.
- `design/system/flai-cli.md`, `design/system/metrics.md` where the events are named, the conventions that name Claude Code's tools, and `docs/users/flai.md` change with the stories that build each piece.

## Alternatives considered

- A per-harness implementation of each mechanism: three guards, three permission handlers, three readers, each with its own policy copy that drifts.
- Refusing every harness but Claude Code: keeps today's code, forgoes other providers' models and self-hosted gateways, which S-0339 was asked to reach.
- A loop of flai's own as the one neutral harness: every contract becomes a function call, but it rebuilds what every harness already gives; kept as a later option.
