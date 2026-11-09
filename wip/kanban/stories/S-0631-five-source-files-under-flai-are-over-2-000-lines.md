---
id: S-0631
type: story
nature: improvement
title: Five source files under flai are over 2,000 lines
status: ready
owner: alex
created: 2026-11-07T23:40:53Z
updated: 2026-11-09T23:22:43Z
transitions:
  - to: ready
    at: 2026-11-09T22:05:09Z
    by: alex
tags: [cli]
topics: [cli, testing, docs]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_run.go, flai/internal/hostapi/writes_actions.go, flai/internal/hostapi/writes_spec.go, flai/internal/hostapi/writes_items.go, flai/internal/hostapi/writes_streams.go, flai/internal/hostapi/writes_host.go, flai/internal/hostapi/writes_describe.go, flai/internal/hostapi/writes_methods.go, flai/internal/hostapi/writes_outcome.go, flai/internal/hostapi/settings.go, flai/internal/serve/agents.go, flai/internal/serve/agents_state.go, flai/internal/serve/agents_launch.go, flai/internal/serve/agents_restart.go, flai/internal/serve/agents_start.go, flai/internal/serve/agents_activity.go, flai/internal/mcpserver/server.go, flai/cmd/serve_actions.go, flai/cmd/serve_agent.go, flai/cmd/serve_checks.go, flai/cmd/serve_components.go, flai/cmd/serve_harness.go, flai/cmd/serve_provider.go, flai/cmd/serve_notify.go, flai/cmd/serve_mcp.go, flai/cmd/serve_view.go, flai/cmd/serve_agent_run.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_strategic.go, flai/internal/manifest/manifest_roles.go, flai/internal/manifest/manifest_accept.go, flai/internal/manifest/manifest_kluster.go, flai/internal/manifest/manifest_tests.go, flai/internal/check/check.go, flai/internal/check/check_layout.go, flai/internal/check/check_items.go, flai/internal/check/check_docs.go, flai/internal/check/check_roles.go, flai/internal/check/check_records.go, flai/internal/check/sourcesize_test.go, flai/internal/guard/guard.go, flai/internal/guard/guard_tables.go, flai/internal/guard/guard_rules.go, flai/internal/guard/guard_refusals.go, flai/internal/guard/guard_shell.go, flai/internal/watch/watch.go, design/system/flai-cli.md, design/system/server-performance.md, design/system/design-system-tool.md, design/system/agent-adapters.md, design/system/strategic-agents.md, design/system/agent-coordination.md, design/system/flaiover-dashboard.md, design/system/agent-prose.md, design/system/release-signing.md, design/issues/I-0298-five-source-files-under-flai-are-over-2-000-lines.md, design/issues/summary.md]
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
    - kind: orchestrator
      seconds: 241
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 31
          output: 302
          cache_read: 1942029
          cache_write: 130331
          cost: 0.5513
        - model: claude-sonnet-5-5
          input: 1
          output: 11
          cache_read: 14232
          cache_write: 8148
          cost: 0.0166
    - kind: analyzer
      seconds: 41
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3102
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 14
          output: 2840
          cache_read: 759049
          cache_write: 13854
          cost: 0.278
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: orchestrator
    at: 2026-11-09T21:59:47Z
  value: 75
  by: planner-S-0631
  at: 2026-11-09T22:03:45Z
forecast:
  duration: 3h
  delivery: 2026-11-10T03:36:00Z
  basis: "Its own forecast of 3h; 5th in the pull order with an in-progress limit of 5, behind S-0633, S-0636, S-0638, S-0640, S-0625, S-0630 and S-0628."
  by: flai
  at: 2026-11-09T23:22:43Z
finalized:
  by: alex
  at: 2026-11-09T22:05:00Z
---
# S-0631 Five source files under flai are over 2,000 lines

## Goal

This story remediates [I-0298](../../../design/issues/I-0298-five-source-files-under-flai-are-over-2-000-lines.md), "Five source files under flai are over 2,000 lines", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0298 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0298 is closed with `flai issue close I-0298 --reason` saying what fixed it

