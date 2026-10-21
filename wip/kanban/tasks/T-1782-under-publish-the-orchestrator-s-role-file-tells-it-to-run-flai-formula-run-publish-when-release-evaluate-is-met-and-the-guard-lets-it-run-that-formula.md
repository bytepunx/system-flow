---
id: T-1782
type: task
nature: improvement
title: Under publish, the orchestrator's role file tells it to run flai formula run publish when release_evaluate is met, and the guard lets it run that formula
status: backlog
parent: S-0423
owner: alex
created: 2026-10-16T05:48:37Z
updated: 2026-10-21T13:26:42Z
transitions: []
stream: S-0423
tags: [cli, template]
touches: [design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/cmd/orchestrate_permissions_test.go, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go, flai/cmd/formula.go, flai/cmd/formula_test.go]
after: [T-1780]
---
# T-1782 Under publish, the orchestrator's role file tells it to run flai formula run publish when release_evaluate is met, and the guard lets it run that formula

## Work

In `design/workflows/roles/orchestrator.md`, template and here:

- list `formula run publish` under the `publish` permission;
- tell the orchestrator to run it when `release_evaluate` is met, and to log the instance's ID with `activity_log`;
- keep `release_publish` for a project with no `publish` formula.

Add a case to `orchestrate_permissions_test.go`: the command passes with `publish` on and is refused with it off, and `formula run` of any other formula is refused.

Waits for the formula file. Runs beside the publish step task, since they share no path.

## Done when

`go test ./cmd/...` passes. The existing orchestrator permission cases are unchanged.

## Notes

Replanned 2026-10-21 by planner-S-0423 on TH-0471's answer 2: the role file is `design/workflows/roles/orchestrator.md`, in both copies, as S-0414 writes it, rewritten here as approved on TH-0583; S-0419 reads the orchestrator's prompt from it, so no `.claude/agents` file changes. The guard's role path (S-0416, S-0422's T-2012) allows `formula run <name>` for the formulas a role file names, so `flai/internal/guard/roles.go` is touched in case the `publish` permission on that line needs code. flai should hold the write itself too, as it holds the orchestrator's other writes: `flai formula run` under `FLAI_ROLE=orchestrate` checks `OrchestratorPermits(PermitPublish)` and that the formula is `publish` in `flai/cmd/formula.go`; the story's agent adds that bullet to the Work when it pulls the story, unless TH-0583 settles it first. The orchestrator starts the instance with `flai formula run publish --ready` (T-2012).
