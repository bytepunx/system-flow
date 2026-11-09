---
id: S-0613
type: story
nature: improvement
title: "flai guard refuses the orchestrator a read or a permitted flai command when the command carries a stderr redirection such as 2>&1 or 2>/dev/null"
status: done
owner: alex
created: 2026-11-07T23:38:58Z
updated: 2026-11-09T09:47:02Z
transitions:
  - to: ready
    at: 2026-11-09T08:19:39Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-09T08:22:35Z
    by: agent-S-0613
  - to: review
    at: 2026-11-09T09:44:31Z
    by: agent-S-0613
  - to: done
    at: 2026-11-09T09:47:02Z
    by: orchestrator
tags: [flai, guard]
topics: [guard, orchestration]
touches: [flai/internal/guard/guard.go, flai/internal/guard/roles.go, flai/internal/guard/orchestrate_test.go, flai/internal/guard/roles_test.go, design/system/strategic-agents.md, docs/users/flai.md, design/issues/I-0230-flai-guard-refuses-the-orchestrator-a-read-or-a-permitted-flai-command-when-the-command-carries-a-stderr-redirection-such-as-2-1-or-2-dev-null.md, design/issues/summary.md, flai/internal/guard/call.go, flai/internal/guard/hostconfig.go, flai/internal/guard/hostconfig_test.go, flai/internal/guard/shared.go, design/issues/I-0318-flai-guard-refuses-the-orchestrator-flai-serve-status-a-read-because-its-reads.md, design/issues/I-0311-teststartdevserverstopsvitewhenitscontextends-fails-under-load-the-fake-vite.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4934
  priced_by: harness
  turns:
    - day: 2026-11-09
      ceremony: 3
      hand_edits: 2
      work: 124
  models:
    - model: claude-opus-5-5
      input: 336
      output: 44946
      cache_read: 46417540
      cache_write: 1242097
      cost: 17.3461
  strategic:
    - kind: planner
      seconds: 383
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 60
          output: 20987
          cache_read: 4123800
          cache_write: 159198
          cost: 2.5183
    - kind: orchestrator
      seconds: 546
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 81
          output: 489
          cache_read: 4033856
          cache_write: 323784
          cost: 1.1306
        - model: claude-sonnet-5-5
          input: 8
          output: 852
          cache_read: 77781
          cache_write: 38240
          cost: 0.0861
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 33
          output: 3103
          cache_read: 238317
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2841
          cache_read: 759050
          cache_write: 13855
          cost: 0.278
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h
    by: flai
    at: 2026-11-07T23:38:58Z
  value: 150
  by: planner-S-0613
  at: 2026-11-09T08:10:15Z
forecast:
  duration: 40m
  delivery: 2026-11-09T09:33:00Z
  basis: "Three comparable flai guard remediations took 31 to 38 minutes of agent time (S-0368 1845 s, S-0568 1927 s, S-0556 2291 s); flai's 18m from size 10 is below them, and delivery is flai's 09:11Z moved by the 22 minutes added."
  by: planner-S-0613
  at: 2026-11-09T08:13:57Z
finalized:
  by: orchestrator
  at: 2026-11-09T08:17:38Z
---
# S-0613 flai guard refuses the orchestrator a read or a permitted flai command when the command carries a stderr redirection such as 2>&1 or 2>/dev/null

## Goal

This story remediates [I-0230](../../../design/issues/I-0230-flai-guard-refuses-the-orchestrator-a-read-or-a-permitted-flai-command-when-the-command-carries-a-stderr-redirection-such-as-2-1-or-2-dev-null.md), "flai guard refuses the orchestrator a read or a permitted flai command when the command carries a stderr redirection such as 2>&1 or 2>/dev/null", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0230 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0230 is closed with `flai issue close I-0230 --reason` saying what fixed it

## Tasks
- T-2706 The orchestrator's and every role file's flai checks judge a command by its arguments, its redirections and their targets left out
- T-2707 The design and the user guide say a redirection neither allows nor refuses a strategic agent's command
- T-2708 Close I-0230 with the reason S-0613's fix gives

## Notes

The analyzer found I-0230: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0230. time_lost_per_cycle 1h carried over from I-0230's Impact section. Its Impact time_lost_per_cycle was taken rather than the 7m derived from its cost and count.

Strategic usage carried over from I-0230: analyzer 42s, 0.8156 USD.

### Planning

Proposed solution, from I-0230's instances and the code: `commands` in `flai/internal/guard/guard.go` returns each simple command's words with a redirection such as `2>&1` kept as a word, and `Guard.refuse` judges those words for the orchestrator (`Guard.orchestrate`) and for a role file (`Guard.asRole` in `flai/internal/guard/roles.go`). So `positionals` counts `2>&1` as a story and `ordersByPolicy` takes `flai order --by throughput --apply 2>&1` for a placement by hand, and `subcommands` reads `2>/dev/null` as the subcommand of `flai board --json 2>/dev/null`. `parse` already splits each command's arguments without its redirections (S-0368); the fix judges those instead (T-2706).

Touches, each a file, none a folder:

- `flai/internal/guard/guard.go`, `flai/internal/guard/roles.go`: layout, the code the cause lives in (`commands`, `refuse`, `orchestrate`, `asRole`); `flai touches suggest S-0613 flai/internal/guard/guard.go` lists `roles.go` at 8% co-change.
- `flai/internal/guard/orchestrate_test.go`, `flai/internal/guard/roles_test.go`: co-change (14% and 8%) and the criterion's reproducing test. `guard_test.go` (67%) is left out: its S-0368 tests of `commands` stay as they are.
- `design/system/strategic-agents.md`, `docs/users/flai.md`: design, the sections that describe the guard's checks of a strategic agent (§ Its permissions and the guard; § Role files); `documentation.md` asks for the user guide in the same story.
- The I-0230 issue file and `design/issues/summary.md`: the second criterion, which `flai issue close` writes.
- Left out from co-change: `flai/cmd/guard.go` and `docs/users/flai-reference.md`, since no flag or help text changes.

Tasks and layers: layer 1 is T-2706 (the fix and its tests); layer 2 is T-2707 (the docs) and T-2708 (closing I-0230), which both wait for T-2706 and share no path, so they run together.

Forecast 40m, adjusted from flai's 18m: three comparable guard remediations took 31 to 38 minutes of agent time (S-0368, S-0568, S-0556). Delivery is flai's 09:11Z moved by the 22 minutes added.

Cost of delay value 150 USD a week, as `flai cod S-0613` gives it from the operator's input of 1h lost per cycle at 150 USD an hour, one cycle a week. It stands: I-0230's Impact counts 52 refusals of permitted calls over two weeks, and each refused `flai order --by <policy> --apply` leaves the ready column unordered until the next run.

No `after`: the story reads no other story's output.

### Accepted by the orchestrator

- Verified: bdb6cd4aad84e6f67eeae79df78b6d08d70c8fe4
- At: 2026-11-09T09:47:02Z

Verdict: meets all criteria
- 1: flai/internal/guard/guard.go, flai/internal/guard/call.go, flai/internal/guard/hostconfig.go, flai/internal/guard/roles.go, flai/internal/guard/shared.go, flai/internal/guard/orchestrate_test.go, flai/internal/guard/roles_test.go, flai/internal/guard/hostconfig_test.go, design/system/strategic-agents.md, docs/users/flai.md
- 2: design/issues/I-0230-flai-guard-refuses-the-orchestrator-a-read-or-a-permitted-flai-command-when-the-command-carries-a-stderr-redirection-such-as-2-1-or-2-dev-null.md, design/issues/summary.md
