---
id: TH-1067
title: "Plan of S-0721: three tasks in two layers; cost of delay inputs needed"
anchor:
  path: wip/kanban/stories/S-0721-every-mcp-server-flai-builds-resolves-all-its-tools-json-schemas-again-which-is.md
  item: S-0721
status: resolved
participants: [planner-S-0721, orchestrator, alex]
created: 2026-11-11T05:29:38Z
updated: 2026-11-11T05:32:07Z
---

# TH-1067 Plan of S-0721: three tasks in two layers; cost of delay inputs needed

On wip/kanban/stories/S-0721-every-mcp-server-flai-builds-resolves-all-its-tools-json-schemas-again-which-is.md.

## Entries

### 2026-11-11T05:29:38Z planner-S-0721
**Question: cost of delay inputs.** S-0721 has no cost of delay inputs and no epic, so `flai cod S-0721` cannot work out a value.

Recommended: `--time-lost-per-cycle 10m` alone, with no revenue and no penalty. The figure is my estimate: each close-out that runs `./internal/mcpserver` under `-race` spends about 600 s of CPU, and three at once took the host to a load of 84 and slowed every agent on it. Answer with your inputs, or confirm this one.

**Plan.**

| Layer | Task | Touches | After |
|---|---|---|---|
| 1 | T-3083 Set one package-level `mcp.SchemaCache` on both MCP servers, with a test that their tool schemas are unchanged | `flai/internal/mcpserver/server.go`, `folder.go`, `schema_cache_test.go` | none |
| 1 | T-3084 Say in go-libraries.md's go-sdk row that flai shares one SchemaCache | `design/tech/go-libraries.md` | none |
| 2 | T-3085 Measure the package's CPU before and after, with and without `-race`, and record it in the narrative | `wip/agents/S-0721.md` | T-3083 |

**Changes to the story.**

- Touches: added `flai/internal/mcpserver/schema_cache_test.go` and `wip/agents/S-0721.md`. The three declared touches are kept.
- Forecast: 30m, raised from flai's 13m for the timed `-race` runs.
- The reasons are under `### Planning` in the Notes.

**Assumptions.**

- Criterion 2's test goes in a new file, to stay off `server_test.go` and `folder_test.go`.
- `permission_prompt`'s pre-defined `permissionInputSchema` map is the only schema that is not inferred. T-3083 checks how the SDK caches it before relying on the cache.
- The before figure is measured again at the story's base, on the same host as the after figure, rather than taken from the Notes. The host's load moves the figures.
- No user or operator docs change: the cache changes no command or setting.

**Note.** The story's agent, agent-S-0721, started at 05:27 while this plan was being written. It may review these tasks, as `work-management.md` says. The start prompt named a thread `asked`, but no such thread exists on S-0721. I planned the story from its own words.

### 2026-11-11T05:30:33Z orchestrator (recommendation)
**Recommendation:** set `time_lost_per_cycle: 10m` on S-0721, as the planner recommends:

```bash
flai edit S-0721 --time-lost-per-cycle 10m
```

## The figure

- **The goal is time lost.** The cost is host CPU and slowed agents during `-race` close-outs, which is agent time, not revenue or a penalty.
- **Nothing gives a reason for an alternative.** No spend comes from the slow tests, so no penalty applies.
- **It matches your answer on TH-0777 (D-0441)**, where you set 10m on S-0550 for the same kind of cost.

## Why this needs you

S-0721 is in progress, not in the backlog, and has no epic. I give cost of delay inputs only to a backlog story, so I did not set it. Confirming this entry does not write the input. Running the command does.

## The plan: approved

- **The tasks fit the criteria.**
  - Criteria 1 and 2: T-3083.
  - Criterion 3: T-3085, after T-3083.
  - Criterion 4: T-3084.
- **The touches fit.** The new `schema_cache_test.go` keeps off the test files other stories change often. `wip/agents/S-0721.md` carries criterion 3's figures.
- **The forecast fits.** The 30m over flai's 13m states its reason: the timed `-race` runs.
- **The assumptions hold.** T-3083 checks how the SDK caches `permissionInputSchema` before relying on the cache, as criterion 2 asks.

The thread stays open for the cost of delay input.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/digest/D-0441-s-0550-plan-eight-tasks-in-four-layers-cost-of-delay-inputs-needed.md

### 2026-11-11T05:30:53Z alex
Confirmed the recommendation of 2026-11-11T05:30:33Z orchestrator.

Source: design/digest/D-0441-s-0550-plan-eight-tasks-in-four-layers-cost-of-delay-inputs-needed.md

### 2026-11-11T05:31:09Z orchestrator
Resolved: Plan approved; alex confirmed the recommendation and set time_lost_per_cycle 10m on S-0721 at 2026-11-11T05:30:40Z.

### 2026-11-11T05:31:16Z planner-S-0721
The operator set `time_lost_per_cycle: 10m` on S-0721 at 05:30:40Z. `flai cod S-0721` works it out as 25 USD a week, and I wrote that value. Nothing else on this plan waits on an answer, so the thread can be resolved.

### 2026-11-11T05:32:07Z orchestrator
Resolved: Plan approved (T-3083, T-3084, T-3085 in two layers); alex set time_lost_per_cycle 10m on S-0721, and the planner wrote its value of 25 USD a week. Nothing on the thread waits on an answer.
