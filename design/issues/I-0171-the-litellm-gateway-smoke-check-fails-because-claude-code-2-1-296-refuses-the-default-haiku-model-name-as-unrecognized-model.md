---
id: I-0171
title: The LiteLLM gateway smoke check fails because Claude Code 2.1.296 refuses the default haiku model name as unrecognized_model
class: blocker
status: open
count: 2
cost: 15m
first_reported: 2026-10-20T20:07:32Z
last_reported: 2026-10-21T01:31:48Z
updated: 2026-10-21T01:31:48Z
---

# I-0171 The LiteLLM gateway smoke check fails because Claude Code 2.1.296 refuses the default haiku model name as unrecognized_model

## Description
The LiteLLM gateway smoke check fails because Claude Code 2.1.296 refuses the default haiku model name as unrecognized_model

## Instances

### 2026-10-20T20:07:32Z
Story: S-0450.
S-0450's close-out stopped at the smoke tier: TestGatewayLiteLLM/unified (flai/internal/serve/gateway_test.go:557) got `[claude-code:unrecognized_model] {"model":"openrouter-messages/nvidia/nemotron-3-super-120b-a12b:free"}` from Claude Code 2.1.296, the session ended in 0s with no model call, so flai-check.md was never written and no guard refusal was seen. It failed the same way when run again alone. The model is litellmHaikuDefault in gateway_test.go and LITELLM_HAIKU_MODEL in scripts/litellm.sh. S-0450 changes nothing under flai/internal/serve; every other tier passed. Every close-out on this host with LITELLM_BASE_URL and LITELLM_API_KEY set stops here until the model name or the check changes.

Correction, 2026-10-20: the diagnosis above is wrong. `unrecognized_model` is only a warning Claude Code prints for a model it has no metadata for, on passing runs too. The cause is upstream: OpenRouter answers 404 for `nvidia/nemotron-3-super-120b-a12b:free` through the proxy, as S-0451 found on TH-0529, and I-0170 records the same cause. `TestGatewayLiteLLM` passes with `FLAI_TEST_LITELLM_HAIKU=openrouter-messages/nvidia/nemotron-3-ultra-550b-a55b:free`. On the operator's word on TH-0530, S-0450's close-out ran with that model set for the run, while `.flai-cache/litellm/litellm.env` still names the retired one.

### 2026-10-21T01:31:48Z
Story: S-0454.
S-0454's close-out passed every tier but smoke, where TestGatewayLiteLLM/unified failed: Claude Code 2.1.296 refused the built-in default openrouter-messages/nvidia/nemotron-3-super-120b-a12b:free as unrecognized_model. Ran it again with FLAI_TEST_LITELLM_HAIKU set to nemotron-3-ultra for that run only, as TH-0534 settled for S-0445.

## Remediation

Story S-0455 remediates this issue, created from it at 2026-10-21T00:29:30Z.
