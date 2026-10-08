---
id: ADR-0132
title: "A story's cost comes from its provider's spend log when it has a provider, from the harness's own figure when it has none, and is estimated otherwise, with the source named on its usage as priced_by"
status: proposed
date: 2026-10-08
supersedes: []
superseded_by: []
refines: [ADR-0051]
---

# ADR-0132 A story's cost comes from its provider's spend log when it has a provider, from the harness's own figure when it has none, and is estimated otherwise, with the source named on its usage as priced_by

## Context

Work items record the tokens and cost their agents spent, measured from the agents' logs and never from a price table ([ADR-0051](0051-work-items-record-the-tokens-and-cost-their-agents-spent-measured-from-the.md)). The one source of a reported cost today is Claude Code's `result` event, whose `total_cost_usd` and `modelUsage` are its own estimate at Anthropic's prices. Over a gateway that estimate is wrong: LiteLLM charges by its model cost map and the deployment's price, OpenRouter by the provider it routed to, and only the gateway knows what was charged. Codex CLI and Goose report tokens and no cost at all. The finding of S-0339, [agent-adapters.md](../system/agent-adapters.md), lists the sources: the harness's figure, the gateway per call (LiteLLM's `x-litellm-response-cost` header, OpenRouter's `usage.cost`), the gateway after the fact (LiteLLM `/spend/logs` by virtual key, OpenRouter `/api/v1/key` and `/api/v1/generation`), a price table, and the blended rate flai estimates from reported totals.

## Decision

A story's cost comes from its provider's spend log when it has a provider, from the harness's own figure when it has none, and is estimated otherwise, with the source named on its usage as priced_by.

- Each provider `api` ([ADR-0129](0129-a-story-s-agent-names-a-provider-beside-its-harness-and-model-and-the-providers.md)) gets a spend reader in flai: for LiteLLM, `/spend/logs` filtered by the key `key_env` names; for OpenRouter, `/api/v1/key` and `/api/v1/generation`. A key per project, or per agent, joins the gateway's spend to a story; `flai serve` reads it when a run ends, as it measures the log today.
- When the story has a provider, its cost is the gateway's figure for the run's calls, and the harness's own cost is kept as tokens only. When it has none, the harness's figure stands, as today. When neither gives one, the cost is estimated from `Rates`, as today, and marked so.
- `usage` on an item gains `priced_by`, one of `gateway`, `harness`, or `estimate`, summed up the hierarchy as the lowest-confidence source among its children. `flai stats` and the dashboard show it beside the cost.
- No price table is added to flai: a table dates, and the gateways already price every call.

## Consequences

- `design/system/metrics.md`, the contract between `flai stats` and the dashboard, gains `priced_by` and the rule for its sum; the metrics tests and the dashboard's usage views change with it.
- A run that was started before a provider was named, or whose gateway log is unreachable when the run ends, is priced by the harness or estimated and says so; flai reads the gateway again on the next measurement.
- The host needs read access to the gateway's spend endpoint with the same key the agent uses, or an admin key the provider entry names; the entry gains that name where the two differ.

## Alternatives considered

- The harness's figure everywhere: wrong by the provider's margin over a gateway, and absent on harnesses that report tokens only.
- A price table in flai, LiteLLM's public model cost map or one of flai's own: against ADR-0051's choice of measured over tabulated, dated as soon as a price moves, and still blind to a gateway's own margin and fallbacks.
- Per-call cost from the response headers: exact, but only a loop of flai's own sees the response; every other harness hides it.
