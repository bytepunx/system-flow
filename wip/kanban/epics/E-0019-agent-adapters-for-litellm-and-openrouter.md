---
id: E-0019
type: epic
nature: feature
title: Agent adapters for LiteLLM and OpenRouter
status: ready
owner: alex
created: 2026-10-08T08:33:28Z
updated: 2026-10-08T08:59:39Z
transitions:
  - to: ready
    at: 2026-10-08T08:38:12Z
    by: alex
tags: [cli, template]
topics: [cli, agents]
touches: [flai/internal/harness, flai/internal/serve, flai/internal/guard, flai/internal/usage, flai/internal/mcpserver/permission.go, flai/internal/protected/protected.go, flai/internal/manifest/agent.go, flai/internal/hostapi/settings.go, design/system/metrics.md, design/system/flai-cli.md, docs/operators/settings.md, docs/users/flai.md]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 278
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 64
          output: 1111
          cache_read: 13917152
          cache_write: 12708
          cost: 3.4322
cost_of_delay:
  inputs:
    penalty_per_week: 500
    by: alex
    at: 2026-10-08T08:45:23Z
  value: 500
  by: planner-E-0019
  at: 2026-10-08T08:56:09Z
---
# E-0019 Agent adapters for LiteLLM and OpenRouter

## Outcome

A story's agent runs on Claude Code over a gateway, OpenRouter first and then a LiteLLM proxy, with the model's provider named apart from its harness and model, every mechanism flai relies on (the guard, the permission prompt, the roles, the usage measurement) meeting the harness through neutral contracts, and the story's cost read from the gateway that charged it. Built on the finding of S-0339, [agent-adapters.md](../../../design/system/agent-adapters.md), and the decisions on TH-0368: [ADR-0129](../../../design/adrs/0129-a-story-s-agent-names-a-provider-beside-its-harness-and-model-and-the-providers.md), [ADR-0130](../../../design/adrs/0130-every-harness-meets-flai-through-neutral-contracts-an-adapter-s-capabilities-a.md), [ADR-0131](../../../design/adrs/0131-a-harness-without-a-guard-hook-runs-a-story-s-agent-but-none-of-its-roles.md), and [ADR-0132](../../../design/adrs/0132-a-story-s-cost-comes-from-its-provider-s-spend-log-when-it-has-a-provider-from.md).

- [ ] A story's `agent` names a `provider`, in the manifest's defaults and in story front matter, and the host holds a `providers` map (`api`, `base_url`, `key_env`, `models`) set with `flai serve agent provider` and the dashboard's settings; a project's manifest overrides an entry's `api`, `base_url`, and `models`, never `key_env` (ADR-0129).
- [ ] The `claude-code` adapter derives `ANTHROPIC_BASE_URL`, `ANTHROPIC_AUTH_TOKEN`, an empty `ANTHROPIC_API_KEY`, and the `ANTHROPIC_DEFAULT_*_MODEL` aliases from the provider, and `flai serve` refuses a start whose key variable is unset or whose `api` the harness cannot speak, naming the setting; no key is written to a log, a work item, or the start record.
- [ ] A story's agent runs through OpenRouter's `/api/v1/messages` with `flai guard` and `permission_prompt` on, and the first story records what OpenRouter forwards of Claude Code's headers and what its spend log reports beside `total_cost_usd`, before the contracts below are fixed.
- [ ] A story's agent runs the same way through a LiteLLM proxy's `/v1/messages`, with the same record.
- [ ] The adapter contract carries `Capabilities` and a usage `Reader`, the guard decides on a neutral `Call` through `Decide` with an input reader per harness and holds a call to ask the owner, and the protected list is flai's policy with each harness's own files in it (ADR-0130); `.claude/agents/` stays the definitions' source.
- [ ] A harness whose capabilities say it cannot run the guard, or cannot hold a call to ask, runs a story's agent but no roles unless the host sets `agent.harnesses.<name>.guard: none`, and a protected write on it is refused unless `auto-approve` is on, each refusal naming what lifts it (ADR-0131).
- [ ] A story with a provider is priced from the gateway's spend log (OpenRouter `/api/v1/key` and `/api/v1/generation`; LiteLLM `/spend/logs`), `usage` carries `priced_by` (`gateway`, `harness`, `estimate`) summed up the hierarchy, `metrics.md` defines it, and `flai stats` and the dashboard show it (ADR-0132).
- [ ] The OpenRouter and LiteLLM integration and smoke tests run only where the gateway's key is configured and are skipped, saying so, where it is not, so that no close-out or CI run without a key fails for want of it.
- [ ] `docs/operators/settings.md` indexes every new setting and its test passes, `docs/users/flai.md` and `design/system/flai-cli.md` describe providers and what a harness without a guard loses, and the template carries the new settings and conventions.

