---
id: TH-0548
title: "S-0456 plan: three tasks in two layers, waits for S-0455; recommend a fallback list of free models that fails, never skips, when the whole list is refused"
anchor:
  path: wip/kanban/stories/S-0456-the-litellm-gateway-smoke-check-s-default-free-model-is-refused-by-claude-code-as-unrecognized-model-so-every-close-out-on-this-host-fails-smoke.md
  item: S-0456
status: resolved
participants: [planner-S-0456, orchestrator]
created: 2026-10-21T06:05:00Z
updated: 2026-10-21T06:05:19Z
---

# TH-0548 S-0456 plan: three tasks in two layers, waits for S-0455; recommend a fallback list of free models that fails, never skips, when the whole list is refused

On wip/kanban/stories/S-0456-the-litellm-gateway-smoke-check-s-default-free-model-is-refused-by-claude-code-as-unrecognized-model-so-every-close-out-on-this-host-fails-smoke.md.

## Entries

### 2026-10-21T06:05:00Z planner-S-0456
**Recommended fix:** make the LiteLLM check's default an ordered list of free models, verified live, and run the check on the first one the proxy answers. Fail, and never skip, when the proxy refuses every model on the list. Confirm this, or answer "skip" if you would rather the check skip with a logged reason when every model is refused.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1944 TestGatewayLiteLLM falls back through a list of free models | `flai/internal/serve/gateway_test.go`, `scripts/litellm.sh` | none |
| 2 | T-1945 The design and the contributors guide describe the fallback | `design/system/agent-adapters.md`, `docs/contributors/index.md` | T-1944 |
| 2 | T-1946 The gateway check passes at smoke, and I-0170 is closed | `design/issues/I-0170-…md`, `design/issues/summary.md` | T-1944 |

T-1945 and T-1946 have no path in common, so they can run together.

## Assumptions

- S-0451 has already moved the default to nemotron-3-ultra. S-0455, in progress, makes the check probe its model and fail at once on a 404. What I-0170 still asks for is a default that survives the next retirement, so this story builds on S-0455's probe. The story has `after: S-0455`, since the two stories change the same five files.
- The proxy routes `openrouter-messages/*` by wildcard, so adding a model to the list needs no change to the proxy.
- At least one more free OpenRouter model passes the whole session (the Write, the refused sub-agent `item_edit`, the alias, and a result) at $0. T-1944 finds it live. If none does, the list stays at one model and the story's agent says so on this thread.
- A slug the operator sets that is not on the list is never replaced. The check probes it alone and fails as S-0455 made it.
- Failing rather than skipping keeps a close-out from passing with no gateway check and nothing to say so. The orchestrator called a skip an unrecorded decision on TH-0534.

## Figures

- Forecast: 1h, raised from flai's 15m to cover the live runs through the proxy. Delivery 2026-10-21T21:58:00Z.
- Cost of delay: 37.50 USD a week, from flai's input of 15m a cycle. Kept as given.
- Touches: six files, no folder. Where each one comes from is in the story's Notes under Planning.

Nothing in this plan proposes splitting, merging, or dropping a task.

### 2026-10-21T06:05:17Z orchestrator
**Approved, as drafted.**

- **The fix fits the goal.** I-0170 asks for a default that survives the next model retirement. An ordered, live-verified list of free models, with the check running on the first one the proxy answers, does that on top of S-0455's probe. A slug the operator sets is never replaced.
- **Fail, never skip: agreed.** That is the check's behaviour today, so keeping it needs no new decision. A skip would be one: a close-out that passes with no gateway check, which TH-0534 already treated as unrecorded. If the operator wants a skip, they can say so on the story before T-1944 starts.
- **Tasks and layers.** T-1944 first. Then T-1945 (docs) and T-1946 (smoke, and closing I-0170) run in parallel, since they share no path.
- **`after: S-0455`: right.** Both change the same files.
- **Figures.** Six file touches. The 1h forecast against flai's 15m is justified by the live proxy runs. The 37.50 USD a week value is kept as flai's 15m input gives it.
- **If no second free model passes the full session:** the list stays at one, as you say. The story's agent records that, and it still meets criterion 1 for this host.

### 2026-10-21T06:05:19Z orchestrator
Resolved: Plan approved as drafted under plan_backlog_stories: three tasks in two layers, after S-0455, 1h forecast, and 37.50 USD a week. "Fail, never skip" keeps the check's existing behaviour.
