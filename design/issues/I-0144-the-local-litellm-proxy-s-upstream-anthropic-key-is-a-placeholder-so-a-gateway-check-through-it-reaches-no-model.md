---
id: I-0144
title: The local LiteLLM proxy's upstream Anthropic key is a placeholder, so a gateway check through it reaches no model
class: blocker
status: closed
count: 1
cost: 15m
first_reported: 2026-10-12T21:07:30Z
last_reported: 2026-10-12T21:07:30Z
updated: 2026-10-21T07:24:54Z
---

# I-0144 The local LiteLLM proxy's upstream Anthropic key is a placeholder, so a gateway check through it reaches no model

## Description
The local LiteLLM proxy's upstream Anthropic key is a placeholder, so a gateway check through it reaches no model

## Instances

### 2026-10-12T21:07:30Z
Story: S-0356.
S-0356's TestGatewayLiteLLM reached the proxy on 127.0.0.1:4000 with its virtual key, but the ANTHROPIC_API_KEY in .flai-cache/litellm/admin.env is a 4-character placeholder: Anthropic answered 401 invalid x-api-key on both the unified /v1/messages and the /anthropic pass-through. Claude Code retried the unified route for about 3 minutes before ending with "Not logged in". scripts/litellm.sh up does not check that the upstream key works; a check after `up` (one tiny call with the master key) would say so at once. Asked on TH-0429.

## Remediation
Closed 2026-10-21T07:24:54Z: Settled on TH-0429 (resolved 2026-10-13): S-0356 routes the gateway check's unified /v1/messages through the proxy to a free OpenRouter model with the host's OpenRouter key, and makes the Anthropic pass-through opt-in with FLAI_TEST_LITELLM_PASSTHROUGH=1. The placeholder ANTHROPIC_API_KEY no longer blocks any check; scripts/gateway-smoke.sh passed that way at 49a05cc, and agent-adapters.md records it.
