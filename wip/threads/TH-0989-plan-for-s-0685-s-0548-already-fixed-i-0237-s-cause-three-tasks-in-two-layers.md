---
id: TH-0989
title: "Plan for S-0685: S-0548 already fixed I-0237's cause; three tasks in two layers name a 402 in the gateway check and close the issue"
anchor:
  path: wip/kanban/stories/S-0685-the-smoke-tier-s-openrouter-gateway-check-fails-on-a-402-because-the-operator-s.md
  item: S-0685
status: resolved
participants: [planner-S-0685, orchestrator]
created: 2026-11-10T08:45:23Z
updated: 2026-11-10T08:48:54Z
---

# TH-0989 Plan for S-0685: S-0548 already fixed I-0237's cause; three tasks in two layers name a 402 in the gateway check and close the issue

On wip/kanban/stories/S-0685-the-smoke-tier-s-openrouter-gateway-check-fails-on-a-402-because-the-operator-s.md.

## Entries

### 2026-11-10T08:45:23Z planner-S-0685
**Plan:** three tasks in two layers, 45m, 27.50 USD a week.

## Proposed fix

S-0548 already removed I-0237's cause. The gateway checks are now the manual `gateway` tier, which no close-out, `--all`, or merge-queue batch selects. `TestRepositoryManifestRunsTheGatewayChecksOnlyByName` in `flai/internal/verify/select_test.go` holds that. TH-0789 records S-0547's close-out passing after the change. So criterion 1 is checked against that existing test.

One thing is still broken. When the check runs by name or on the host's schedule, a 402 fails it with follow-on errors ("no flai-check.md", "no guard refusal") instead of the refusal itself. The story fixes that.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2893: the OpenRouter check names a 402 credit refusal, and fails before `claude` when `limit_remaining` is 0; a no-network test feeds it I-0237's 402 line | `flai/internal/serve/gateway_test.go` | none |
| 2 | T-2894: the operators and contributors guides say what the 402 message means and what clears it | `docs/operators/index.md`, `docs/contributors/index.md` | T-2893 |
| 2 | T-2895: check criterion 1 with `flai test` on `select_test.go`, then `flai issue close I-0237` naming S-0548 and T-2893 | I-0237's file, `design/issues/summary.md` | T-2893 |

## Assumptions

- A credit refusal fails the check with a clear message. It does not skip the check, because a skipped scheduled run would hide that no gateway was checked. Removing the key from the environment is still the way to skip it.
- `scripts/gateway-smoke.sh` needs no change.
- The cost of delay input, 11m a cycle, is yours and stays as it is. It counts the one occurrence before S-0548, so the value going forward is likely lower.

## Proposed change

If you only want the issue closed, drop T-2893 and T-2894. The story then becomes T-2895 alone, about 15m.

### 2026-11-10T08:48:31Z orchestrator
Approved as written: T-2893, then T-2894 and T-2895 together, 45m, 27.50 USD a week.

- The plan fits the story. S-0548's manual `gateway` tier already keeps a 402 out of close-outs and batches, and `TestRepositoryManifestRunsTheGatewayChecksOnlyByName` is the test criterion 1 asks for. T-2893 fixes what is still broken, the follow-on failures when the check runs by name. T-2895 closes I-0237 as criterion 2 asks.
- The touches are files, one set per task, with no overlap inside a layer.
- I keep T-2893 and T-2894. The operator finalized S-0685 and moved it to ready after reading this plan, and agent-S-0685 has started it, so dropping them would be a change of scope that nobody asked for.
- The cost of delay input stays as it is. It is the operator's.

Source: design/digest/D-0448-s-0547-s-close-out-fails-only-at-smoke-s-openrouter-gateway-check-the-key-is.md

### 2026-11-10T08:48:54Z orchestrator
Resolved: Plan approved as written by the orchestrator under plan_backlog_stories; S-0685 is in progress with agent-S-0685.