## Stories

To be drafted by the planner from the finding. Order: the provider split and OpenRouter first, with the first story settling the two open checks (headers forwarded, cost reported) against OpenRouter; then the neutral contracts inside the Claude Code adapter; then LiteLLM; then cost from the spend log and `priced_by`; the harness-without-a-guard rule lands with the capabilities; documentation and the template with each story.
- S-0349 A story's agent names a provider from the host's providers map, and a project's manifest may override the entry's api, base_url, and models
- S-0350 flai serve starts a claude-code agent through its provider, with the key read from the variable the host names and never written down
- S-0351 A claude-code agent runs through OpenRouter with the guard and permission_prompt on, checked by tests that run only where its key is set, and what OpenRouter forwards and charges is recorded
- S-0352 Each adapter states its capabilities and reads its harness's log into neutral usage events, and flai measures usage from those events
- S-0353 flai guard decides on a neutral call through Decide, and Claude Code's hook input is one reader of it
- S-0354 Asking the owner before a protected write is one hold-and-ask that the guard's ask verdict and permission_prompt share, and the protected list is flai's, kept per harness
- S-0355 A harness that cannot run the guard or hold a call to ask runs a story's agent but none of its roles, unless the host sets guard none for it, and each refusal names what lifts it
- S-0356 A claude-code agent runs through a LiteLLM proxy as it does through OpenRouter, checked by tests that run only where the proxy and its key are set, and what LiteLLM forwards and logs is recorded
- S-0357 A story's agent run through a provider is priced from the gateway's spend log, and an item's usage records priced_by, summed up the hierarchy as its least certain source
- S-0358 flai stats, flai show, and the dashboard show what priced each cost, and metrics.md defines priced_by and its sum
- S-0359 The dashboard's settings show and set the host's providers, and a story's and the project's agent fields name a provider
- S-0360 flai serve checks each new claude version through the provider a project's agents name, so a host that reaches Claude only through a gateway is not reported failing

## Notes

- The operator's answers, TH-0368 on 2026-10-08: the split with project overrides of non-secret values; option A alone, Claude Code over a gateway, no Codex CLI, OpenCode, Goose, or loop of flai's own in this epic; a harness without a guard as ADR-0131; cost from the gateway as ADR-0132; OpenRouter first, its tests gated on a configured key.
- The planner reads first: `design/system/agent-adapters.md`, its `## Abstractions` and `## Decision`; the four ADRs; `flai/internal/harness/adapters.go`, `flai/internal/guard/guard.go`, `flai/internal/mcpserver/permission.go`, `flai/internal/usage/log.go`, `flai/internal/manifest/agent.go`.
- What the finding leaves **to check**, under `## LiteLLM and OpenRouter`: whether each gateway's `/v1/messages` forwards every beta header and body field Claude Code sends; what `total_cost_usd` says against the gateway's spend log; whether the model aliases resolve; whether LiteLLM reads `x-claude-code-session-id` so that spend can be joined per session. The first story answers them against OpenRouter.
- Spend: an OpenRouter key with a few dollars of credit, and a local LiteLLM proxy, free, configured by the operator; cap the trial with `max_budget_usd` on the story's agent and a limit on the key.
- Claude models only: Anthropic and OpenRouter both say Claude Code over a gateway is supported for Claude models alone. Non-Claude models need another harness, which the finding recommends as Codex CLI, for a later epic.
- Cost of delay inputs are the operator's.