## Tasks
- T-2794 Split flai/internal/hostapi/writes.go by responsibility, cutting itemSpecs in three
- T-2795 Split flai/internal/serve/agents.go by responsibility
- T-2796 Split flai/cmd/serve_actions.go into one file per flai serve subcommand family
- T-2797 Split flai/internal/manifest/manifest.go into one file per manifest block
- T-2798 Split flai/internal/check/check.go into its runner and one file per family of rules
- T-2799 Split flai/internal/guard/guard.go into its decision, its rule tables, and its shell parsing
- T-2800 Give the nolint directive in flai/internal/watch/watch.go its reason
- T-2801 Test that no non-test Go file under flai is over 2,000 lines
- T-2802 Point the design documents at the files the split code moved to, and close I-0298

## Notes

The analyzer found I-0298: the analysis the goal links gives its evidence.

Strategic usage carried over from I-0298: analyzer 41s, 0.8156 USD.

### Planning

Proposed fix: split each Go source file over 2,000 lines by responsibility, within its own package. Move declarations without changing what they do. Then add a test that fails when any non-test Go file under `flai/` is over 2,000 lines.

Six files are over the limit now, not five. All five have grown since the analysis, and `flai/internal/guard/guard.go` is at 2,051 lines.

| File | Lines on 2026-11-09 |
|------|---------------------|
| `flai/internal/hostapi/writes.go` | 3,626 |
| `flai/internal/serve/agents.go` | 3,298 |
| `flai/cmd/serve_actions.go` | 2,754 |
| `flai/internal/manifest/manifest.go` | 2,449 |
| `flai/internal/check/check.go` | 2,235 |
| `flai/internal/guard/guard.go` | 2,051 |

The test would fail on `guard.go`, so the story splits it too.

The test covers non-test files only. Seven `_test.go` files are over 2,000 lines, and splitting them is proposed on the plan's thread as a follow-up story.

The issue's nolint without a reason is now at `flai/internal/watch/watch.go:94`, `//nolint:nilerr`. The story gives it its reason.

Touches, and where each came from:

- Declared: none. The story declared no touches.
- Design (I-0298 and its analysis): the five files the issue names, and `flai/internal/watch/watch.go`.
- Layout (the explorer's read of each package):
  - `flai/internal/guard/guard.go`.
  - The new files each split creates, named as its task proposes. A task's touches widen to the names it ends up using when it closes.
  - `flai/cmd/serve_mcp.go`, an existing file that takes the MCP subcommands from `serve_actions.go`.
  - `flai/internal/check/sourcesize_test.go`, the new size test.
  - The two Go comments that name moved code by file: `flai/internal/hostapi/settings.go:615` and `flai/internal/mcpserver/server.go:1041`.
  - The design documents that name these files by path: `design/system/flai-cli.md`, `server-performance.md`, `design-system-tool.md`, `agent-adapters.md`, `strategic-agents.md`, `agent-coordination.md`, `flaiover-dashboard.md`, `agent-prose.md`, and `release-signing.md`.
  - The issue file and `design/issues/summary.md`, which `flai issue close` writes.
- Co-change: nothing taken. `flai touches suggest` lists the user docs and tests that features beside these files changed. A move that keeps behaviour changes neither. Every split keeps its file's name, so `flai/internal/verify/select_test.go:554`, which stats `writes.go`, still passes.
- Accepted ADRs and digest entries that name these files stay as they are.

No folder touches: each task names the files it creates.

Figures:

- Forecast 3h. flai forecast gave 1h38m. Its size counts touches, not lines: about 16,400 lines move across six packages, `itemSpecs` (`writes.go:989-2464`) is one 1,476-line function to cut into three, and each split runs its package's tests and lint.
- Cost of delay 75 USD a week, as flai cod works out from the `time_lost_per_cycle` of 30m that TH-0944 set: 30m per 168h cycle at 150 USD an hour.

No `after`: the story reads nothing another open story writes.
