---
id: I-0170
title: The LiteLLM gateway smoke check's default free model is refused by Claude Code as unrecognized_model, so every close-out on this host fails smoke
class: blocker
status: open
count: 1
cost: 15m
first_reported: 2026-10-20T19:17:06Z
last_reported: 2026-10-20T19:17:06Z
updated: 2026-10-20T19:17:06Z
---

# I-0170 The LiteLLM gateway smoke check's default free model is refused by Claude Code as unrecognized_model, so every close-out on this host fails smoke

## Description
The LiteLLM gateway smoke check's default free model is refused by Claude Code as unrecognized_model, so every close-out on this host fails smoke

## Instances

### 2026-10-20T19:17:06Z
Story: S-0451.
S-0451's close-out failed only at smoke, in TestGatewayLiteLLM/unified (flai/internal/serve/gateway_test.go:557-573): Claude Code 2.1.296 answered [claude-code:unrecognized_model] for openrouter-messages/nvidia/nemotron-3-super-120b-a12b:free, litellmHaikuDefault, and ended in 1 turn without the Write or the sub-agent's item_edit. S-0451 changes only flai/internal/itemnew and docs. Setting FLAI_TEST_LITELLM_HAIKU to a model Claude Code accepts, or changing litellmHaikuDefault, would clear it; I-0144 is the proxy's placeholder upstream key, a different cause.

## Remediation
