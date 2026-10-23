---
id: D-0336
title: The LiteLLM gateway smoke check fails because Claude Code 2.1.296 refuses the default haiku model name as unrecognized_model
kind: cause
sources: [I-0171]
topics: [flai, cli, go, code, all]
paths: [scripts/litellm.sh, scripts/README.md, flai/internal/serve/gateway_test.go, docs/contributors/index.md, design/system/agent-adapters.md, design/issues/I-0171-the-litellm-gateway-smoke-check-fails-because-claude-code-2-1-296-refuses-the-default-haiku-model-name-as-unrecognized-model.md, design/issues/summary.md, design/issues/I-0173-testachildisstoppedonlyonceitsprocesshasended-fails-under-load-when-the-stubborn-child-dies-on-sigterm-before-it-ignores-it.md, design/issues/I-0174-the-merge-queue-s-batch-runs-smoke-with-no-close-out-story-so-scripts-check-sh-runs-an-unscoped-flai-check-strict-and-fails-every-batch-on-findings-outside-its-stories.md]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0336 The LiteLLM gateway smoke check fails because Claude Code 2.1.296 refuses the default haiku model name as unrecognized_model

I-0171 was closed on 2026-10-21.

The cause:

```text
The LiteLLM gateway smoke check fails because Claude Code 2.1.296 refuses the default haiku model name as unrecognized_model
```

The fix:

```text
S-0455. The cause was OpenRouter's 404 for the retired free nemotron-3-super, not Claude Code's unrecognized_model warning. T-1941: scripts/litellm.sh keeps LITELLM_HAIKU_MODEL in admin.env only when the operator sets it, and drops a retired default on any run, so a host follows the script's current default (nemotron-3-ultra, S-0451). T-1942: TestGatewayLiteLLM probes its model through the proxy before the Claude Code session and fails naming the slug, its source, and the upstream answer; TestLiteLLMProbeNamesTheModel reproduces the 404. On 2026-10-21 the host's files were rewritten on the operator's word (TH-0547), and the gateway check passed from the host's own litellm.env.
```

Source: I-0171.
