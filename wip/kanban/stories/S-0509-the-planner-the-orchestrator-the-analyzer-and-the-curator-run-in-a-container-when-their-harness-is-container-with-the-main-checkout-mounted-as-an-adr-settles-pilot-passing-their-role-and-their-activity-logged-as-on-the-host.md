---
id: S-0509
type: story
nature: feature
title: The planner, the orchestrator, the analyzer, and the curator run in a container when their harness is container, with the main checkout mounted as an ADR settles, pilot passing their role, and their activity logged as on the host
status: backlog
parent: E-0029
owner: alex
created: 2026-10-26T04:44:48Z
updated: 2026-10-26T04:57:43Z
transitions: []
tags: [cli, docker]
topics: [harness, docker, planning]
touches: [flai/internal/harness/container.go, flai/internal/harness/container_test.go, flai/internal/pilot/run.go, flai/internal/pilot/run_test.go, flai/internal/serve/plan.go, flai/internal/serve/plan_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/analyze.go, flai/internal/serve/analyze_test.go, flai/internal/serve/curate.go, flai/internal/serve/curate_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, design/adrs/README.md, design/adrs, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, design/issues/summary.md, design/issues]
after: [S-0508]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 2h
  delivery: 2026-10-27T23:25:00Z
  basis: "Its own forecast of 2h; 28th in the pull order with an in-progress limit of 5, behind S-0498, S-0504, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0499, S-0502, S-0503, S-0505, S-0506, S-0507 and S-0508."
  by: flai
  at: 2026-10-26T04:57:43Z
---
# S-0509 The planner, the orchestrator, the analyzer, and the curator run in a container when their harness is container, with the main checkout mounted as an ADR settles, pilot passing their role, and their activity logged as on the host

## Goal

E-0029 asks that flai start the right image as a container for every kind of agent, the strategic agents among them. A story's agent runs in its worktree, and S-0508 mounts that alone; the planner, the orchestrator, the analyzer, and the curator run in the main checkout (`flai/internal/serve/plan.go`, `orchestrate.go`, `analyze.go`, `curate.go` start them in `entry.Root`), where the planner and the orchestrator write only through flai, the analyzer writes its report under `design/analysis/`, and the curator writes the digest through `flai digest`. ADR-0203 mounts no main checkout, so an ADR must settle what a strategic agent's container mounts: the main checkout read-write at its host path, with the writes it makes being the ones the role's guard allows today and the main checkout's `.git` inside the mount, or the main checkout read-only with every write through the HTTP MCP server, which the analyzer's report and the curator's `flai digest` do not reach today. This story writes that ADR, makes the four start paths take the container adapter with the mounts the ADR names and `FLAI_ROLE` and `--agent` passed through pilot's start request, keeps their activity documents and journal entries as they are, and proves it with one planner run in a container that plans a story of this repository.

## Acceptance criteria

- [ ] An ADR refining ADR-0203 settles the mounts of a strategic agent's container and the rule for its writes, choosing between the main checkout read-write at its host path and read-only with writes over MCP, with the reason, the risk each carries, and what the guard refuses inside the container for each role
- [ ] `plan.go`, `orchestrate.go`, `analyze.go`, and `curate.go` start their agent through `harness.For` as today and, when the harness is `container`, the adapter's `Start` gives `docker run` with the mounts the ADR names for a strategic role in place of a story's worktree, the role's cache folder `.flai-cache/agents/<role>` for the session, and pilot's start request carrying `Role`, `Item`, and `Focus` so the inner harness is started with `--agent <role>` as on the host
- [ ] The guard runs inside the container on each role's calls as it does on the host, holding the planner to planning, the orchestrator to its permissions, the analyzer to its report, and the curator to `flai digest`, with the role read from the request as `flai guard` reads it today, and a test per role checks the argv and the mounts
- [ ] A planner run in a container plans one story of this repository and `wip/agents/planner.md` gains its activity entry with duration and cost measured from the container's log as from a host run; the orchestrator's, the analyzer's, and the curator's start paths are tested with a fake adapter, and each limitation met is an issue under `design/issues`
- [ ] `design/system/strategic-agents.md` (each role's "Its agent" section), `design/system/flai-cli.md`, `docs/users/flai.md`, and `docs/operators/index.md` say how a strategic agent runs in a container and what it mounts

## Tasks
- T-2158 An ADR settles what a strategic agent's container mounts and the rule for its writes
- T-2163 The container adapter mounts a strategic role's checkout and cache folder as the ADR says, and pilot's request carries Role, Item, and Focus so the inner harness starts with --agent
- T-2170 plan.go, orchestrate.go, analyze.go, and curate.go start their agent through the container adapter when its harness is container, the orchestrator's container stopped by name, with a test per role on a fake adapter
- T-2175 A planner run in a container plans one story of this repository with its activity logged, each limitation is an issue, and the docs say how a strategic agent runs in a container

## Notes

Waits on S-0508 for the adapter. The orchestrator is long-running: a container of its own for as long as the `orchestrate` host action is on, stopped by name as a story's is. `flai serve`'s claude check (`claudecheck.go`) reads the claude-code adapter directly; whether it probes the image's `claude` or the host's is for the ADR to say.

### Planning

Planned by planner-E-0029 on 2026-10-26 with the epic.

Touches. `flai touches suggest S-0509`, seeded with `flai/internal/serve/plan.go`, `orchestrate.go`, `analyze.go`, and `curate.go` (16 commits), names `serve/agents.go`, `serve/agents_test.go`, and `plan_test.go`, each kept (co-change), and `cmd/serve_actions.go`, `hostapi/writes.go`, `serve/activity.go`, `serve/serve.go`, `mcpserver/*`, and `flaiover/src/lib/server/agent.ts`, left out because the story adds no setting, no hostapi method, and no dashboard change; `activity.go`'s reader is read through the inner harness's by T-2170 without changing the file unless the task finds it must, in which case `flai task done` widens the touches. The rest are the tasks' (layout):

- `flai/internal/harness/container.go`, `container_test.go`, `flai/internal/pilot/run.go`, `run_test.go`: S-0508's and S-0507's files, which gain the role's mounts and `--agent`. T-2163.
- `flai/internal/serve/orchestrate_test.go`, `analyze_test.go`, `curate_test.go`: the per-role tests beside their start paths. T-2170.
- `design/system/strategic-agents.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/operators/index.md` (design, criterion 5). T-2158, T-2175.
- `design/adrs/README.md` and the folder `design/adrs` (layout, folder kept): `flai adr new` in T-2158 names the ADR's file, which no task can name; the folder is kept for that file alone.
- `design/issues/summary.md` and the folder `design/issues` (layout, folder kept): `flai issue new` in T-2175 writes issue files no task can name; the folder is kept for them alone. `flai task done` widens the touches with the committed paths (ADR-0144).

Forecast and cost of delay: see the second Planning paragraph, written once the figures are set.
