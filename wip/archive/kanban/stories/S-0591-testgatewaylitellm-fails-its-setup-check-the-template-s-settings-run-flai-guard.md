---
id: S-0591
type: story
nature: remediation
title: "TestGatewayLiteLLM fails its setup check: the template's settings run flai guard 6 times, and it wants 4"
status: done
owner: alex
created: 2026-11-07T15:15:23Z
updated: 2026-11-09T07:21:55Z
transitions:
  - to: ready
    at: 2026-11-07T21:49:23Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-07T23:44:33Z
    by: agent-S-0591
  - to: review
    at: 2026-11-08T00:00:13Z
    by: system-flow
  - to: done
    at: 2026-11-09T07:21:55Z
    by: alex
tags: [flai, test]
touches: [flai/internal/serve/gateway_test.go, design/issues/I-0303-testgatewaylitellm-fails-its-setup-check-the-template-s-settings-run-flai-guard.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 976
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 4
      hand_edits: 2
      work: 114
    - day: 2026-11-08
      ceremony: 1
      work: 1
  models:
    - model: claude-opus-5-5
      input: 254
      output: 19040
      cache_read: 31069355
      cache_write: 319486
      cost: 9.1516
  strategic:
    - kind: planner
      seconds: 199
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 42
          output: 12652
          cache_read: 2578669
          cache_write: 155960
          cost: 2.0166
    - kind: orchestrator
      seconds: 119
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 45
          output: 305
          cache_read: 1674241
          cache_write: 123485
          cost: 0.4543
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-07T15:15:23Z
  value: 12.5
  by: planner-S-0591
  at: 2026-11-07T21:47:34Z
forecast:
  duration: 12m
  delivery: 2026-11-07T23:54:00Z
  basis: "Its own forecast of 12m; 1st in the pull order with an in-progress limit of 5, with nothing ahead of it."
  by: flai
  at: 2026-11-07T23:38:49Z
finalized:
  by: orchestrator
  at: 2026-11-07T21:49:11Z
---
# S-0591 TestGatewayLiteLLM fails its setup check: the template's settings run flai guard 6 times, and it wants 4

## Goal

This story remediates [I-0303](../../../design/issues/I-0303-testgatewaylitellm-fails-its-setup-check-the-template-s-settings-run-flai-guard.md), "TestGatewayLiteLLM fails its setup check: the template's settings run flai guard 6 times, and it wants 4". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0303 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0303 is closed with `flai issue close I-0303 --reason` saying what fixed it

## Tasks
- T-2592 guardSettings counts and rewrites only the hooks' flai guard invocations, with a test on the template's settings that -short runs
- T-2593 Close I-0303 with the reason naming the guardSettings fix and its test

## Notes

Cost of delay inputs set by flai from I-0303. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T10:19:48Z, 0.2 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from the instance and the code. `guardSettings` in `flai/internal/serve/gateway_test.go` counts every `flai guard` substring in `template/root/.claude/settings.json` and wants 4. S-0518 (commit ece3f9e69) added the message `flai guard did not run (exit %s)` to the two refusing hooks, so the file now holds four invocations and two messages: 6. The same `strings.ReplaceAll` also rewrites the two messages to name the built flai's path. The settings are right; the test's count is wrong. So the fix is in the test: count and rewrite only the hooks' invocations of `flai guard` (the command, not the text of the message), and add a test that runs without `-short` skipping it, on the template's real settings, so the ordinary Go tier catches the next drift instead of the manual `gateway` tier. TestGatewayOpenRouter goes through the same `checkGateway` and `guardSettings`, so it fails the same way and this fixes it too. Changing the settings' message instead was rejected: it is the guard's own error text, which `flai/cmd/guard_test.go` checks, and a change to `template/` needs a changelog line and a template release for a test's sake.

Touches:

| Path | From | Why |
|------|------|-----|
| `flai/internal/serve/gateway_test.go` | layout | holds `guardSettings` and the gateway checks; the new test goes beside them |
| `design/issues/I-0303-testgatewaylitellm-fails-its-setup-check-the-template-s-settings-run-flai-guard.md` | design (criterion 2) | `flai issue close` writes it |
| `design/issues/summary.md` | design (criterion 2) | `flai issue close` regenerates it |

`flai touches suggest` on the test file and the template's settings listed `.claude/settings.json`, `flai/cmd/guard_test.go`, `flai/cmd/guard.go`, `scripts/litellm.sh`, the docs, and `template/CHANGELOG.md` by co-change; none is kept, since the fix changes neither the settings, the guard, nor what an operator or user sees.

Forecast: see the plan's thread for the review of flai's figure; the basis written with the forecast says why it stands or moved.

Cost of delay: 12.50 USD a week, flai's figure from the operator's input of 5m lost per 168h cycle at 150 USD an hour, kept: the gateway tier is manual and runs weekly at most, so one failure a cycle is what the input says.
