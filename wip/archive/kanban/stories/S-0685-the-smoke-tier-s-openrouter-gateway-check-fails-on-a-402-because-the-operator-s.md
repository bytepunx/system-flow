---
id: S-0685
type: story
nature: remediation
title: The smoke tier's OpenRouter gateway check fails on a 402 because the operator's OpenRouter key is out of credits, so every close-out and batch that runs smoke fails
status: done
owner: alex
created: 2026-11-10T08:37:11Z
updated: 2026-11-10T09:33:15Z
transitions:
  - to: ready
    at: 2026-11-10T08:47:51Z
    by: alex
  - to: in-progress
    at: 2026-11-10T08:48:20Z
    by: agent-S-0685
  - to: review
    at: 2026-11-10T09:28:20Z
    by: agent-S-0685
  - to: done
    at: 2026-11-10T09:33:15Z
    by: orchestrator
tags: [gateway, tests]
touches: [flai/internal/serve/gateway_test.go, docs/operators/index.md, docs/contributors/index.md, design/issues/I-0237-the-smoke-tier-s-openrouter-gateway-check-fails-on-a-402-because-the-operator-s.md, design/issues/summary.md, design/issues/I-0334-testgatewaystreamisjudged-and-testgatewaymissingnameswhatismissing-run-in-no.md, design/issues/I-0337-flai-check-finds-adr-index-outside-the-story-at-close-out.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2412
  priced_by: harness
  turns:
    - day: 2026-11-10
      ceremony: 3
      hand_edits: 1
      work: 90
  models:
    - model: claude-opus-5-5
      input: 196
      output: 23317
      cache_read: 27779932
      cache_write: 325841
      cost: 8.6298
  strategic:
    - kind: orchestrator
      seconds: 85
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 34
          output: 219
          cache_read: 1270853
          cache_write: 153217
          cost: 0.3857
cost_of_delay:
  inputs:
    time_lost_per_cycle: 11m
    by: flai
    at: 2026-11-10T08:37:11Z
  value: 27.5
  by: planner-S-0685
  at: 2026-11-10T08:45:09Z
forecast:
  duration: 45m
  delivery: 2026-11-10T12:28:00Z
  basis: flai forecast's 16m (131 s per unit over 50 done remediation stories, size 7) raised to 45m because a change under flai/ runs the integration and smoke tiers in its close-out, which took 3 to 10 minutes on S-0547
  by: planner-S-0685
  at: 2026-11-10T08:45:09Z
finalized:
  by: alex
  at: 2026-11-10T08:47:48Z
---
# S-0685 The smoke tier's OpenRouter gateway check fails on a 402 because the operator's OpenRouter key is out of credits, so every close-out and batch that runs smoke fails

## Goal

This story remediates [I-0237](../../../design/issues/I-0237-the-smoke-tier-s-openrouter-gateway-check-fails-on-a-402-because-the-operator-s.md), "The smoke tier's OpenRouter gateway check fails on a 402 because the operator's OpenRouter key is out of credits, so every close-out and batch that runs smoke fails". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0237 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0237 is closed with `flai issue close I-0237 --reason` saying what fixed it

## Tasks
- T-2893 The OpenRouter gateway check says the key is out of credit when OpenRouter answers 402, before the checks that follow from it
- T-2894 The operators and contributors guides say what the gateway check prints on a 402 and what to do about it
- T-2895 Close I-0237 naming S-0548's manual gateway tier and T-2893's credit message as what fixed it

## Notes

Cost of delay inputs set by flai from I-0237. time_lost_per_cycle 11m: 15m per occurrence × 1 occurrence ÷ 1.4 cycles of 168h (first reported 2026-10-31T10:10:59Z, 9.9 days before this story).

### Planning

Proposed fix. S-0548 already removed I-0237's cause. The gateway checks are now the manual `gateway` tier in `system-flow.yaml`, which no path, close-out, `--all`, or merge-queue batch selects. `TestRepositoryManifestRunsTheGatewayChecksOnlyByName` in `flai/internal/verify/select_test.go` holds that, and TH-0789 records S-0547's close-out passing after it. What remains is the check itself, run by name or on the host's schedule: on a 402 it fails with follow-on errors, not with the refusal. T-2893 makes it name the refusal. T-2894 documents that. T-2895 checks criterion 1 against the existing test and closes I-0237.

Touches:

| Touch | From |
|-------|------|
| `flai/internal/serve/gateway_test.go` | layout: `checkGateway`, `keyNow`, and `TestGatewayOpenRouter` live there (T-2893) |
| `docs/operators/index.md` | co-change with `gateway_test.go` and `scripts/gateway-smoke.sh` (2 of 11 commits); its gateway section (T-2894) |
| `docs/contributors/index.md` | co-change (5 of 11 commits); its "Gateway checks" section (T-2894) |
| `design/issues/I-0237-…md` | design: criterion 2, `flai issue close` (T-2895) |
| `design/issues/summary.md` | layout: `flai issue close` updates the summary (T-2895) |

Left out of co-change: `scripts/litellm.sh` and `design/system/agent-adapters.md`, which the LiteLLM check and the adapters change, not a 402 message. `scripts/gateway-smoke.sh` needs no change: it already skips a gateway whose key is unset.

Forecast: `flai forecast` gave 16m from 131 s per unit over 50 done remediation stories, size 7. I raised it to 45m. A change under `flai/` runs integration and smoke in its close-out, which took 3 to 10 minutes on S-0547. Delivery was moved later by the same 29 minutes.

Cost of delay: 27.50 USD a week, from `flai cod`, kept as computed from the operator's input, 11m lost per 168h cycle at 150 USD an hour. The input counts the one occurrence on 2026-10-31. No occurrence has followed since S-0548, so the real figure going forward is likely lower. The input is the operator's, so I left it as it is.

No `after`: the story reads S-0548's tier, which is already on main.

### Accepted by the orchestrator

- Verified: 9441bbbed91a631b585fa7f8477fbf80df720f6e
- At: 2026-11-10T09:33:15Z

Verdict: pass
- 1: flai/internal/serve/gateway_test.go, docs/operators/index.md, docs/contributors/index.md
- 2: design/issues/I-0237-the-smoke-tier-s-openrouter-gateway-check-fails-on-a-402-because-the-operator-s.md, design/issues/summary.md
