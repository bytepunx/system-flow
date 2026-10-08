---
id: TH-0368
title: "S-0339: the harness/provider/model split, which adapters first, and four decisions for the LiteLLM and OpenRouter epic"
anchor:
  path: wip/kanban/stories/S-0339-explore-adapters-for-litellm-and-openrouter.md
  item: S-0339
status: resolved
participants: [agent-S-0339, orchestrator, alex]
created: 2026-10-08T08:04:25Z
updated: 2026-10-08T08:37:05Z
---

# TH-0368 S-0339: the harness/provider/model split, which adapters first, and four decisions for the LiteLLM and OpenRouter epic

On wip/kanban/stories/S-0339-explore-adapters-for-litellm-and-openrouter.md.

## Entries

### 2026-10-08T08:04:25Z agent-S-0339
Recommendation: adopt the five answers below as recommended, and I will record each as an ADR and create the epic.

The finding is `design/system/agent-adapters.md` on `story/S-0339` (worktree `.flai-cache/worktrees/S-0339`, commit c3a7acb4). `## Today` lists where flai is tied to Claude Code per concern; `## LiteLLM and OpenRouter` sets out the two gateways and six harnesses (Claude Code over a gateway, Codex CLI, OpenCode, Goose, aider, a loop of flai's own) with what their documentation does not settle; `## Abstractions` gives the neutral contracts, the options compared, and the recommendation.

Two facts that shape everything: Anthropic and OpenRouter both say Claude Code over a gateway is supported for Claude models only, so non-Claude models need another harness; and only the gateway knows what a call cost, since Claude Code's `total_cost_usd` is its own estimate at Anthropic's prices.

## Questions

1. **The split.** Add `agent.provider` beside `harness` and `model`, naming an entry in a `providers` map (`api`, `base_url`, `key_env` naming the key's environment variable, optional model aliases; never a key). Recommended: the map lives on the host (`~/.flai/config.json`), as `agent.harnesses` does, because base URLs and variable names are host-bound. Alternative: in the manifest with a host override.
2. **Which adapters, in what order.** Recommended: (0) a spike story that runs `claude -p` through LiteLLM's `/v1/messages`, its `/anthropic` pass-through, and OpenRouter's `/api/v1/messages` with the guard and `permission_prompt` on, settling the beta-header and cost questions; (A) the provider split and Claude Code over a gateway, with the neutral contracts (`Capabilities`, a usage `Reader` per harness, `Call`/`Decide` in the guard, hold-and-ask); (B) a Codex CLI adapter, which unlocks OpenAI-shaped models through either gateway. OpenCode on demand; Goose and a loop of flai's own not now.
3. **A harness without a guard hook or a permission handler.** Recommended: its story's agent may run, its roles may not, unless the operator sets `agent.harnesses.<name>.guard: none` on the host; a protected write there is refused unless `auto-approve` is on. Alternative: refuse such a harness outright.
4. **Cost when a harness logs none, or estimates at the wrong prices.** Recommended: read the gateway's spend log (LiteLLM `/spend/logs` by virtual key, OpenRouter `/api/v1/key` and `/generation`) when the story has a provider, the harness's figure when it has none, an estimate otherwise, with the source named on the item's `usage` (`priced_by`). No price table. This changes `metrics.md`, so it is an ADR.
5. **A paid trial before the epic.** The spike needs a LiteLLM proxy (free, local) and an OpenRouter key with a few dollars of credit. Recommended: yes, both, inside the spike story; nothing is spent in S-0339.

Answer by number; "as recommended" covers them all. I will write one ADR per decision, refining ADR-0037, 0038, 0051, 0060, 0065, 0086, 0106, and 0124 as the finding's table says, add `## Decision` to the finding, and create the epic with no stories, for the planner.

### 2026-10-08T08:04:48Z orchestrator (recommendation)
Recommendation for the operator: answer "as recommended" on 1 to 4, and decide 5 yourself. All five are decisions not yet recorded, and 5 is spend, so they are yours. I am only recommending.

1. **Provider map on the host: yes.** ADR-0038 § Decision, which ADR-0041 leaves standing on this point, holds that "what runs, and what the agent may do, are the operator's alone", set on the host. `agent.harnesses` already lives there. Base URLs and key variable names are host-bound in the same way, and keeping them out of the manifest keeps them out of git. The manifest-with-override alternative would let anyone who edits the manifest point a story at a different endpoint.
2. **Order: spike (0), then A (provider split with Claude Code over a gateway and the neutral contracts), then B (Codex CLI): yes.** The spike settles the two facts the finding leaves open, the beta headers and who reports cost, before the contracts are fixed. A keeps today's harness and its guard, so it carries the least risk.
3. **A harness with no guard hook: off unless the host opts in, as recommended.** This follows the same ADR-0038 rule, permissions being the operator's. ADR-0060 is the guard, so roles stay refused. Refusing such a harness outright would rule out Codex CLI-style harnesses until each gets a hook.
4. **Cost from the gateway's spend log, with `priced_by` on usage: yes.** It needs an ADR, since `metrics.md` is the contract with the dashboard and changes only with one.
5. **Paid trial: your call.** It needs an OpenRouter key with a few dollars of credit, plus a local LiteLLM proxy, which is free. If you approve, cap it in the spike story, for example with `max_budget_usd` on its agent and a spend limit on the OpenRouter key.

For the agent: until the operator answers, it can draft the ADRs as `proposed` without deciding them.

Source: design/adrs/0038-flai-serve-starts-a-story-s-own-agent-through-an-adapter-with-what-the-operator.md § Decision

### 2026-10-08T08:29:00Z agent-S-0339
While this awaits you, the first four decisions are drafted as proposed ADRs on `story/S-0339`, as the orchestrator suggested: ADR-0129 (the split, `providers` on the host), ADR-0130 (the neutral contracts), ADR-0131 (a harness without a guard hook), ADR-0132 (cost from the gateway's spend log, `priced_by`). `design/system/agent-adapters.md` has a `## Decision` section linking them. On "as recommended" I set each to accepted and create the epic; on another answer I rewrite the ADR concerned. Question 5, the paid trial, is recorded in the epic's notes either way.

### 2026-10-08T08:30:29Z alex
1. take the recommendation but allow project specific overrides for the configuration values that aren't secrets 2. A 3. recommendation 4. yes 5. begin work to support open router but make it so that the openrouter integration/smoke tests do not work if there isn't a configured key for it.

### 2026-10-08T08:31:44Z agent-S-0339
Recorded. How I read each answer, so you can correct it here before acceptance:

1. ADR-0129 accepted, with a clause: a project's manifest may override an entry's non-secret values, `api`, `base_url`, and `models`; `key_env` stays the host's, since it names a variable of the host's environment, and the host must have the entry.
2. Option A only: the provider split and Claude Code over a gateway, with the neutral contracts (ADR-0130). No Codex, OpenCode, Goose, or loop of flai's own in this epic. The spike's two checks (what the gateway forwards, what its spend log says) become the first story's work against OpenRouter.
3. ADR-0131 accepted as recommended.
4. ADR-0132 accepted.
5. The epic starts with OpenRouter; LiteLLM follows. Its integration and smoke tests run only where an OpenRouter key is configured and do nothing where it is not, so no close-out or CI run without the key fails for want of it. I read "do not work" as "are skipped", not "fail".

Creating the epic next, with no stories, for the planner.

### 2026-10-08T08:36:04Z orchestrator (recommendation)
Recommendation for the operator: confirm the agent's readings of your five answers (08:31:44Z), or correct them here, then resolve this thread. It is the only thing blocking S-0339's acceptance.

- On 5, "skipped, saying so", not "fail", is the reading I recommend. A test that fails without a key would fail every close-out and CI run on a host without OpenRouter credit, which is exactly what S-0340 is removing for GitHub.
- S-0339 is in review. Verify passed at its head, ce4a34de. The verifier matched all three criteria:
  - E-0019, with nine criteria including key-gated tests;
  - ADR-0129 to ADR-0132, with ADR-0129 letting a manifest override `api`, `base_url` and `models` while `key_env` stays on the host;
  - `design/system/agent-adapters.md`, which has a Decision section and is indexed.
- No secret appears anywhere.
- `flai accept --dry-run` names one blocker: this thread is open.

Once it is resolved, I accept S-0339 and start E-0019's planner, which I have held back until the finding and the ADRs are on main.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md

### 2026-10-08T08:37:05Z alex
Resolved: S-0339 was accepted
