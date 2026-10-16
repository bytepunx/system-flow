---
id: T-1710
type: task
nature: feature
title: A gateway test runs a story's agent on OpenCode through OpenRouter only where the key and the known opencode are present, and asserts a gateway-priced cost
status: backlog
parent: S-0410
owner: alex
created: 2026-10-16T05:19:05Z
updated: 2026-10-16T05:19:05Z
transitions: []
stream: S-0410
tags: [cli]
touches: [flai/internal/serve/gateway_test.go]
after: [T-1709]
---
# T-1710 A gateway test runs a story's agent on OpenCode through OpenRouter only where the key and the known opencode are present, and asserts a gateway-priced cost

## Work

- Add `TestGatewayOpenRouterOpenCode` in `flai/internal/serve/gateway_test.go`, beside `TestGatewayOpenRouter`: skip under `testing.Short()`, and with a message naming what is missing when `OPENROUTER_API_KEY` is unset or `opencode` is missing or of another version, through `gatewayMissing`.
- Otherwise run one short prompt on a cheap OpenRouter model through the adapter, measure it, and assert `priced_by: gateway` and a cost above zero.
- Waits for T-1709: it measures through the serve path that task builds.

## Done when

The test passes with the key and the binary and is skipped, saying which is missing, without either.

## Notes
