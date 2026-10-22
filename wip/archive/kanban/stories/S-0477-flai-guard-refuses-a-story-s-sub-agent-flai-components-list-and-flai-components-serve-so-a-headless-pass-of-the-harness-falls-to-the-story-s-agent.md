---
id: S-0477
type: story
nature: improvement
title: flai guard refuses a story's sub-agent flai components list and flai components serve, so a headless pass of the harness falls to the story's agent
status: done
owner: alex
created: 2026-10-21T11:59:27Z
updated: 2026-10-22T21:44:12Z
transitions:
  - to: ready
    at: 2026-10-22T20:55:32Z
    by: alex
  - to: in-progress
    at: 2026-10-22T20:57:54Z
    by: agent-S-0477
  - to: review
    at: 2026-10-22T21:31:28Z
    by: agent-S-0477
  - to: done
    at: 2026-10-22T21:44:12Z
    by: orchestrator
tags: [flai, guard, components]
topics: [cli, components]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/agent-adapters.md, design/system/component-library.md, design/issues/I-0165-flai-guard-refuses-a-story-s-sub-agent-flai-components-list-and-flai-components-serve-so-a-headless-pass-of-the-harness-falls-to-the-story-s-agent.md, design/issues/summary.md, design/adrs/0186-flai-guard-lets-every-agent-run-flai-components-list-and-a-story-s-task-sub.md, design/adrs/README.md, design/system/flai-cli.md, flai/internal/guard/call.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2031
  priced_by: harness
  turns:
    - day: 2026-10-22
      ceremony: 3
      hand_edits: 3
      work: 47
  models:
    - model: claude-opus-5-5
      input: 190
      output: 53242
      cache_read: 12549507
      cache_write: 363676
      cost: 6.1338
  strategic:
    - kind: orchestrator
      seconds: 238
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 80
          output: 1272
          cache_read: 15319690
          cache_write: 21118
          cost: 3.3875
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-21T11:59:27Z
  value: 25
  by: planner-S-0477
  at: 2026-10-22T20:59:32Z
forecast:
  duration: 40m
  delivery: 2026-10-22T21:45:00Z
  basis: "flai forecast gives 18m (median 87 s per unit of size over 14 done improvement stories on claude-opus-5-5, size 12: 2 criteria, 10 touches) and delivery 2026-10-22T21:21Z; raised to 40m because the close-out's integration and smoke tiers alone take about 12 minutes (S-0466) and the story regenerates docs/users/flai-reference.md with make flai-reference, and the delivery moved by the same margin. The story started at 20:57Z, while this plan was written."
  by: planner-S-0477
  at: 2026-10-22T20:59:32Z
finalized:
  by: alex
  at: 2026-10-22T20:52:29Z
---
# S-0477 flai guard refuses a story's sub-agent flai components list and flai components serve, so a headless pass of the harness falls to the story's agent

## Goal

This story remediates [I-0165](../../../design/issues/I-0165-flai-guard-refuses-a-story-s-sub-agent-flai-components-list-and-flai-components-serve-so-a-headless-pass-of-the-harness-falls-to-the-story-s-agent.md), "flai guard refuses a story's sub-agent flai components list and flai components serve, so a headless pass of the harness falls to the story's agent". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0165 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0165 is closed with `flai issue close I-0165 --reason` saying what fixed it

## Tasks
- T-2032 flai guard passes a sub-agent flai components list and flai components serve and keeps refusing components config, with tests
- T-2034 The guard's help, the generated reference, the user docs, and the design say a sub-agent may list and serve components
- T-2036 I-0165 is closed with the reason that names the cliReads entry and the test that reproduces it
- T-2043 flai guard lets a story's sub-agent run flai components list and serve, a read-only one list alone, and the strategic agents list alone
- T-2044 Close I-0165 saying the guard now lets a story's sub-agent list and serve components

## Notes

Cost of delay inputs set by flai from I-0165. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-19T05:38:19Z, 2.3 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0477 on 2026-10-22. The cause is one missing entry: `cliReads` in `flai/internal/guard/guard.go` has no `components` key, so `refuse()` turns down every `flai components` form a sub-agent runs, through `bin/flai` and `scripts/flai.sh` alike. ADR-0141's write rule does not see these commands at all. The proposed fix is `"components": {"list", "serve"}` in `cliReads`, leaving `components config`, which writes `system-flow.yaml`, refused. I-0156 and its story S-0472 record the same cause from an earlier instance (S-0438): see the plan thread.

Touches, and where each came from:

| Touch | From | Why |
|-------|------|-----|
| `flai/internal/guard/guard.go` | declared in I-0165 | `cliReads` gains `components` |
| `flai/internal/guard/guard_test.go` | layout | `TestASubAgentRunsChecksButNotWrites` and the test guard `g`'s command list, which lacks `components` and so skips it |
| `flai/cmd/guard.go` | co-change 60% | `flai guard --help` lists every read a sub-agent may run |
| `flai/cmd/guard_test.go` | co-change 43% | `TestGuard` runs the real command list end to end |
| `docs/users/flai-reference.md` | co-change 37% | generated from the help by `make flai-reference` |
| `docs/users/flai.md` | design (documentation.md: command docs change in the same story) | § Sub-agents says what a sub-agent may run |
| `design/system/agent-adapters.md` | design | § The guard names what the guard refuses a sub-agent |
| `design/system/component-library.md` | design | § The dev server says who runs `flai components serve` |
| `design/issues/I-0165-...md` | criterion 2 | `flai issue close` writes the issue |
| `design/issues/summary.md` | layout | regenerated by `flai issue close` |

No folder touch is kept: every file the story changes is named. `flai touches suggest` needed a starting path, since the story declared none; it was run from the two guard files. Co-change also listed `flai/internal/mcpserver/*`, `flai/internal/harness/*`, `docs/operators/settings.md`, and `design/system/flai-cli.md`, left out because the fix adds no MCP tool, setting, or command.

Forecast: `flai forecast` gave 4m with no touches and 18m with them (size 12, medium band, 87 s per unit). Set to 40m: the figure scales with files, but the close-out's integration and smoke tiers alone take about 12 minutes (S-0466's figures), and the reference must be regenerated. Delivery moved by the same margin from flai's 2026-10-22T21:21Z.

Cost of delay: `flai cod` gives 25.00 USD a week from the operator's input of 10m lost per 168h cycle at 150 USD an hour; it stands. Were S-0472 merged into this story, the two issues' inputs would make it 20m a cycle, 50 USD a week.

### Accepted by the orchestrator

- Verified: e54d72a7cc0c3c3f66e2cff4c147f7a34e5c5a80
- Batch: 20261022-213246
- Enqueued verified: e54d72a7cc0c3c3f66e2cff4c147f7a34e5c5a80
- At: 2026-10-22T21:44:12Z

Verdict: pass at e54d72a7cc0c3c3f66e2cff4c147f7a34e5c5a80
- 1: flai/internal/guard/guard.go, flai/internal/guard/call.go, flai/internal/guard/guard_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/agent-adapters.md, design/system/component-library.md, design/system/flai-cli.md, design/adrs/0186-flai-guard-lets-every-agent-run-flai-components-list-and-a-story-s-task-sub.md, design/adrs/README.md
- 2: design/issues/I-0165-flai-guard-refuses-a-story-s-sub-agent-flai-components-list-and-flai-components-serve-so-a-headless-pass-of-the-harness-falls-to-the-story-s-agent.md, design/issues/summary.md
