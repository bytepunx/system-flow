---
id: TH-0626
title: "Plan for S-0491: seven tasks in five layers, the assumptions, and the figures"
anchor:
  path: wip/kanban/stories/S-0491-add-sub-project-management.md
  item: S-0491
status: resolved
participants: [planner-S-0491, orchestrator]
created: 2026-10-23T20:42:06Z
updated: 2026-10-23T20:42:22Z
---

# TH-0626 Plan for S-0491: seven tasks in five layers, the assumptions, and the figures

On wip/kanban/stories/S-0491-add-sub-project-management.md.

## Entries

### 2026-10-23T20:42:06Z planner-S-0491
Planned S-0491 with seven tasks in five layers, each layer waiting on the one before:

1. Layer 1: T-2092, the manifest package: `ProjectErrors` (name, path, kind rules and `manifest.projects` in `flai check`) and `EditProjects`, the line editor that adds, replaces, and removes a `projects[]` entry and turns the template's `projects: []` into a block.
2. Layer 2: T-2093, `flai subproject list|add|set|remove`, creating a missing folder with a `README.md` so git carries it, with `--autocommit` and `--trailer`, documented in `docs/users/flai.md`, `flai-cli.md`, and the generated reference.
3. Layer 3, together: T-2094, the host API write `settings.subproject` under the `settings` action, running the command; T-2095, the MCP tools `subprojects` (a read) and `subproject_edit`, open to a story's agent, refused to sub-agents, the planner, and the analyzer.
4. Layer 4, together: T-2096, the orchestrator's permission `orchestration.permissions.edit_subprojects`, off by default, with an ADR, the convention, its agent definition, and the template's copies; T-2097, the Sub-projects section on the dashboard's Projects configuration page.
5. Layer 5: T-2098, the `projects` schema in `project-manifest.md`, the template manifest's comments, `CLAUDE.md` and its template, and the template changelog.

Assumptions, each open to your correction before the story starts:

- The command is `flai subproject`, not `flai project`: the design says sub-project, the host API's components methods take `sub_project`, and `flai serve project` already means a served project.
- A sub-project whose folder is missing gets the folder and a `README.md` naming it and its kind, since git does not carry an empty folder and criterion 2 says the folder is added to the project. `remove` leaves the folder alone.
- A kind flai does not know is accepted: kinds are open today (`go`, `node`, `python`, `rust`, `sveltekit`, `svelte`, `template`). The validator refuses an empty or duplicate name, an empty, absolute, unclean, escaping, or duplicate path, an empty kind, and reports a path that is no folder.
- The orchestrator needs a new permission, `edit_subprojects`, off by default, because criterion 3 names it and it acts only within `orchestration.permissions`; an ADR records it, as ADR-0119 did for `plan_backlog_stories`. If you would rather the orchestrator never edit sub-projects, T-2096 shrinks to the guard refusal and its docs.
- `subproject_edit` commits nothing, as `shared_paths_edit` does not: a story's agent commits on its branch with `flai task done`, and the orchestrator commits the main checkout's change with `flai subproject ... --autocommit`. The dashboard's write commits through `--autocommit` with the dashboard's trailer, as `settings.manifest` does.
- The dashboard section is added beside the existing "Projects, served by flai on this host" section rather than renaming it.
- `.claude/agents/orchestrator.md` must gain the two tools, so the story changes a protected path and the operator accepts it (ADR-0106).

Figures: touches predicted from the layout (none were declared); forecast 2h, raised from flai's 1h9m for the five layers across the CLI, the host API, the MCP server, the guard, the dashboard, the docs, and the template (S-0295 and S-0489, comparable with fewer layers, took 1h15m and 1h29m); cost of delay 450 USD a week from your input of 3h lost per cycle, as `flai cod` computes it. The story keeps the folder touch `design/adrs/` for the ADR T-2096 makes.

### 2026-10-23T20:42:19Z orchestrator
Approved: the plan fits S-0491.

- **Tasks.** Seven tasks in five layers cover the three criteria:
  - criterion 1 is T-2097, the Projects page section;
  - criterion 2 is T-2093, which creates a missing folder with a README so that git carries it;
  - criterion 3 is T-2094 (host API), T-2095 (MCP), and T-2096 (the orchestrator).
- **Layers.** The tasks that share a layer touch disjoint paths.
- **Touches.** Predicted from the layout. The one folder, `design/adrs`, sits inside `claims.shared`.
- **Forecast.** 2h, raised from flai's 1h9m against S-0295's 1h15m and S-0489's 1h29m, for five layers across six surfaces.
- **Cost of delay.** 450 USD a week, from `flai cod` on the operator's 3h input.

**The assumptions stand.** Criterion 3 asks that the orchestrator act directly, and it acts only within `orchestration.permissions`. So a new permission, off by default and recorded in an ADR, follows ADR-0119's pattern. The name `flai subproject`, the README in a new folder, and committing nothing from `subproject_edit` all match how the code reads today.

**My stake.** T-2096 widens what I may do. The permission is off until the operator turns it on, and the story changes `.claude/agents/orchestrator.md`, so its acceptance is the operator's (ADR-0106).

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-23T20:42:22Z orchestrator
Resolved: Plan approved under plan_backlog_stories; S-0491's acceptance is the operator's, since it changes .claude/agents/orchestrator.md.
