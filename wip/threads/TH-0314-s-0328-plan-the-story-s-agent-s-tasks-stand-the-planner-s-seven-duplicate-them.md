---
id: TH-0314
title: "S-0328 plan: the story's agent's tasks stand; the planner's seven duplicate them"
anchor:
  path: wip/kanban/stories/S-0328-add-permission-and-ability-to-orchestrate-to-trigger-planner.md
  item: S-0328
status: resolved
participants: [planner-S-0328, alex]
created: 2026-10-07T19:42:22Z
updated: 2026-10-07T20:16:32Z
---

# TH-0314 S-0328 plan: the story's agent's tasks stand; the planner's seven duplicate them

On wip/kanban/stories/S-0328-add-permission-and-ability-to-orchestrate-to-trigger-planner.md.

## Entries

### 2026-10-07T19:42:22Z planner-S-0328
Recommendation: keep the story's agent's plan (T-1174, T-1176 to T-1179) and cancel the planner's seven tasks, T-1168 to T-1173 and T-1175, as duplicates. Before cancelling them, the agent should widen its tasks' touches to cover the five files listed below.

S-0328 was already in progress when I was started. Its agent wrote its own tasks while I wrote mine, so the two plans overlap. I set the story's touches, its forecast (1h45m, delivery 2026-10-07T22:00Z), and its cost of delay value (800 USD a week, from your `penalty_per_week: 800`). My reasons are under `### Planning` in its Notes.

## The agent's plan (live)

| Layer | Task | After |
|-------|------|-------|
| 1 | T-1174 ADR (done), T-1176 permission, story candidates, flai's checks (in progress) | none |
| 2 | T-1177 guard, T-1178 prompt, definition, and convention, T-1179 design, guides, and changelog | T-1176 |

## The planner's plan (duplicates)

| Layer | Task | After | Duplicates |
|-------|------|-------|------------|
| 1 | T-1168 ADR | none | T-1174 |
| 2 | T-1169 manifest permission, T-1170 story candidates | T-1168 | T-1176 |
| 3 | T-1171 orchestrator plans a story, flai and guard | T-1169, T-1170 | T-1176, T-1177 |
| 4 | T-1172 orchestrator sets cost of delay inputs | T-1171 | T-1176, T-1177 |
| 5 | T-1173 prompt, definition, and convention; T-1175 design and guides | T-1171, T-1172 | T-1178, T-1179 |

## Files my tasks name that the agent's tasks do not

- `flai/internal/manifest/settings_test.go`
- `flai/cmd/plan_candidates_test.go`
- `flai/internal/mcpserver/items_write_test.go`
- `flai/cmd/edit_test.go`
- `flai/internal/workitem/promotable.go`: `Repo.OrchestratorPermits`, if flai's own check of the cost of delay edit goes there

## Assumptions and points for you

- **Cancelling.** I did not cancel any task. The planner never cancels without asking. The pulling agent may cancel tasks it would plan differently, and should say why in `## Decisions`.
- **Cost of delay.** T-1174's ADR lets the orchestrator answer and resolve a planner's threads, and set a backlog story's cost of delay, under `plan_backlog_stories` alone. My T-1172 had recommended also requiring `answer_threads: autonomous`. The ADR also narrows two rules: the convention's rule that money questions escalate to you, and its rule that the orchestrator never resolves a thread it did not open. Confirm that this is what criterion 3 meant.
- **Off by default.** `system-flow.yaml` is not changed. Turning the permission on is yours.
- **`.claude/` files.** T-1178 changes `.claude/agents/orchestrator.md`, so the branch is accepted by you only (ADR-0106).

### 2026-10-07T20:16:32Z alex
Resolved.
