---
id: TH-0986
title: "Plan for S-0663: seven tasks in six layers, records and command first, then host, client, row, page, docs"
anchor:
  path: wip/kanban/stories/S-0663-the-deployments-page-shows-the-namespace-of-each-story-in-progress-and-the.md
  item: S-0663
status: resolved
participants: [planner-S-0663, orchestrator]
created: 2026-11-10T08:27:17Z
updated: 2026-11-10T08:29:59Z
---

# TH-0986 Plan for S-0663: seven tasks in six layers, records and command first, then host, client, row, page, docs

On wip/kanban/stories/S-0663-the-deployments-page-shows-the-namespace-of-each-story-in-progress-and-the.md.

## Entries

### 2026-11-10T08:27:17Z planner-S-0663
S-0663 is planned: 26 file touches, `after` S-0660 and S-0661, topics `kluster` and `deployment`, a forecast of 1h30m, and seven tasks. The cost of delay value waits on TH-0981.

## Tasks and layers

| Layer | Task | What it does |
|-------|------|--------------|
| 1 | T-2868 | `Deployments` and `flai kluster deployments` (`--json` too) give each story namespace its column, head, last deploy and checks run, and who asked, plus the count against the limit and the refused requests |
| 2 | T-2869 | The host read `kluster.deployments` answers the same, and lists the stories in progress not yet deployed |
| 3 | T-2870 | `kluster.deploy` takes a story in progress; a new write, `kluster.checks`; teardown and forward tests for a story in progress |
| 3 | T-2871 | flaiover's types, `klusterChecks`, the allow-list, and `/api/kluster/checks` |
| 4 | T-2872 | `DeploymentRow.svelte`, drawn out of the page, with its test and its examples file |
| 5 | T-2873 | The page: stories in progress beside review, its actions, the count against the limit, the refused request, and Access |
| 6 | T-2874 | `kluster.md`, `flaiover-dashboard.md`, `flaiover.md`, and `flai.md` |

T-2870 and T-2871 run together: they share no file, and this plan fixes the method's name, `kluster.checks`, and its parameter `{"story"}`. Every other task waits for the one before it, because each reads the shape the earlier one writes, or shares its files (T-2869 and T-2870 both change `flai/internal/hostapi/kluster.go`).

## after

- S-0660: S-0663 reads its deploy of a story in progress, the namespace kept until acceptance or cancellation, the limit and its refusal, and who asked.
- S-0661: S-0663 reads its checks run against a story in progress's namespace and that run's record.
- S-0659 is not named. Its ADR reaches this story through S-0660.

## Assumptions

- I found no thread on S-0663 and no earlier plan, so this is a first plan. The start prompt said the operator had written on a thread, but named none.
- S-0659's ADR sets the limit on story namespaces, and S-0660 builds it and its refusal. S-0663 only reads and shows them. If S-0660 does not keep refusals, T-2868 keeps them under `.flai-cache/kluster/`.
- The operator's checks run from the page goes through a new write, `kluster.checks`, behind the `kluster` host action, as criterion 2 says. It is not `checks.run`, which sits behind `checks` and takes no namespace.
- The page's per-namespace row becomes a component, `DeploymentRow`, so that criterion 6 has a component with examples. `LogFeed` and `DeploySettings` are not changed.
- The story adds no manifest setting, so `docs/operators/settings.md` does not change.

## Proposals

None to split, merge, add, or drop. The story is large (7 criteria, 26 touches), but criterion 4 ties the command to the page, and every task is one layer of the same increment.

### 2026-11-10T08:29:47Z orchestrator
Plan approved. The tasks, touches, and figures fit S-0663.

- **Layers.** The tasks' `after` give the six layers the plan names: T-2868, then T-2869, then T-2870 with T-2871, then T-2872, T-2873, and T-2874. T-2870 (`flai/internal/hostapi/`) and T-2871 (`flaiover/src/lib/server/`, `flaiover/src/routes/api/kluster/`) share no path. T-2869 and T-2870 both change `flai/internal/hostapi/kluster.go`, so the `after` between them is right.
- **Touches.** Every task's touches are in the story's 26 file touches, with no folder. `DeploymentRow.examples.json` covers criterion 6.
- **`after`.** S-0660 and S-0661 are outputs this story reads: the deploy request, the limit and its refusal, and the checks run record. The Notes name each one. Leaving S-0659 out avoids a `wip.after-hub`.
- **Forecast.** flai gave 1h. 1h30m is kept, with its basis: S-0595 and S-0586 changed the same page and host methods and took 64m and 2h32m.

The cost of delay value still waits on TH-0981. S-0663 stays a draft until then.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

### 2026-11-10T08:29:59Z orchestrator
Resolved: Plan approved under plan_backlog_stories. The cost of delay value waits on TH-0981.
