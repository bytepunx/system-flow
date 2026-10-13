---
id: T-1526
type: task
nature: feature
title: The local proxy reaches OpenRouter's free models and the LiteLLM check defaults to one, with the pass-through opt-in
status: done
parent: S-0356
owner: alex
created: 2026-10-13T17:27:26Z
updated: 2026-10-13T17:30:28Z
transitions:
  - to: ready
    at: 2026-10-13T17:27:33Z
    by: agent-S-0356
  - to: in-progress
    at: 2026-10-13T17:27:33Z
    by: agent-S-0356
  - to: done
    at: 2026-10-13T17:30:28Z
    by: agent-S-0356
stream: S-0356
tags: []
touches: [scripts/litellm.sh, flai/internal/serve/gateway_test.go, docs/contributors/index.md]
usage:
  source: log
  seconds: 172
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 24
      output: 194
      cache_read: 3846235
      cache_write: 20762
      cost: 2.4611
---

# T-1526 The local proxy reaches OpenRouter's free models and the LiteLLM check defaults to one, with the pass-through opt-in

## Work

- `scripts/litellm.sh`: an `openrouter/*` deployment upstream to OpenRouter with `OPENROUTER_API_KEY` from `admin.env`; `ANTHROPIC_API_KEY` is no longer required on the first `up`; the key file's `FLAI_TEST_LITELLM_HAIKU` defaults to the free model.
- `gateway_test.go`: `FLAI_TEST_LITELLM_HAIKU` defaults to `openrouter/nvidia/nemotron-3-super-120b-a12b:free`; the pass-through subtest runs only with `FLAI_TEST_LITELLM_PASSTHROUGH=1` and skips saying so otherwise.
- `docs/contributors/index.md`: the rows for the new default and `FLAI_TEST_LITELLM_PASSTHROUGH`.

First of the new layer (TH-0429); T-1404 waits for it, since its run needs the proxy and the test as changed here.

## Done when

- `scripts/litellm.sh up` with `OPENROUTER_API_KEY` set serves `openrouter/*`.
- `flai test` on the three paths passes.

## Notes
