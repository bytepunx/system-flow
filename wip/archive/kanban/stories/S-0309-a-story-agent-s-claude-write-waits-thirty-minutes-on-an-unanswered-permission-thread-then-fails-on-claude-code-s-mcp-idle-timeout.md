---
id: S-0309
type: story
nature: improvement
title: A story agent's .claude/ write waits thirty minutes on an unanswered permission thread, then fails on Claude Code's MCP idle timeout
status: done
owner: alex
created: 2026-10-07T06:48:44Z
updated: 2026-10-08T06:12:43Z
transitions:
  - to: ready
    at: 2026-10-08T04:23:09Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-08T05:52:20Z
    by: agent-S-0309
  - to: review
    at: 2026-10-08T06:11:28Z
    by: agent-S-0309
  - to: done
    at: 2026-10-08T06:12:43Z
    by: orchestrator
tags: []
topics: [cli, agents]
touches: [flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/delegation.md, docs/users/flai.md, design/system/flai-cli.md, design/adrs, design/adrs/README.md, design/issues/I-0103-a-story-agent-s-claude-write-waits-thirty-minutes-on-an-unanswered-permission-thread-then-fails-on-claude-code-s-mcp-idle-timeout.md, design/issues/summary.md, design/issues/I-0108-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1160
  turns:
    - day: 2026-10-08
      ceremony: 1
      hand_edits: 1
      work: 35
  models:
    - model: claude-opus-5-5
      input: 200
      output: 62366
      cache_read: 8819149
      cache_write: 402844
      cost: 5.6425
  strategic:
    - kind: planner
      seconds: 296
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 178
          output: 32
          cache_read: 1205366
          cache_write: 89708
          cost: 0.2551
        - model: claude-opus-5-5
          input: 58
          output: 463
          cache_read: 2537009
          cache_write: 106840
          cost: 1.0999
    - kind: orchestrator
      seconds: 1564
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 103
          output: 1567
          cache_read: 14417689
          cache_write: 52480
          cost: 3.5658
        - model: claude-sonnet-5-5
          input: 12
          output: 62
          cache_read: 162614
          cache_write: 47869
          cost: 0.1814
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-10-07T06:48:44Z
  value: 75
  by: planner-S-0309
  at: 2026-10-07T23:28:00Z
forecast:
  duration: 35m
  delivery: 2026-10-08T05:40:00Z
  basis: "Its own forecast of 35m; 2nd in the pull order with an in-progress limit of 3, behind S-0232, S-0318, S-0320 and S-0319."
  by: flai
  at: 2026-10-08T05:02:40Z
finalized:
  by: orchestrator
  at: 2026-10-08T04:23:05Z
---
# S-0309 A story agent's .claude/ write waits thirty minutes on an unanswered permission thread, then fails on Claude Code's MCP idle timeout

## Goal

This story remediates [I-0103](../../../design/issues/I-0103-a-story-agent-s-claude-write-waits-thirty-minutes-on-an-unanswered-permission-thread-then-fails-on-claude-code-s-mcp-idle-timeout.md), "A story agent's .claude/ write waits thirty minutes on an unanswered permission thread, then fails on Claude Code's MCP idle timeout". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0103 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0103 is closed with `flai issue close I-0103 --reason` saying what fixed it

## Tasks
- T-1274 Record the remedy for I-0103 in an ADR refining ADR-0086 and ADR-0097
- T-1275 permission_prompt bounds its wait, keeps an unanswered thread open, and takes its answer on the retry
- T-1276 The start prompt and delegation.md tell the agent to retry a protected write after its thread is answered
- T-1277 Document permission_prompt's bounded wait and retry in the user guide and the CLI design
- T-1278 Close I-0103 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0103. time_lost_per_cycle 30m: 30m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T04:55:43Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Touches, and where each came from:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/mcpserver/permission.go` | layout | `askOperator` and `awaitAnswer` hold the call with no bound and settle the thread as refused when it ends |
| `flai/internal/mcpserver/permission_test.go` | co-change | changed in 5 of 5 commits with `permission.go`; the test that reproduces I-0103 goes here |
| `flai/internal/harness/harness.go` | layout | the start prompt's `.flai-cache/` and `cp` workaround for the thirty-minute wait |
| `flai/internal/harness/harness_test.go` | layout | checks the prompt's wording on that workaround |
| `design/conventions/delegation.md` | design | its project additions carry the same workaround |
| `docs/users/flai.md` | design | § Writes to paths Claude Code protects says the agent waits until you answer |
| `design/system/flai-cli.md` | design | the `flai mcp` row describes `permission_prompt` |
| `design/adrs/` | design | the new ADR refining ADR-0086 and ADR-0097; its number and slug are unknown until written |
| `design/adrs/README.md` | design | the ADR index |
| `design/issues/I-0103-…md`, `design/issues/summary.md` | criteria | written by `flai issue close` |

`flai touches suggest S-0309` declared nothing, so it ran from `permission.go`; `folder.go` (2 of 5 commits) was left out because the fix does not change how a folder server finds the project.

The one folder touch kept is `design/adrs/`, because the ADR's file name is not known until it is written. The template's `delegation.md` lacks the workaround text, so it is not touched. `wait_for_events`' `endWhy` in `flai/internal/mcpserver/server.go` already counts an open permission thread as the agent's own question, so `server.go` is not touched.

Forecast: flai gave 17m. It is raised to 35m from S-0299 (31m) and S-0257 (36m), the permission stories closest in scope: this one adds an ADR, the bounded wait and the retry with tests, the prompt, and two documents. Delivery moves by the 18m added, to 2026-10-08T05:45Z.

Cost of delay: 75 USD a week as `flai cod` computes it from the operator's input (30m lost per 168h cycle at 150 USD an hour). It stands: one instance so far, but every protected write made while the operator is away hits it.

### Accepted by the orchestrator

- Verified: 32badc936063f9af700e0359241115a1210d6d3e
- At: 2026-10-08T06:12:43Z

Verdict: accept. flai verify passed every step at the branch head 32badc93, and the verifier matched both criteria to the diff.
- 1: design/adrs/0124-permission-prompt-holds-a-write-at-most-four-minutes-then-refuses-it-and-leaves.md, design/adrs/README.md, flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/delegation.md, docs/users/flai.md, design/system/flai-cli.md
- 2: design/issues/I-0103-a-story-agent-s-claude-write-waits-thirty-minutes-on-an-unanswered-permission-thread-then-fails-on-claude-code-s-mcp-idle-timeout.md, design/issues/summary.md
