---
id: D-0256
title: "S-0339: the harness/provider/model split, which adapters first, and four decisions for the LiteLLM and OpenRouter epic: take the recommendation but allow project specific overrides for the…"
kind: recommendation
sources: [TH-0368]
topics: [cli, template, all]
paths: [design/system/agent-adapters.md, design/system/README.md, design/adrs, design/adrs/README.md, design/issues/I-0123-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md, design/issues/summary.md]
date: 2026-10-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0256 S-0339: the harness/provider/model split, which adapters first, and four decisions for the LiteLLM and OpenRouter epic: take the recommendation but allow project specific overrides for the…

The operator, alex, answered TH-0368 at 2026-10-08T08:30:29Z.

What was asked:

```text
Recommendation: adopt the five answers below as recommended, and I will record each as an ADR and create the epic.
```

The recommendation it overruled:

```text
Recommendation for the operator: answer "as recommended" on 1 to 4, and decide 5 yourself. All five are decisions not yet recorded, and 5 is spend, so they are yours. I am only recommending.

1. **Provider map on the host: yes.** ADR-0038 § Decision, which ADR-0041 leaves standing on this point, holds that "what runs, and what the agent may do, are the operator's alone", set on the host. `agent.harnesses` already lives there. Base URLs and key variable names are host-bound in the same way, and keeping them out of the manifest keeps them out of git. The manifest-with-override alternative would let anyone who edits the manifest point a story at a different endpoint.
2. **Order: spike (0), then A (provider split with Claude Code over a gateway and the neutral contracts), then B (Codex CLI): yes.** The spike settles the two facts the finding leaves open, the beta headers and who reports cost, before the contracts are fixed. A keeps today's harness and its guard, so it carries the least risk.
3. **A harness with no guard hook: off unless the host opts in, as recommended.** This follows the same ADR-0038 rule, permissions being the operator's. ADR-0060 is the guard, so roles stay refused. Refusing such a harness outright would rule out Codex CLI-style harnesses until each gets a hook.
4. **Cost from the gateway's spend log, with `priced_by` on usage: yes.** It needs an ADR, since `metrics.md` is the contract with the dashboard and changes only with one.
5. **Paid trial: your call.** It needs an OpenRouter key with a few dollars of credit, plus a local LiteLLM proxy, which is free. If you approve, cap it in the spike story, for example with `max_budget_usd` on its agent and a spend limit on the OpenRouter key.

For the agent: until the operator answers, it can draft the ADRs as `proposed` without deciding them.

Source: design/adrs/0038-flai-serve-starts-a-story-s-own-agent-through-an-adapter-with-what-the-operator.md § Decision
```

The operator's answer:

```text
1. take the recommendation but allow project specific overrides for the configuration values that aren't secrets 2. A 3. recommendation 4. yes 5. begin work to support open router but make it so that the openrouter integration/smoke tests do not work if there isn't a configured key for it.
```

Source: TH-0368.
