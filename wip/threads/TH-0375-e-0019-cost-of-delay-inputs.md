---
id: TH-0375
title: E-0019 cost of delay inputs
anchor:
  path: wip/kanban/epics/E-0019-agent-adapters-for-litellm-and-openrouter.md
  item: E-0019
status: answered
participants: [planner-E-0019, alex]
created: 2026-10-08T08:39:05Z
updated: 2026-10-08T08:46:18Z
---

# TH-0375 E-0019 cost of delay inputs

On wip/kanban/epics/E-0019-agent-adapters-for-litellm-and-openrouter.md.

## Entries

### 2026-10-08T08:39:05Z planner-E-0019
E-0019 has no cost of delay inputs, so `flai cod` refuses it and every story under it. The inputs are yours.

**Recommendation:** `penalty_per_week: 250` USD, set with `flai edit E-0019 --penalty-per-week 250`.

Why 250:

- Today every agent runs on Claude Code against Anthropic direct. The epic adds a choice of gateway. It does not unblock work already on the board.
- OpenRouter charges the provider's rate plus a 5.5% fee on credits (agent-adapters.md § OpenRouter), so Claude over it is not cheaper. The gain is a self-hosted LiteLLM in front of Bedrock or Vertex credits, and a key and spend log per project.
- That puts it well below E-0018 (1000 a week) and above E-0015 (25 a week).

Alternatives:

1. `penalty_per_week: 1000`, if gateway credits or a spend cap per project are pressing now.
2. `penalty_per_week: 100`, if this is exploratory and can wait behind E-0015 and E-0018.
3. `revenue_per_week: <n>`, if you expect a weekly saving from cheaper credits: the saving in USD.

I am drafting the stories and tasks meanwhile. Each story's cost of delay value waits for this answer.

### 2026-10-08T08:46:18Z alex
I set the penalty to 500 because with the Claude subscription cap, the only way to continue each week is to pay for very expensive API tokens.
