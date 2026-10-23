---
id: D-0287
title: "S-0356: the LiteLLM proxy's upstream Anthropic key is a placeholder, so no run reaches the model: LiteLLM is not present to test an alternative way to spend money on Claude since the user could use…"
kind: recommendation
sources: [TH-0429]
topics: [agents, testing, cli, flai, go, code, all]
paths: [flai/internal/serve/gateway_test.go, scripts/gateway-smoke.sh, docs/contributors/index.md, docs/operators/index.md, design/system/agent-adapters.md, design/issues/I-0144-the-local-litellm-proxy-s-upstream-anthropic-key-is-a-placeholder-so-a-gateway-check-through-it-reaches-no-model.md, design/issues/summary.md, design/issues/I-0125-a-settled-conversation-between-two-stories-agents-keeps-waking-each-in-turn-because-every-reply-awaits-the-other.md, scripts/litellm.sh]
date: 2026-10-13
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0287 S-0356: the LiteLLM proxy's upstream Anthropic key is a placeholder, so no run reaches the model: LiteLLM is not present to test an alternative way to spend money on Claude since the user could use…

The operator, alex, answered TH-0429 at 2026-10-13T16:43:18Z.

What was asked:

```text
**Recommendation:** put a real Anthropic key in the proxy, then answer here. I will run the check again and record it.
```

The recommendation it overruled:

```text
**Recommendation:** put a real Anthropic key in the proxy, as S-0356's agent proposes. Pass it only in the environment: `ANTHROPIC_API_KEY=… scripts/litellm.sh up`. Then answer here so the agent reruns `TestGatewayLiteLLM`.

This is your call, not mine. It calls a paid service with a credential only you hold, and the safety rules need your confirmation for both. The spend is small and capped:

- The virtual key's budget stops it at $5.
- One run through each route is a handful of calls.
- Criterion 4 of S-0356 needs a reply from the model to be checked at all.

Do not paste the key into this thread or any work item. Secrets go in the environment only, and `.flai-cache/litellm/admin.env` is git-ignored.

**Alternative:** if you would rather not spend on Anthropic directly, say so. The agent will then:

- record both routes as checked up to the upstream;
- leave criterion 4 unticked and say why in S-0356's notes.

The story could still be accepted with that criterion out of scope only by your decision.

Source: design/conventions/safety.md § Rules
```

The operator's answer:

```text
LiteLLM is not present to test an alternative way to spend money on Claude since the user could use Claude directly as they can today in system-flow. Can we identify a free/near-free alternative for testing that the LiteLLM integration works? We don't need anything robust, just a task that can be performed to ensure the integration is working across the board as expected.
```

Source: TH-0429.
