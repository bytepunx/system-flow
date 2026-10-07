---
id: S-0293
type: story
nature: improvement
title: flai stats classifies a story run's tool calls, so the ceremony turns E-0017 removes are measured per story and over time
status: done
parent: E-0017
owner: alex
created: 2026-10-06T11:37:08Z
updated: 2026-10-07T15:03:53Z
transitions:
  - to: ready
    at: 2026-10-07T09:19:04Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-07T09:22:20Z
    by: system-flow
  - to: review
    at: 2026-10-07T15:03:06Z
    by: agent-S-0293
  - to: done
    at: 2026-10-07T15:03:53Z
    by: orchestrator
tags: [cli, metrics]
topics: [automation, conventions]
touches: [flai/internal/usage/log.go, flai/internal/usage/usage.go, flai/internal/usage/log_test.go, flai/internal/metrics/usage.go, flai/internal/metrics/usage_test.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go, design/system/metrics.md, design/adrs, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, flai/internal/usage/turns.go, flai/internal/usage/turns_test.go, flai/internal/usage/usage_test.go, flai/internal/workitem/fields_test.go, flai/internal/workitem/front-matter-fields.txt, flai/internal/workitem/store_test.go, flai/internal/workitem/usage.go, flai/internal/workitem/usage_test.go, flai/internal/metrics/metrics.go, flai/internal/metrics/turns.go, flai/internal/metrics/turns_test.go, design/system/work-hierarchy.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2215
  models:
    - model: claude-opus-5-5
      input: 364
      output: 151484
      cache_read: 30171700
      cache_write: 939315
      cost: 15.6892
  strategic:
    - kind: orchestrator
      seconds: 919
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 94
          output: 1506
          cache_read: 13915819
          cache_write: 34023
          cost: 3.6396
        - model: claude-sonnet-5-5
          input: 5
          output: 28
          cache_read: 55169
          cache_write: 19291
          cost: 0.05
cost_of_delay:
  value: 59.02
  by: planner-E-0017
  at: 2026-10-06T11:37:27Z
forecast:
  duration: 40m
  delivery: 2026-10-07T14:20:00Z
  basis: "Its own forecast of 40m; 18th in the pull order with an in-progress limit of 3, behind S-0215, S-0246, S-0275, S-0216, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0279, S-0280, S-0287, S-0288, S-0289, S-0290 and S-0291."
  by: flai
  at: 2026-10-07T08:59:42Z
finalized:
  by: alex
  at: 2026-10-07T02:18:48Z
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
- [x] `flai stats` reports, per story and per day of the window, the main agent's turns in each class above, as text and `--json`, from the run logs `usage` already reads
- [x] The classes are defined in `design/system/metrics.md`, with an ADR, since that file is the contract with the dashboard
- [x] Run over the logs to 2026-10-04, the counts are compared with the epic's evidence, and the story records the comparison and explains each class that differs by more than ten percent.
- [x] `design/system/flai-cli.md` and the user guide describe the report

## Tasks
- T-1157 Measuring a story's usage classifies its own agent's turns per day into usage.turns
- T-1158 flai stats reports the story agents' turns per class, per story and per day of the window, as text and --json
- T-1159 metrics.md defines the turn classes with an ADR, and the CLI design and the user guide describe the report
- T-1160 The turn counts over the logs to 2026-10-04 are compared with E-0017's evidence in the story's notes

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

### Comparison with E-0017's evidence

Measured with this branch's flai at T-1157 (2fd91f0a), read-only, over every story log `~/.flai/serve/agents` holds (`flai serve agent usage --all --json --config ~/.flai/config.json`), summing `usage.turns` over the days to 2026-10-04: 88 stories, 8,504 turns (the epic: 108 runs, a mean of 80 turns, about 8,640).

| Class | Turns to 2026-10-04 | E-0017's figure | Difference |
|-------|--------------------:|-----------------|-----------:|
| ceremony | 748 (8.8% of turns) | "turns whose every call is pure ceremony are 9% of all turns", about 770 | −3% |
| test runs | 1,240 | 1,449 hand test runs | −14% |
| empty wakes | 197 | 650 turns woke with nothing to do | −70% |
| hand edits | 431 | 253 | +70% |
| work | 5,888 (69.2%) | none | |

Where the figures differ by more than ten percent, the epic counted something other than a turn of the class:

- Ceremony: the epic's "about 2,900 ceremony turns" counts ceremony commands inside every turn, not turns. Counted that way over the same logs, the commands the epic lists (commit, sync, move, log, touches, check, inbox) come to 3,097, +7%. A turn that only runs them is the class; its 9% agrees.
- Empty wakes: 650 is exactly the number of `wait_for_events` calls of the story agents in those logs, so the epic counted every wait. ADR-0105 counts a wake empty only when it timed out with no events, no changed paths, and no end: 211 calls, 197 turns. The other 439 woke on a change, usually another story's.
- Hand edits: of the 431, 288 are Python scripts that rewrite a narrative's `## Current state` and `## Next steps` or a story's criteria, and 35 are `sed -i`. The epic's 253 is closer to the Edit and Write calls and `sed -i` alone.
- Test runs: matching test tools anywhere in a command's text, as a one-off search does, finds about 290 more, which explains the epic's higher figure; but those commands name `make test` or `go test` in a heredoc or a quoted string (a task body, a narrative, a commit message) and run nothing. Of the rest, 47 are close-outs and 11 are `flai test`, which are the commands that replace hand runs.

### Accepted by the orchestrator

- Verified: d483224102818da8bb8f18a75d3d7818f2c5675b
- At: 2026-10-07T15:03:53Z

Verdict: Pass. All four criteria are met and ticked, with ADR-0116 behind the metrics.md change, verified at d483224102818da8bb8f18a75d3d7818f2c5675b.

- 1: flai/cmd/stats.go, flai/internal/metrics/turns.go, flai/internal/metrics/metrics.go, flai/internal/usage/turns.go, flai/internal/usage/log.go, flai/internal/usage/usage.go, flai/internal/workitem/usage.go, flai/internal/workitem/front-matter-fields.txt
- 2: design/system/metrics.md, design/adrs/0116-when-flai-measures-a-story-s-usage-from-its-logs-it-classifies-each-turn-of-the.md, design/adrs/README.md, design/system/work-hierarchy.md
- 3: flai/internal/usage/turns.go, flai/internal/usage/turns_test.go
- 4: design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md
