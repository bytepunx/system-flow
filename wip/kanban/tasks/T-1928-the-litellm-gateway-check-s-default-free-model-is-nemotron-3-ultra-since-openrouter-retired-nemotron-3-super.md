---
id: T-1928
type: task
nature: feature
title: The LiteLLM gateway check's default free model is nemotron-3-ultra, since OpenRouter retired nemotron-3-super
status: done
parent: S-0451
owner: alex
created: 2026-10-20T22:20:12Z
updated: 2026-10-20T22:21:22Z
transitions:
  - to: ready
    at: 2026-10-20T22:20:22Z
    by: agent-S-0451
  - to: in-progress
    at: 2026-10-20T22:20:23Z
    by: agent-S-0451
  - to: done
    at: 2026-10-20T22:21:22Z
    by: agent-S-0451
stream: S-0451
tags: []
touches: [flai/internal/serve/gateway_test.go, scripts/litellm.sh, docs/operators/index.md, docs/contributors/index.md]
usage:
  source: log
  seconds: 59
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 12
      output: 2474
      cache_read: 748126
      cache_write: 29634
      cost: 0.428
---
# T-1928 The LiteLLM gateway check's default free model is nemotron-3-ultra, since OpenRouter retired nemotron-3-super

## Work

Set `litellmHaikuDefault` in `flai/internal/serve/gateway_test.go` to `openrouter-messages/nvidia/nemotron-3-ultra-550b-a55b:free`, since OpenRouter answers 404 for the retired `nvidia/nemotron-3-super-120b-a12b:free` and fails every close-out at smoke (I-0170). The operator chose this on TH-0529. Waits for nothing: T-1912 and T-1913 are done.

## Done when

- `TestGatewayLiteLLM` passes with the new default.
- The close-out passes smoke.

## Notes
