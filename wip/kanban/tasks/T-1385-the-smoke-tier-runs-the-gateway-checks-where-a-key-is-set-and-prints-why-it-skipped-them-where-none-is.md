---
id: T-1385
type: task
nature: feature
title: The smoke tier runs the gateway checks where a key is set and prints why it skipped them where none is
status: backlog
parent: S-0351
owner: alex
created: 2026-10-08T08:45:10Z
updated: 2026-10-08T08:45:10Z
transitions: []
stream: S-0351
tags: [cli]
touches: [scripts/gateway-smoke.sh, scripts/smoke.sh, docs/contributors/index.md]
after: [T-1384]
---
# T-1385 The smoke tier runs the gateway checks where a key is set and prints why it skipped them where none is

## Work

- `scripts/gateway-smoke.sh`, POSIX `sh` under `set -eu`: for each gateway whose key variable is set, run `go test -count=1 -v -run '^TestGateway<Name>$' ./internal/serve/` from `flai/`; for each whose key is unset, print `smoke: gateway <name> skipped: <VAR> is not set`. OpenRouter is the one gateway now; the script lists the gateways so that S-0351's LiteLLM successor adds one line.
- `scripts/smoke.sh` calls it after the installer step, under its own `echo` line.
- `docs/contributors/index.md`: how to run the checks, the variables (`OPENROUTER_API_KEY`, `FLAI_TEST_OPENROUTER_HAIKU`, `FLAI_TEST_GATEWAY_BUDGET`), and that a key should carry its own credit limit.

It waits for T-1384, whose test it runs.

## Done when

- `scripts/smoke.sh` on a host without the key passes and prints the skip line.
- `flai test scripts/ docs/contributors/index.md` passes.

## Notes

Layer 2 of S-0351.
