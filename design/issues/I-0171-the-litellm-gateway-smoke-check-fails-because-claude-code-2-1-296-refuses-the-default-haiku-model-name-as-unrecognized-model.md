---
id: I-0171
title: The LiteLLM gateway smoke check fails because Claude Code 2.1.296 refuses the default haiku model name as unrecognized_model
class: blocker
status: open
count: 1
cost: 15m
first_reported: 2026-10-20T20:07:32Z
last_reported: 2026-10-20T20:07:32Z
updated: 2026-10-20T20:07:32Z
---

# I-0171 The LiteLLM gateway smoke check fails because Claude Code 2.1.296 refuses the default haiku model name as unrecognized_model

## Description
The LiteLLM gateway smoke check fails because Claude Code 2.1.296 refuses the default haiku model name as unrecognized_model

## Instances

### 2026-10-20T20:07:32Z
Story: S-0450.
S-0450's close-out stopped at the smoke tier: TestGatewayLiteLLM/unified (flai/internal/serve/gateway_test.go:557) got `[claude-code:unrecognized_model] {"model":"openrouter-messages/nvidia/nemotron-3-super-120b-a12b:free"}` from Claude Code 2.1.296, the session ended in 0s with no model call, so flai-check.md was never written and no guard refusal was seen. It failed the same way when run again alone. The model is litellmHaikuDefault in gateway_test.go and LITELLM_HAIKU_MODEL in scripts/litellm.sh. S-0450 changes nothing under flai/internal/serve; every other tier passed. Every close-out on this host with LITELLM_BASE_URL and LITELLM_API_KEY set stops here until the model name or the check changes.

## Remediation
