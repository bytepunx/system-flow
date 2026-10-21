---
id: S-0456
type: story
nature: remediation
title: The LiteLLM gateway smoke check's default free model is refused by Claude Code as unrecognized_model, so every close-out on this host fails smoke
status: ready
owner: alex
created: 2026-10-21T02:08:41Z
updated: 2026-10-21T07:33:37Z
transitions:
  - to: ready
    at: 2026-10-21T07:28:20Z
    by: alex
tags: [flai, serve, tests]
touches: [flai/internal/serve/gateway_test.go, scripts/litellm.sh, design/system/agent-adapters.md, docs/contributors/index.md, design/issues/I-0170-the-litellm-gateway-smoke-check-s-default-free-model-is-refused-by-claude-code-as-unrecognized-model-so-every-close-out-on-this-host-fails-smoke.md, design/issues/summary.md]
after: [S-0455]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 9911
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 167
          output: 3051
          cache_read: 28729430
          cache_write: 46798
          cost: 7.2145
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-21T02:08:41Z
  value: 37.5
  by: planner-S-0456
  at: 2026-10-21T06:04:37Z
forecast:
  duration: 1h
  delivery: 2026-10-21T08:41:00Z
  basis: "Its own forecast of 1h; 2nd in the pull order with an in-progress limit of 5, behind S-0453, S-0464 and S-0463."
  by: flai
  at: 2026-10-21T07:33:37Z
finalized:
  by: alex
  at: 2026-10-21T05:58:31Z
---
# S-0456 The LiteLLM gateway smoke check's default free model is refused by Claude Code as unrecognized_model, so every close-out on this host fails smoke

## Goal

This story remediates [I-0170](../../../design/issues/I-0170-the-litellm-gateway-smoke-check-s-default-free-model-is-refused-by-claude-code-as-unrecognized-model-so-every-close-out-on-this-host-fails-smoke.md), "The LiteLLM gateway smoke check's default free model is refused by Claude Code as unrecognized_model, so every close-out on this host fails smoke". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0170 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0170 is closed with `flai issue close I-0170 --reason` saying what fixed it

## Tasks
- T-1944 TestGatewayLiteLLM falls back through a list of free models when the proxy answers its default with a refusal
- T-1945 The design and the contributors guide say the LiteLLM check falls back through its list of free models
- T-1946 The gateway check passes at smoke on this host with the fallback, and I-0170 is closed

## Notes

Cost of delay inputs set by flai from I-0170. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-20T19:17:06Z, 0.3 days before this story; under one cycle counts as one).

### Planning

The fix proposed, from I-0170's corrected instance and TH-0529: the cause is OpenRouter retiring a free model, which the proxy answers with a 404, not Claude Code refusing it. S-0451 moved the default to `nemotron-3-ultra-550b-a55b:free`. S-0455, in progress, makes the check probe its model first and fail at once, naming the slug, on a refusal. What I-0170 still asks for is a default that survives the next retirement. T-1944 makes the default an ordered list of free models, verified live, and runs the check on the first the proxy answers. An operator's own slug is not replaced. When the whole list is refused the check fails, never skips, so a close-out never passes without the gateway check and nothing saying so. T-1945 describes it, and T-1946 shows it at smoke on this host and closes I-0170.

`after: S-0455`: its T-1942 adds the probe to `gateway_test.go`, and it changes `scripts/litellm.sh`, `design/system/agent-adapters.md`, `docs/contributors/index.md`, and `design/issues/summary.md`, all of which this story changes again.

Touches, all files, no folder:

| Touch | Source | Task |
|-------|--------|------|
| `flai/internal/serve/gateway_test.go` | design (I-0170 names `TestGatewayLiteLLM` in it) | T-1944 |
| `scripts/litellm.sh` | co-change (4 of 8 commits) | T-1944 |
| `design/system/agent-adapters.md` | co-change (2 of 8 commits) | T-1945 |
| `docs/contributors/index.md` | co-change (3 of 8 commits) | T-1945 |
| `design/issues/I-0170-…md` | criterion 2 | T-1946 |
| `design/issues/summary.md` | layout (kept with every issue change) | T-1946 |

`flai touches suggest` listed nothing at first, since the story declared no touches. Run from `gateway_test.go` and the issue file, it gave the three co-change paths above and `docs/operators/index.md` (2 of 8). I left that one out: the fallback changes no operator setting.

Forecast: flai gave 5m with no touches, then 15m with them (110 s per unit of size over 24 done remediation stories, size 8). Raised to 1h. T-1944 tries free models live through the proxy, about 150 s a run, and T-1946 runs the smoke check live again. S-0455, with the same kind of live runs, was forecast 45m. Delivery 2026-10-21T21:58:00Z, flai's 21:13 plus the 45m added; it is 29th in the pull order.

Cost of delay: 37.50 USD a week, as `flai cod` works it out from flai's input, 15m lost per 168h cycle at 150 USD an hour. Kept as given. The input understates the cost while no fallback exists, since a retirement stops every close-out on this host at smoke, but the input is the operator's.
