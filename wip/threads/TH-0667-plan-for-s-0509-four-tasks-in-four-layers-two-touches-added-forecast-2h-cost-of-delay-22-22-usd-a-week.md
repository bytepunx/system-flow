---
id: TH-0667
title: "Plan for S-0509: four tasks in four layers, two touches added, forecast 2h, cost of delay 22.22 USD a week"
anchor:
  path: wip/kanban/stories/S-0509-the-planner-the-orchestrator-the-analyzer-and-the-curator-run-in-a-container-when-their-harness-is-container-with-the-main-checkout-mounted-as-an-adr-settles-pilot-passing-their-role-and-their-activity-logged-as-on-the-host.md
  item: S-0509
status: open
participants: [planner-S-0509]
created: 2026-10-26T12:42:34Z
updated: 2026-10-26T12:42:34Z
---

# TH-0667 Plan for S-0509: four tasks in four layers, two touches added, forecast 2h, cost of delay 22.22 USD a week

On wip/kanban/stories/S-0509-the-planner-the-orchestrator-the-analyzer-and-the-curator-run-in-a-container-when-their-harness-is-container-with-the-main-checkout-mounted-as-an-adr-settles-pilot-passing-their-role-and-their-activity-logged-as-on-the-host.md.

## Entries

### 2026-10-26T12:42:34Z planner-S-0509
S-0509 was planned with E-0029; this run revisited its four tasks against the story's outcome, wrote the cost of delay value the epic's input now allows, and kept the forecast. No task is split, merged, or dropped.

Tasks, order, and layers, each a layer of its own, since each needs what the one before settles:

1. T-2158: the ADR refining ADR-0203 on what a strategic agent's container mounts and the rule for its writes, with a paragraph per role in `design/system/strategic-agents.md`. Waits for nothing.
2. T-2163: `container.go` mounts the main checkout as the ADR says, uses `.flai-cache/agents/<role>`, and pilot's request carries `Role`, `Item`, and `Focus` so the inner harness starts with `--agent`; a test per role on argv and mounts. After T-2158, which settles read-write or read-only.
3. T-2170: `plan.go`, `orchestrate.go`, `analyze.go`, `curate.go` pass the adapter the role's cache folder, `FLAI_ROLE`, and the request; the orchestrator's container is stopped by name; a test per role on a fake adapter. After T-2163, which gives the adapter its role path.
4. T-2175: one planner run in a container plans a story of this repository, logged in `wip/agents/planner.md`; each limitation is an issue; the four documents. After T-2170.

Changed on this run:

- `flai/internal/serve/activity.go` and `claudecheck.go` added to T-2170's and the story's touches: T-2170's Work names both (the reader of a container run's log, the claude probe the ADR settles). Nothing in progress or in review touches them.
- Forecast kept at 2h over flai's 49m, with the reason in the story's `### Planning`; delivery 2026-10-28T06:53Z, flai's playout plus what the 2h adds.
- Cost of delay value 22.22 USD a week, flai's figure: 2h of the 13h30m over E-0029's six open stories, from the 150 USD a week you set on the epic.

Assumptions:

- S-0508 lands `container.go` with a path for a `Request` that carries a `Role`, and S-0507 lands `pilot/run.go`; `after: [S-0508]` stands and S-0507 is reached through it.
- The ADR chooses between read-write at the host path and read-only with writes over MCP; T-2158 recommends read-write, and T-2163 follows whichever is chosen, so the story's agent reads the ADR before T-2163.
- `FLAI_ROLE` reaches the guard inside the container by `-e`, as `runMarks` in `agents.go` carry it to a host run today.
- The proof run in T-2175 is one planner session on a backlog story of this repository, charged as any planner run, and needs the image S-0505 and S-0507 build and Docker on the host.

Considered and not proposed: splitting T-2175's documentation from the proof run. The documentation says what the proof run found, so one task keeps them together; the story's agent may split it when it pulls the story.

Nothing here waits on an answer.
