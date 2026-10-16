---
id: E-0023
type: epic
nature: feature
title: An OpenCode adapter runs a story's agent on OpenCode with the guard, the permission prompt, the roles, and usage measured, through its plugin hooks and its server's events, so a non-Claude model works a story end to end
status: backlog
owner: alex
created: 2026-10-16T04:28:00Z
updated: 2026-10-16T05:20:15Z
transitions: []
tags: [cli, template]
topics: [cli, agents, template]
touches: [flai/internal/harness, flai/internal/guard, flai/internal/usage, flai/internal/mcpserver, flai/internal/serve, flai/internal/manifest, design/system/agent-adapters.md, design/system/flai-cli.md, docs/operators/index.md, docs/users/flai.md, template]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 1022
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 150
          output: 2868
          cache_read: 62277197
          cache_write: 26387
          cost: 15.6189
---
# E-0023 An OpenCode adapter runs a story's agent on OpenCode with the guard, the permission prompt, the roles, and usage measured, through its plugin hooks and its server's events, so a non-Claude model works a story end to end

## Outcome

flai drives one harness fully, Claude Code, and runs anything else through the `command` escape hatch, which declares no capability, so a story on it has no guard, no permission prompt, no roles, and no measured usage ([ADR-0130](../../../design/adrs/0130-every-harness-meets-flai-through-neutral-contracts-an-adapter-s-capabilities-a.md), [ADR-0131](../../../design/adrs/0131-a-harness-without-a-guard-hook-runs-a-story-s-agent-but-none-of-its-roles.md)). E-0019 built the neutral contracts inside the Claude Code adapter and left Codex CLI, OpenCode, and Goose out on the operator's word (TH-0368). Gas Town ships eleven runtime presets, Paperclip about ten, Vibe Kanban ten. The finding in [agent-adapters.md § OpenCode](../../../design/system/agent-adapters.md#opencode) says OpenCode has what the contracts need: `opencode run` headless with `--session`, `--agent`, `--model provider/model`, and `--format json`; plugins whose `tool.execute.before` hook can refuse a call and whose `permission.asked` event can be answered; agents as markdown files with their own model and permissions; a headless server with an OpenAPI spec and an event stream; providers that reach OpenRouter and any OpenAI-compatible gateway. This epic builds the second real adapter on it, proves the contracts against a second harness, and gets a non-Claude model working a story end to end.

- [ ] A story's `agent.harness` may be `opencode`. The adapter's default host runs `opencode run` headless in the story's worktree with the harness prompt, the story's model as `provider/model`, `--format json`, and resumes a session with `--session`; `flai serve agent harness opencode` sets the program and arguments as it does for `claude-code`.
- [ ] flai's MCP server reaches the agent as a local MCP entry in a project-local `opencode.json` that flai writes or merges at start, with `flai mcp --agent <name>`, and the instruction files are read: OpenCode loads `AGENTS.md` and falls back to `CLAUDE.md`; the template settles which it ships, and `flai check` warns when neither is there.
- [ ] The guard runs: a plugin flai ships into `.opencode/plugins/` calls `flai guard` on `tool.execute.before` with a neutral `Call` and throws on refusal, and the adapter's `Capabilities` say `Guard: true` only once the first story has checked whether the hook can tell a sub-agent's call from the primary's; if it cannot, the roles stay off under ADR-0131 and `agent-adapters.md` records why.
- [ ] The permission prompt holds and asks: `permission.asked` reaches flai's hold-and-ask (`permission_prompt`, [ADR-0124](../../../design/adrs/0124-permission-prompt-holds-a-write-at-most-four-minutes-then-refuses-it-and-leaves.md)) and the plugin replies with the owner's answer, or the adapter runs `--auto` behind deny rules for the protected list; the first story settles which and records it.
- [ ] The roles run as OpenCode agents: `explore`, `verify`, and the planner's sub-agents are generated into `.opencode/agents/` from the definitions in `.claude/agents/`, which stay the source (ADR-0130), each with the model the story's `agent.roles` names and read-only permissions where the definition is read-only.
- [ ] Usage is measured: the adapter's `Reader` turns OpenCode's JSON events (`step_finish` tokens and cost) into neutral usage events, a story on OpenCode without a provider is `priced_by: harness`, and one with a provider is priced from the gateway's spend log as a Claude Code story is ([ADR-0132](../../../design/adrs/0132-a-story-s-cost-comes-from-its-provider-s-spend-log-when-it-has-a-provider-from.md)), with the provider entry derived from the host's `providers` map ([ADR-0129](../../../design/adrs/0129-a-story-s-agent-names-a-provider-beside-its-harness-and-model-and-the-providers.md)).
- [ ] The agent's status reaches the dashboard: the adapter's stream for the story's card and page comes from `--format json`, or from `opencode serve` and its event stream when the first story finds that better; the choice is an ADR.
- [ ] A non-Claude model works one story of this repository end to end on OpenCode through OpenRouter, with the guard and the permission prompt on, and the story's record names the model, what broke, and what it cost.
- [ ] OpenCode's integration and smoke tests run only where an `opencode` binary of a known version is installed and are skipped, saying so, where it is not; `docs/operators/index.md` § Harnesses, `docs/users/flai.md`, `design/system/agent-adapters.md`, and `design/system/flai-cli.md` describe the adapter and what a story on it has and lacks, and the template carries the plugin and the `opencode.json` skeleton.

