---
id: S-0293
type: story
nature: improvement
title: flai stats classifies a story run's tool calls, so the ceremony turns E-0017 removes are measured per story and over time
status: backlog
parent: E-0017
owner: alex
created: 2026-10-06T11:37:08Z
updated: 2026-10-06T23:03:32Z
transitions: []
tags: [cli, metrics]
topics: [automation, conventions]
touches: [flai/internal/usage/log.go, flai/internal/usage/usage.go, flai/internal/usage/log_test.go, flai/internal/metrics/usage.go, flai/internal/metrics/usage_test.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go, design/system/metrics.md, design/adrs, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  value: 59.02
  by: planner-E-0017
  at: 2026-10-06T11:37:27Z
forecast:
  duration: 40m
  delivery: 2026-10-07T10:56:00Z
  basis: "Its own forecast of 40m; 40th in the pull order with an in-progress limit of 3, behind S-0301, S-0300, S-0228, S-0261, S-0269, S-0270, S-0271, S-0212, S-0213, S-0214, S-0215, S-0216, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0245, S-0246, S-0251, S-0254, S-0265, S-0272, S-0273, S-0274, S-0275, S-0277, S-0279, S-0280, S-0281, S-0286, S-0287, S-0288, S-0289, S-0290 and S-0291."
  by: flai
  at: 2026-10-06T23:03:32Z
---
# S-0293 flai stats classifies a story run's tool calls, so the ceremony turns E-0017 removes are measured per story and over time

## Goal

E-0017's evidence came from a one-off classification of every main-agent tool call in the logs under `~/.flai/serve/agents/` on 2026-10-04. It found about 2,900 ceremony turns, 1,449 hand test runs, 650 empty wakes, and 253 hand edits of narratives, task bodies, and criteria. The epic is done when a story agent's run shows no turn that only commits, syncs, moves, logs, ticks, or reads a test log. Nothing measures that today. `flai stats` reads the run logs it already reads for usage, and puts each main-agent turn in a class:

- ceremony: only flai commands that a story-loop command replaces
- test run: a test, lint, or format command run by hand
- empty wake: a `wait_for_events` that returned nothing
- hand edit: of a narrative, a criterion, or a task body
- work: everything else

It reports the counts per story and per day, so each E-0017 story's saving shows against the baseline, and the analyzer can read it.

## Acceptance criteria
- [ ] `flai stats` reports, per story and per day of the window, the main agent's turns in each class above, as text and `--json`, from the run logs `usage` already reads
- [ ] The classes are defined in `design/system/metrics.md`, with an ADR, since that file is the contract with the dashboard
- [ ] Run over the logs to 2026-10-04, the counts agree with the epic's evidence within ten percent, and the story records the comparison
- [ ] `design/system/flai-cli.md` and the user guide describe the report

## Tasks

## Notes

Added by the planner for E-0017: the epic's done condition needs a measure, and the baseline should be taken before the other stories land. S-0272's fourth criterion, counting empty wakes, is one class of this report; the plan thread proposes moving it here.

### Planning

Touches, from layout and design:

- `flai/internal/usage/log.go`, `usage.go`, `log_test.go`: layout. The run logs are parsed here for tokens and cost; the turn's tool calls are read from the same lines.
- `flai/internal/metrics/usage.go`, `usage_test.go`, `flai/cmd/stats.go`, `flai/cmd/check_stats_test.go`: layout. Usage metrics and the stats command.
- `design/system/metrics.md`, `design/adrs`: design. `CLAUDE.md` says metric definitions change only with an ADR.
- `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`: design. The stats command is described there.

Left out: `design/system/flaiover-dashboard.md` and the operator docs, which `flai touches suggest` lists (76 and 72 of 299 commits). A chart of the classes is a later story, if the operator wants one.

Forecast: flai gave 24m (89 s per unit over 21 done large improvement stories, times size 16). I raised it to 40m for two things: calibrating the classes against the epic's hand count within ten percent, and an ADR. The delivery is flai's, 2026-10-07T02:46Z, moved by the added 16m.

Cost of delay: 59.02 USD a week, flai's share of E-0017's 900 USD a week by forecast duration, unadjusted. It removes no turns itself, so a turn-based share would be zero. Without it, though, the epic's done condition cannot be checked. It sits on top of the seven turn-based values, which sum to 900.
