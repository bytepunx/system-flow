---
id: T-1406
type: task
nature: feature
title: usage.Spend is the spend reader's contract, and OpenRouter's reads a run's cost from /api/v1/generation or /api/v1/key
status: backlog
parent: S-0357
owner: alex
created: 2026-10-08T08:50:37Z
updated: 2026-10-08T08:50:37Z
transitions: []
stream: S-0357
tags: [cli]
touches: [flai/internal/usage/spend.go, flai/internal/usage/spend_test.go, flai/internal/usage/openrouter.go, flai/internal/usage/openrouter_test.go]
---
# T-1406 usage.Spend is the spend reader's contract, and OpenRouter's reads a run's cost from /api/v1/generation or /api/v1/key

## Work

- `spend.go`: `Spend interface{ Run(ctx, RunRef) (map[string]Model, error) }`, where `RunRef` holds the provider, the key's variable name, the session, the run's window, and its message IDs, and the result is the cost and tokens per model the gateway charged. `For(api, base_url)` picks the reader: OpenRouter by its host, LiteLLM otherwise for `anthropic-messages`.
- `openrouter.go`: by the run's message IDs through `/api/v1/generation?id=`, summing `total_cost` per model, as S-0351 found the IDs join; where they do not, `/api/v1/key`'s `usage` before and after the run, with the reason in the result. The key is read from its variable as the request is made and is never logged.
- Tests serve recorded OpenRouter responses with `httptest`: a join by ID, a fallback by key, a 401, and a timeout.

First layer.

## Done when

- `flai test flai/internal/usage/` passes.

## Notes

Layer 1 of S-0357. Read `agent-adapters.md` § Checked against OpenRouter before writing the join.
