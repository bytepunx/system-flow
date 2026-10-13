---
id: TH-0429
title: "S-0356: the LiteLLM proxy's upstream Anthropic key is a placeholder, so no run reaches the model"
anchor:
  path: wip/kanban/stories/S-0356-a-claude-code-agent-runs-through-a-litellm-proxy-as-it-does-through-openrouter-checked-by-tests-that-run-only-where-the-proxy-and-its-key-are-set-and-what-litellm-forwards-and-logs-is-recorded.md
  item: S-0356
status: answered
participants: [agent-S-0356, orchestrator, alex]
created: 2026-10-12T21:05:51Z
updated: 2026-10-13T17:26:32Z
---

# TH-0429 S-0356: the LiteLLM proxy's upstream Anthropic key is a placeholder, so no run reaches the model

On wip/kanban/stories/S-0356-a-claude-code-agent-runs-through-a-litellm-proxy-as-it-does-through-openrouter-checked-by-tests-that-run-only-where-the-proxy-and-its-key-are-set-and-what-litellm-forwards-and-logs-is-recorded.md.

## Entries

### 2026-10-12T21:05:51Z agent-S-0356
**Recommendation:** put a real Anthropic key in the proxy, then answer here. I will run the check again and record it.

```sh
ANTHROPIC_API_KEY=sk-ant-... scripts/litellm.sh up
```

`up` writes the new key into `.flai-cache/litellm/admin.env` and recreates the container. It keeps the virtual key and its $5 budget.

## What happened

- `TestGatewayLiteLLM` ran once through each route on 2026-10-12. Claude Code reached the proxy, and the virtual key was accepted.
- The proxy's upstream `ANTHROPIC_API_KEY` in `admin.env` is a 4-character placeholder. Anthropic answered `401 invalid x-api-key` on both routes:
  - unified `/v1/messages`: Claude Code retried for 3 minutes, 11 calls, then ended with "Not logged in".
  - `/anthropic` pass-through: one call, then the same end.
- Nothing was spent: the key's spend stayed $0.

## Already learned, without a model reply

- On the unified route, every spend-log row's `session_id` equals Claude Code's session ID.
- On the pass-through, LiteLLM makes up its own `session_id`. Claude Code's session ID appears only in `end_user`.
- Only the master key reads `/spend/logs`; a virtual key gets 401.

## What waits on this

T-1404 needs a reply from the model: the run that answers which fields reach the model, whether the alias resolves, and how the spend compares with `total_cost_usd`. The S-0356 test code and the smoke step do not wait; I am finishing them now.

If you would rather not spend on Anthropic directly, say so. I will then record the two routes as checked up to the upstream, leave criterion 4 unticked, and say why in S-0356's notes.

### 2026-10-13T04:34:43Z orchestrator (recommendation)
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

### 2026-10-13T16:43:18Z alex
LiteLLM is not present to test an alternative way to spend money on Claude since the user could use Claude directly as they can today in system-flow. Can we identify a free/near-free alternative for testing that the LiteLLM integration works? We don't need anything robust, just a task that can be performed to ensure the integration is working across the board as expected.

### 2026-10-13T16:48:09Z agent-S-0356
**Recommendation:** route the proxy's unified route to one of OpenRouter's free models, using the OpenRouter key this host already has. Drop the `/anthropic` pass-through from the default run. Confirm, or say otherwise, and I will build it.

## What I checked (free, 2026-10-13)

- OpenRouter lists 14 free models that support tools. There is no local Ollama on this host.
- I ran a one-step `claude -p` straight through OpenRouter's Anthropic endpoint, asking for one `Write`. It cost $0.
  - `nvidia/nemotron-3-super-120b-a12b:free`: wrote the file; `is_error` false, 2 turns.
  - `poolside/laguna-s-2.1:free`: wrote the file.
  - `google/gemma-4-31b-it:free`: did not.

## The change

- `scripts/litellm.sh` gains an `openrouter/*` deployment upstream to OpenRouter, keyed by `OPENROUTER_API_KEY` in `admin.env`. Its first `up` no longer requires `ANTHROPIC_API_KEY`.
- `FLAI_TEST_LITELLM_HAIKU` defaults to `openrouter/nvidia/nemotron-3-super-120b-a12b:free`. That sends Claude Code through LiteLLM's translating `/v1/messages`, which is the part worth testing.
- The `/anthropic` pass-through forwards only to Anthropic, so no free model can sit behind it. It runs only when `FLAI_TEST_LITELLM_PASSTHROUGH=1` names a real Anthropic upstream, and skips saying so otherwise.
- The run still checks the guard, `permission_prompt`, the alias, the session join, and the spend log. The spend will read $0, so it shows only that rows are written and joined, not a price comparison.

## What it changes in S-0356

- **Criteria 1 and 4:** the pass-through becomes opt-in. agent-adapters.md records it as not run, with the reason.
- **Touches:** `scripts/litellm.sh` joins the story's touches.

## Risks

- Anthropic does not support Claude Code on non-Claude models.
- A free model may fumble the harder step: the sub-agent's `item_edit`, which the guard must refuse.
- Free models are rate-limited, and the list changes.
- If the model fails a step, I record what it did and do not count that as a LiteLLM defect.

## Alternative

LiteLLM's `mock_response` costs nothing, but it makes no tool calls. It would prove only that the proxy answers, not that the integration works.

### 2026-10-13T16:48:30Z orchestrator (recommendation)
**Recommendation:** confirm the agent's plan. Route the unified `/v1/messages` to `openrouter/nvidia/nemotron-3-super-120b-a12b:free` with the host's existing OpenRouter key, and make the `/anthropic` pass-through opt-in. It answers your ask: a near-free task that runs the integration end to end. It checks the guard, `permission_prompt`, the alias, the session join, and the spend-log rows. It spends nothing.

This changes S-0356's scope, so it is your decision, not mine:

- Criterion 1 says the test runs "once through `/anthropic`". It would run that route only when `FLAI_TEST_LITELLM_PASSTHROUGH=1` names a real Anthropic upstream, and skip saying so otherwise. Please reword criterion 1 to match, or tell the agent to.
- Criterion 4 would record the pass-through as not run, with the reason. Its "gateway's spend against `total_cost_usd`" reads $0, so it shows that rows are written and joined, not a price comparison.

**What it cannot show:** the design records that Anthropic does not support Claude Code on non-Claude models. So a free model's run proves what LiteLLM forwards and logs, not how a Claude model behaves behind it. A failed step by the free model is recorded as the model's, not as a LiteLLM defect, as the agent proposes.

**Alternative:** `mock_response` costs nothing but makes no tool calls. It would prove only that the proxy answers, so it does not meet criterion 1.

Source: design/system/agent-adapters.md § Claude Code over a gateway

### 2026-10-13T17:26:32Z alex
Confirmed the recommendation of 2026-10-13T16:48:30Z orchestrator.

Source: design/system/agent-adapters.md § Claude Code over a gateway
