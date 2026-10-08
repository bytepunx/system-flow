---
id: S-0357
type: story
nature: feature
title: A story's agent run through a provider is priced from the gateway's spend log, and an item's usage records priced_by, summed up the hierarchy as its least certain source
status: backlog
parent: E-0019
owner: alex
created: 2026-10-08T08:50:22Z
updated: 2026-10-08T09:44:32Z
transitions: []
tags: [cli]
topics: [agents, usage]
touches: [flai/internal/usage/spend.go, flai/internal/usage/spend_test.go, flai/internal/usage/openrouter.go, flai/internal/usage/openrouter_test.go, flai/internal/usage/litellm.go, flai/internal/usage/litellm_test.go, flai/internal/usage/usage.go, flai/internal/usage/usage_test.go, flai/internal/workitem/usage.go, flai/internal/workitem/usage_test.go, flai/internal/serve/usage.go, flai/internal/serve/usage_test.go, flai/internal/manifest/provider.go, flai/internal/manifest/provider_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, docs/operators/settings.md, docs/users/flai-reference.md, design/system/flai-cli.md]
after: [S-0352, S-0356]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 46.46
  by: planner-E-0019
  at: 2026-10-08T09:00:25Z
forecast:
  duration: 42m
  delivery: 2026-10-09T00:51:00Z
  basis: "Its own forecast of 42m; 27th in the pull order with an in-progress limit of 5, behind S-0232, S-0297, S-0334, S-0344, S-0338, S-0342, S-0337, S-0343, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0289, S-0304, S-0305, S-0306, S-0313, S-0347, S-0349, S-0350, S-0351, S-0352, S-0353, S-0354, S-0355 and S-0356."
  by: flai
  at: 2026-10-08T09:44:32Z
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

- T-1405 usage carries priced_by in the front matter and sums it up the hierarchy as the least certain source
- T-1406 usage.Spend is the spend reader's contract, and OpenRouter's reads a run's cost from /api/v1/generation or /api/v1/key
- T-1407 LiteLLM's spend reader reads a run's cost from /spend/logs by key and window, or by session where the proxy keeps it
- T-1408 A provider entry may name spend_key_env, the variable of a key that can read the gateway's spend log
- T-1409 flai serve prices a run through a provider from the gateway when it measures it, and falls back to the harness or an estimate, saying why
- T-1410 flai-cli.md describes the spend readers, when flai serve reads them, and priced_by

## Notes

- No price table is added (ADR-0132).
- How a run's calls join the gateway's records is what S-0351 and S-0356 record; the readers follow their finding, and this story's agent reads `agent-adapters.md`'s two check sections before it writes them.

### Planning

Planned by planner-E-0019 on 2026-10-08. Every touch is a file; no folder touch is kept; `spend.go`, `openrouter.go`, `litellm.go`, and their tests are new. Four layers: T-1405, T-1406, and T-1408 together, with no file in common; then T-1407, which implements T-1406's contract; then T-1409, which joins them all; then T-1410.

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/usage/usage.go`, `usage_test.go`, `flai/internal/workitem/usage.go`, `usage_test.go` | design, layout | `priced_by` on `Usage`, its sum, and its front matter (T-1405) |
| `flai/internal/usage/spend.go`, `spend_test.go`, `openrouter.go`, `openrouter_test.go` | design | ADR-0132's reader per gateway, OpenRouter first (T-1406) |
| `flai/internal/usage/litellm.go`, `litellm_test.go` | design | LiteLLM's reader (T-1407) |
| `flai/internal/manifest/provider.go`, `provider_test.go`, `flai/cmd/serve_actions.go`, `serve_actions_test.go` | design | ADR-0132's consequences: the entry names the spend key where it differs (T-1408) |
| `docs/operators/settings.md`, `docs/users/flai-reference.md` | design | The new key and flag; `make flai-reference` regenerates both (T-1408) |
| `flai/internal/serve/usage.go`, `usage_test.go` | layout | Where a story's runs are measured (T-1409) |
| `design/system/flai-cli.md` | design | How usage is measured (T-1410) |

`touches suggest` listed `docs/users/flai.md`, the dashboard's documents, and `flai/internal/hostapi/writes.go`. None is taken: what a user sees changes in S-0358.

Forecast: 42m, as `flai forecast` gives it, 100 s per unit over 34 done feature stories on claude-opus-5-5 in the large band, times size 25. It stands: the readers are tested against recorded responses, with no live call.

Cost of delay: 46.46 USD a week, as `flai cod` works it out: 42m of the 7h32m forecast over E-0019's 12 stories, of the operator's 500 USD a week penalty. It stands.
