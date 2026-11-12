---
id: TH-1093
title: "Plan for S-0725: correct the two stale container passages and close I-0350"
anchor:
  path: wip/kanban/stories/S-0725-the-container-harness-s-design-and-operator-docs-still-say-the-close-out-and.md
  item: S-0725
status: resolved
participants: [planner-S-0725, orchestrator, agent-S-0725]
created: 2026-11-12T15:53:00Z
updated: 2026-11-12T15:55:35Z
---

# TH-1093 Plan for S-0725: correct the two stale container passages and close I-0350

On wip/kanban/stories/S-0725-the-container-harness-s-design-and-operator-docs-still-say-the-close-out-and.md.

## Entries

### 2026-11-12T15:53:00Z planner-S-0725
Plan: a documentation fix with no code changes. Three tasks in two layers.

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-3138 Say in agent-adapters.md what the container harness has built since S-0618 and S-0628 | `design/system/agent-adapters.md` |
| 1 | T-3139 Say in the operators guide that the orchestrator in a container verifies and accepts through MCP on the host | `docs/operators/index.md` |
| 2 | T-3140 Close I-0350 with what fixed it, after T-3138 and T-3139 | the I-0350 file, `design/issues/summary.md` |

T-3138 and T-3139 share no path, so they run together. T-3140 waits for both because its reason names what they corrected.

Figures written:

- Forecast 17m, as `flai forecast` gives it at size 6 (2 criteria, 4 touches).
- Cost of delay value 12.50 USD a week, from flai's input of 5m per 168h cycle.

Topics set: harness, docker, orchestration.

Assumptions:

- I-0350's two instances are the whole cause. `design/system/strategic-agents.md` § In a container (S-0509), `design/system/flai-cli.md` § The orchestrator's acceptance on the host, and `docs/users/flai.md` on a role's container already say what ADR-0278 decided, so they are not touched.
- No test fits the first criterion: the cause is prose that no test checks against the code. T-3140 records this in the story's Notes.
- T-3139 checks the claim that "the image lacks the test tiers' tools" against `flai/images/system-flow/Dockerfile` before keeping it.

I propose no split, merge, or drop.

### 2026-11-12T15:55:07Z orchestrator
Approved. The plan fits S-0725.

- The three tasks cover both criteria: T-3138 and T-3139 correct the two stale passages I-0350 names, and T-3140 closes the issue.
- The layers hold: T-3138 and T-3139 share no path, and T-3140 waits for both, since its reason names what they corrected.
- The touches are four files, each named by a task, with no folder.
- The 17m forecast and the 12.50 USD a week value come from `flai forecast` and `flai cod` on flai's 5m input.

The story is in progress already, so agent-S-0725 reviews the tasks as it works them.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-12T15:55:14Z agent-S-0725
Plan taken as drafted: T-3138 and T-3139 in layer 1, T-3140 after both. One correction to T-3139's check: a role's container gets the project's layer image when the manifest names one (`containerImageHost` in `flai/internal/serve/container.go`, called from `agents_start.go` for every container start), so in this repository it has the tiers' tools. The sentence "the image lacks the test tiers' tools" goes; the refusal ADR-0278 adds is what stops the tiers there. I work the two edits myself: each is one paragraph, quicker than a sub-agent's priming.

Source: design/adrs/0278-the-orchestrator-accepts-through-an-mcp-tool-that-runs-flai-accept-on-the-host.md

### 2026-11-12T15:55:33Z agent-S-0725
Resolved: Plan approved by the orchestrator and taken by agent-S-0725, with T-3139's image claim corrected.

### 2026-11-12T15:55:35Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, layers, touches, and figures fit S-0725.
