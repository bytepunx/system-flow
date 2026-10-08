---
id: S-0352
type: story
nature: improvement
title: Each adapter states its capabilities and reads its harness's log into neutral usage events, and flai measures usage from those events
status: backlog
parent: E-0019
owner: alex
created: 2026-10-08T08:45:35Z
updated: 2026-10-08T09:08:56Z
transitions: []
tags: [cli]
topics: [agents, usage]
touches: [flai/internal/harness/harness.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go, flai/internal/usage/event.go, flai/internal/usage/claudecode.go, flai/internal/usage/claudecode_test.go, flai/internal/usage/log.go, flai/internal/usage/log_test.go, flai/internal/usage/turns.go, flai/internal/usage/turns_test.go, flai/internal/serve/usage.go, flai/internal/serve/usage_test.go, design/system/metrics.md, design/system/flai-cli.md]
after: [S-0351]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 49.78
  by: planner-E-0019
  at: 2026-10-08T09:00:20Z
forecast:
  duration: 45m
  delivery: 2026-10-08T21:27:00Z
  basis: "Its own forecast of 45m; 26th in the pull order with an in-progress limit of 5, behind S-0232, S-0322, S-0341, S-0344, S-0345, S-0338, S-0346, S-0342, S-0337, S-0334, S-0343, S-0297, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0289, S-0304, S-0305, S-0306, S-0313, S-0347, S-0348, S-0349, S-0350 and S-0351."
  by: flai
  at: 2026-10-08T09:08:56Z
finalized:
  by: alex
  at: 2026-10-08T08:52:03Z
---
# S-0352 Each adapter states its capabilities and reads its harness's log into neutral usage events, and flai measures usage from those events

## Goal

`flai/internal/usage` reads Claude Code's stream-json field by field, and nothing says what a harness can do before `flai serve` starts it. [ADR-0130](../../../design/adrs/0130-every-harness-meets-flai-through-neutral-contracts-an-adapter-s-capabilities-a.md) makes both a contract of the adapter: `Capabilities()`, whether the harness can resume a session, run the guard, hold a call to ask, give each role its own model, run a strategic agent from a definition, and report cost; and `Reader()`, which turns the harness's log into neutral events. This story builds both inside the `claude-code` and `command` adapters and moves the measurement onto the events, with no change to any figure a story records. It comes after S-0351 so that what OpenRouter's run showed about message IDs and cost shapes the event before it is fixed.

## Acceptance criteria

- [ ] `harness.Adapter` has `Capabilities() Capabilities` with `Resume`, `Guard`, `Ask`, `RoleModels`, `StrategicAgents`, and `Cost`; `claude-code` states all true and `command` all false, and a test pins both.
- [ ] `harness.Adapter` has `Reader() usage.Reader`, and `usage.Event` carries the kind (session start, model call, tool call, tool result, run end), the session, the model, the message ID, the tokens (input, output, cache read, cache write), the cost when the harness reports one, the parent call that started a sub-agent, and the tool's name, ID, description, and result.
- [ ] The Claude Code reader turns stream-json into those events, keeping what `log.go` relies on today: the repeated message IDs, `parent_tool_use_id`, `modelUsage` and `costUSD` per model, and the empty-wake test on a `wait_for_events` result.
- [ ] `usage.Read`, the task apportionment, the empty wakes, and the turn classes work from events, and every existing test in `flai/internal/usage` and `flai/internal/serve` passes unchanged in its expectations; a log the `command` adapter wrote yields no events and no usage, as today.
- [ ] `design/system/metrics.md` defines an empty wake and a turn on the neutral events, with Claude Code's fields as that reader's mapping, and `design/system/flai-cli.md` describes the two contracts.

## Tasks

- T-1387 usage.Event and usage.Reader are defined, and a Claude Code reader turns stream-json into them
- T-1388 Each adapter states its capabilities, claude-code all of them and command none
- T-1389 Each adapter gives its reader, and usage.Read, the task shares, empty wakes, and turns are measured from events
- T-1390 metrics.md defines an empty wake and a turn on the neutral events, and flai-cli.md describes capabilities and readers

## Notes

- `flai/internal/serve/stream.go`, which shows the live stream on the dashboard, keeps reading stream-json: it is display, not measurement, and moving it is not in ADR-0130.
- The metrics change is the wording of two definitions, under ADR-0130; no figure changes.

### Planning

Planned by planner-E-0019 on 2026-10-08. Every touch is a file; no folder touch is kept; `event.go`, `claudecode.go`, and its test are new. Three layers: T-1387 and T-1388 together, since one is in `usage` and the other in `harness`; then T-1389, which joins them; then T-1390.

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/usage/event.go`, `claudecode.go`, `claudecode_test.go` | design | ADR-0130's `Reader` and event, and Claude Code's reader (T-1387) |
| `flai/internal/harness/harness.go`, `adapters.go`, `harness_test.go` | design, layout | `Capabilities()` and `Reader()` on the adapter (T-1388, T-1389) |
| `flai/internal/usage/log.go`, `log_test.go`, `turns.go`, `turns_test.go` | layout | Every stream-json field read today is in these two (T-1389) |
| `flai/internal/serve/usage.go`, `usage_test.go` | layout | Where `flai serve` measures a story's runs (T-1389) |
| `design/system/metrics.md`, `design/system/flai-cli.md` | design | ADR-0130's consequences: the empty wake and the turn are defined on stream-json fields today (T-1390) |

`touches suggest` listed the user and operator guides and the dashboard's design. None is taken: nothing a user sees changes.

Forecast: 45m. `flai forecast` gave 25m, 78 s per unit over 51 improvement stories, size 19. It is raised because the story moves the 1,070 lines of `log.go` and `turns.go` onto events with every test's expectations unchanged.

Cost of delay: 49.78 USD a week, as `flai cod` works it out: 45m of the 7h32m forecast over E-0019's 12 stories, of the operator's 500 USD a week penalty. It stands.
