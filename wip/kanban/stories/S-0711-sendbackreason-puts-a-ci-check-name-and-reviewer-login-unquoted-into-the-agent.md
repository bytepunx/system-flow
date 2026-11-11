---
id: S-0711
type: story
nature: improvement
title: sendBackReason puts a CI check name and reviewer login unquoted into the agent's prompt
status: ready
owner: alex
created: 2026-11-10T09:56:05Z
updated: 2026-11-11T16:15:10Z
transitions:
  - to: ready
    at: 2026-11-11T16:15:03Z
    by: orchestrator
tags: []
topics: [forge, orchestration]
touches: [flai/internal/serve/forgestatus.go, flai/internal/serve/forgestatus_test.go, flai/internal/serve/agents_test.go, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0340-sendbackreason-puts-a-ci-check-name-and-reviewer-login-unquoted-into-the-agent.md, design/issues/summary.md]
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
      seconds: 61
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 21
          output: 3693
          cache_read: 1647235
          cache_write: 35083
          cost: 0.6841
    - kind: orchestrator
      seconds: 42
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 16
          output: 137
          cache_read: 709554
          cache_write: 115442
          cost: 0.2307
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1m
    by: orchestrator
    at: 2026-11-11T09:08:43Z
  value: 2.5
  by: planner-S-0711
  at: 2026-11-11T09:09:25Z
forecast:
  duration: 1h
  delivery: 2026-11-11T18:10:00Z
  basis: "Its own forecast of 1h; 4th in the pull order with an in-progress limit of 3, behind S-0715, S-0707 and S-0714."
  by: flai
  at: 2026-11-11T16:15:10Z
finalized:
  by: alex
  at: 2026-11-11T09:03:41Z
---
# S-0711 sendBackReason puts a CI check name and reviewer login unquoted into the agent's prompt

## Goal

This story remediates [I-0340](../../../design/issues/I-0340-sendbackreason-puts-a-ci-check-name-and-reviewer-login-unquoted-into-the-agent.md), "sendBackReason puts a CI check name and reviewer login unquoted into the agent's prompt". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0340 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0340 is closed with `flai issue close I-0340 --reason` saying what fixed it

## Tasks
- T-3118 sendBackReason quotes each check name and reviewer login in a code span with its control characters dropped
- T-3119 Document the quoted send-back reason and close I-0340

## Notes

### Planning

The cause: `sendBackReason` in `flai/internal/serve/forgestatus.go` writes a failed check's name and a reviewer's login as they are into the move's reason, and `sentBackFor` in `flai/internal/serve/agents_restart.go` hands that reason to the agent's start prompt as `harness.SentBack.What`. Both strings come from the forge. The proposed remedy is the one the orchestrator recommended on TH-0500, recorded as D-0505: quote each name in a code span and drop its control characters, as the thread body's `mdlint.Fenced` quote does.

Touches, and where each came from:

| Touch | From |
|-------|------|
| `flai/internal/serve/forgestatus.go` | design: I-0340 and D-0505 name `sendBackReason` there |
| `flai/internal/serve/forgestatus_test.go` | co-change: 4 of 5 commits to `forgestatus.go` changed it; it asserts the reason's words |
| `flai/internal/serve/agents_test.go` | layout: its send-back tests expect the reason's words in the prompt |
| `design/system/flai-cli.md` | design: the send back's reason example, § The send back |
| `docs/users/flai.md` | design: the user guide's send-back step 2 quotes the reason |
| the I-0340 file and `design/issues/summary.md` | criterion 2: closing the issue |

`flai/internal/serve/agents_restart.go` is left out: it passes the reason on and needs no change. `flai touches suggest` also listed `forge_test.go`, `forge.go`, and four files under `flai/internal/forge/` at 40 to 60%; none builds or asserts the reason, so they are left out.

Figures:

- Forecast 1h, delivery 2026-11-11T22:47Z. flai forecast gave 2m from size 2 with no touches. I raised it for two serial tasks and the integration and smoke tiers a close-out of a `flai/` change runs, in line with the 1h to 1h30m of comparable remediation stories planned today.
- Cost of delay 2.50 USD a week, as flai cod works it out from the input set on TH-1079, `time_lost_per_cycle: 1m`. I left it unchanged: the risk is low, since only someone who can push sets a check name.
- No `after`: S-0711 reads no other open story's output. No folder touch.
