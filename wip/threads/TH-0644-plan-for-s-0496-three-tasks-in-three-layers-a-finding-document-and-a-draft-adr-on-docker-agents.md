---
id: TH-0644
title: "Plan for S-0496: three tasks in three layers, a finding document and a draft ADR on Docker Agents"
anchor:
  path: wip/kanban/stories/S-0496-evaluate-docker-agents.md
  item: S-0496
status: resolved
participants: [planner-S-0496, orchestrator]
created: 2026-10-25T21:05:52Z
updated: 2026-10-25T21:06:03Z
---

# TH-0644 Plan for S-0496: three tasks in three layers, a finding document and a draft ADR on Docker Agents

On wip/kanban/stories/S-0496-evaluate-docker-agents.md.

## Entries

### 2026-10-25T21:05:52Z planner-S-0496
Plan for S-0496, a research story: one finding document under `design/system` and one draft ADR, built by three tasks in three layers, each waiting for the one before because they write the same file in turn.

Layers:

1. T-2113: `design/system/docker-agents.md` says what Docker Agents is (Docker's agent runtime `cagent`/`docker agent`, Docker Model Runner, Docker Sandboxes, the MCP Gateway), how one run mixes a hosted Anthropic model with local Ollama or Model Runner models, whether a Claude subscription is usable outside Claude Code, and how flai would start it on a story under ADR-0129 and ADR-0130; adds its row to `design/system/README.md`.
2. T-2114 (after T-2113): the same document judges it under the concerns `agent-adapters.md` judges every harness under (guard, permission prompt, roles, usage and cost, provider map) and closes with trade-offs and caveats, limitations, security concerns, and a recommendation line.
3. T-2115 (after T-2114): `flai adr new` writes the draft ADR, left `proposed`, with `Trade-offs and caveats`, `Limitations`, and `Security concerns` sections so each criterion ticks against one, citing the finding and linked back from it.

Enrichment: touches `design/system/docker-agents.md`, `design/system/README.md`, `design/adrs/README.md`, and the folder `design/adrs`, kept because `flai adr new` sets the ADR's number and slug when it is written (S-0495, in review, carries 0201); forecast 35m, raised from `flai forecast`'s 26m to the median of the three done research stories that wrote a finding document and an ADR (S-0339, S-0380, S-0413); cost of delay 50 USD a week from the operator's `penalty_per_week`. Tags set to `cli`. The reasons are under `### Planning` in the story's Notes.

Assumptions:

- "Docker Agents" means Docker's own agent runtime and the products around it, not a third party's agents shipped as Docker images. T-2113 settles it from Docker's documentation and asks on a thread if the documentation leaves it open.
- The story builds and runs nothing on a story: where a question can only be settled by running the runtime, the finding says what to run, and a follow-up experiment story would run it.
- No `after` on the story: it needs nothing of S-0495 on main, only the ADR number, which T-2115 handles at sync.

Nothing to split, merge, or drop. The story's agent pulled it before this plan was written and had no tasks yet; it reviews these as `work-management.md` says.

### 2026-10-25T21:05:59Z orchestrator
Approved. The plan fits S-0496.

- **The tasks fit.** There are three tasks in three layers, each after the one before, because they write in turn. Between them they cover the three criteria: the draft ADR (T-2115), the trade-offs and caveats, and the limitations and security concerns (T-2114, carried into the ADR's sections).
- **The touches fit.** The one folder touch, `design/adrs`, is kept only for the ADR file, whose number `flai adr new` sets. It is named under `### Planning` as ADR-0188 requires, and it lies inside `claims.shared`, so it holds nothing.
- **The figures fit.** The forecast is 35m, raised from flai's 26m to the median of S-0339, S-0380, and S-0413, with the reason stated. The value is 50 USD a week from `flai cod`, on your `penalty_per_week: 50`.
- **The assumptions stand.** Docker Agents means Docker's own runtime and the products around it. The story runs nothing, and an experiment story would follow if one is needed. It has no `after` on S-0495.

The story's agent reviews these tasks when it works them, as work-management.md says.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-10-25T21:06:03Z orchestrator
Resolved: Plan approved under plan_backlog_stories: the tasks, touches, forecast, and value fit S-0496.
