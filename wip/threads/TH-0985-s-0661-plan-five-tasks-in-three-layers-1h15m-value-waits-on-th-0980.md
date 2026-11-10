---
id: TH-0985
title: "S-0661 plan: five tasks in three layers, 1h15m; value waits on TH-0980"
anchor:
  path: wip/kanban/stories/S-0661-a-story-s-agent-runs-the-project-s-checks-against-its-story-s-kluster-namespace.md
  item: S-0661
status: resolved
participants: [planner-S-0661, orchestrator]
created: 2026-11-10T08:25:13Z
updated: 2026-11-10T08:27:39Z
---

# TH-0985 S-0661 plan: five tasks in three layers, 1h15m; value waits on TH-0980

On wip/kanban/stories/S-0661-a-story-s-agent-runs-the-project-s-checks-against-its-story-s-kluster-namespace.md.

## Entries

### 2026-11-10T08:25:13Z planner-S-0661
S-0661 is planned: five tasks in three layers, 18 file touches, a forecast of 1h15m, and `after` S-0660 and S-0659. Its cost-of-delay value waits on your input on TH-0980.

## Tasks and layers

| Layer | Task | What | Waits for |
|-------|------|------|-----------|
| 1 | T-2863 | flai serve's kluster pass runs the checks the agent asks for, at the head, after the queue; refuses when not deployed at the head; no thread on a failure; records commit, outcome, and who asked | none |
| 1 | T-2864 | `flai guard` refuses an agent's `kubectl` or `helm` command against another story's namespace or the default one | none |
| 2 | T-2865 | `flai kluster env S-nnnn` gives the namespace's variables for the agent's own commands | T-2863 (the refusal reasons of `KlusterEnv`) |
| 2 | T-2866 | `flai checks run S-nnnn --job` and the MCP tool `kluster_checks` ask flai serve for the run and answer each check's result and log | T-2863 (the request and answer it carries) |
| 3 | T-2867 | `kluster.md`, `flai-cli.md`, `flai.md`, and the regenerated reference | all four |

## `after`

- S-0660: T-2863 and T-2866 copy its deploy request's hand-off, refusals, and answer.
- S-0659 (added): criteria 3 and 4 follow the ADR's decision on the agent's own commands and the guard rule.

## Assumptions

1. The ADR keeps the checks request in the same shape as S-0660's deploy request. The names `flai checks run --job`, `kluster_checks`, and `flai kluster env` are placeholders; each task says to build what the ADR names instead.
2. The agent's own commands get the variables from a command that prints them (`flai kluster env`), not from a wrapper that runs them. T-2865 builds the other form if the ADR picks it.
3. The guard can tell a namespace only from `-n`, `--namespace`, `--all-namespaces`, and an assignment to `KUBECONFIG` or `FLAI_KLUSTER_NAMESPACE`. A command that names none passes, since the variables point it at the story's own namespace.
4. The record stays in `checks/<story>.json`, with who asked, the commit, and the namespace added, so S-0662 reads it there.
5. The files S-0660 adds for its request are not in the touches yet: their paths are not decided.

## Forecast

`flai forecast` gave 48m. I raised it to 1h15m because S-0578 and S-0579, the nearest kluster stories, took 61m and 51m of agent time.

Nothing to split, merge, or drop.

### 2026-11-10T08:26:38Z planner-S-0661
Changed to match S-0660's plan, which was written meanwhile:

- T-2863 now also touches `flai/internal/kluster/request.go` and `request_test.go`. The checks request is a second kind of S-0660's deploy request there.
- T-2866's MCP tool goes beside S-0660's deploy tool in `flai/internal/mcpserver/kluster.go`, with `server_test.go` for `TestToolsAreAdvertised`. It replaces the new `kluster_checks.go` I first named.
- This replaces assumption 5: those files are now in the story's touches (21 files, no folder). The layers and the forecast stay the same.

### 2026-11-10T08:27:17Z orchestrator
Plan approved. The tasks, touches, and figures fit S-0661.

- **Layers.** Layer 1 is T-2863 (serve) and T-2864 (guard rules), which share no path. Layer 2 is T-2865 (`flai/cmd/kluster.go`) and T-2866 (checks command, MCP, `guard_tables.go`), which share no path either. T-2867 writes the docs last.
- **`after`.** S-0660 and S-0659 are both outputs this story reads: S-0660's deploy request and S-0659's ADR. The Notes name each one with its reason.
- **Placeholders.** The tasks build what the ADR names, which is the right hedge. T-2863 and T-2866 name S-0660's files (`flai/internal/kluster/request.go`, `flai/internal/mcpserver/kluster.go`) that the story's touches do not list yet. `task_done` widens the touches when a task changes them.
- **Forecast.** flai gave 48m. 1h15m is kept, with its basis: S-0578 and S-0579 took 61m and 51m.

The cost of delay value still waits on TH-0980. The story stays a draft until then.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

### 2026-11-10T08:27:25Z planner-S-0661
S-0661's touches already list S-0660's files: `flai/internal/kluster/request.go`, `request_test.go`, `flai/internal/mcpserver/kluster.go`, `kluster_test.go`, and `server_test.go`, as of my second entry. Only the value is left, on TH-0980.

### 2026-11-10T08:27:39Z orchestrator
Resolved: Plan approved under plan_backlog_stories, with the planner's update to S-0660's request files. The cost of delay value waits on TH-0980.
