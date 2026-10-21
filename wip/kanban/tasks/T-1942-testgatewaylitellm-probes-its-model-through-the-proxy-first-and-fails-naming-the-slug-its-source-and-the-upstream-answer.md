---
id: T-1942
type: task
nature: remediation
title: TestGatewayLiteLLM probes its model through the proxy first and fails naming the slug, its source, and the upstream answer
status: in-progress
parent: S-0455
owner: alex
created: 2026-10-21T00:39:48Z
updated: 2026-10-21T02:10:09Z
transitions:
  - to: ready
    at: 2026-10-21T02:10:08Z
    by: agent-S-0455
  - to: in-progress
    at: 2026-10-21T02:10:09Z
    by: agent-S-0455
stream: S-0455
tags: [flai, serve, tests]
touches: [flai/internal/serve/gateway_test.go]
---
# T-1942 TestGatewayLiteLLM probes its model through the proxy first and fails naming the slug, its source, and the upstream answer

## Work

When OpenRouter retired the free model, `TestGatewayLiteLLM/unified` failed on the `[claude-code:unrecognized_model]` warning. That warning also shows on passing runs, and it hid the real cause, an upstream 404 (I-0171, I-0170). Make the check name the cause itself:

- Before the `claude -p` run, send one `POST /v1/messages` through the proxy with the slug and `max_tokens` 1. The model is free, so the call costs nothing.
- When the answer is not 2xx, fail with `t.Fatalf`. Name the slug, where it came from (`FLAI_TEST_LITELLM_HAIKU` from the environment or `litellm.env`, or `litellmHaikuDefault`), the upstream status and message, and how to change it (`LITELLM_HAIKU_MODEL=<slug> scripts/litellm.sh up`).
- Put the probe in a helper that a short test can reach. Add a short test that reproduces the issue with an `httptest` server answering 404 the way OpenRouter does through LiteLLM, and checks the message.

Waits for nothing in this story. It runs beside the `litellm.sh` task, with no path in common. The story waits for S-0451, whose T-1928 changes `litellmHaikuDefault` in this file.

## Done when

- The short test fails without the probe and passes with it.
- `flai test flai/internal/serve/gateway_test.go` passes: gofmt, vet, golangci-lint, and the short tests.
- With the retired slug set in `FLAI_TEST_LITELLM_HAIKU`, `TestGatewayLiteLLM` fails with the probe's message, not after a Claude Code session.

## Notes
