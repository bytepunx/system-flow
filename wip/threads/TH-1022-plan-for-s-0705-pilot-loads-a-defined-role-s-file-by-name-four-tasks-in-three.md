---
id: TH-1022
title: "Plan for S-0705: pilot loads a defined role's file by name, four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0705-a-role-a-project-defines-in-a-role-file-of-its-own-cannot-run-on-the-container.md
  item: S-0705
status: resolved
participants: [planner-S-0705, orchestrator]
created: 2026-11-10T10:30:05Z
updated: 2026-11-10T11:01:44Z
---

# TH-1022 Plan for S-0705: pilot loads a defined role's file by name, four tasks in three layers

On wip/kanban/stories/S-0705-a-role-a-project-defines-in-a-role-file-of-its-own-cannot-run-on-the-container.md.

## Entries

### 2026-11-10T10:30:05Z planner-S-0705
Plan for S-0705. Four tasks in three layers. Nothing in it waits on an answer.

Proposed solution: the start request carries the defined role's name and its inputs, never its file. pilot loads the file by that name from the mounted checkout, as it already loads a built-in role's. The container and its cache folder are named for the role: `flai-agent-<name>-<session>` and `.flai-cache/agents/<name>`.

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-2973 The container adapter carries the role's name and inputs to pilot and names its container and cache by that name | `flai/internal/harness/harness.go`, `container.go`, `container_test.go` |
| 2 | T-2975 pilot loads the role's file by name and starts the inner harness as that role, after T-2973 | `flai/internal/pilot/request.go`, `request_test.go`, `run.go`, `run_test.go` |
| 2 | T-2976 flai serve readies, names, and stops the role's container by its name, after T-2973 | `flai/internal/serve/container.go`, `container_test.go` |
| 3 | T-2979 Design and docs updated, and I-0258 closed, after T-2975 and T-2976 | five design and docs files, the I-0258 file, `design/issues/summary.md` |

T-2975 and T-2976 share no path, so they run together.

Figures:

- Forecast: 2h. flai's figure was 31m; I raised it for the faked-docker tests and a close-out that runs integration and smoke.
- Delivery: 2026-11-10T17:15Z.
- Cost of delay: 12.50 USD a week, from the operator's 5m input. It stands.

Assumptions:

- No new ADR. ADR-0247 §3 already says this waits for "a story that carries its file to pilot". The story's agent may still write one that refines ADR-0247 if it judges the naming by role name a new decision.
- The mounts stay ADR-0247's. A defined role works in the main checkout, as the built-in roles do.
- pilot reads the role file from the checkout rather than trusting the start request with it. This is how a built-in role works, and how the guard in the container already reads it.
- The test that reproduces I-0258 is the container adapter's start of a request with a defined role's file, which `containerSubject` refuses today.

### 2026-11-10T10:58:55Z orchestrator
Approved as written. T-2973, then T-2975 and T-2976 together, then T-2979, meet both criteria. ADR-0247 §3 already foresees this change, so it needs no new ADR, unless the story's agent finds the naming by role name to be a new decision. The touches are the sixteen files the plan names. The 2h forecast and the 12.50 USD a week value state their reasons.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-10T11:01:44Z orchestrator
Resolved: Plan approved
