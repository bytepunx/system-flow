---
id: T-1944
type: task
nature: remediation
title: TestGatewayLiteLLM falls back through a list of free models when the proxy answers its default with a refusal
status: done
parent: S-0456
owner: alex
created: 2026-10-21T06:03:28Z
updated: 2026-10-21T07:47:56Z
transitions:
  - to: ready
    at: 2026-10-21T07:29:51Z
    by: agent-S-0456
  - to: in-progress
    at: 2026-10-21T07:29:51Z
    by: agent-S-0456
  - to: done
    at: 2026-10-21T07:47:56Z
    by: agent-S-0456
stream: S-0456
tags: [flai, serve, tests]
touches: [flai/internal/serve/gateway_test.go, scripts/litellm.sh]
---
# T-1944 TestGatewayLiteLLM falls back through a list of free models when the proxy answers its default with a refusal

## Work

I-0170 asks for a fallback or a skip when the proxy answers 404, because OpenRouter can retire a free model at any time. S-0455's probe now fails the check at once on such an answer, naming the slug. This task makes the default survive a retirement instead:

- Replace the single `litellmHaikuDefault` with an ordered list of free OpenRouter slugs, `openrouter-messages/<slug>`, the current default first. The proxy routes `openrouter-messages/*` by wildcard, so a slug on the list needs no proxy change.
- Find at least one more free model that passes `TestGatewayLiteLLM` live, at $0 of key spend: the protected Write, the sub-agent's refused `item_edit`, the alias, and a result. A model that answers but fails the session does not go on the list. Record each slug tried and its result in this task's notes.
- When the unified case's slug is on the list, whether from the default or from `FLAI_TEST_LITELLM_HAIKU` set by `scripts/litellm.sh` to its default, probe the slugs in order and run the session on the first that answers 2xx. Log each slug passed over and the proxy's answer.
- A slug the operator chose, one not on the list, is probed alone and fails as S-0455 made it, naming the slug and its source.
- When every slug on the list is refused, fail with `t.Fatal`, naming each slug and its answer, and how to name another model. Do not skip: a skip would let every close-out pass with no gateway check and nothing to say so.
- Keep `HAIKU_DEFAULT` in `scripts/litellm.sh` on the list, and say in its comment that it must be.
- Extend the short test with the `httptest` stand-in: the first slug answers 404 as OpenRouter does through LiteLLM and the second answers 200, so the pick is the second, which reproduces I-0170; every slug refused fails with each named; an operator's slug refused is not replaced.

Waits for nothing in this story. The story waits for S-0455, whose T-1942 adds the probe in this file and whose T-1941 changes `scripts/litellm.sh`.

## Done when

- The short test fails on the code before this task and passes after it.
- `flai test flai/internal/serve/gateway_test.go scripts/litellm.sh` passes.
- With the first slug on the list set to a retired one, `TestGatewayLiteLLM` passes on the next slug, and its log names the one passed over.
- Every slug on the list passed `TestGatewayLiteLLM` live, recorded in this task's notes.

## Notes

Live runs of `TestGatewayLiteLLM` through `scripts/gateway-smoke.sh` on 2026-10-21, each slug as `openrouter-messages/<slug>`, every one at $0.000000 of key spend:

| Slug | Result | Test time |
|------|--------|-----------|
| `nvidia/nemotron-3-ultra-550b-a55b:free` | passed, picked from the key file's default | 146.3 s |
| `minimax/minimax-m3.1-flash:free` | passed | 148.1 s |
| `nvidia/nemotron-3.5-lightning:free` | passed | 161.2 s |
| retired `nvidia/nemotron-3-super-120b-a12b:free` put first for one run | passed over on its 404, logged as `passed over`; the run passed on nemotron-3-ultra | 147.9 s |

Probes only, not run whole: `poolside/laguna-s-2.1:free` and `cohere/north-mini-code:free` answered 200; `nvidia/nemotron-3.5-super-vl:free` 502, `google/gemma-4-31b-it:free` 429, `thinkingmachines/inkling:free` 403.

The first run also ran `TestGatewayOpenRouter`, as `OPENROUTER_API_KEY` was set in the sub-agent's environment, and spent $0.051729 on the OpenRouter key; the later runs unset it.
