---
id: T-1710
type: task
nature: feature
title: A gateway test runs a story's agent on OpenCode through OpenRouter only where the key and the known opencode are present, and asserts a gateway-priced cost
status: done
parent: S-0410
owner: alex
created: 2026-10-16T05:19:05Z
updated: 2026-11-07T21:38:55Z
transitions:
  - to: ready
    at: 2026-11-07T21:31:40Z
    by: agent-S-0410
  - to: in-progress
    at: 2026-11-07T21:31:41Z
    by: agent-S-0410
  - to: done
    at: 2026-11-07T21:38:55Z
    by: agent-S-0410
stream: S-0410
tags: [cli]
touches: [flai/internal/serve/gateway_test.go, scripts/gateway-smoke.sh]
after: [T-1709]
usage:
  source: log
  seconds: 434
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 36
      output: 12527
      cache_read: 5607031
      cache_write: 98113
      cost: 1.9884
---
# T-1710 A gateway test runs a story's agent on OpenCode through OpenRouter only where the key and the known opencode are present, and asserts a gateway-priced cost

## Work

- Add `TestGatewayOpenRouterOpenCode` in `flai/internal/serve/gateway_test.go`, beside `TestGatewayOpenRouter`: skip under `testing.Short()`, and with a message naming what is missing when `OPENROUTER_API_KEY` is unset or `opencode` is missing or of another version, through `gatewayMissing`.
- Otherwise run one short prompt on a cheap OpenRouter model through the adapter, measure it, and assert `priced_by: gateway` and a cost above zero.
- Waits for T-1709: it measures through the serve path that task builds.

## Done when

The test passes with the key and the binary and is skipped, saying which is missing, without either.

## Notes
