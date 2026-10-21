---
id: T-1943
type: task
nature: remediation
title: The host's LiteLLM files take the current model, the gateway check passes at smoke, and I-0171 is closed
status: backlog
parent: S-0455
owner: alex
created: 2026-10-21T00:40:05Z
updated: 2026-10-21T00:40:05Z
transitions: []
stream: S-0455
tags: [litellm, issues]
touches: [design/system/agent-adapters.md, design/issues/I-0171-the-litellm-gateway-smoke-check-fails-because-claude-code-2-1-296-refuses-the-default-haiku-model-name-as-unrecognized-model.md, design/issues/summary.md]
after: [T-1941, T-1942]
---
# T-1943 The host's LiteLLM files take the current model, the gateway check passes at smoke, and I-0171 is closed

## Work

Show the fix on this host and close the issue. It waits for T-1941 and T-1942: the run needs both the script that drops the retired model and the probe that names one.

- Ask the operator on a thread on S-0455 before running `scripts/litellm.sh up`, since it rewrites the host's `.flai-cache/litellm/admin.env` and `litellm.env`. Then run it and check that both files name the current default, `openrouter-messages/nvidia/nemotron-3-ultra-550b-a55b:free` as S-0451 set it.
- Run `TestGatewayLiteLLM` through the smoke tier with no `FLAI_TEST_LITELLM_HAIKU` set by hand, and record the result.
- In `design/system/agent-adapters.md`, under the LiteLLM gateway check, say that the check probes its model first, and that `litellm.sh` keeps only a model the operator chose.
- Close I-0171 with `flai issue close I-0171 --reason`, naming T-1941 and T-1942. Keep `design/issues/summary.md` current.

## Done when

- The smoke tier passes on this host with the host's own `litellm.env`.
- `agent-adapters.md` describes the probe and the model rule, and the markdown lint passes on it.
- I-0171 is closed with its reason, and `flai check --strict` is clean.

## Notes
