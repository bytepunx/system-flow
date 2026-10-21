---
id: S-0455
type: story
nature: remediation
title: The LiteLLM gateway smoke check fails because Claude Code 2.1.296 refuses the default haiku model name as unrecognized_model
status: ready
owner: alex
created: 2026-10-21T00:29:30Z
updated: 2026-10-21T00:41:56Z
transitions:
  - to: ready
    at: 2026-10-21T00:37:27Z
    by: alex
tags: []
touches: [scripts/litellm.sh, scripts/README.md, flai/internal/serve/gateway_test.go, docs/contributors/index.md, design/system/agent-adapters.md, design/issues/I-0171-the-litellm-gateway-smoke-check-fails-because-claude-code-2-1-296-refuses-the-default-haiku-model-name-as-unrecognized-model.md, design/issues/summary.md]
after: [S-0451]
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
    - kind: planner
      seconds: 222
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 82
          output: 16041
          cache_read: 11067813
          cache_write: 85575
          cost: 3.5607
    - kind: orchestrator
      seconds: 286
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 38
          output: 805
          cache_read: 14332245
          cache_write: 9388
          cost: 3.5954
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-21T00:29:30Z
  value: 37.5
  by: planner-S-0455
  at: 2026-10-21T00:40:41Z
forecast:
  duration: 45m
  delivery: 2026-10-21T01:33:00Z
  basis: "Its own forecast of 45m; 1st in the pull order with an in-progress limit of 5, behind S-0445 and S-0453."
  by: flai
  at: 2026-10-21T00:41:56Z
finalized:
  by: alex
  at: 2026-10-21T00:37:14Z
---
# S-0455 The LiteLLM gateway smoke check fails because Claude Code 2.1.296 refuses the default haiku model name as unrecognized_model

## Goal

This story remediates [I-0171](../../../design/issues/I-0171-the-litellm-gateway-smoke-check-fails-because-claude-code-2-1-296-refuses-the-default-haiku-model-name-as-unrecognized-model.md), "The LiteLLM gateway smoke check fails because Claude Code 2.1.296 refuses the default haiku model name as unrecognized_model". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0171 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0171 is closed with `flai issue close I-0171 --reason` saying what fixed it

## Tasks
- T-1941 scripts/litellm.sh stops freezing its default haiku model into admin.env, and replaces a retired one on up
- T-1942 TestGatewayLiteLLM probes its model through the proxy first and fails naming the slug, its source, and the upstream answer
- T-1943 The host's LiteLLM files take the current model, the gateway check passes at smoke, and I-0171 is closed

## Notes

Cost of delay inputs set by flai from I-0171. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-20T20:07:32Z, 0.2 days before this story; under one cycle counts as one).

### Planning

The fix proposed, from I-0171's corrected instance and TH-0529 and TH-0530: the cause is OpenRouter's 404 for the retired `nvidia/nemotron-3-super-120b-a12b:free`, not Claude Code's `unrecognized_model` warning. S-0451's T-1928, in review, moves the built-in default to `nemotron-3-ultra-550b-a55b:free`. That change alone does not fix this host, because `scripts/litellm.sh` writes the model into `.flai-cache/litellm/admin.env` on the first `up` and carries it into `litellm.env` as `FLAI_TEST_LITELLM_HAIKU`, which beats the default. So T-1941 stops the script freezing its default and drops a retired one. T-1942 makes the check probe its model and name an upstream refusal, with a short test that reproduces it. T-1943 shows the result on this host and closes I-0171.

`after: S-0451`: T-1928 changes `litellmHaikuDefault` in `gateway_test.go` and the default in `scripts/litellm.sh`, both of which this story changes again.

Touches, all files, no folder:

| Touch | Source | Task |
|-------|--------|------|
| `scripts/litellm.sh` | design (I-0171 names `LITELLM_HAIKU_MODEL` in it) | T-1941 |
| `scripts/README.md` | layout (the script's row) | T-1941 |
| `docs/contributors/index.md` | co-change (2 of 5 commits with the two files) | T-1941 |
| `flai/internal/serve/gateway_test.go` | design (I-0171 names `litellmHaikuDefault` in it) | T-1942 |
| `design/system/agent-adapters.md` | co-change (2 of 5 commits) | T-1943 |
| `design/issues/I-0171-…md` | criterion 2 | T-1943 |
| `design/issues/summary.md` | layout (kept with every issue change) | T-1943 |

`flai touches suggest` listed nothing from the story at first, since it declared no touches; run from the two files I-0171 names, it gave the two co-change paths.

Forecast: flai gave 17m (110 s per unit of size over 24 done remediation stories, size 9). Raised to 45m: T-1943 runs the smoke tier live through the proxy and waits on the operator's word before rewriting the host's LiteLLM files. Delivery 2026-10-21T03:00:00Z, after S-0451's acceptance.

Cost of delay: 37.50 USD a week, as `flai cod` works it out from flai's input, 15m lost per 168h cycle at 150 USD an hour. Kept as given. The input likely understates the cost, since every close-out on this host with the proxy set stops at smoke, but the input is the operator's.
