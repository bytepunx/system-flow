---
id: S-0556
type: story
nature: improvement
title: flai guard still lists the verify MCP tool among a sub-agent's reads, though it now asks flai serve for a close-out job that commits in the story's worktree
status: backlog
owner: alex
created: 2026-11-01T10:44:38Z
updated: 2026-11-07T03:43:33Z
transitions: []
tags: [flai, template]
topics: [guard, testing]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/internal/guard/orchestrate_test.go, flai/internal/roles/baseline/orchestrator.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, design/conventions/delegation.md, template/root/design/conventions/delegation.md, template/CHANGELOG.md, design/system/flai-cli.md, design/system/strategic-agents.md, design/issues/I-0240-flai-guard-still-lists-the-verify-mcp-tool-among-a-sub-agent-s-reads-though-it.md, design/issues/summary.md]
agent:
  harness: container
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 586
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 192
          output: 34854
          cache_read: 12498755
          cache_write: 420234
          cost: 5.5392
    - kind: orchestrator
      seconds: 186
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 46
          output: 324
          cache_read: 1798520
          cache_write: 210094
          cost: 0.4791
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: orchestrator
    at: 2026-11-05T01:40:15Z
  value: 12.5
  by: planner-S-0556
  at: 2026-11-05T01:41:54Z
forecast:
  duration: 40m
  delivery: 2026-11-07T08:55:00Z
  basis: "Its own forecast of 40m; 11th in the pull order with an in-progress limit of 5, behind S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0241, S-0543 and S-0555."
  by: flai
  at: 2026-11-07T03:29:11Z
finalized:
  by: orchestrator
  at: 2026-11-05T01:43:38Z
---
# S-0556 flai guard still lists the verify MCP tool among a sub-agent's reads, though it now asks flai serve for a close-out job that commits in the story's worktree

## Goal

This story remediates [I-0240](../../../design/issues/I-0240-flai-guard-still-lists-the-verify-mcp-tool-among-a-sub-agent-s-reads-though-it.md), "flai guard still lists the verify MCP tool among a sub-agent's reads, though it now asks flai serve for a close-out job that commits in the story's worktree". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0240 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0240 is closed with `flai issue close I-0240 --reason` saying what fixed it

## Tasks
- T-2459 Give the orchestrator the MCP verify tool by name, in its role file and the guard, before verify leaves the reads
- T-2460 Refuse every sub-agent, and every strategic agent but the orchestrator, the MCP verify tool and flai verify --job
- T-2461 Say in the design, the delegation convention, and the template's changelog that the guard refuses a sub-agent the close-out job
- T-2462 Close I-0240 with the reason naming the guard's refusal of verify and flai verify --job

## Notes

I-0240's cost "0m" is not a duration longer than zero, so no time_lost_per_cycle was derived from it.

### Planning

The proposed fix, which T-2459 records in the narrative first: take `verify` out of `guard.MCPReads`, refuse a sub-agent `flai verify --job`, and keep the tool for the story's agent, whose own calls the guard passes, and for the orchestrator, whose role file then names it. D-0446 recorded why S-0546 kept the tool among the reads, and anticipated this narrowing: the two calls S-0546 kept it for stay, the orchestrator's verify of a story in review and a reader's `flai verify --last`. Cite D-0446 in the narrative's `## Decisions`. The plan was approved on TH-0827.

Touches, all files, no folder:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/guard/guard.go` | design | I-0240 names `guard.MCPReads` and `waitsOnSubAgents`; `flagReads` holds the `flai verify` rule; `MCPOrchestrates` mirrors the orchestrator's tools |
| `flai/internal/guard/guard_test.go` | layout | holds the test that lets a verifier call MCP `verify`, to be turned around |
| `flai/internal/guard/orchestrate_test.go` | layout | the orchestrator's guard tests |
| `flai/internal/roles/baseline/orchestrator.md`, `design/workflows/roles/orchestrator.md`, `template/root/design/workflows/roles/orchestrator.md` | layout | the role file's `tools:` list does not name `verify`, which the orchestrator reaches only through `MCPReads` today (`asRole` in `roles.go`) |
| `design/conventions/delegation.md`, `template/root/design/conventions/delegation.md` | design | D-0446 names the delegation convention; its rule "unless your prompt asks" becomes the guard's refusal |
| `template/CHANGELOG.md` | co-change | the template's role file and convention change |
| `design/system/flai-cli.md` | design | the `flai guard` row lists `verify` among `guard.MCPReads` |
| `design/system/strategic-agents.md` | co-change | 5% of commits on these paths; lists the orchestrator's tools and its guard |
| `design/issues/I-0240-...md`, `design/issues/summary.md` | declared by criterion 2 | `flai issue close` writes both |

`flai touches suggest` also ranked `docs/users/flai.md`, `docs/operators/index.md`, and other hub documents by co-change. None of them says a sub-agent may call `verify`, so none is a touch.

No `after`: the story reads nothing an open story writes. S-0546, whose MCP `verify` this builds on, is done.

Forecast: `flai forecast` gave 22m. I raised it to 40m. The story first proposes narrowing a rule D-0446 recorded, and it keeps three copies of the orchestrator's role file identical. Delivery moved by the same 18m.

Cost of delay: 12.50 USD a week, as `flai cod` works it out. The input is the `time_lost_per_cycle` of 5m that the orchestrator set on TH-0827, from the recommendation there: 5m per 168h cycle at 150 USD an hour. It stands unadjusted: no instance has cost time yet, and the risk is a stray commit that takes a few minutes to undo.
