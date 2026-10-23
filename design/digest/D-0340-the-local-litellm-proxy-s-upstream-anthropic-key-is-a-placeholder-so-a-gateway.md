---
id: D-0340
title: The local LiteLLM proxy's upstream Anthropic key is a placeholder, so a gateway check through it reaches no model
kind: cause
sources: [I-0144]
topics: [agents, testing, cli, flai, go, code, all]
paths: [flai/internal/serve/gateway_test.go, scripts/gateway-smoke.sh, docs/contributors/index.md, docs/operators/index.md, design/system/agent-adapters.md, design/issues/I-0144-the-local-litellm-proxy-s-upstream-anthropic-key-is-a-placeholder-so-a-gateway-check-through-it-reaches-no-model.md, design/issues/summary.md, design/issues/I-0125-a-settled-conversation-between-two-stories-agents-keeps-waking-each-in-turn-because-every-reply-awaits-the-other.md, scripts/litellm.sh]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0340 The local LiteLLM proxy's upstream Anthropic key is a placeholder, so a gateway check through it reaches no model

I-0144 was closed on 2026-10-21.

The cause:

```text
The local LiteLLM proxy's upstream Anthropic key is a placeholder, so a gateway check through it reaches no model
```

The fix:

```text
Settled on TH-0429 (resolved 2026-10-13): S-0356 routes the gateway check's unified /v1/messages through the proxy to a free OpenRouter model with the host's OpenRouter key, and makes the Anthropic pass-through opt-in with FLAI_TEST_LITELLM_PASSTHROUGH=1. The placeholder ANTHROPIC_API_KEY no longer blocks any check; scripts/gateway-smoke.sh passed that way at 49a05cc, and agent-adapters.md records it.
```

Source: I-0144.
