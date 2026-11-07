---
id: S-0589
type: story
nature: remediation
title: A story already in ready when the agent host action is turned on gets no agent until someone presses Start agent
status: in-progress
owner: alex
created: 2026-11-07T09:44:04Z
updated: 2026-11-07T10:07:35Z
transitions:
  - to: ready
    at: 2026-11-07T10:02:33Z
    by: alex
  - to: in-progress
    at: 2026-11-07T10:06:00Z
    by: agent-S-0589
tags: []
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, docs/operators/index.md, design/system/flai-cli.md, design/issues/I-0300-a-story-already-in-ready-when-the-agent-host-action-is-turned-on-gets-no-agent.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-07T09:44:04Z
  value: 25
  by: planner-S-0589
  at: 2026-11-07T10:07:35Z
forecast:
  duration: 1h30m
  delivery: 2026-11-07T12:00:00Z
  basis: "flai forecast's 24m (117 s per unit over 40 medium remediation stories, size 12) raised to 1h30m: the cause is unconfirmed, and serve-level loop tests are slow on this host (S-0560)."
  by: planner-S-0589
  at: 2026-11-07T10:07:35Z
finalized:
  by: alex
  at: 2026-11-07T10:02:06Z
---
# S-0589 A story already in ready when the agent host action is turned on gets no agent until someone presses Start agent

## Goal

This story remediates [I-0300](../../../design/issues/I-0300-a-story-already-in-ready-when-the-agent-host-action-is-turned-on-gets-no-agent.md), "A story already in ready when the agent host action is turned on gets no agent until someone presses Start agent". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0300 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0300 is closed with `flai issue close I-0300 --reason` saying what fixed it

## Tasks
- T-2573 Reproduce I-0300 in flai serve's tests: a story already in ready gets no agent after the agent action is turned on
- T-2574 Start a story already in ready on the first look after the agent action is turned on for its project
- T-2575 Say in the operators guide and flai-cli.md that a story already in ready gets its agent once the agent action is turned on
- T-2576 Close I-0300 with flai issue close, naming the cause and the fix

## Notes

Cost of delay inputs set by flai from I-0300. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T08:31:47Z, 0.1 days before this story; under one cycle counts as one).

### Planning

The proposed solution: reproduce the hold in a test first (T-2573), because the issue's guess does not match the code. `(*launcher).look` in `flai/internal/serve/agents.go` runs on every change under `kanban/` and `threads/` and every minute, and reads the `agent` action on each look. So something else kept S-0001 waiting from 08:26 to 08:32. Then make turning the action on start the stories already in ready (T-2574). The three candidate causes are listed in T-2573. The scratch project is gone (TH-0863), so its `serve/agents.json` cannot be read.

Touches, by source:

| Touch | Source |
|-------|--------|
| `flai/internal/serve/agents.go` | layout: the launcher's `look` and `readyStories` |
| `flai/internal/serve/serve.go` | layout: the project loop's minute look and the watcher's paths |
| `flai/internal/serve/agents_test.go`, `flai/internal/serve/serve_test.go` | layout: the launcher's and the loop's tests, for the reproducing test |
| `flai/cmd/serve_actions.go` | layout: `flai serve enable agent` and `agentConfig` |
| `flai/internal/hostapi/writes.go` | layout: the `ActionAgent` comment and description, which say only "when a story becomes ready" |
| `flai/cmd/serve_actions_test.go`, `flai/internal/hostapi/writes_test.go` | co-change: 60 and 78 of the commits that changed the touches above |
| `docs/operators/index.md`, `design/system/flai-cli.md` | design: § "Starting an agent when a story becomes ready", and the `flai serve agent` row |
| `design/issues/I-0300-…md`, `design/issues/summary.md` | goal: criterion 2, `flai issue close` |

No folder touch, and no `after`: the story reads nothing another open story writes. S-0587 (the default harness) is about a different gate. S-0001 started on the operator's Start agent (TH-0860), so its harness was set.

Figures:

- Forecast: flai gave 24m (117 s per unit over 40 medium remediation stories, size 12). I raised it to 1h30m, delivery 2026-11-07T12:00:00Z. The cause is not confirmed yet. The serve-level loop tests are slow on this host, near their timeouts (S-0560).
- Cost of delay: 25 USD a week, as flai cod gave it from the operator's input of 10m lost per 168h cycle at 150 USD an hour. I kept it. The issue says every import formula run hits this, but the time lost is the operator's input, not mine.
