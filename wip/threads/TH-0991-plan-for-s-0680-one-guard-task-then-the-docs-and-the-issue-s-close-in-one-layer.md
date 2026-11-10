---
id: TH-0991
title: "Plan for S-0680: one guard task, then the docs and the issue's close in one layer"
anchor:
  path: wip/kanban/stories/S-0680-flai-guard-refuses-the-orchestrator-flai-serve-status-a-read-because-its-reads.md
  item: S-0680
status: open
participants: [planner-S-0680]
created: 2026-11-10T08:47:45Z
updated: 2026-11-10T08:47:45Z
---

# TH-0991 Plan for S-0680: one guard task, then the docs and the issue's close in one layer

On wip/kanban/stories/S-0680-flai-guard-refuses-the-orchestrator-flai-serve-status-a-read-because-its-reads.md.

## Entries

### 2026-11-10T08:47:45Z planner-S-0680
S-0680 is planned as three tasks in two layers.

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-2901 The guard counts the read forms of flai serve among the reads | `flai/internal/guard/guard_tables.go`, `guard.go`, `hostconfig.go`, `guard_test.go`, `hostconfig_test.go`, `orchestrate_test.go` |
| 2 | T-2902 The design and the user guide say the read forms of flai serve are open to every agent (after T-2901) | `design/system/strategic-agents.md`, `docs/users/flai.md` |
| 2 | T-2903 Close I-0318 with what fixed it (after T-2901) | `design/issues/I-0318-…-reads.md`, `design/issues/summary.md` |

T-2902 and T-2903 share no path, so they run together.

Proposed fix: add `serve` to `flagReads` with a form function that passes only the read forms ADR-0220 lists. Those are `status`, `actions`, `journal`, `agent` bare or `show`, `agent harness` and `agent provider` bare or with a name alone, `checks show`, `components show`, `forge` bare or with a forge alone, `notify show`, `mcp show`, and `import list`. Every form that starts, stops, or writes the host stays refused.

Figures: forecast 45m, delivery 2026-11-10T12:19:00Z. flai gave 22m; I raised it for the exhaustive table test and the integration and smoke tiers. Cost of delay value 5 USD a week, from the operator's 2m-per-cycle input.

Assumptions:
- These read forms open to every agent the guard holds to its reads: the orchestrator, the planner, the analyzer, the curator, defined roles, and every sub-agent. I-0318 calls that scope the operator's choice. I took the story's creation from the issue as that choice, because `askForHost` already tells all of those agents the forms are open. Alternative: open `flai serve status` to the orchestrator alone, through `design/workflows/roles/orchestrator.md`. That would replace T-2901's table entry with a role-file entry, and the sub-agent tests would go.
- No new ADR. Earlier read additions (S-0428 digest, S-0491 subproject, S-0612 hostapi) took none. The story's agent may add one if it judges ADR-0220's note on sub-agents needs superseding. `design/adrs` is a shared path, so that would hold nothing.
- No `after` on another story, and no folder touch.
