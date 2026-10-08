---
id: S-0357
type: story
nature: feature
title: A story's agent run through a provider is priced from the gateway's spend log, and an item's usage records priced_by, summed up the hierarchy as its least certain source
status: backlog
parent: E-0019
owner: alex
created: 2026-10-08T08:50:22Z
updated: 2026-10-08T08:52:43Z
transitions: []
tags: [cli]
topics: [agents, usage]
touches: [flai/internal/usage/spend.go, flai/internal/usage/spend_test.go, flai/internal/usage/openrouter.go, flai/internal/usage/openrouter_test.go, flai/internal/usage/litellm.go, flai/internal/usage/litellm_test.go, flai/internal/usage/usage.go, flai/internal/usage/usage_test.go, flai/internal/workitem/usage.go, flai/internal/workitem/usage_test.go, flai/internal/serve/usage.go, flai/internal/serve/usage_test.go, flai/internal/manifest/provider.go, flai/internal/manifest/provider_test.go, docs/operators/settings.md, design/system/flai-cli.md]
after: [S-0352, S-0356]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
finalized:
  by: alex
  at: 2026-10-08T08:52:43Z
---
# S-0357 A story's agent run through a provider is priced from the gateway's spend log, and an item's usage records priced_by, summed up the hierarchy as its least certain source

## Goal

Over a gateway, Claude Code's `total_cost_usd` is its own estimate at Anthropic's prices, not what the gateway charged. [ADR-0132](../../../design/adrs/0132-a-story-s-cost-comes-from-its-provider-s-spend-log-when-it-has-a-provider-from.md) takes the cost of a run through a provider from the gateway's spend log, keeps the harness's figure when there is no provider, and estimates otherwise, naming the source on the item's `usage` as `priced_by`. This story builds the spend readers, OpenRouter's from `/api/v1/generation` and `/api/v1/key` and LiteLLM's from `/spend/logs`, the pricing when `flai serve` measures a run, and `priced_by` in the front matter. It joins a run's calls to the gateway as S-0351 and S-0356 found they can be joined. Showing `priced_by` in `flai stats` and the dashboard is the next story.

## Acceptance criteria

- [ ] `usage.Spend` is the interface a spend reader meets, and one exists per gateway: OpenRouter's reads `/api/v1/generation` by the run's generation IDs, or `/api/v1/key` before and after where IDs do not join, and LiteLLM's reads `/spend/logs` by key and the run's window, or by session where S-0356 found the session ID; each is tested against a recorded response served by `httptest`.
- [ ] When `flai serve` measures a run whose agent had a provider, the run's cost per model is the gateway's figure and `priced_by` is `gateway`; the harness's own cost is kept out of the total. With no provider, the harness's figure stands and `priced_by` is `harness`. When neither gives one, the cost is estimated from `Rates` and `priced_by` is `estimate`, as `estimated` says today.
- [ ] A gateway that cannot be reached, or answers nothing for the run, leaves the run priced by the harness or estimated, says so in `flai serve`'s log, and is read again at the next measurement.
- [ ] `usage.priced_by` is written to an item's front matter and summed up the hierarchy as the least certain source among its parts (`estimate` below `harness` below `gateway`), and `estimated` keeps its meaning; tests cover the sum.
- [ ] A provider entry may name `spend_key_env`, the variable of a key that can read the spend log where the agent's own key cannot; `docs/operators/settings.md` indexes it and `TestSettingsIndexIsComplete` passes; no key is logged.
- [ ] `design/system/flai-cli.md` describes the readers and when each runs.

## Tasks

Drafted by the planner; see the children.
- T-1405 usage carries priced_by in the front matter and sums it up the hierarchy as the least certain source
- T-1406 usage.Spend is the spend reader's contract, and OpenRouter's reads a run's cost from /api/v1/generation or /api/v1/key
- T-1407 LiteLLM's spend reader reads a run's cost from /spend/logs by key and window, or by session where the proxy keeps it
- T-1408 A provider entry may name spend_key_env, the variable of a key that can read the gateway's spend log
- T-1409 flai serve prices a run through a provider from the gateway when it measures it, and falls back to the harness or an estimate, saying why
- T-1410 flai-cli.md describes the spend readers, when flai serve reads them, and priced_by

## Notes

- No price table is added (ADR-0132).
- How a run's calls join the gateway's records is what S-0351 and S-0356 record; the readers follow their finding, and this story's agent reads `agent-adapters.md`'s two check sections before it writes them.
