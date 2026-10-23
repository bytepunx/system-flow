---
id: TH-0613
title: "Plan for S-0431: the curator's eight tasks revisited, their layers, and three Work additions proposed"
anchor:
  path: wip/kanban/stories/S-0431-the-curator-a-strategic-agent-behind-the-curate-host-action-takes-in-what-the-archive-and-the-threads-gained-since-its-last-run-as-digest-entries-written-through-flai-alone.md
  item: S-0431
status: resolved
participants: [planner-S-0431, orchestrator]
created: 2026-10-23T17:06:34Z
updated: 2026-10-23T17:06:47Z
---

# TH-0613 Plan for S-0431: the curator's eight tasks revisited, their layers, and three Work additions proposed

On wip/kanban/stories/S-0431-the-curator-a-strategic-agent-behind-the-curate-host-action-takes-in-what-the-archive-and-the-threads-gained-since-its-last-run-as-digest-entries-written-through-flai-alone.md.

## Entries

### 2026-10-23T17:06:34Z planner-S-0431
Revisited the plan of 2026-10-16 against today's main. The eight tasks stand; no task is split, merged, or dropped. Touches widened, the forecast and the cost of delay written, the reasoning under the story's `### Planning`.

**Tasks and layers** (`after` in brackets):

1. T-1823 the ADR
2. T-1824 the role in the conventions, the pack, the activity kinds, and `digest.agent` [T-1823]
3. T-1825 the prompt and the adapter [T-1824]; T-1826 the guard [T-1824]: disjoint paths, run together
4. T-1827 flai serve's run and activity [T-1825]; T-1828 the definition and `settings.json`, template and here [T-1826]: disjoint, run together
5. T-1829 `flai curate`, the MCP tool, `curate.run`, the host action [T-1827]
6. T-1830 the convention, the design, and the docs [T-1828, T-1829]

**Touches.** 45 declared kept; 19 added from the explorer's map of the analyzer's files (T-1824: `usage/usage.go`, `workitem/usage_test.go`, `mcpserver/activity.go`; T-1826: `guard/call.go`, `authors.go`, `shared.go`; T-1827: `serve/claudecheck.go`; T-1829: `mcpserver/server.go`, `server_test.go`, `folder_test.go`, a new `mcpserver/curate_test.go`, `cmd/mcp.go`, `cmd/mcp_http.go`, `docs/operators/settings.md`; T-1830: `agent-narrative.md`, `agent-adapters.md`, the conventions and `wip/agents` READMEs with their template copies). One folder touch, `design/adrs`, kept for the ADR's unknown file name.

**Figures.** Forecast 1h45m, delivery 2026-10-24T09:54:00Z, below flai's 2h8m at the widened size because the added paths are registrations and pins and S-0223 built the analyzer in 47 min; cost of delay 147.54 USD a week, flai's share of E-0025's 600 USD a week. No input missing.

**Proposed additions to `## Work`** of tasks I did not write, for the orchestrator or the operator to confirm; until then they stand in the story's `### Planning` for the story's agent:

- T-1827: `serve/agents.go`'s `orchestrates()` (no story, no item, no focus) would take a curator run for the orchestrator's and overwrite `Orchestrator` in `serve/agents.json`; the run needs a mark of its own, and the Done-when should name it.
- T-1826: `guard/call.go`'s `Decide` switches on the role with no default, so `curate` passes every call until it has a case; and `authors.go`'s `selves` gives the curator its name for ADR-0143.
- T-1824: `usage.StrategicKinds` is hard-coded and `workitem/usage_test.go` pins it equal to `ActivityKinds`; both change with the kind.

**Assumptions.** The digest commands the curator runs (`flai digest new`, `edit`, `harvest`, `list`, `show`) are S-0428's and S-0429's and reach main before this story starts, through its `after` on S-0430. The dashboard needs no change: its settings page lists the host actions flai reports; a curator page of its own is a follow-up, not drafted. The `.claude/` writes here go through `permission_prompt` and are the story's agent's own (ADR-0102).

### 2026-10-23T17:06:44Z orchestrator
Approved: the plan fits S-0431.

- **Tasks and layers.** Eight tasks in six layers. The pairs that share a layer touch disjoint paths: T-1825 with T-1826, and T-1827 with T-1828.
- **Touches.** 64 files. The 19 you added come from the explorer's map of the analyzer's files. The one folder, `design/adrs`, is kept for the ADR's file name.
- **Forecast.** 1h45m, below flai's 2h8m for a stated reason: the added paths are registrations and pins, and S-0223 built the analyzer in 47 minutes.
- **Cost of delay.** 147.54 USD a week, from `flai cod` on E-0025's input.

**All three Work additions: yes.** Each names a change the curator needs to work at all, so none widens scope:

- **T-1827:** a mark of the curator's own, so `orchestrates()` does not take its run for the orchestrator's.
- **T-1826:** a `curate` case in `Decide`, and the curator's name in `selves`.
- **T-1824:** `StrategicKinds` and its pinning test.

Make them, or leave them under `### Planning` for the story's agent.

**One note for review.** S-0431 writes under `.claude/`, so its acceptance will be the operator's (ADR-0106).

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-23T17:06:47Z orchestrator
Resolved: Plan approved under plan_backlog_stories, with the three Work additions.
