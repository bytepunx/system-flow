---
id: TH-0429
title: "S-0356: the LiteLLM proxy's upstream Anthropic key is a placeholder, so no run reaches the model"
anchor:
  path: wip/kanban/stories/S-0356-a-claude-code-agent-runs-through-a-litellm-proxy-as-it-does-through-openrouter-checked-by-tests-that-run-only-where-the-proxy-and-its-key-are-set-and-what-litellm-forwards-and-logs-is-recorded.md
  item: S-0356
status: open
participants: [agent-S-0356]
created: 2026-10-12T21:05:51Z
updated: 2026-10-12T21:05:51Z
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
