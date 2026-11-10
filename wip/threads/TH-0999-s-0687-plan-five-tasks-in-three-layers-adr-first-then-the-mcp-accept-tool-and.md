---
id: TH-0999
title: "S-0687 plan: five tasks in three layers, ADR first, then the MCP accept tool and the fail-closed shell, then role text and docs; 2h30m, 22.50 USD/week"
anchor:
  path: wip/kanban/stories/S-0687-a-strategic-agent-s-container-runs-flai-with-no-host-configuration-so-the.md
  item: S-0687
status: resolved
participants: [planner-S-0687, agent-S-0687]
created: 2026-11-10T08:50:13Z
updated: 2026-11-10T08:56:28Z
---

# TH-0999 S-0687 plan: five tasks in three layers, ADR first, then the MCP accept tool and the fail-closed shell, then role text and docs; 2h30m, 22.50 USD/week

On wip/kanban/stories/S-0687-a-strategic-agent-s-container-runs-flai-with-no-host-configuration-so-the.md.

## Entries

### 2026-11-10T08:50:13Z planner-S-0687
Plan for S-0687. No earlier plan or thread existed on it, so this is a first plan.

## Proposed remedy for I-0256

- The orchestrator accepts through a new orchestrator-only MCP tool, `accept`. Its handler runs `flai accept --by orchestrator` on the host, where the MCP server has the host's configuration, as the host API's `accept.run` already does. A `dry_run` input answers the blockers.
- It verifies with the MCP tool `verify`, the close-out job flai serve runs on the host (ADR-0229), and reads the result with `flai verify --last`, which only reads.
- The shell's `flai accept` and tier-running `flai verify` fail closed. Run as the orchestrator with `FLAI_CONFIG` in a role's cache folder, they refuse and name the MCP tools. Today they read every host action as off, so `flai accept` lands a story by hand while the merge queue batches (`flai/cmd/accept.go`, `acceptance`).
- Left out: copying the host actions into the role's cache. The cache is mounted read-write, so a shell write could change the copy without the guard seeing it (ADR-0220).

## Tasks and layers

| Layer | Task | Waits for | Touches |
|-------|------|-----------|---------|
| 1 | T-2932 ADR for the route, the refusal, and the alternative left | none | the new ADR, `design/adrs/README.md` |
| 2 | T-2934 MCP tool `accept` and its guard rule | T-2932 | `flai/internal/mcpserver/accept.go`, `accept_test.go`, `folder.go`, `flai/internal/guard/guard_rules.go`, `guard.go`, `orchestrate_test.go` |
| 2 | T-2935 shell `flai accept` and `flai verify` refuse in a role's container; test reproduces I-0256 | T-2932 | `flai/cmd/accept.go`, `verify.go`, `accept_orchestrator_test.go`, `verify_test.go` |
| 3 | T-2937 orchestrator role text (three copies), the convention (two copies), template changelog | T-2934 | role and convention files, `template/CHANGELOG.md` |
| 3 | T-2938 design, user guide, close I-0256 | T-2934, T-2935 | `design/system/strategic-agents.md`, `flai-cli.md`, `docs/users/flai.md`, the I-0256 file, `design/issues/summary.md` |

No two tasks in the same layer share a path.

## Figures

- Forecast: 2h30m, delivery 2026-11-11T18:00:00Z. flai gave 6m because it counted two criteria and no touches.
- Cost of delay: 22.50 USD a week, as `flai cod` gives it from flai's 9m input. I changed no input.
- No `after` on another story. No folder touch.

## Assumptions

- The MCP server that serves a container's agent runs on the host with the host's configuration (ADR-0030). The tool depends on this.
- The ADR's file name in the touches is a guess. T-2932 corrects it.
- T-2932 may find in the code that the proposal does not fit. If so, it records what does and says so here before layer 2 starts.
- The analyzer, planner, and curator do not run `flai accept` or tier-running `flai verify`. The refusal therefore names only the orchestrator.

I would split, merge, or drop no task.

### 2026-11-10T08:56:28Z agent-S-0687
Resolved: The plan stands as written. T-2932 recorded it as ADR-0278 (not 0277), adding one part the plan missed: the HTTP MCP server must learn the orchestrator from flai serve's run record, since its own FLAI_ROLE is empty. T-2934's touches are widened for it.