## Stories

To be drafted by the planner. Order: the adapter and the headless run first, with the MCP entry, so a story runs with no guard and the findings to check are settled; then the guard plugin and the capabilities it proves; then the permission prompt; then the roles; then usage and pricing; then the status stream and its ADR; then the end-to-end story on a non-Claude model; documentation and the template with each story.
- S-0406 A story's agent may run on OpenCode headless with flai's MCP server, no guard, and no roles, and the findings the OpenCode section leaves to check are settled against a real opencode binary
- S-0407 flai guard judges an OpenCode agent's calls through a plugin on tool.execute.before, and the opencode adapter states the guard only where the plugin can tell a sub-agent's call
- S-0408 An OpenCode agent's protected write is asked of the story's owner through flai's hold-and-ask from the plugin's permission.asked, or refused by deny rules under --auto, as an ADR settles
- S-0409 The explorer, the verifier, and the strategic agents run on OpenCode as agents flai generates from .claude/agents, each on the model the story's roles name and read-only where the definition is
- S-0410 A story on OpenCode has its usage read from OpenCode's JSON events, priced by the harness or its gateway's spend log, and reaches its provider through an entry derived from the host's
- S-0411 A story's card and page show its OpenCode agent's status from a stream the adapter reads, from --format json or opencode serve's events as an ADR decides
- S-0412 A non-Claude model works one story of this repository end to end on OpenCode through OpenRouter, with the guard and the permission prompt on, and the results record what broke and what it cost

## Notes

- From the comparison of system-flow with OpenCode, Gas Town, Paperclip, and Vibe Kanban on 2026-10-16: OpenCode's hook surface is closer to Claude Code's than Codex's is, which makes it the cheapest second harness to test ADR-0130's contracts against.
- The planner reads first: `design/system/agent-adapters.md` § Today, § OpenCode, § Abstractions, and § Decision; ADR-0129 to ADR-0132; `flai/internal/harness/adapters.go` and `harness.go`; `flai/internal/guard/guard.go`; `flai/internal/mcpserver/permission.go`; `flai/internal/usage/log.go`.
- What the finding leaves to check, and the first story answers: whether a plugin can tell a sub-agent's call from the primary's; whether an `ask` in headless mode can be answered from outside; what fields `step_finish` carries; whether `--attach` to a running `opencode serve` is the better way to run many stories on one host.
- Spend: an OpenRouter key with a few dollars of credit for the end-to-end story; cap it with a limit on the key. OpenCode is MIT-licensed; its binary is installed by the operator on the host, not by flai.
- Later harnesses, Gemini CLI, Copilot CLI, Codex CLI, Goose, are not in this epic; each follows the same contracts once this one proves them.
