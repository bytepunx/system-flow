---
id: S-0509
type: story
nature: feature
title: The planner, the orchestrator, the analyzer, and the curator run in a container when their harness is container, with the main checkout mounted as an ADR settles, pilot passing their role, and their activity logged as on the host
status: backlog
parent: E-0029
owner: alex
created: 2026-10-26T04:44:48Z
updated: 2026-10-28T23:13:41Z
transitions: []
tags: [cli, docker]
topics: [harness, docker, planning]
touches: [flai/internal/harness/container.go, flai/internal/harness/container_test.go, flai/internal/pilot/run.go, flai/internal/pilot/run_test.go, flai/internal/serve/plan.go, flai/internal/serve/plan_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/analyze.go, flai/internal/serve/analyze_test.go, flai/internal/serve/curate.go, flai/internal/serve/curate_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/activity.go, flai/internal/serve/claudecheck.go, design/adrs/README.md, design/adrs, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, design/issues/summary.md, design/issues]
after: [S-0508]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 190
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 546
          output: 14923
          cache_read: 3320733
          cache_write: 204659
          cost: 8.1773
    - kind: orchestrator
      seconds: 127
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 34
          output: 265
          cache_read: 1884393
          cache_write: 125641
          cost: 0.4414
cost_of_delay:
  value: 22.22
  by: planner-S-0509
  at: 2026-10-26T12:41:56Z
forecast:
  duration: 2h
  delivery: 2026-10-29T15:12:00Z
  basis: "Its own forecast of 2h; 17th in the pull order with an in-progress limit of 5, behind S-0507, S-0535, S-0532, S-0537, S-0534, S-0536, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0424, S-0425, S-0506 and S-0508."
  by: flai
  at: 2026-10-28T23:13:41Z
finalized:
  by: orchestrator
  at: 2026-10-28T08:03:28Z
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

Planned again by planner-S-0509 on 2026-10-26, once the operator had set E-0029's cost of delay input (TH-0657); first planned by planner-E-0029 with the epic, whose four tasks stand.

Touches. Every touch declared before this run is kept. `flai touches suggest S-0509` over the declared paths (1213 of 2427 commits changed them) ranks `docs/users/flai-reference.md`, `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`, `docs/operators/settings.md`, `flai/cmd/serve_actions.go`, `flai/internal/hostapi/writes.go`, and `flai/internal/harness/harness.go` with its test highest (co-change); each is left out, since the story adds no flag, no setting, no host API method, and no dashboard change, and the adapter contract in `harness.go` is S-0507's and S-0508's to change. Where each touch stands:

- `flai/internal/harness/container.go`, `container_test.go`, `flai/internal/pilot/run.go`, `run_test.go` (layout): S-0508's and S-0507's files, not on main yet, which is why `after` names S-0508. T-2163.
- `flai/internal/serve/plan.go`, `orchestrate.go`, `analyze.go`, `curate.go` and their `_test.go` (declared): the four start paths call `harness.For` and `startOf` today. `agents.go` and `agents_test.go` (declared, co-change): `startOf` and the stop of a run. T-2170.
- `flai/internal/serve/activity.go` and `claudecheck.go` (layout, added on this run to T-2170 and the story): T-2170's Work names both, the reader of a container run's log and the claude probe the ADR settles; the epic's plan left them to `flai task done` to widen, but a touch declared before the work is what holds an overlap truthfully.
- `design/system/strategic-agents.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/operators/index.md` (design, criterion 5). T-2158, T-2175.
- `design/adrs/README.md` and the folder `design/adrs` (layout, folder kept): `flai adr new` in T-2158 names the ADR's file, which no task can name; the folder is kept for that file alone, and it is in `claims.shared`, so it holds nothing.
- `design/issues/summary.md` and the folder `design/issues` (layout, folder kept): `flai issue new` in T-2175 writes issue files no task can name; the folder is kept for them alone, and it is in `claims.shared`. `flai task done` widens the touches with the committed paths (ADR-0144).

Forecast. `flai forecast S-0509` gives 49m (median 107 s per unit over 93 done feature stories on claude-opus-5-5 in the large band; size 27: 5 criteria, 22 touches) and delivery 2026-10-28T05:42Z, 27th in the pull order with an in-progress limit of 5. Kept at 2h, as the epic's plan set it: the ADR settles a real choice on the main checkout's mount, four start paths change with a test each, and the proof is a planner session inside a container that the story's agent waits for. Delivery 2026-10-28T06:53Z: flai's playout, plus the 1h11m the 2h adds in the story's own place.

Cost of delay. `flai cod S-0509` gives 22.22 USD a week: the story's share of E-0029's `penalty_per_week` of 150 USD, set by the operator on 2026-10-26 (TH-0657), 2h of the 13h30m forecast over the epic's six open stories without inputs. Written as given: nothing on the story argues for inputs of its own.
