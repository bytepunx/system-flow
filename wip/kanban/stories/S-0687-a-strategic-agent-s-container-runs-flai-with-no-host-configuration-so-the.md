---
id: S-0687
type: story
nature: remediation
title: A strategic agent's container runs flai with no host configuration, so the orchestrator's flai accept and flai verify inside it see no host actions and lack the test tiers' tools
status: review
owner: alex
created: 2026-11-10T08:37:14Z
updated: 2026-11-10T09:49:33Z
transitions:
  - to: ready
    at: 2026-11-10T08:49:53Z
    by: alex
  - to: in-progress
    at: 2026-11-10T08:50:09Z
    by: agent-S-0687
  - to: review
    at: 2026-11-10T09:49:33Z
    by: agent-S-0687
tags: [flai]
topics: [orchestration, review]
touches: [design/adrs/README.md, flai/internal/mcpserver/accept.go, flai/internal/mcpserver/accept_test.go, flai/internal/mcpserver/folder.go, flai/internal/guard/guard_rules.go, flai/internal/guard/guard.go, flai/internal/guard/orchestrate_test.go, flai/cmd/accept.go, flai/cmd/verify.go, flai/cmd/accept_orchestrator_test.go, flai/cmd/verify_test.go, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/roles/baseline/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0256-a-strategic-agent-s-container-runs-flai-with-no-host-configuration-so-the.md, design/issues/summary.md, design/adrs/0278-the-orchestrator-accepts-through-an-mcp-tool-that-runs-flai-accept-on-the-host.md, flai/cmd/mcp_http.go, flai/cmd/mcp.go, flai/internal/mcpserver/server.go, flai/internal/guard/guard_refusals.go, flai/internal/guard/call.go, flai/internal/hostapi/writes_outcome.go, flai/internal/roles/conditions.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder_test.go, flai/cmd/mcp_http_test.go, flai/internal/harness/harness_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3577
  priced_by: harness
  turns:
    - day: 2026-11-10
      ceremony: 5
      hand_edits: 3
      work: 132
  models:
    - model: claude-opus-5-5
      input: 586
      output: 132682
      cache_read: 101035678
      cache_write: 2228394
      cost: 35.2759
  strategic:
    - kind: orchestrator
      seconds: 128
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 30
          output: 261
          cache_read: 1022854
          cache_write: 68753
          cost: 0.2976
        - model: claude-sonnet-5-5
          input: 10
          output: 64
          cache_read: 137139
          cache_write: 50413
          cost: 0.1404
cost_of_delay:
  inputs:
    time_lost_per_cycle: 9m
    by: flai
    at: 2026-11-10T08:37:14Z
  value: 22.5
  by: planner-S-0687
  at: 2026-11-10T08:49:05Z
forecast:
  duration: 2h30m
  delivery: 2026-11-10T18:17:00Z
  basis: "Its own forecast of 2h30m; 28th in the pull order with an in-progress limit of 5, behind S-0685, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663, S-0665, S-0666, S-0667, S-0668, S-0669, S-0670, S-0671, S-0672, S-0673, S-0674, S-0675, S-0676, S-0677, S-0678, S-0679, S-0680, S-0681, S-0682, S-0683, S-0684 and S-0686."
  by: flai
  at: 2026-11-10T08:49:28Z
finalized:
  by: alex
  at: 2026-11-10T08:49:49Z
---
# S-0687 A strategic agent's container runs flai with no host configuration, so the orchestrator's flai accept and flai verify inside it see no host actions and lack the test tiers' tools

## Goal

This story remediates [I-0256](../../../design/issues/I-0256-a-strategic-agent-s-container-runs-flai-with-no-host-configuration-so-the.md), "A strategic agent's container runs flai with no host configuration, so the orchestrator's flai accept and flai verify inside it see no host actions and lack the test tiers' tools". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0256 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0256 is closed with `flai issue close I-0256 --reason` saying what fixed it

## Tasks
- T-2932 An ADR records that the orchestrator accepts and verifies through MCP calls that run on the host, and that its shell's accept and verify fail closed in a container
- T-2934 The orchestrator's MCP tool accept runs flai accept --by orchestrator on the host, with a dry run that answers the blockers
- T-2935 The orchestrator's shell flai accept and tier-running flai verify refuse inside a role's container, naming the MCP tools, with a test that reproduces I-0256
- T-2937 The orchestrator's role text and the strategic-agents convention, in every copy, accept through the MCP tool accept and verify through the MCP tool verify
- T-2938 The design and the user guide describe the orchestrator's acceptance through the host in a container, and I-0256 is closed

## Notes

Cost of delay inputs set by flai from I-0256. time_lost_per_cycle 9m: 10m per occurrence × 1 occurrence ÷ 1.1 cycles of 168h (first reported 2026-11-02T19:06:36Z, 7.6 days before this story).

### Planning

The remedy proposed, which the first task records as an ADR before any code: the orchestrator accepts through a new orchestrator-only MCP tool, `accept`, whose handler runs `flai accept --by orchestrator` on the host, where the MCP server runs with the host's configuration, as the host API's `accept.run` already does (`flai/internal/hostapi/writes_items.go`); it verifies through the MCP tool `verify`, the close-out job `flai serve` runs on the host (ADR-0229), and reads the result with `flai verify --last`, which only reads. Its shell's `flai accept` and tier-running `flai verify` then fail closed: run as the orchestrator with a `FLAI_CONFIG` in the role's cache folder (`.flai-cache/agents/<role>`), they refuse and name the MCP tools, instead of seeing every host action off and landing a story by hand while the merge queue batches (`flai/cmd/accept.go`, `acceptance`; `flai/cmd/verify.go`, `hostActionOn`). Copying the host actions into the role's cache was weighed and left: the cache is mounted read-write, so a shell write could turn the copy's actions on or off unseen by the guard, against ADR-0220.

Touches, each from the code layout an explorer read, since the story declared none and `flai touches suggest` had nothing to start from:

- Design: `design/adrs/README.md` and the new ADR (its file name predicted from the planned title; the story's agent renames the touch if the number or slug differs); `design/system/strategic-agents.md` (§ In a container, § Accepting a story), `design/system/flai-cli.md` (the strategic container section that states the `FLAI_CONFIG` gap).
- Layout: `flai/internal/mcpserver/accept.go` and `accept_test.go` (new), `flai/internal/mcpserver/folder.go` (registers the tools), `flai/internal/guard/guard_rules.go` and `guard.go` (the orchestrator's MCP tools and `accept_reviews`), `flai/internal/guard/orchestrate_test.go`; `flai/cmd/accept.go`, `flai/cmd/verify.go`, `flai/cmd/accept_orchestrator_test.go`, `flai/cmd/verify_test.go`.
- Layout: the orchestrator's role text in its three copies, `design/workflows/roles/orchestrator.md`, `template/root/design/workflows/roles/orchestrator.md`, `flai/internal/roles/baseline/orchestrator.md`; the convention in both copies, `design/conventions/strategic-agents.md` and `template/root/design/conventions/strategic-agents.md`; `template/CHANGELOG.md` for the template's change.
- Design: `docs/users/flai.md` (the MCP tools table, § When the orchestrator accepts, the container harness section).
- Criteria: `design/issues/I-0256-…md` and `design/issues/summary.md`, which `flai issue close` writes.

No folder touch is kept. No `after`: the story reads nothing another open story writes.

Figures:

- Forecast: adjusted from flai's 6m to 2h30m. flai sized the story by two criteria and no touches; the plan's five tasks over 23 files, three copies of the role text among them, are the shape of past orchestrator stories that took two to three hours. Delivery 2026-11-11T18:00:00Z allows for its 29th place in the pull order with an in-progress limit of 5.
- Cost of delay: 22.50 USD a week stands as `flai cod` gives it, from flai's input of 9m lost per 168h cycle at 150 USD an hour. One occurrence, seen in the design rather than a live run, gives no reason to move it.
