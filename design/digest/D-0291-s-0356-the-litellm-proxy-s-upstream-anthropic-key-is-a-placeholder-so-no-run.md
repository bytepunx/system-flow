---
id: D-0291
title: "S-0356: the LiteLLM proxy's upstream Anthropic key is a placeholder, so no run reaches the model: Recommendation: confirm the agent's plan"
kind: answer
sources: [TH-0429]
topics: [agents, testing, cli, flai, go, code, all]
paths: [flai/internal/serve/gateway_test.go, scripts/gateway-smoke.sh, docs/contributors/index.md, docs/operators/index.md, design/system/agent-adapters.md, design/issues/I-0144-the-local-litellm-proxy-s-upstream-anthropic-key-is-a-placeholder-so-a-gateway-check-through-it-reaches-no-model.md, design/issues/summary.md, design/issues/I-0125-a-settled-conversation-between-two-stories-agents-keeps-waking-each-in-turn-because-every-reply-awaits-the-other.md, scripts/litellm.sh]
date: 2026-10-13
status: absorbed
absorbed_by: "design/system/agent-adapters.md#checked-against-litellm"
asked_again: []
cited: []
---

# D-0291 S-0356: the LiteLLM proxy's upstream Anthropic key is a placeholder, so no run reaches the model: Recommendation: confirm the agent's plan

The operator, alex, answered TH-0429 at 2026-10-13T17:26:32Z.

What was asked:

```text
**Recommendation:** put a real Anthropic key in the proxy, then answer here. I will run the check again and record it.
```

The answer:

```text
**Recommendation:** confirm the agent's plan. Route the unified `/v1/messages` to `openrouter/nvidia/nemotron-3-super-120b-a12b:free` with the host's existing OpenRouter key, and make the `/anthropic` pass-through opt-in. It answers your ask: a near-free task that runs the integration end to end. It checks the guard, `permission_prompt`, the alias, the session join, and the spend-log rows. It spends nothing.

This changes S-0356's scope, so it is your decision, not mine:

- Criterion 1 says the test runs "once through `/anthropic`". It would run that route only when `FLAI_TEST_LITELLM_PASSTHROUGH=1` names a real Anthropic upstream, and skip saying so otherwise. Please reword criterion 1 to match, or tell the agent to.
- Criterion 4 would record the pass-through as not run, with the reason. Its "gateway's spend against `total_cost_usd`" reads $0, so it shows that rows are written and joined, not a price comparison.

**What it cannot show:** the design records that Anthropic does not support Claude Code on non-Claude models. So a free model's run proves what LiteLLM forwards and logs, not how a Claude model behaves behind it. A failed step by the free model is recorded as the model's, not as a LiteLLM defect, as the agent proposes.

**Alternative:** `mock_response` costs nothing but makes no tool calls. It would prove only that the proxy answers, so it does not meet criterion 1.

Source: design/system/agent-adapters.md § Claude Code over a gateway
```

Source: TH-0429.
