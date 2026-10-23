---
id: ADR-0195
title: "With edit_subprojects on, the orchestrator adds, changes, and removes a sub-project that a story or a thread asks for, with flai subproject and --autocommit on the main checkout"
status: accepted
date: 2026-10-23
supersedes: []
superseded_by: []
refines: [ADR-0060, ADR-0087]
topics: [cli, template, orchestration]
---

# ADR-0195 With edit_subprojects on, the orchestrator adds, changes, and removes a sub-project that a story or a thread asks for, with flai subproject and --autocommit on the main checkout

## Context

S-0491 gives the manifest's sub-projects, `projects` in `system-flow.yaml`, commands of their own: `flai subproject list`, `add`, `set`, and `remove`, the MCP tools `subprojects` and `subproject_edit`, and the host API write `settings.subproject` behind the dashboard's projects page. A story's agent changes them on its branch, in its worktree, and its close-out commits the change. Every sub-agent is refused them ([ADR-0060](0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md)), and so are the planner and the analyzer. The orchestrator ([ADR-0087](0087-flai-serve-runs-one-orchestrator-per-project-behind-the-orchestrate-host-action.md)) acts only within `orchestration.permissions`, each off by default, and no permission named the sub-projects, so flai guard refused it the edit with no permission to name. The story asks that the orchestrator, as well as a story's agent, can handle a sub-project directly.

Two facts shape how. The orchestrator works the main checkout and has no story, so `subproject_edit`, which commits nothing, would leave the manifest and a new folder's README.md uncommitted on main. And flai guard refuses the orchestrator every git command that writes, so it could not commit them afterwards.

## Decision

`orchestration.permissions.edit_subprojects`, a boolean, off by default, lets the orchestrator add, change, and remove the manifest's sub-projects.

1. **What it allows.** With the permission on, flai guard lets the orchestrator run `flai subproject add`, `set`, and `remove`, and call `subproject_edit`. With it off, it refuses each and names `orchestration.permissions.edit_subprojects` as the permission it needs. Listing them, with `subprojects` or `flai subproject list`, stays a read open to the orchestrator whatever its permissions.
2. **How it commits.** The orchestrator makes the change with the shell, as `flai subproject <op> ... --autocommit --trailer "<trailer>"`, in the main checkout. The command commits `system-flow.yaml` and the README.md of a folder it created, and nothing else, in one step, as the dashboard's `settings.subproject` does. It does not use `subproject_edit`, which leaves the change uncommitted on main where the orchestrator cannot commit it, so its agent definition lists `subprojects` but not `subproject_edit`. The guard still holds `subproject_edit` to the same permission, so a harness that offers it the tool changes nothing about what the operator allowed.
3. **When.** The orchestrator changes a sub-project only when a story or a thread asks for that change, and never on its own judgement: a story whose criteria or notes name the sub-project to add, change, or remove, or a thread that asks for it. It logs each change with `activity_log`, naming the story or thread that asked and the commit, and it answers the thread when one asked.

## Consequences

- An operator who turns the permission on lets the orchestrator change the manifest on main without review, and adding a sub-project creates its folder there. Its risk on the settings page says so.
- A story's agent still makes the change its own story needs, on its branch. The orchestrator's change is for what the work around stories needs: a sub-project a ready story will work in, or one a thread asks for.
- A revert of the one commit reverts the manifest and the folder's README.md together.

## Alternatives considered

- Letting the orchestrator call `subproject_edit` and commit afterwards: flai guard refuses it every git write, and widening that for one tool would let it commit anything on main.
- Giving `subproject_edit` a commit option, or having it commit whenever it writes in the main checkout: a change to the tool for one caller, outside S-0491's plan, and the operator's own session may not want its main-checkout change committed for it.
- Folding the sub-projects into an existing permission such as `plan_backlog_epics`: an operator who allowed planning would silently allow manifest edits.
