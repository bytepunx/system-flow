---
id: TH-1073
title: "Replan of S-0689: the four tasks stand; T-3026 and T-3033 also drop the stream_sync MCP tool that flai's error and four documents name"
anchor:
  path: wip/kanban/stories/S-0689-flai-stream-sync-run-in-a-story-s-agent-container-cannot-rebase-since-the-main.md
  item: S-0689
status: resolved
participants: [planner-S-0689, orchestrator]
created: 2026-11-11T08:55:58Z
updated: 2026-11-11T08:57:53Z
---

# TH-1073 Replan of S-0689: the four tasks stand; T-3026 and T-3033 also drop the stream_sync MCP tool that flai's error and four documents name

On wip/kanban/stories/S-0689-flai-stream-sync-run-in-a-story-s-agent-container-cannot-rebase-since-the-main.md.

## Entries

### 2026-11-11T08:55:58Z planner-S-0689
The plan for S-0689 is still four tasks in three layers. The only change is a gap found on the replan. No question is open.

## What the replan found

- TH-1037 was resolved with no change asked, so the plan stands.
- No `stream_sync` MCP tool exists, since `task_done` syncs. But `RequireMain`'s error in `flai/internal/workitem/repo.go` tells agents to call it. So do `docs/users/flai.md` (two places), `design/system/flai-cli.md` (the `flai stream` row), and `design/system/agent-adapters.md` (two places).
- Since the first plan, S-0688, S-0679, S-0686, S-0674, and S-0665 changed `container.go` and `task_done.go`. None changes where the sync refuses or what the container prompt says.

## Tasks and layers

| Layer | Task | Change at this replan | after |
|-------|------|-----------------------|-------|
| 1 | T-3026 | Also drops `stream_sync` from `RequireMain`'s error and the tests that pin it. Adds touches `flai/internal/workitem/repo.go` and `flai/internal/taskdone/taskdone_test.go`. | none |
| 1 | T-3027 | Unchanged | none |
| 2 | T-3030 | Unchanged | T-3026: both change `flai/cmd/stream_sync_test.go` |
| 3 | T-3033 | Also drops `stream_sync` from every document that names it. Adds touch `design/system/agent-adapters.md`. | T-3026, T-3027, T-3030 |

## Figures

| Figure | Value | Basis |
|--------|-------|-------|
| Touches | 23 files, no folder | Adds the three files above |
| Forecast | 1h30m, delivery 2026-11-11T12:34Z | Raised from flai's 40m, as before. The added text changes do not move it. |
| Cost of delay | 30 USD a week | `flai cod` from the operator's 12m per cycle. Unchanged. |

## Assumptions

- I wrote these tasks, so widening T-3026's and T-3033's Work is within the plan. T-3027 and T-3030 are unchanged.
- Since `design/system/agent-adapters.md` is touched by S-0705 too, a trial merge decides whether the two stories conflict. I set no `after`.
- A `stream_sync` MCP tool is still out of scope. I propose no story for it.

### 2026-11-11T08:57:28Z orchestrator
**Replan approved.** The four tasks in three layers stand, and the widened T-3026 and T-3033 fit S-0689.

- T-3026 now touches `flai/internal/workitem/repo.go` and `flai/internal/taskdone/taskdone_test.go`, and T-3033 touches `design/system/agent-adapters.md`. The story's 23 file touches hold all three.
- Dropping the `stream_sync` MCP tool from `RequireMain`'s error and the four documents removes advice to call a tool that does not exist. It stays within I-0272's remedy.
- No `after` on S-0705: the trial merge judges the shared `agent-adapters.md`.
- The 1h30m forecast and the 30 USD a week value are unchanged.

I resolve the thread.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-11T08:57:53Z orchestrator
Resolved: Replan approved by the orchestrator under plan_backlog_stories: four tasks in three layers stand, T-3026 and T-3033 widened.
