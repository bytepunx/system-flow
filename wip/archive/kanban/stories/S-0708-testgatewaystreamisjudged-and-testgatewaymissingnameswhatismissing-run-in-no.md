---
id: S-0708
type: story
nature: remediation
title: "TestGatewayStreamIsJudged and TestGatewayMissingNamesWhatIsMissing run in no tier, because go-test and integration skip ^TestGateway and the gateway tier runs only the two paid checks by name"
status: done
owner: alex
created: 2026-11-10T09:55:53Z
updated: 2026-11-11T01:43:36Z
transitions:
  - to: ready
    at: 2026-11-10T10:29:38Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-10T10:52:56Z
    by: agent-S-0708
  - to: review
    at: 2026-11-10T19:36:12Z
    by: system-flow
  - to: done
    at: 2026-11-11T01:43:36Z
    by: alex
blocked:
  - from: 2026-11-10T11:10:53Z
    until: 2026-11-10T11:14:20Z
    reason: awaiting the operator on TH-1049
  - from: 2026-11-10T11:14:57Z
    until: 2026-11-10T19:29:38Z
    reason: awaiting the operator on TH-1049
tags: [flai, go, testing]
topics: [testing, code]
touches: [flai/internal/serve/gateway_test.go, design/conventions/code-quality.md, docs/contributors/index.md, design/issues/I-0334-testgatewaystreamisjudged-and-testgatewaymissingnameswhatismissing-run-in-no.md, design/issues/summary.md, design/issues/I-0229-two-stories-accepted-apart-leave-main-red-s-0535-reworded-to-backlog-s-condition-and-s-0416-s-testadefinedroleisheldtoitsfile-still-expects-the-old-words.md, design/issues/I-0345-wait-for-events-fails-with-no-front-matter-while-another-story-s-checked.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1505
  priced_by: harness
  turns:
    - day: 2026-11-10
      ceremony: 3
      hand_edits: 5
      work: 135
  models:
    - model: claude-opus-5-5
      input: 328
      output: 40350
      cache_read: 45879026
      cache_write: 753105
      cost: 15.8797
  strategic:
    - kind: orchestrator
      seconds: 241
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 39
          output: 221
          cache_read: 1457435
          cache_write: 159774
          cost: 0.4434
        - model: claude-sonnet-5-5
          input: 6
          output: 48
          cache_read: 49038
          cache_write: 29325
          cost: 0.0589
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-10T09:55:53Z
  value: 25
  by: planner-S-0708
  at: 2026-11-10T10:27:11Z
forecast:
  duration: 45m
  delivery: 2026-11-10T15:36:00Z
  basis: "flai forecast's 16m (134 s per unit over 54 medium remediation stories, size 7) raised to 45m for an AST guard test and a close-out that runs integration and smoke, about 12 minutes, since the story touches flai/; delivery moved by the same 29m"
  by: planner-S-0708
  at: 2026-11-10T10:27:11Z
finalized:
  by: orchestrator
  at: 2026-11-10T10:28:40Z
---
# S-0708 TestGatewayStreamIsJudged and TestGatewayMissingNamesWhatIsMissing run in no tier, because go-test and integration skip ^TestGateway and the gateway tier runs only the two paid checks by name

## Goal

This story remediates [I-0334](../../../design/issues/I-0334-testgatewaystreamisjudged-and-testgatewaymissingnameswhatismissing-run-in-no.md), "TestGatewayStreamIsJudged and TestGatewayMissingNamesWhatIsMissing run in no tier, because go-test and integration skip ^TestGateway and the gateway tier runs only the two paid checks by name". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0334 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0334 is closed with `flai issue close I-0334 --reason` saying what fixed it

## Tasks
- T-2961 The two no-network gateway tests leave the TestGateway prefix, and a test fails on a TestGateway function gateway-smoke.sh does not run
- T-2962 code-quality.md and the contributors guide say the TestGateway prefix names only the paid checks the gateway tier runs
- T-2964 I-0334 is closed with the reason that the no-network gateway tests left the TestGateway prefix and a test guards it

## Notes

Cost of delay inputs set by flai from I-0334. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T09:01:59Z, 0 days before this story; under one cycle counts as one).

### Planning

Proposed remedy: keep the `-skip '^TestGateway'` of the go-test tier (`system-flow.yaml`) and of `scripts/integration.sh` as it is, and keep the `TestGateway` prefix for the paid checks alone. Rename the two no-network tests in `flai/internal/serve/gateway_test.go` out of the prefix, as S-0685 named `TestCreditRefusalIsFoundInAGatewayRun`. Add a test, named outside the prefix, that parses the serve package's test files and fails when a `TestGateway*` function is not one `scripts/gateway-smoke.sh` runs by name. On main's names it fails, naming the two tests, which reproduces I-0334; a future no-network test named `TestGateway*` fails it too.

The alternative the issue names, narrowing the skip to `^TestGateway(OpenRouter|LiteLLM)$` in the go-test tier and `scripts/integration.sh`, was not taken. It changes the manifest, the script, and the pin in `flai/internal/verify/select_test.go`, and a new paid check named `TestGateway*` would then run in every close-out and batch and spend on the gateway's key. The hand-run advice in D-0451, `-skip '^TestGateway'`, stays right under the remedy.

Touches, each a file:

| Touch | Source | Why |
| --- | --- | --- |
| `flai/internal/serve/gateway_test.go` | design (I-0334) | the two renames and the guard test |
| `design/conventions/code-quality.md` | layout | its project additions name the tiers and the paid gateway tests; they gain the naming rule |
| `docs/contributors/index.md` | layout | its gateway section says the go-test and integration tiers skip `^TestGateway`; it gains the naming rule |
| `design/issues/I-0334-testgatewaystreamisjudged-and-testgatewaymissingnameswhatismissing-run-in-no.md` | layout | `flai issue close` writes it |
| `design/issues/summary.md` | layout | `flai issue close` writes it |

`flai touches suggest` listed co-change paths at 5 percent or less, none of which the remedy changes, so none was added. No folder touch. No `after`: the story reads nothing another open story writes.

Forecast: `flai forecast` gave 16m (134 s per unit of size over 54 done medium remediation stories, size 7). Raised to 45m: the guard test parses Go files and the script, and a story that touches `flai/` runs integration and smoke in its close-out, about 12 minutes. Delivery moved by the same 29m, to 2026-11-10T15:36:00Z.

Cost of delay: `flai cod` gives 25 USD a week from flai's input of 10m lost per 168h cycle at 150 USD an hour. It stands: one occurrence, and the cost is a regression in stream parsing or skip messages passing every close-out unseen.
