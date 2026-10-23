---
id: D-0343
title: The LiteLLM gateway smoke check's default free model is refused by Claude Code as unrecognized_model, so every close-out on this host fails smoke
kind: cause
sources: [I-0170]
topics: [flai, cli, go, code, all]
paths: [flai/internal/serve/gateway_test.go, scripts/litellm.sh, design/system/agent-adapters.md, docs/contributors/index.md, design/issues/I-0170-the-litellm-gateway-smoke-check-s-default-free-model-is-refused-by-claude-code-as-unrecognized-model-so-every-close-out-on-this-host-fails-smoke.md, design/issues/summary.md, design/issues/I-0176-a-usage-write-to-a-story-overwrites-its-move-to-in-progress-made-meanwhile-so-a-started-story-reads-as-ready-again.md, design/issues/I-0174-the-merge-queue-s-batch-runs-smoke-with-no-close-out-story-so-scripts-check-sh-runs-an-unscoped-flai-check-strict-and-fails-every-batch-on-findings-outside-its-stories.md, design/issues/I-0178-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-node-modules-without-it-so-every-other-close-out-fails-testeverylistedflaiovercomponentcompilesintheharness.md, design/issues/I-0158-testasettledreplystartstheagentwaitingonitonce-in-flai-internal-serve-fails-now-and-then-in-the-close-out-s-integration-tier-on-a-branch-that-does-not-touch-serve.md]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0343 The LiteLLM gateway smoke check's default free model is refused by Claude Code as unrecognized_model, so every close-out on this host fails smoke

I-0170 was closed on 2026-10-21.

The cause:

```text
The LiteLLM gateway smoke check's default free model is refused by Claude Code as unrecognized_model, so every close-out on this host fails smoke

Corrected 2026-10-20 on TH-0529: Claude Code does not refuse the model. OpenRouter's upstream answers 404 ("Function ... Not found for account") for `nvidia/nemotron-3-super-120b-a12b:free` through the proxy, and Claude Code reports that as `unrecognized_model`. It prints `unrecognized_model` for any model it has no metadata for, even on a passing run. `TestGatewayLiteLLM` passes with `FLAI_TEST_LITELLM_HAIKU=openrouter-messages/nvidia/nemotron-3-ultra-550b-a55b:free`, at $0 of key spend. A free model can be retired at any time, so the default needs a fallback or a skip when the proxy answers 404.
```

The fix:

```text
Fixed by S-0456 on S-0455's probe. TestGatewayLiteLLM no longer rests on one free model: it probes litellmHaikuDefaults (nemotron-3-ultra, minimax-m3.1-flash, nemotron-3.5-lightning, each passed the whole check live at $0 on 2026-10-21) in order and runs on the first the proxy answers, logging each passed over, so an OpenRouter retirement no longer fails smoke. An operator's own slug is probed alone and never replaced. The check fails, naming each slug and its answer, only when the whole list is refused, and never skips. TestLiteLLMPicksTheFirstFreeModelTheProxyAnswers reproduces the retired-first case against a stand-in proxy.
```

Source: I-0170.
