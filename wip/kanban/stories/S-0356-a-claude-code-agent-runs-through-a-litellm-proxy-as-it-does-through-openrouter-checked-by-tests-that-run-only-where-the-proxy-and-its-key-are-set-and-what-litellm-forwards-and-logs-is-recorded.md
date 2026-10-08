---
id: S-0356
type: story
nature: feature
title: A claude-code agent runs through a LiteLLM proxy as it does through OpenRouter, checked by tests that run only where the proxy and its key are set, and what LiteLLM forwards and logs is recorded
status: backlog
parent: E-0019
owner: alex
created: 2026-10-08T08:49:40Z
updated: 2026-10-08T08:52:29Z
transitions: []
tags: [cli]
topics: [agents, testing]
touches: [flai/internal/serve/gateway_test.go, scripts/gateway-smoke.sh, docs/contributors/index.md, docs/operators/index.md, design/system/agent-adapters.md]
after: [S-0351]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
finalized:
  by: alex
  at: 2026-10-08T08:52:29Z
---
# S-0356 A claude-code agent runs through a LiteLLM proxy as it does through OpenRouter, checked by tests that run only where the proxy and its key are set, and what LiteLLM forwards and logs is recorded

## Goal

S-0351 checks Claude Code against OpenRouter. This story runs the same check against a LiteLLM proxy the operator runs locally, through both of its Anthropic routes: the unified `/v1/messages`, which translates and may drop fields, and the `/anthropic` pass-through, which forwards untranslated. It records what each forwards, whether LiteLLM's spend log can be joined to a session through `x-claude-code-session-id`, and what `/spend/logs` reports beside `total_cost_usd` (agent-adapters.md § LiteLLM). The tests run only where the proxy's URL and a key are set.

## Acceptance criteria

- [ ] `TestGatewayLiteLLM` in `flai/internal/serve` runs the S-0351 check once through the proxy's unified `/v1/messages` and once through `/anthropic`, each with a provider entry for that route, when `LITELLM_BASE_URL` and `LITELLM_API_KEY` are set and `claude` is on `PATH`, and skips naming what is missing otherwise; the integration tier passes on a host without them.
- [ ] The test reads `/spend/logs` for the key after each run and logs the gateway's spend, and whether a row carries the session's ID, beside `total_cost_usd`.
- [ ] `scripts/gateway-smoke.sh` runs it where the variables are set and prints the skip line where they are not; the smoke tier passes on a host without them.
- [ ] `design/system/agent-adapters.md` gains a LiteLLM section beside S-0351's, with the date, the proxy's version, and the answers for each route: fields forwarded, aliases resolved, session join, and the gateway's spend against `total_cost_usd`. It says which route a provider entry should use.
- [ ] `docs/operators/index.md` shows a provider entry for a local LiteLLM proxy and the virtual key it needs; `docs/contributors/index.md` lists the variables.

## Tasks

Drafted by the planner; see the children.
- T-1402 TestGatewayLiteLLM runs the gateway check through the proxy's unified and pass-through routes and reads its spend log
- T-1403 The smoke tier runs the LiteLLM check where the proxy and its key are set, and the contributor guide lists its variables
- T-1404 Run the LiteLLM check against the operator's proxy and record what each route forwarded and logged, and the provider entry to use

## Notes

- The proxy is the operator's to run, free, from Docker or `litellm --config`, with a virtual key carrying a `max_budget` (E-0019's notes). Without it in the agent's environment, the story's agent asks on a thread.
- The epic orders LiteLLM after the neutral contracts. This story needs only S-0351's test, so it may run beside S-0352 to S-0354; the pull order is the operator's.
