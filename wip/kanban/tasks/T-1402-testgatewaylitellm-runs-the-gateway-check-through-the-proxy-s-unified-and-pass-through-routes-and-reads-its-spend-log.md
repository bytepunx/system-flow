---
id: T-1402
type: task
nature: feature
title: TestGatewayLiteLLM runs the gateway check through the proxy's unified and pass-through routes and reads its spend log
status: backlog
parent: S-0356
owner: alex
created: 2026-10-08T08:49:47Z
updated: 2026-10-08T08:49:47Z
transitions: []
stream: S-0356
tags: [cli]
touches: [flai/internal/serve/gateway_test.go]
---
# T-1402 TestGatewayLiteLLM runs the gateway check through the proxy's unified and pass-through routes and reads its spend log

## Work

- In `gateway_test.go`, `TestGatewayLiteLLM`: skip unless `LITELLM_BASE_URL` and `LITELLM_API_KEY` are set and `claude` is on `PATH`, naming what is missing; skip under `-short`.
- Two subtests, `unified` with `base_url` the proxy's root and `pass-through` with `<root>/anthropic`, each with `key_env: LITELLM_API_KEY` and `models.haiku` from `FLAI_TEST_LITELLM_HAIKU`, default `anthropic/claude-haiku-4-5`, through S-0351's check: the protected `Write`, the guard's refusal, the alias, and a `result`.
- After each, `GET /spend/logs` with the key, filtered to the run's window: `t.Log` the spend, the tokens, and whether any row's metadata carries the session ID Claude Code sent, beside `total_cost_usd`. A missing row is logged, not failed.

First layer.

## Done when

- Without the variables, the test skips naming them.
- With them, both subtests pass and log both costs.
- `flai test flai/internal/serve/` passes.

## Notes

Layer 1 of S-0356.
