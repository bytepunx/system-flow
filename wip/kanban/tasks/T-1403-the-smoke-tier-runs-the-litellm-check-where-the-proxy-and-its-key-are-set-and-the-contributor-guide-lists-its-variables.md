---
id: T-1403
type: task
nature: feature
title: The smoke tier runs the LiteLLM check where the proxy and its key are set, and the contributor guide lists its variables
status: backlog
parent: S-0356
owner: alex
created: 2026-10-08T08:49:53Z
updated: 2026-10-08T08:49:53Z
transitions: []
stream: S-0356
tags: [cli]
touches: [scripts/gateway-smoke.sh, docs/contributors/index.md]
after: [T-1402]
---
# T-1403 The smoke tier runs the LiteLLM check where the proxy and its key are set, and the contributor guide lists its variables

## Work

- `scripts/gateway-smoke.sh`: add LiteLLM to its list of gateways, run when both `LITELLM_BASE_URL` and `LITELLM_API_KEY` are set, with the skip line naming the one missing otherwise.
- `docs/contributors/index.md`: `LITELLM_BASE_URL`, `LITELLM_API_KEY`, `FLAI_TEST_LITELLM_HAIKU`, and a virtual key with a `max_budget`.

It waits for T-1402, whose test it runs.

## Done when

- `scripts/smoke.sh` on a host without the variables passes and prints both skip lines.
- `flai test scripts/ docs/contributors/index.md` passes.

## Notes

Layer 2 of S-0356.
