---
id: T-1406
type: task
nature: feature
title: usage.Spend is the spend reader's contract, and OpenRouter's reads a run's cost from /api/v1/generation or /api/v1/key
status: done
parent: S-0357
owner: alex
created: 2026-10-08T08:50:37Z
updated: 2026-10-13T19:25:35Z
transitions:
  - to: ready
    at: 2026-10-13T19:17:08Z
    by: agent-S-0357
  - to: in-progress
    at: 2026-10-13T19:17:08Z
    by: agent-S-0357
  - to: done
    at: 2026-10-13T19:25:35Z
    by: agent-S-0357
stream: S-0357
tags: [cli]
touches: [flai/internal/usage/spend.go, flai/internal/usage/spend_test.go, flai/internal/usage/openrouter.go, flai/internal/usage/openrouter_test.go, flai/internal/usage/testdata/openrouter]
usage:
  source: log
  seconds: 507
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 25
      output: 13202
      cache_read: 1620907
      cache_write: 49294
      cost: 0.8846
---
# T-1406 usage.Spend is the spend reader's contract, and OpenRouter's reads a run's cost from /api/v1/generation or /api/v1/key

## Work

- `spend.go`, the contract, in `flai/internal/usage` and with no import of `manifest`:
  - `RunRef`: the run's harness session IDs (`Sessions`), its window (`From`, `To`), the gateway's generation IDs where the harness surfaces them (`Generations`, none from Claude Code), and `KeyBefore *float64`, the key's usage read when the run started.
  - `Spent`: `Cost` in US dollars that the gateway charged for the run, `Calls` it counted where it says, and `By`, how it joined (`generation`, `key`, `session`, `window`).
  - `Spend interface { Run(ctx context.Context, run RunRef) (*Spent, error) }`, and `ErrNoSpend`, which a reader returns, wrapped with why, when the gateway answers but holds nothing for the run.
  - `KeyReader interface { KeyUsage(ctx context.Context) (float64, error) }`, the reading a reader takes when the run starts, for a gateway that joins by the key's usage before and after.
  - `For(baseURL, key string, client *http.Client) Spend`: OpenRouter's reader when the base URL's host is `openrouter.ai` or below it, LiteLLM's otherwise, nil when the host is `api.anthropic.com` (no spend log; the harness's figure stands). LiteLLM's constructor is T-1407's: here `For` calls `NewLiteLLM(baseURL, key, client)`, and T-1406 leaves a minimal `litellm.go` out: write `For` so it compiles with only OpenRouter, by returning nil for LiteLLM until T-1407 adds it, and say so in a comment T-1407 removes. The key is a value the caller looked up; it is sent as a bearer and never logged or put in an error.
- `openrouter.go`: `NewOpenRouter(baseURL, key, client)`, whose base URL is the provider's (`https://openrouter.ai/api` or `.../api/v1`: normalize to the API root).
  - With `Generations`, `GET /api/v1/generation?id=<id>` each, summing `data.total_cost`; `By` is `generation`.
  - Otherwise, with `KeyBefore`, `GET /api/v1/key`'s `data.usage` now, taking the highest of three reads (S-0351: copies lag one another), minus `KeyBefore`; `By` is `key`. A delta of zero or less is `ErrNoSpend` (the run has not shown yet; the next measurement reads again).
  - Neither: `ErrNoSpend`, saying the run has no generation IDs and no reading from its start.
  - `KeyUsage` takes the highest of three reads of `/api/v1/key`.
  - A non-2xx answer is an error naming the endpoint and status, never the key.
- Tests serve recorded OpenRouter responses with `httptest`, shaped as `flai/internal/serve/gateway_test.go`'s `openRouterKeyInfo` and generation lookups read them: a join by generation ID, a join by the key before and after with lagging copies, no reading (`ErrNoSpend`), a 401, a context timeout, and `For` choosing by host.

First layer.

## Done when

- `flai test flai/internal/usage/` passes.

## Notes

Layer 1 of S-0357. Read `design/system/agent-adapters.md` § Checked against OpenRouter: Claude Code's message IDs are not OpenRouter's generation IDs, so a Claude Code run joins by the key. Meets part of criterion 1.
